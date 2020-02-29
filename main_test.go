package main

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest"
)

func TestBetweenTime(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("00:00->23:59", now)
	if got != true {
		t.Errorf("Should be within time")
	}
}

func TestOutsideTime(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("00:00->15:00", now)
	if got == true {
		t.Errorf("Should be whithout time")
	}
}

func TestOutsideTimeInsideMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("00:00->15:00;January", now)
	if got != true {
		t.Errorf("Should be whithin time")
	}
}

func TestDayOfTheWeek(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("Wednesday", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestOutsideDayOfTheWeek(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("Thursday", now)
	if got == true {
		t.Errorf("Should be correct for date")
	}
}

func TestDayAndMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("January 01", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestOutsideDayAndMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("January 02", now)
	if got == true {
		t.Errorf("Should be correct for date")
	}
}

func TestMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("January", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestMonthMultipleLastCorrect(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February;January", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestMonthMultipleFirstCorrect(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February;January", now)
	if got != true {
		t.Errorf("Should be correct for Month")
	}
}

func TestOutsideMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February", now)
	if got == true {
		t.Errorf("Should be incorrect for Month")
	}
}

func TestInsideCustomWeekDay(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February", now)
	if got == true {
		t.Errorf("Should be correct for date")
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_getSubscriptions(t *testing.T) {
	type args struct {
		auth autorest.Authorizer
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSubscriptions(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSubscriptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evaluateStatus(t *testing.T) {
	type args struct {
		auth         autorest.Authorizer
		subscription string
		wg           *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluateStatus(tt.args.auth, tt.args.subscription, tt.args.wg)
		})
	}
}

func Test_getResource(t *testing.T) {
	type args struct {
		resource string
	}
	tests := []struct {
		name string
		args args
		want *AzureResource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResource(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isOn(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOn(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("isOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isOnStatuses(t *testing.T) {
	type args struct {
		statuses *[]compute.InstanceViewStatus
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOnStatuses(tt.args.statuses); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("isOnStatuses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isWithinTime(t *testing.T) {
	type args struct {
		schedule string
		now      time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWithinTime(tt.args.schedule, tt.args.now); got != tt.want {
				t.Errorf("isWithinTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evaluateTimeRange(t *testing.T) {
	type args struct {
		timeSchedule string
		now          time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateTimeRange(tt.args.timeSchedule, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateTimeRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("evaluateTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newAuthorizer(t *testing.T) {
	tests := []struct {
		name    string
		want    *autorest.Authorizer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newAuthorizer()
			if (err != nil) != tt.wantErr {
				t.Errorf("newAuthorizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAuthorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}
