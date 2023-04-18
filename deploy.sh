#!/bin/bash

# Set the AWS region and profile
aws configure set region us-east-1
aws configure set profile default

# Set the stack name using the AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)

TIMESTAMP=$(date +%Y%m%d%H%M%S)
STACK_NAME="txn-summary-stack-${AWS_ACCOUNT_ID}-${TIMESTAMP}"
BUCKET_NAME="transactions-summary-bucket-${AWS_ACCOUNT_ID}"

# Get the VPC ID and subnet IDs
VPC_ID=$(aws ec2 describe-vpcs --filters Name=isDefault,Values=true --query 'Vpcs[0].VpcId' --output text)
SUBNET1_ID=$(aws ec2 describe-subnets --filters Name=vpc-id,Values=$VPC_ID --query 'Subnets[0].SubnetId' --output text)
SUBNET2_ID=$(aws ec2 describe-subnets --filters Name=vpc-id,Values=$VPC_ID --query 'Subnets[1].SubnetId' --output text)

# Set the csv file path
CSV_FILE_PATH="mock-data/txns2.csv"

# Build the SAM application
sam build

# Deploy the SAM application
sam deploy --guided --stack-name $STACK_NAME --capabilities CAPABILITY_IAM --parameter-overrides VpcId=$VPC_ID PublicSubnet1=$SUBNET1_ID PublicSubnet2=$SUBNET2_ID

# Upload the CSV file to the S3 bucket
aws s3 cp $CSV_FILE_PATH s3://$BUCKET_NAME/txns.csv
