package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

func main() {
	sleep, _ := strconv.Atoi(os.Getenv("SLEEP"))
	if sleep == 0 {
		sleep = 5
	}
	sess := session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
	}))

	svc := sts.New(sess)
	s3svc := s3.New(sess)
	for {
		result, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				fmt.Println(aerr.Error())
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		}
		fmt.Println(result)

		input := &s3.ListObjectsInput{
			Bucket:  aws.String("exp-1-sys-prom-prometheus-storage"),
			MaxKeys: aws.Int64(1),
		}

		s3result, err := s3svc.ListObjects(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		}
		fmt.Println(s3result)
		time.Sleep(time.Duration(sleep*1000) * time.Millisecond)
	}
}
