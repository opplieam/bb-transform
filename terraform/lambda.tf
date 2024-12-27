variable "lambda_function_name" {
  type = string
  description = "Lambda function name"
}

variable "BUYBETTER_DEV_SUPABASE_DSN" {
  type = string
  description = "Suppabase DSN environment variable"
}

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.17.0"

  function_name = var.lambda_function_name
  description   = "Transform category data from Supabase to datasets for machine learning"
  handler       = "bootstrap"
  runtime       = "provided.al2"
  publish       = true

  memory_size = 128
  timeout     = 20

  # Env
  environment_variables = {
    BUYBETTER_DEV_SUPABASE_DSN: var.BUYBETTER_DEV_SUPABASE_DSN
  }

  create_package         = false
  local_existing_package = "./dummy/function.zip"

  # Lambda Concurrency
#   reserved_concurrent_executions = 2

  # CloudWatch
  cloudwatch_logs_retention_in_days = 30

  # Allow Lambda to be triggered by SQS
  allowed_triggers = {
    SQS = {
      principal  = "sqs.amazonaws.com"
      source_arn = module.transform_queue.queue_arn
    }
  }

  # IAM role attachments
  attach_policies    = true
  policies          = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"]
  attach_policy_statements = true
  policy_statements = {
    sqs = {
      effect = "Allow"
      actions = [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
      ]
      resources = [module.transform_queue.queue_arn]
    }
  }

  tags = {
    Environment = "prod"
    Project     = "transform-category"
  }
}

resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = module.transform_queue.queue_arn
  function_name    = module.lambda_function.lambda_function_arn
  enabled          = true
  batch_size       = 1 # Adjust as needed
}


output "lambda_function_arn" {
  description = "The ARN of the Lambda Function"
  value       = module.lambda_function.lambda_function_arn
}

output "lambda_function_name" {
  description = "The name of the Lambda Function"
  value       = module.lambda_function.lambda_function_name
}