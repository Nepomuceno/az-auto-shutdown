package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// ActionType Which action should be takem
type ActionType int32

const (
	// TURNOFF Shutdown the target vm
	TURNOFF ActionType = 0
	// START Start the target vm
	START ActionType = 1
	// DONOTHING do nothing
	DONOTHING ActionType = 3
)

// AzureResource Parts of azure resource identification
type AzureResource struct {
	Subscription  string
	ResourceGroup string
	Provider      string
	Type          string
	Name          string
	SubType       string
	SubName       string
}

var (
	resouceIDRegex   = regexp.MustCompile(`(?m)\/subscriptions\/(?P<subscription>[^\/]+)\/resourceGroups\/(?P<resourceGroup>[^\/]+)\/providers\/(?P<resourceProvider>[^\/]+)\/(?P<resouceType>[^\/]+)\/(?P<resourceName>[^\/]+)(\/)?(?P<resourceSubtype>[^\/]+)?(\/)?(?P<resourceSubtypeName>[^\/]+)?$`)
	weekdayFunctions = map[string]func(now time.Time) bool{
		"Weekdays": func(now time.Time) bool {
			return now.Weekday() != time.Sunday && now.Weekday() != time.Saturday
		},
		"Weekends": func(now time.Time) bool {
			return (now.Weekday() == time.Sunday || now.Weekday() == time.Saturday)
		},
		"WorkingHours": func(now time.Time) bool {
			start, _ := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", now.Format("2006-01-02"), "06:00"))
			end, _ := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", now.Format("2006-01-02"), "20:00"))
			return (now.Before(end) && now.After(start))
		},
	}
)

func main() {
	// create an authorizer from env vars or Azure Managed Service Idenity
	log.Println("Starting app Press CTRL+C to end.")
	authorizer, err := newAuthorizer()
	if err != nil || authorizer == nil {
		log.Fatalln("Impossible to authenticate")
	}
	var interval = 300
	intervalSrt, intervalConfigured := os.LookupEnv("CHECK_SECONDS_INTERVAL")
	if intervalConfigured {
		interval, err = strconv.Atoi(intervalSrt)
		if err != nil {
			log.Println("CHECK_SECONDS_INTERVAL is not a valid integer")
			interval = 300
		}
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				subs, err := getSubscriptions(*authorizer)
				if err != nil {
					log.Panic(err)
				}
				var wg sync.WaitGroup
				wg.Add(len(subs))
				for _, sub := range subs {
					go evaluateStatus(*authorizer, sub, &wg)
				}
				wg.Wait()
				fmt.Println("Tick at", t)
			}
		}
	}()
	<-done
	log.Println("End of schedule")
}

func getSubscriptions(auth autorest.Authorizer) ([]string, error) {
	var subs []string
	client := subscription.NewSubscriptionsClient()
	client.Authorizer = auth
	result, err := client.ListComplete(context.Background())
	if err != nil {
		return nil, err
	}
	for result.NotDone() {
		subs = append(subs, *result.Value().SubscriptionID)
		result.Next()
	}
	return subs, nil
}

func evaluateStatus(auth autorest.Authorizer, subscription string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Evaluating status for: %s", subscription)
	resourceClient := resources.NewClient(subscription)
	computeClient := compute.NewVirtualMachinesClient(subscription)
	computeClient.Authorizer = auth
	resourceClient.Authorizer = auth
	listResources, err := resourceClient.ListComplete(context.Background(), "resourceType eq 'Microsoft.Compute/virtualMachines'", "", to.Int32Ptr(1000))
	if err != nil {
		log.Fatal(err)
	}
	for listResources.NotDone() {
		res := listResources.Value()
		resID := getResource(*res.ID)
		if res.Tags["AutoShutDown-OFF"] != nil || res.Tags["AutoShutDown-ON"] != nil {
			vm, err := computeClient.Get(context.Background(), resID.ResourceGroup, resID.Name, compute.InstanceView)
			if err != nil {
				log.Println(err)
			} else {
				isOn := isOnStatuses(vm.InstanceView.Statuses)
				if isOn == nil {
					log.Printf("Machine transitioning: %s/%s", resID.ResourceGroup, resID.Name)
				} else {
					if res.Tags["AutoShutDown-OFF"] != nil {
						insideSchedule := isWithinTime(*res.Tags["AutoShutDown-OFF"], time.Now())
						if insideSchedule && *isOn {
							_, err := computeClient.Deallocate(context.Background(), resID.ResourceGroup, resID.Name)
							log.Printf("VM: %s/%s DEALLOCATED", resID.ResourceGroup, resID.Name)
							if err != nil {
								log.Println(err)
							}
						}
					} else if res.Tags["AutoShutDown-ON"] != nil {
						insideSchedule := isWithinTime(*res.Tags["AutoShutDown-ON"], time.Now())
						if insideSchedule && !*isOn {
							_, err := computeClient.Start(context.Background(), resID.ResourceGroup, resID.Name)
							log.Printf("VM: %s/%s STARTED", resID.ResourceGroup, resID.Name)
							if err != nil {
								log.Println(err)
							}
						}
					}
				}
			}
		} else {
			log.Printf("Tag not found for: %s/%s", resID.ResourceGroup, resID.Name)
		}
		listResources.Next()
	}

}

func getResource(resource string) *AzureResource {
	matches := resouceIDRegex.FindStringSubmatch(resource)
	result := &AzureResource{}
	if len(matches) > 1 {
		result.Subscription = matches[1]
		result.ResourceGroup = matches[2]
		result.Provider = matches[3]
		result.Type = matches[4]
		result.Name = matches[5]
		result.SubType = matches[7]
		result.SubName = matches[9]
	}
	return result
}

func isOn(status string) *bool {
	switch status {
	case "PowerState/starting":
		return nil
	case "PowerState/running":
		return to.BoolPtr(true)
	case "PowerState/stopping":
		return nil
	case "PowerState/stopped":
		return to.BoolPtr(false)
	case "PowerState/deallocating":
		return nil
	case "PowerState/deallocated":
		return to.BoolPtr(false)
	default:
		return nil
	}
}

func isOnStatuses(statuses *[]compute.InstanceViewStatus) *bool {
	for _, status := range *statuses {
		if strings.Contains(*status.Code, "PowerState") {
			return isOn(*status.Code)
		}
	}
	return nil
}

func isWithinTime(schedule string, now time.Time) bool {
	schedules := strings.Split(schedule, ";")
	result := false
	for _, timeSchedule := range schedules {
		timeSchedule = strings.Trim(timeSchedule, " \t\n")
		if strings.Contains(timeSchedule, "->") {
			isIn, err := evaluateTimeRange(timeSchedule, now)
			if err != nil {
				continue
			}
			result = result || isIn
		}
		if timeSchedule == now.Weekday().String() {
			result = true
		}
		if timeSchedule == now.Month().String() {
			result = true
		}
		if len(strings.Split(timeSchedule, " ")) == 2 {
			paths := strings.Split(timeSchedule, " ")
			if paths[0] == now.Month().String() {
				if paths[1] == now.Format("02") {
					result = true
				}
			}
		}
		if f, exists := weekdayFunctions[timeSchedule]; exists {
			result = result || f(now)
		}
	}
	return result
}

func evaluateTimeRange(timeSchedule string, now time.Time) (bool, error) {
	timerange := strings.Split(timeSchedule, "->")
	if len(timerange) == 2 {
		start, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", now.Format("2006-01-02"), timerange[0]))
		if err != nil {
			log.Printf("Incorrect time range %s", timeSchedule)
			return false, err
		}
		end, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", now.Format("2006-01-02"), timerange[1]))
		if err != nil {
			log.Printf("Incorrect time range %s", timeSchedule)
			return false, err
		}
		if end.Before(start) {
			end = end.AddDate(0, 0, 1)
		}
		if now.Before(end) && now.After(start) {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("Incorrect time range %s", timeSchedule)

}

func newAuthorizer() (*autorest.Authorizer, error) {
	// Carry out env var lookups
	_, clientIDExists := os.LookupEnv("AZURE_CLIENT_ID")
	_, tenantIDExists := os.LookupEnv("AZURE_TENANT_ID")
	_, fileAuthSet := os.LookupEnv("AZURE_AUTH_LOCATION")

	// Execute logic to return an authorizer from the correct method
	if clientIDExists && tenantIDExists {
		log.Println("Logging from environment")
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		return &authorizer, err
	} else if fileAuthSet {
		log.Println("Logging from file")
		authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
		return &authorizer, err
	} else {
		log.Println("Logging from CLI")
		authorizer, err := auth.NewAuthorizerFromCLI()
		return &authorizer, err
	}
}
