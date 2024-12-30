variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-2"
}

# Get current AWS account ID
data "aws_caller_identity" "current" {}

# IAM User for GitHub Actions
resource "aws_iam_user" "github_actions" {
  name = "github-actions-lambda"

  tags = {
    Description = "IAM user for GitHub Actions Lambda deployment"
    Environment = "prod"
  }
}

# Access key for the IAM user
resource "aws_iam_access_key" "github_actions" {
  user = aws_iam_user.github_actions.name
}

# IAM Policy
resource "aws_iam_policy" "lambda_deploy" {
  name        = "github-actions-lambda-deploy"
  description = "Policy for GitHub Actions to deploy Lambda functions"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:UpdateFunctionCode",
          "lambda:PublishVersion",
          "lambda:GetFunction"
        ]
        Resource = "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${var.lambda_function_name}"
      }
    ]
  })
}

# Attach policy to user
resource "aws_iam_user_policy_attachment" "lambda_deploy" {
  user       = aws_iam_user.github_actions.name
  policy_arn = aws_iam_policy.lambda_deploy.arn
}

output "github_actions_access_key_id" {
  value     = aws_iam_access_key.github_actions.id
  sensitive = true
}

output "github_actions_secret_access_key" {
  value     = aws_iam_access_key.github_actions.secret
  sensitive = true
}