# Azure Auto Shutdown

The goal of the project it is to create a tool that automaticaly spins up and shutdown your Azure VMs in a set of scenarios.

## How to use it

### Tag your VMS

The first step to use `az-auto-shutdown` it is to add the tags in your VMs you can add the tags:

- **AutoShutDown-OFF**: To specify a time where you want your VMs to be **OFF**. Outside of this time they will be turned ON.
- **AutoShutDown-ON**: To specify a time where you want your VMs to be **ON**. Outside of this time they will be turned OFF.

You can tag you machines with a series of values separated by ; if any of this values is match the action will be taken.

The values suported for the tag are 

- Time: `O1:02->15:04` this will match against any time between 1:02AM and 3:04PM Times are always in 24h format.
- Weekday: `Thursday` this will match against every time at Thursday. Weekdays need to be in English
- Month: `February` This will match against every day in February. Months need to be in English
- Date: `December 25` This will match agaist any time of the day December 25

You can combine the tags for example and `AutoShutDown-OFF` tag with the value `20:00->06:00;Saturday;Sunday;December 25` for a VM that will be off from 8PM to 6AM every day and off every weekend plus off also on December 25


## Container Configuration

### Required Settings

You are required to pass your azure authentication as environment variables to the container. The application will check all the subscription that that user has access to.

You van find mor details on how to use the environment variables to authenticate on azure on [Azure Docs](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization)

A sample configuration would have the following env variables

```
AZURE_TENANT_ID=YOUR_TENAT_ID
AZURE_CLIENT_ID=YOUR_SERVICE_PRINCIPAL_ID
AZURE_CLIENT_SECRET=YOUR_SERVICE_PRINCIPAL_PASSWORD
```

### Optional Settings

By default the VMs are checked every 5 minutes but you can make that more or less frequent dependent on your requirementes changing the variables `CHECK_SECONDS_INTERVAL` with the interval that makes sense for you.


## Advanced Tag Usages

There are the normal tag usages discussed at [Tag Your Vms](#tag-your-vms) but you can have way more advanced Tags on your VMs.

First how the Tags do work:

`AutoShutDown-ON` will be evaluated only if there is no `AutoShutDown-OFF` tag on the VM. Meaning that you can use a combination of multiple tags acress your vms but if you use ON and OFF on the same VM only the OFF configuratioon will be evaluated. This is done to avoid


### Sample Scenarios

Description | Tag Key | Tag value
----------- | ----------- | -----------
Shut down from 10PM to 6 AM UTC every day (different format, same result as above) | AutoShutDown-OFF | 22:00->06:00
Shut down from 8PM to 12AM and from 2AM to 7AM UTC every day (bringing online from 12-2AM for maintenance in between) | AutoShutDown-OFF | 20:00->00:00;02:00->07:00
Shut down all day Saturday and Sunday (midnight to midnight) | AutoShutDown-OFF | Saturday;Sunday
Shut down from 2AM to 7AM UTC every day and all day on weekends | AutoShutDown-OFF | 02:00->07:00;Saturday;Sunday
Shut down on Christmas Day and New Year’s Day | AutoShutDown-OFF | December 25, January 1 
Shut down from 2AM to 7AM UTC every day,and all day on weekends, and on Christmas Day | AutoShutDown-OFF | 02:00->07:00;Saturday;Sunday;December 25
Shut down always – I don’t want this VM online, ever | AutoShutDown-OFF | 00:00->23:59