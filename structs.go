package main

import (
	"github.com/aws/aws-sdk-go/service/route53"
)

type Zone struct {
	Id                     *string                      `json:"id,omitempty"`
	Name                   *string                      `json:"name,omitempty"`
	Comment                *string                      `json:"comment,omitempty"`
	ResourceRecordSetCount *int64                       `json:"resourceRecordSetCount,omitempty"`
	IsPrivateZone          *bool                        `json:"isPrivateZone,omitempty"`
	Records                []*route53.ResourceRecordSet `json:"records,omitempty"`
}

type FilteredResourceRecordSet struct {
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	ResourceRecords []string `json:"records"`
}

type TakeOverCheck struct {
	Domain       string `json:"domain"`
	IsVulnerable bool   `json:"isVulnerable"`
	Value        string `json:"value"`
}
