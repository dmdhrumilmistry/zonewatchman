package main

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getElasticIps(sess *session.Session) []*string {
	var ips []*string
	// Create EC2 client
	client := ec2.New(sess)

	// Get a list of AWS regions

	regions, err := getRegions(client)
	if err != nil {
		fmt.Println("[!] Error getting regions:", err)
		return nil
	}

	// Iterate over regions
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, region := range regions {
		wg.Add(1)

		go func(sess *session.Session, region string) {
			defer wg.Done()
			// Describe Elastic IP addresses in the region
			regionIps, err := describeElasticIPs(sess, region)

			if err != nil {
				fmt.Println("[!] Error describing Elastic IPs:", err)
			}

			// acquire lock to avoid race conditions
			mu.Lock()

			// Print the public IPs
			ips = append(ips, regionIps...)
			mu.Unlock()

		}(sess, region)
	}
	wg.Wait()

	return ips
}

// getRegions retrieves a list of AWS regions
func getRegions(client *ec2.EC2) ([]string, error) {
	resp, err := client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, region := range resp.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

// describeElasticIPs retrieves Elastic IP addresses in the specified region
func describeElasticIPs(sess *session.Session, region string) ([]*string, error) {
	client := ec2.New(sess, aws.NewConfig().WithRegion(region))

	resp, err := client.DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	var ips []*string
	for _, address := range resp.Addresses {
		ips = append(ips, address.PublicIp)
	}

	return ips, nil
}

func IsAwsIp(awsIps []*string, ip string) bool {
	for _, awsIp := range awsIps {
		if *awsIp == ip {
			return true
		}
	}

	return false
}
