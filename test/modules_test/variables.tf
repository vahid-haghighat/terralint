// Variables for the complex module test

variable "aws_region" {
  description = "The AWS region to deploy resources in"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "The environment (dev, staging, prod)"
  type        = string
  default     = "dev"
  
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "project_name" {
  description = "The name of the project"
  type        = string
  default     = "terraform-test"
}

variable "owner" {
  description = "The owner of the resources"
  type        = string
  default     = "terraform"
}

variable "account_id" {
  description = "The AWS account ID"
  type        = string
  default     = "123456789012"
}

variable "region" {
  description = "The AWS region for the backend"
  type        = string
  default     = "us-west-2"
}

variable "kms_key_id" {
  description = "The KMS key ID for the backend"
  type        = string
  default     = null
}

variable "state_role_arn" {
  description = "The IAM role ARN for the backend"
  type        = string
  default     = null
}

variable "prod_role_arn" {
  description = "The IAM role ARN for production"
  type        = string
  default     = "arn:aws:iam::123456789012:role/terraform-prod"
}

variable "non_prod_role_arn" {
  description = "The IAM role ARN for non-production"
  type        = string
  default     = "arn:aws:iam::123456789012:role/terraform-non-prod"
}

variable "external_id" {
  description = "The external ID for the IAM role"
  type        = string
  default     = "terraform-external-id"
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default     = {
    ManagedBy = "Terraform"
    Project   = "terraform-test"
  }
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "az_count" {
  description = "The number of availability zones to use"
  type        = number
  default     = 3
  
  validation {
    condition     = var.az_count > 0 && var.az_count <= 6
    error_message = "AZ count must be between 1 and 6."
  }
}

variable "allowed_ips" {
  description = "List of allowed IP addresses for SSH access"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "create_database" {
  description = "Whether to create a database"
  type        = bool
  default     = true
}

variable "db_engine" {
  description = "The database engine to use"
  type        = string
  default     = "postgres"
}

variable "db_engine_version" {
  description = "The database engine version to use"
  type        = string
  default     = "13.4"
}

variable "db_instance_class" {
  description = "The database instance class to use for each environment"
  type        = map(string)
  default     = {
    dev     = "db.t3.small"
    staging = "db.t3.medium"
    prod    = "db.m5.large"
  }
}

variable "db_allocated_storage" {
  description = "The amount of storage to allocate for the database for each environment"
  type        = map(number)
  default     = {
    dev     = 20
    staging = 50
    prod    = 100
  }
}

variable "db_max_allocated_storage" {
  description = "The maximum amount of storage to allocate for the database for each environment"
  type        = map(number)
  default     = {
    dev     = 100
    staging = 200
    prod    = 500
  }
}

variable "db_name" {
  description = "The name of the database"
  type        = string
  default     = "appdb"
}

variable "db_username" {
  description = "The username for the database"
  type        = string
  default     = "dbadmin"
  sensitive   = true
}

variable "db_password" {
  description = "The password for the database"
  type        = string
  default     = "dbpassword123!"
  sensitive   = true
}

variable "db_port" {
  description = "The port for the database"
  type        = number
  default     = 5432
}

variable "db_parameters" {
  description = "Additional database parameters"
  type        = list(object({
    name  = string
    value = string
  }))
  default     = []
}

variable "services" {
  description = "The services to deploy"
  type        = map(object({
    name           = string
    container_port = number
    host_port      = optional(number)
    protocol       = optional(string)
    cpu            = optional(number)
    memory         = optional(number)
    desired_count  = optional(number)
    max_count      = optional(number)
    health_check   = optional(object({
      path                = string
      interval            = number
      timeout             = number
      healthy_threshold   = number
      unhealthy_threshold = number
    }))
    environment    = optional(list(object({
      name  = string
      value = string
    })))
    secrets        = optional(list(object({
      name      = string
      valueFrom = string
    })))
    tags           = optional(map(string))
  }))
  default     = {
    api = {
      name           = "api"
      container_port = 8080
    }
    web = {
      name           = "web"
      container_port = 80
    }
  }
}

variable "s3_buckets" {
  description = "The S3 buckets to create"
  type        = list(object({
    name                = string
    acl                 = optional(string)
    versioning_enabled  = optional(bool)
    lifecycle_rules     = optional(list(object({
      id                       = string
      enabled                  = bool
      prefix                   = optional(string)
      expiration_days          = optional(number)
      noncurrent_version_expiration_days = optional(number)
      transition_days          = optional(number)
      transition_storage_class = optional(string)
    })))
    server_side_encryption = optional(bool)
    cors_enabled           = optional(bool)
    cors_rules             = optional(list(object({
      allowed_headers = list(string)
      allowed_methods = list(string)
      allowed_origins = list(string)
      expose_headers  = optional(list(string))
      max_age_seconds = optional(number)
    })))
    logging_enabled        = optional(bool)
    logging_target_bucket  = optional(string)
    logging_target_prefix  = optional(string)
    tags                   = optional(map(string))
  }))
  default     = [
    {
      name = "assets"
    },
    {
      name = "logs"
    }
  ]
}

variable "create_cloudfront" {
  description = "Whether to create a CloudFront distribution"
  type        = bool
  default     = true
}

variable "domain_name" {
  description = "The domain name for the application"
  type        = string
  default     = "example.com"
}

variable "acm_certificate_arn" {
  description = "The ARN of the ACM certificate for the CloudFront distribution"
  type        = string
  default     = "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"
}

variable "create_dns_records" {
  description = "Whether to create DNS records"
  type        = bool
  default     = true
}

variable "route53_zone_id" {
  description = "The Route53 zone ID"
  type        = string
  default     = "Z1234567890ABCDEFGHIJ"
}

variable "sns_topic_arn" {
  description = "The ARN of the SNS topic for notifications"
  type        = string
  default     = null
}