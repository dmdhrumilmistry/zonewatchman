package main

import (
	"fmt"
	"sync"
)

func testDnsAtypeEntries(resourceSet []*FilteredResourceRecordSet, awsElasticIps []*string) []*TakeOverCheck {
	var checks []*TakeOverCheck
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, resourceRecord := range resourceSet {
		for _, ip := range resourceRecord.ResourceRecords {
			wg.Add(1)

			go func(resourceRecord *FilteredResourceRecordSet, ip string) {
				defer wg.Done()
				err := makeGetRequest(false, ip, 80) // http
				if err != nil {
					fmt.Printf("[!!] %s Host is not reachable\n", ip)
				}

				IsVulnerable := !IsAwsIp(awsElasticIps, ip)

				mu.Lock()
				defer mu.Unlock()

				checks = append(checks, &TakeOverCheck{
					Domain:       resourceRecord.Name,
					IsVulnerable: IsVulnerable,
					Value:        ip,
				})
			}(resourceRecord, ip)
		}
	}
	wg.Wait()

	return checks
}
