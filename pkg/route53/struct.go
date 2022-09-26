package route53

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"log"
)

var domain = "opssec.in."

func GetConfig() aws.Config {
	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}
	return cfg
}

func GetRoute53Service(cfg aws.Config) *route53.Client {
	return route53.NewFromConfig(cfg)
}

func GetRecord(svc *route53.Client, hostedZone, record string) bool {
	out, err := svc.ListResourceRecordSets(
		context.Background(),
		&route53.ListResourceRecordSetsInput{HostedZoneId: aws.String(hostedZone)},
	)
	if err != nil {
		return false
	}
	for _, n := range out.ResourceRecordSets {
		if n.Type == types.RRTypeTxt {
			if *n.Name == record {
				return true
			}
		}
	}
	return false
}

func SaveRecord(svc *route53.Client, hostedZone, record string, values []string) types.ChangeStatus {
	var rrs []types.ResourceRecord
	for _, v := range values {
		rrs = append(rrs, types.ResourceRecord{Value: aws.String(v)})
	}

	action := types.ChangeActionUpsert
	if !GetRecord(svc, hostedZone, record) {
		action = types.ChangeActionCreate
	}

	changeBatch := types.ChangeBatch{
		Changes: []types.Change{
			{
				Action: action,
				ResourceRecordSet: &types.ResourceRecordSet{
					TTL:             aws.Int64(1),
					Name:            aws.String(record),
					Type:            types.RRTypeTxt,
					ResourceRecords: rrs,
				},
			},
		},
	}

	out, err := svc.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZone),
		ChangeBatch:  &changeBatch,
	})
	if err != nil {
		log.Fatal(err)
	}
	return out.ChangeInfo.Status
}
