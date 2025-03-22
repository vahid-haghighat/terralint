// Main module file that references other modules and demonstrates complex module usage

// Provider configuration with complex expressions
provider "aws" {
  region = var.aws_region
  
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
}

// Terraform configuration
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
}

// Local values with complex expressions
locals {
  // Environment-specific settings
  env_config = {
    dev = {
      instance_type = "t3.small"
      min_capacity  = 1
      max_capacity  = 3
      storage_gb    = 20
      multi_az      = false
    }
    staging = {
      instance_type = "t3.medium"
      min_capacity  = 2
      max_capacity  = 5
      storage_gb    = 50
      multi_az      = true
    }
    prod = {
      instance_type = "m5.large"
      min_capacity  = 3
      max_capacity  = 10
      storage_gb    = 100
      multi_az      = true
    }
  }
  
  // Get current environment config with fallback to dev
  current_env_config = lookup(local.env_config, var.environment, local.env_config.dev)
  
  // Computed values
  vpc_name = "${var.project_name}-${var.environment}"
  
  // Subnet calculations
  availability_zones = slice(data.aws_availability_zones.available.names, 0, var.az_count)
  
  // Complex subnet CIDR calculations
  public_subnets = [
    for i, az in local.availability_zones :
    cidrsubnet(var.vpc_cidr, 8, i)
  ]
  
  private_subnets = [
    for i, az in local.availability_zones :
    cidrsubnet(var.vpc_cidr, 8, i + length(local.availability_zones))
  ]
  
  database_subnets = [
    for i, az in local.availability_zones :
    cidrsubnet(var.vpc_cidr, 8, i + (2 * length(local.availability_zones)))
  ]
  
  // Tags with complex merging
  common_tags = merge(
    var.common_tags,
    {
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "Terraform"
    }
  )
  
  // Resource naming with conditional logic
  resource_name_prefix = var.environment == "prod" ? "prod-${var.project_name}" : "${var.environment}-${var.project_name}"
  
  // Complex map transformation
  service_config = {
    for service_key, service in var.services :
    service_key => {
      name           = service.name
      container_port = service.container_port
      host_port      = lookup(service, "host_port", service.container_port)
      protocol       = lookup(service, "protocol", "tcp")
      cpu            = lookup(service, "cpu", 256)
      memory         = lookup(service, "memory", 512)
      desired_count  = lookup(service, "desired_count", local.current_env_config.min_capacity)
      max_count      = lookup(service, "max_count", local.current_env_config.max_capacity)
      health_check   = lookup(service, "health_check", {
        path                = "/health"
        interval            = 30
        timeout             = 5
        healthy_threshold   = 3
        unhealthy_threshold = 3
      })
      environment = concat(
        [
          {
            name  = "ENVIRONMENT"
            value = var.environment
          },
          {
            name  = "LOG_LEVEL"
            value = var.environment == "prod" ? "INFO" : "DEBUG"
          }
        ],
        lookup(service, "environment", [])
      )
      secrets = lookup(service, "secrets", [])
      tags = merge(
        local.common_tags,
        {
          Service = service.name
        },
        lookup(service, "tags", {})
      )
    }
  }
}

// Data sources
data "aws_availability_zones" "available" {}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

// VPC Module
module "vpc" {
  source = "./modules/vpc"
  
  vpc_name             = local.vpc_name
  vpc_cidr             = var.vpc_cidr
  availability_zones   = local.availability_zones
  public_subnets       = local.public_subnets
  private_subnets      = local.private_subnets
  database_subnets     = local.database_subnets
  enable_nat_gateway   = var.environment != "dev"
  single_nat_gateway   = var.environment != "prod"
  enable_vpn_gateway   = var.environment == "prod"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = local.common_tags
}

// Security Module
module "security" {
  source = "./modules/security"
  
  environment      = var.environment
  vpc_id           = module.vpc.vpc_id
  vpc_cidr         = var.vpc_cidr
  allowed_ips      = var.allowed_ips
  bastion_enabled  = var.environment != "dev"
  
  tags = local.common_tags
  
  depends_on = [module.vpc]
}

// Database Module with conditional creation
module "database" {
  source = "./modules/database"
  count  = var.create_database ? 1 : 0
  
  identifier           = "${local.resource_name_prefix}-db"
  engine               = var.db_engine
  engine_version       = var.db_engine_version
  instance_class       = lookup(var.db_instance_class, var.environment, "db.t3.small")
  allocated_storage    = lookup(var.db_allocated_storage, var.environment, 20)
  max_allocated_storage = lookup(var.db_max_allocated_storage, var.environment, 100)
  storage_encrypted    = true
  
  name                 = var.db_name
  username             = var.db_username
  password             = var.db_password
  port                 = var.db_port
  
  vpc_security_group_ids = [module.security.database_security_group_id]
  subnet_ids             = module.vpc.database_subnet_ids
  
  multi_az             = local.current_env_config.multi_az
  backup_retention_period = var.environment == "prod" ? 30 : 7
  backup_window        = "03:00-04:00"
  maintenance_window   = "sun:04:30-sun:05:30"
  
  parameters = concat(
    [
      {
        name  = "character_set_server"
        value = "utf8"
      },
      {
        name  = "character_set_client"
        value = "utf8"
      }
    ],
    var.db_parameters
  )
  
  tags = merge(
    local.common_tags,
    {
      Name = "${local.resource_name_prefix}-db"
    }
  )
  
  depends_on = [
    module.vpc,
    module.security
  ]
}

// ECS Cluster Module
module "ecs" {
  source = "./modules/ecs"
  
  cluster_name = "${local.resource_name_prefix}-cluster"
  vpc_id       = module.vpc.vpc_id
  subnet_ids   = module.vpc.private_subnet_ids
  
  services = {
    for service_key, service in local.service_config :
    service_key => {
      name           = service.name
      container_port = service.container_port
      host_port      = service.host_port
      protocol       = service.protocol
      cpu            = service.cpu
      memory         = service.memory
      desired_count  = service.desired_count
      max_count      = service.max_count
      health_check   = service.health_check
      environment    = service.environment
      secrets        = service.secrets
      tags           = service.tags
      
      // Add database connection if database exists
      database_connection = var.create_database ? {
        host     = module.database[0].endpoint
        port     = var.db_port
        name     = var.db_name
        username = var.db_username
        password = var.db_password
      } : null
    }
  }
  
  load_balancer_security_group_id = module.security.load_balancer_security_group_id
  service_security_group_id       = module.security.service_security_group_id
  
  tags = local.common_tags
  
  depends_on = [
    module.vpc,
    module.security,
    module.database
  ]
}

// S3 Bucket Module with dynamic creation
module "s3_buckets" {
  source   = "./modules/s3"
  for_each = { for bucket in var.s3_buckets : bucket.name => bucket }
  
  bucket_name          = "${local.resource_name_prefix}-${each.value.name}"
  acl                  = lookup(each.value, "acl", "private")
  versioning_enabled   = lookup(each.value, "versioning_enabled", false)
  lifecycle_rules      = lookup(each.value, "lifecycle_rules", [])
  server_side_encryption = lookup(each.value, "server_side_encryption", true)
  
  // Conditional CORS configuration
  cors_enabled = lookup(each.value, "cors_enabled", false)
  cors_rules   = lookup(each.value, "cors_rules", [])
  
  // Conditional logging configuration
  logging_enabled = lookup(each.value, "logging_enabled", false)
  logging_target_bucket = lookup(each.value, "logging_target_bucket", null)
  logging_target_prefix = lookup(each.value, "logging_target_prefix", null)
  
  tags = merge(
    local.common_tags,
    {
      Name = "${local.resource_name_prefix}-${each.value.name}"
    },
    lookup(each.value, "tags", {})
  )
}

// CloudFront Distribution Module with conditional creation
module "cloudfront" {
  source = "./modules/cloudfront"
  count  = var.create_cloudfront ? 1 : 0
  
  distribution_name = "${local.resource_name_prefix}-distribution"
  
  origins = concat(
    [
      {
        domain_name = module.ecs.load_balancer_dns_name
        origin_id   = "ecs-load-balancer"
        custom_origin_config = {
          http_port              = 80
          https_port             = 443
          origin_protocol_policy = "https-only"
          origin_ssl_protocols   = ["TLSv1.2"]
        }
      }
    ],
    [
      for bucket_key, bucket in module.s3_buckets :
      {
        domain_name = bucket.bucket_regional_domain_name
        origin_id   = "s3-${bucket_key}"
        s3_origin_config = {
          origin_access_identity = "origin-access-identity/cloudfront/${bucket.cloudfront_access_identity_path}"
        }
      }
    ]
  )
  
  default_cache_behavior = {
    target_origin_id       = "ecs-load-balancer"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }
  
  ordered_cache_behaviors = [
    for bucket_key, bucket in module.s3_buckets :
    {
      path_pattern           = "/static/*"
      target_origin_id       = "s3-${bucket_key}"
      viewer_protocol_policy = "redirect-to-https"
      allowed_methods        = ["GET", "HEAD", "OPTIONS"]
      cached_methods         = ["GET", "HEAD"]
      compress               = true
      min_ttl                = 0
      default_ttl            = 86400
      max_ttl                = 31536000
    }
  ]
  
  aliases = var.environment == "prod" ? [var.domain_name, "www.${var.domain_name}"] : ["${var.environment}.${var.domain_name}"]
  
  viewer_certificate = {
    acm_certificate_arn      = var.acm_certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
  
  geo_restriction = {
    restriction_type = "none"
    locations        = []
  }
  
  tags = local.common_tags
  
  depends_on = [
    module.ecs,
    module.s3_buckets
  ]
}

// Route53 Records Module with conditional creation
module "route53" {
  source = "./modules/route53"
  count  = var.create_dns_records ? 1 : 0
  
  zone_id     = var.route53_zone_id
  domain_name = var.domain_name
  
  records = concat(
    var.create_cloudfront ? [
      {
        name    = var.environment == "prod" ? "" : var.environment
        type    = "A"
        alias   = {
          name                   = module.cloudfront[0].distribution_domain_name
          zone_id                = module.cloudfront[0].distribution_hosted_zone_id
          evaluate_target_health = false
        }
      },
      {
        name    = var.environment == "prod" ? "www" : "www-${var.environment}"
        type    = "A"
        alias   = {
          name                   = module.cloudfront[0].distribution_domain_name
          zone_id                = module.cloudfront[0].distribution_hosted_zone_id
          evaluate_target_health = false
        }
      }
    ] : [],
    [
      {
        name    = var.environment == "prod" ? "api" : "api-${var.environment}"
        type    = "A"
        alias   = {
          name                   = module.ecs.load_balancer_dns_name
          zone_id                = module.ecs.load_balancer_zone_id
          evaluate_target_health = true
        }
      }
    ]
  )
  
  depends_on = [
    module.ecs,
    module.cloudfront
  ]
}

// Monitoring Module
module "monitoring" {
  source = "./modules/monitoring"
  
  environment         = var.environment
  project_name        = var.project_name
  resource_name_prefix = local.resource_name_prefix
  
  // Conditional alarm thresholds based on environment
  alarm_thresholds = {
    cpu_utilization = var.environment == "prod" ? 80 : 90
    memory_utilization = var.environment == "prod" ? 80 : 90
    response_time = var.environment == "prod" ? 1 : 2
    error_rate = var.environment == "prod" ? 1 : 5
  }
  
  // Resources to monitor
  ecs_cluster_name = module.ecs.cluster_name
  ecs_services     = [for service_key, service in local.service_config : service.name]
  
  db_instance_identifier = var.create_database ? module.database[0].identifier : null
  
  load_balancer_arn = module.ecs.load_balancer_arn
  
  // Notification settings
  sns_topic_arn = var.sns_topic_arn
  
  tags = local.common_tags
  
  depends_on = [
    module.ecs,
    module.database
  ]
}

// Outputs
output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "List of public subnet IDs"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "List of private subnet IDs"
  value       = module.vpc.private_subnet_ids
}

output "database_subnet_ids" {
  description = "List of database subnet IDs"
  value       = module.vpc.database_subnet_ids
}

output "database_endpoint" {
  description = "The database endpoint"
  value       = var.create_database ? module.database[0].endpoint : null
  sensitive   = true
}

output "ecs_cluster_name" {
  description = "The name of the ECS cluster"
  value       = module.ecs.cluster_name
}

output "ecs_service_names" {
  description = "The names of the ECS services"
  value       = module.ecs.service_names
}

output "load_balancer_dns_name" {
  description = "The DNS name of the load balancer"
  value       = module.ecs.load_balancer_dns_name
}

output "s3_bucket_names" {
  description = "The names of the S3 buckets"
  value       = { for key, bucket in module.s3_buckets : key => bucket.bucket_name }
}

output "cloudfront_distribution_domain_name" {
  description = "The domain name of the CloudFront distribution"
  value       = var.create_cloudfront ? module.cloudfront[0].distribution_domain_name : null
}

output "domain_urls" {
  description = "The URLs for the application"
  value = {
    main = var.create_dns_records ? (
      var.environment == "prod" ? "https://${var.domain_name}" : "https://${var.environment}.${var.domain_name}"
    ) : null
    api = var.create_dns_records ? (
      var.environment == "prod" ? "https://api.${var.domain_name}" : "https://api-${var.environment}.${var.domain_name}"
    ) : null
  }
}