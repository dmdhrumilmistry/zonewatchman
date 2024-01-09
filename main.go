package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	// Create a session with default credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	recordSetTypes := []string{"A"}

	// get all zone domains
	hostedZones, err := getAllZoneDomains(sess)
	if err != nil {
		fmt.Println("Error listing resource record sets:", err)
		os.Exit(1)
	}

	publicHostedZones := getPublicZones(hostedZones)
	for _, publicHostedZone := range publicHostedZones {
		fmt.Println(*publicHostedZone.Name)
		fmt.Println("=====================================================")
		filteredRecordSet := filterRecordSetByType(publicHostedZone.Records, recordSetTypes)
		results := testDnsAtypeEntries(filteredRecordSet)

		// fmt.Println(results)
		for _, result := range results {
			fmt.Println(result)
		}

		fmt.Println("=====================================================")
		fmt.Println("=====================================================")
	}

}
