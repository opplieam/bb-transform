# bb-transform: Data Transformation Tool for Machine Learning Dataset Generation

[![Go Report Card](https://goreportcard.com/badge/github.com/opplieam/bb-transform)](https://goreportcard.com/report/github.com/opplieam/bb-transform)

## Overview

![diagram](https://raw.githubusercontent.com/opplieam/bb-transform/refs/heads/main/transform.png)

`bb-transform` is a small but powerful tool designed to transform raw data from a database into a structured machine learning dataset. It is deployed as a serverless function on AWS Lambda and triggered by messages from AWS SQS. The entire AWS infrastructure is provisioned using Terraform, ensuring consistent and reproducible deployments.

This project demonstrates a typical workflow for preparing data for machine learning, leveraging cloud-native technologies for scalability and efficiency.

## Project Structure

```bash
├── cmd
│   └── lambda
│       └── main.go  # Entry point for the AWS Lambda function.
├── internal
│   ├── lambdahandler
│   │   └── lambdahandler.go # Handles SQS events and triggers data transformation.
│   ├── store
│   │   ├── category.go  # Manages interactions with the category data in the database.
│   │   └── db.go  # Handles database connection and configuration.
│   └── transform
│       ├── transform.go
├── terraform  # Terraform configurations for AWS infrastructure.
```

## Features

*   **Database Integration:** Connects to a PostgreSQL database (using Supabase in this example) to retrieve raw data.
*   **Data Transformation:** Processes raw category data and matched category data to create a structured dataset, including:
    *   Handling hierarchical category structures.
    *   Generating train, validation, and test datasets with configurable ratios.
    *   Optional shuffling of data for better model training.
*   **Serverless Deployment:** Deploys as an AWS Lambda function for scalability and cost-effectiveness.
*   **SQS Trigger:** Triggered by messages from an AWS SQS queue, enabling event-driven processing.
*   **Infrastructure as Code:** Uses Terraform to manage and provision all necessary AWS resources.
*   **Local Development:** Supports local development for testing and debugging using environment variable for development env

## Prerequisites

*   **Go:** Version 1.23.3 or higher.
*   **Terraform:** Version 1.10.0 or higher.
*   **AWS Account:** An active AWS account with permissions to create Lambda functions, SQS queues, IAM roles, and other necessary resources.
*   **Supabase Account:** A Supabase project for the PostgreSQL database. Or you can use your own Postgres DB. 
*   **Jet:** For database schema management using `jet`, install it following their official documentation. 
scheme can be found [here](https://github.com/opplieam/bb-admin-api/)
*   **(Optional) Golangci-lint:** For code linting, install `golangci-lint`.

## Setup

### Database migration

Please visit [here](https://github.com/opplieam/bb-admin-api/)

### Environment Variables

Create a `.env` file in the root directory of the project and add the following environment variables:
```bash
BUYBETTER_DEV_SUPABASE_DSN="your_supabase_dsn" # Replace with your Supabase DSN
SQS_QUEUE_URL="your_sqs_queue_url" # Replace with your SQS queue URL
ENV="dev" # Set to "dev" for local development
```

### Terraform Configuration

1. **Backend Configuration (`backend-config.tfvars`):**
   Create a `terraform/backend-config.tfvars` file to configure your Terraform backend (e.g., S3 bucket for storing Terraform state).

    ```
    bucket = "your-terraform-state-bucket" # Replace with your bucket name
    key    = "terraform.tfstate"
    region = "your-aws-region" # Replace with your AWS region
    ```

2. **Variables (`terraform.tfvars`):**
   Create a `terraform/terraform.tfvars` file to define variables for your AWS infrastructure.

    ```
    aws_region            = "us-east-2" # Replace with your desired AWS region
    sqs_queue_name        = "bb-transform-queue"
    sqs_dlq_name          = "bb-transform-dlq"
    lambda_function_name = "transform-category"
    BUYBETTER_DEV_SUPABASE_DSN = "your_supabase_dsn" # Same as in .env
    ```

### Local Development

1. **Database Setup:**
    *   Ensure your Supabase database is set up with the necessary tables.
    *   Use `jet-gen` to generate the necessary Go code for database interactions:

        ```bash
        make jet-gen
        ```

2. **Run Locally:**

    ```bash
    go run ./cmd/lambda/main.go
    ```

### Deployment

1. **Initialize Terraform:**

    ```bash
    make terraform-init
    ```

2. **Plan Infrastructure Changes:**

    ```bash
    make terraform-plan
    ```

3. **Apply Infrastructure Changes:**

    ```bash
    make terraform-apply
    ```

4. **Build the Lambda Function:**

    ```bash
    make build-lambda
    ```

5. **Deploy the Lambda Function:**

    ```bash
    make deploy-lambda
    ```

## Usage

### Triggering Data Transformation

To trigger the data transformation process, send a message to the configured SQS queue with the following JSON payload:

```json
{
    "version": "v1-lambda",  
    "shuffle": true,
    "train_ratio": 60,
    "validate_ratio": 20,
    "test_ratio": 20
}

```

- version: A string representing the version of the dataset (used for cleanup).
- shuffle: A boolean indicating whether to shuffle the data before splitting.
- train_ratio, validate_ratio, test_ratio: Integers (0-100) representing the percentage of data to use for each dataset split. These should add up to 100.

You can use the following command to send a message:

```bash 
make sent-message
```

### Testing
Run unit tests:
```bash
make test
```

### Linting
Run linter:
```bash
make lint
```

### Cleanup
To destroy the AWS infrastructure created by Terraform, run:
```bash
make terraform-destroy
```