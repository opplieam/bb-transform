module "transform_queue_dlq" {
  source  = "terraform-aws-modules/sqs/aws"
  version = "4.2.1"

  name = var.sqs_dlq_name

  tags = {
    Environment = "prod"
    Project     = "transform-category"
  }
}

module "transform_queue" {
  source  = "terraform-aws-modules/sqs/aws"
  version = "4.2.1"

  name = var.sqs_queue_name

  redrive_policy = {
    maxReceiveCount = 3
    deadLetterTargetArn = module.transform_queue_dlq.queue_arn
  }

  tags = {
    Environment = "prod"
    Project     = "transform-category"
  }
}

output "sqs_queue_url" {
  description = "The URL of the SQS Queue"
  value       = module.transform_queue.queue_url
}

output "dlq_queue_url" {
  description = "The URL of the Dead Letter Queue"
  value       = module.transform_queue_dlq.queue_url
}