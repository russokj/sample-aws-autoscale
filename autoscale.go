/*
   Copyright Ken
*/

package main

import (
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/autoscaling"
)

func main() {
    // We can explicitly filter on a specific region or use env vars to do it
    // sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})

    // Load session from shared config
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create new EC2 client object
    ec2Svc := ec2.New(sess)

    // Call to get detailed information on each instance
    result, err := ec2Svc.DescribeInstances(nil)
    if err != nil {
        fmt.Println("Error", err)
    } else {
        fmt.Println("Success", result)
    }

    fmt.Println("\n\n*****\n\n")

    // Create Autoscale object
    asSvc := autoscaling.New(sess)

    // Example of finding all instances (we could filter for a specific one)
    result2, err2 := asSvc.DescribeAutoScalingInstances(nil)
    if err2 != nil {
       fmt.Println(err.Error())
    } else {
       fmt.Println(result2)
    }

    fmt.Println("\n\n*****\n\n")

    // Example of specifying a specific autoscale group
    // (Set input to nil if we want to find all autoscale gropus)
    input := &autoscaling.DescribeAutoScalingGroupsInput{
        AutoScalingGroupNames: []*string{
            aws.String("damian-heptio-k8s-K8sStack-1MZY47ZF88RCO-K8sNodeGroup-1LY9CV0XVVTLK"),
        },
    }

    result3, err3 := asSvc.DescribeAutoScalingGroups(input)
    if err3 != nil {
       fmt.Println(err.Error())
    } else {
       // TODO: Cycle through all instances of the autoscale group
       fmt.Println(result3)
    }

}
