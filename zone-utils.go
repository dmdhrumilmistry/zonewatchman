package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func listHostedZones(sess *session.Session) ([]*route53.HostedZone, error) {
	svc := route53.New(sess)
	input := &route53.ListHostedZonesInput{}

	result, err := svc.ListHostedZones(input)
	if err != nil {
		return nil, err
	}

	return result.HostedZones, nil
}

func listResourceRecordSets(sess *session.Session, hostedZoneID string) ([]*route53.ResourceRecordSet, error) {
	svc := route53.New(sess)
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	result, err := svc.ListResourceRecordSets(input)
	if err != nil {
		return nil, err
	}

	return result.ResourceRecordSets, nil
}

func getAllZoneDomains(sess *session.Session) ([]*Zone, error) {
	// List all hosted zones
	hostedZones, err := listHostedZones(sess)
	if err != nil {
		fmt.Println("Error listing hosted zones:", err)
		os.Exit(1)
		return make([]*Zone, 0), fmt.Errorf("Failed to list hosted zones")
	}

	var zonesData []*Zone

	// append zone data into slice
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, zoneData := range hostedZones {
		// List resource record sets for each hosted zone
		wg.Add(1)

		go func(zone *route53.HostedZone) {
			defer wg.Done()
			resourceRecordSets, err := listResourceRecordSets(sess, aws.StringValue(zone.Id))
			if err != nil {
				fmt.Println("Error listing resource record sets:", err)
				os.Exit(1)
			}

			mu.Lock()
			defer mu.Unlock()

			zonesData = append(zonesData, &Zone{
				Id:                     zone.Id,
				Name:                   zone.Name,
				Comment:                zone.Config.Comment,
				IsPrivateZone:          zone.Config.PrivateZone,
				ResourceRecordSetCount: zone.ResourceRecordSetCount,
				Records:                resourceRecordSets,
			})
		}(zoneData)
	}

	wg.Wait()

	return zonesData, nil
}

func getPublicZones(zones []*Zone) []*Zone {
	var publicZones []*Zone

	for _, zone := range zones {
		// Check if IsPrivateZone is false
		if zone.IsPrivateZone != nil && !*zone.IsPrivateZone {
			publicZones = append(publicZones, zone)
		}
	}
	return publicZones
}

func filterRecordSetByType(resourceRecordSet []*route53.ResourceRecordSet, recordTypes []string) []*FilteredResourceRecordSet {
	var filteredZoneRecordSet []*FilteredResourceRecordSet

	for _, resourceSet := range resourceRecordSet {
		if resourceSet.Type != nil && SearchStringInSlice(recordTypes, *resourceSet.Type) {
			var resourceRecords []string
			for _, resourceRecord := range resourceSet.ResourceRecords {
				resourceRecords = append(resourceRecords, *resourceRecord.Value)
			}

			filteredZoneRecordSet = append(filteredZoneRecordSet, &FilteredResourceRecordSet{
				Name:            *resourceSet.Name,
				Type:            *resourceSet.Type,
				ResourceRecords: resourceRecords,
			})

		}
	}

	return filteredZoneRecordSet
}
