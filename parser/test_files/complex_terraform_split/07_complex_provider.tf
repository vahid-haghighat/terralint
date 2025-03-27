// Provider configuration with complex expressions
provider "aws" {
  region = var.region

  assume_role {
    role_arn = var.environment == "prod" ? var.prod_role_arn : var.non_prod_role_arn
    session_name = "terraform-${var.environment}-${formatdate("YYYY-MM-DD-hh-mm-ss", timestamp())}"
    external_id = var.external_id
  }

  default_tags {
    tags = merge(
      var.common_tags,
      {
        Environment = var.environment
        ManagedBy   = "Terraform"
        Project     = var.project_name
        Owner       = var.owner
        CreatedAt   = formatdate("YYYY-MM-DD", timestamp())
      }
    )
  }

  dynamic "endpoints" {
    for_each = var.custom_endpoints != null ? var.custom_endpoints : {}
    content {
      apigateway     = lookup(endpoints.value, "apigateway", null)
      cloudformation = lookup(endpoints.value, "cloudformation", null)
      cloudwatch     = lookup(endpoints.value, "cloudwatch", null)
      dynamodb       = lookup(endpoints.value, "dynamodb", null)
      ec2            = lookup(endpoints.value, "ec2", null)
      s3             = lookup(endpoints.value, "s3", null)
      sts            = lookup(endpoints.value, "sts", null)
    }
  }
}