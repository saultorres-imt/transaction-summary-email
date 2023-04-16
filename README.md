# Transaction-summary-email
This project processes a CSV file containing account transactions and sends a summary email to the user. The CSV file is read from an AWS S3 bucket. The application is designed to run on AWS Lambda.

**Table of Contents**
- [Prerequisites](#prerequisites)
- [CVS File Format](#csv-file-format)
- [Setup](#setup)
- [Usage](#usage)

## Prerequisites

- [Go](https://golang.org/doc/install) installed
- [AWS CLI](https://aws.amazon.com/cli/) installed and configured with AWS credentials

## CSV File Format

The CSV file should have the following format:

Id,Date,Transaction
0,7/15,+60.5
1,7/28,-10.3
2,8/2,-20.46
3,8/13,+10

Where:
- Id is a unique identifier for the transaction
- Date is the transaction date in the format M/D
- Transaction is the transaction amount, with a '+' sign for credit transactions and a '-' sign for debit transactions

## Setup

1. Clone the repository:
```shell
git clone https://github.com/saultorres-imt/transaction-summary-email.git
cd transaction-summary-email
```

## Usage

Run the following command to verify, Amazon SES, the email addresses that will be used
```bash
aws ses verify-email-identity --email-address youremail@example.com
```