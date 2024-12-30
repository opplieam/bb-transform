variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-2"
}

variable "lambda_function_name" {
  type = string
  description = "Lambda function name"
}

variable "BUYBETTER_DEV_SUPABASE_DSN" {
  type = string
  description = "Suppabase DSN environment variable"
}

variable "sqs_queue_name" {
  type = string
  description = "SQS queue name"
}

variable "sqs_dlq_name" {
  type = string
  description = "SQS dead letter queue name"
}