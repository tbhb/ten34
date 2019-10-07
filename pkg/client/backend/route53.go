package backend

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Route53 implements a key-value store on top of AWS Route53 hosted zones and DNS
type Route53 struct {
	// HostedZoneID is the ID of the Route53 hosted zone
	HostedZoneID string

	// HosteddZoneName is the name of the Route53 hosted zone
	HostedZoneName string

	// URI is the uniform resource identifier of the database
	URI url.URL

	client *route53.Route53
}

// NewRoute53 creates a new Route53 database backend session
func NewRoute53(uri url.URL) *Route53 {
	backend := &Route53{
		HostedZoneID:   "",
		HostedZoneName: "",
		URI:            uri,
	}
	return backend
}

// Setup sets up a backend session
func (be *Route53) Setup(uri url.URL) error {
	be.HostedZoneName = fmt.Sprintf("%s.", be.URI.Host)

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session")
		return err
	}

	be.client = route53.New(sess)

	return nil
}

// CreateDB creates a database
func (be *Route53) CreateDB(db url.URL) error {
	listInput := &route53.ListHostedZonesByNameInput{}
	resp, err := be.client.ListHostedZonesByName(listInput)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, v := range resp.HostedZones {
		hzn := *v.Name
		if hzn == be.HostedZoneName {
			id := *v.Id
			id = strings.TrimPrefix(id, "/hostedzone/")
			return fmt.Errorf("database already exists: %s", id)
		}
	}

	t := time.Now().Unix()
	callerRef := fmt.Sprintf("%v", t)
	createInput := &route53.CreateHostedZoneInput{
		CallerReference: aws.String(callerRef),
		Name:            aws.String(be.HostedZoneName),
	}
	_, err = be.client.CreateHostedZone(createInput)
	if err != nil {
		return fmt.Errorf("failed to create database: %s", be.URI.String())
	}

	return nil
}

// DropDB deletes a database
func (be *Route53) DropDB(db url.URL) error {
	be.setHostedZoneID()
	if be.HostedZoneID == "" {
		return fmt.Errorf("database does not exist: %s", be.URI.Host)
	}

	input := &route53.DeleteHostedZoneInput{
		Id: aws.String(be.HostedZoneID),
	}
	_, err := be.client.DeleteHostedZone(input)
	if err != nil {
		return fmt.Errorf("failed to delete database: %s", be.URI.String())
	}

	return nil
}

// Delete deletes a key from the database
func (be *Route53) Delete(db url.URL, key string) error {
	be.setHostedZoneID()
	if be.HostedZoneID == "" {
		return fmt.Errorf("database does not exist: %s", be.URI.Host)
	}

	dnsKey := fmt.Sprintf("%s.%s", key, be.HostedZoneName)

	listInput := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(be.HostedZoneID),
		StartRecordName: aws.String(dnsKey),
		StartRecordType: aws.String("TXT"),
		MaxItems:        aws.String("1"),
	}
	resp, err := be.client.ListResourceRecordSets(listInput)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	if len(resp.ResourceRecordSets) == 0 {
		return fmt.Errorf("key does not exist: %s", key)
	}

	deleteInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(dnsKey),
						ResourceRecords: resp.ResourceRecordSets[0].ResourceRecords,
						TTL:             aws.Int64(60),
						Type:            aws.String("TXT"),
					},
				},
			},
			Comment: aws.String(fmt.Sprintf("Deleting key %s", key)),
		},
		HostedZoneId: aws.String(be.HostedZoneID),
	}

	_, err = be.client.ChangeResourceRecordSets(deleteInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}
	return nil
}

// Get retrieves keys from the database
func (be *Route53) Get(db url.URL, key string) (string, error) {
	be.setHostedZoneID()

	resolvable := true
	lookupHost := strings.TrimSuffix(be.HostedZoneName, ".")
	_, err := net.LookupNS(lookupHost)
	if err != nil {
		resolvable = false
	}

	dnsKey := fmt.Sprintf("%s.%s", key, be.HostedZoneName)

	if resolvable {
		fmt.Println("DNS is resolvable")
		vals, err := net.LookupTXT(dnsKey)
		if err != nil || len(vals) == 0 {
			return "", fmt.Errorf("key not found: %s", key)
		}
		return vals[0], nil
	}

	fmt.Println("DNS is not resolvable, falling back to ListResourceRecordSets")

	listInput := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(be.HostedZoneID),
		StartRecordName: aws.String(dnsKey),
		StartRecordType: aws.String("TXT"),
		MaxItems:        aws.String("1"),
	}
	resp, err := be.client.ListResourceRecordSets(listInput)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}

	if len(resp.ResourceRecordSets) == 0 {
		return "", fmt.Errorf("key does not exist: %s", key)
	}

	val := *resp.ResourceRecordSets[0].ResourceRecords[0].Value
	val = strings.TrimPrefix(val, "\"")
	val = strings.TrimSuffix(val, "\"")
	return val, nil
}

// Put writes a key-value pair to the database
func (be *Route53) Put(db url.URL, key, val string) error {
	be.setHostedZoneID()
	if be.HostedZoneID == "" {
		return fmt.Errorf("database does not exist: %s", be.URI.Host)
	}

	dnsKey := fmt.Sprintf("%s.%s", key, be.HostedZoneName)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(dnsKey),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(fmt.Sprintf("\"%s\"", val)),
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("TXT"),
					},
				},
			},
			Comment: aws.String(fmt.Sprintf("Updating key %s", key)),
		},
		HostedZoneId: aws.String(be.HostedZoneID),
	}

	_, err := be.client.ChangeResourceRecordSets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case route53.ErrCodeNoSuchHostedZone:
				fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
			case route53.ErrCodeNoSuchHealthCheck:
				fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
			case route53.ErrCodeInvalidChangeBatch:
				fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
			case route53.ErrCodeInvalidInput:
				fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
			case route53.ErrCodePriorRequestNotComplete:
				fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}
	return nil
}

func (be *Route53) setHostedZoneID() error {
	input := &route53.ListHostedZonesByNameInput{}
	resp, err := be.client.ListHostedZonesByName(input)

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, v := range resp.HostedZones {
		hzn := *v.Name
		if hzn == be.HostedZoneName {
			id := *v.Id
			be.HostedZoneID = strings.TrimPrefix(id, "/hostedzone/")
			return nil
		}
	}
	return nil
}
