This repo just makes calls to AWS to get a list of instances and autoscale
groups for a particular IAM user

Requirements:
1. Requires (for now) the AWS GO SDK to be downloaded and available at
~/Documents/Projects/AWS/sdk

2. Requires the IAM role credentials to be specified in aws-env.sh

Note: There are 3 ways for specifying AWS credentials. See
https://docs.aws.amazon.com/sdk-for-go/api (under 'Configuring Credentials')

