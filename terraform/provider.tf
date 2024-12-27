variable "region" {
  type = string
  description = "AWS region"
}

variable "bucket" {
  type = string
  description = "S3 bucket name"
}

variable "key" {
  type = string
  description = "S3 bucket key"
}

terraform {
  required_version = ">= 1.10.0"

  backend "s3" {

  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}