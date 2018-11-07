/*
   Copyright Ken
*/

package main

import (
    "fmt"
    "log"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/autoscaling"
    "github.com/aws/aws-sdk-go/aws/ec2metadata"
)

func main() {
    // Create sessiontly filter on a specific region or use env vars to do it
    sess, err := session.NewSession()
    if err != nil {
        log.Fatal(err.Error())
    }

    /*
     * 1. Create session with the specific AWS account credentials
     * (defaults to current account if not specified; this requires
     * the 'autoscaling:Describe*' action to be specified in the IAM role
     * for the blue master instance)
     */

    // Required - this must be specified in the CRD
    autoscaleName := os.Getenv("AWS_AUTOSCALE_GROUP_NAME")

    if autoscaleName == "" {
        log.Fatal("Must specify the environment variable AWS_AUTOSCALE_GROUP_NAME")
    }

    // Optional - if specified, all must be specified
    region := os.Getenv("AWS_REGION")
    accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
    secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

    if region == "" {
        // Call the metadata server to get our region
        metaDataSvc := ec2metadata.New(sess)
        region, err = metaDataSvc.Region()
        if err == nil {
            sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
        }
    } else {
        sess, err = session.NewSession(&aws.Config{
            Region:      aws.String(region),
            Credentials: credentials.NewStaticCredentials(accessKey, secretKey, "")})
    }
    if err != nil {
        log.Fatal(err.Error())
    }
    fmt.Printf("\nAWS region: %s\n", region)

    /*
     * 2. Get the autoscale group information (list of instances & associated state)
     */
    asInput := &autoscaling.DescribeAutoScalingGroupsInput{
        AutoScalingGroupNames: []*string{aws.String(autoscaleName)},
    }

    asClient := autoscaling.New(sess)
    asDescription, err := asClient.DescribeAutoScalingGroups(asInput)
    if err != nil {
       log.Fatal(err.Error())
    }

    /*
     * 3. Get the IP (int/ext) and instance name  of all instances in the autoscale group
     */
    ec2Client := ec2.New(sess)
    // (should probably verify we only have 1 group)
    for _, instance := range asDescription.AutoScalingGroups[0].Instances {
	fmt.Printf("Instance ID: %s,", *instance.InstanceId)
	fmt.Printf("  Lifecycle State: %s", *instance.LifecycleState)
	fmt.Printf(",  Health State: %s", *instance.HealthStatus)

        // We should add all instance Ids for a single autoscale group to cut down on the api calls
        ec2Input := &ec2.DescribeInstancesInput{
            InstanceIds: []*string{aws.String(*instance.InstanceId)},
        }
        ecDescription, err := ec2Client.DescribeInstances(ec2Input)
        if err == nil {
            if len(ecDescription.Reservations[0].Instances) > 0  && len(ecDescription.Reservations[0].Instances[0].NetworkInterfaces) > 0 {
                fmt.Printf(",  Private IP: %s", *ecDescription.Reservations[0].Instances[0].NetworkInterfaces[0].PrivateIpAddress)
                fmt.Printf(",  Public IP: %s", *ecDescription.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicIp)
                // FIXME(kenr): This assumes the second index.  We need to traverse looking for a tag where Tags[i].Key == 'Name'
                for _, tag := range ecDescription.Reservations[0].Instances[0].Tags {
                    if *tag.Key == "Name" {
                        fmt.Printf(",  Name: %s", *tag.Value)
                    }
                }
            }
        }
        fmt.Printf("\n")
    }
    if err != nil {
        log.Fatal(err.Error())
    }
}
