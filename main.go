package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type Zone struct {
	Id                     string                       `json:"id,omitempty"`
	Name                   string                       `json:"name,omitempty"`
	Comment                string                       `json:"comment,omitempty"`
	ResourceRecordSetCount int64                        `json:"resourceRecordSetCount,omitempty"`
	IsPrivateZone          bool                         `json:"isPrivateZone,omitempty"`
	Records                []*route53.ResourceRecordSet `json:"records,omitempty"`
}

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

	zonesData := make([]*Zone, 0)

	// append zone data into slice
	for _, zone := range hostedZones {
		// List resource record sets for each hosted zone
		resourceRecordSets, err := listResourceRecordSets(sess, aws.StringValue(zone.Id))
		if err != nil {
			fmt.Println("Error listing resource record sets:", err)
			os.Exit(1)
		}

		zonesData = append(zonesData, &Zone{
			Id:                     *zone.Id,
			Name:                   *zone.Name,
			Comment:                *zone.Config.Comment,
			IsPrivateZone:          *zone.Config.PrivateZone,
			ResourceRecordSetCount: *zone.ResourceRecordSetCount,
			Records:                resourceRecordSets,
		})

	}

	return zonesData, nil
}

func main() {
	// Create a session with default credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// get all zone domains
	hostedDomains, err := getAllZoneDomains(sess)
	if err != nil {
		fmt.Println("Error listing resource record sets:", err)
		os.Exit(1)
	}
	for _, domain := range hostedDomains {
		fmt.Println(domain)
	}
}
