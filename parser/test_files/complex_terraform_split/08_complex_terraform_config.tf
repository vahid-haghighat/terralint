// Terraform configuration with complex expressions
terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0.0, < 5.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.1"
    }
  }

  backend "s3" {
    bucket         = "terraform-state-${var.environment}-${var.account_id}"
    key            = "complex-module/terraform.tfstate"
    region         = var.region
    dynamodb_table = "terraform-locks"
    encrypt        = true
    kms_key_id     = var.kms_key_id
    role_arn       = var.state_role_arn
  }

  experiments = [module_variable_optional_attrs]
}