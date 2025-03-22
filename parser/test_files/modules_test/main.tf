// Main Terraform file for testing module structure

// Provider configuration
provider "aws" {
  region = var.region
  
  default_tags {
    tags = {
      ManagedBy = "Terraform"
      Project   = var.project_name
    }
  }
}

// Terraform configuration
terraform {
  required_version = ">= 1.0.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  
  backend "s3" {
    bucket         = "terraform-state-bucket"
    key            = "modules-test/terraform.tfstate"
    region         = "us-west-2"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}

// Local variables
locals {
  environment_code = {
    development = "dev"
    staging     = "stg"
    production  = "prod"
  }
  
  env_code = lookup(local.environment_code, var.environment, "dev")
  
  name_prefix = "${var.project_name}-${local.env_code}"
  
  common_tags = {
    Environment = var.environment
    Project     = var.project_name
    ManagedBy   = "Terraform"
    Owner       = var.owner
  }
}

// VPC Module
module "vpc" {
  source = "modules/vpc"
  
  vpc_name       = "${local.name_prefix}-vpc"
  vpc_cidr       = var.vpc_cidr
  azs            = var.availability_zones
  private_subnets = var.private_subnet_cidrs
  public_subnets  = var.public_subnet_cidrs
  
  enable_nat_gateway = var.enable_nat_gateway
  single_nat_gateway = var.environment != "production"
  
  enable_vpn_gateway = var.enable_vpn_gateway
  
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-vpc"
    }
  )
}

// Security group for the application
module "security" {
  source = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"
  
  name        = "${local.name_prefix}-sg"
  description = "Security group for ${var.project_name} application"
  vpc_id      = module.vpc.vpc_id
  
  ingress_with_cidr_blocks = [
    {
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      description = "HTTP"
      cidr_blocks = "0.0.0.0/0"
    },
    {
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      description = "HTTPS"
      cidr_blocks = "0.0.0.0/0"
    }
  ]
  
  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      description = "Allow all outbound traffic"
      cidr_blocks = "0.0.0.0/0"
    }
  ]
  
  tags = local.common_tags
}

// EC2 instances
resource "aws_instance" "app" {
  count = var.instance_count
  
  ami           = var.ami_id
  instance_type = var.instance_type
  subnet_id     = element(module.vpc.private_subnets, count.index % length(module.vpc.private_subnets))
  
  vpc_security_group_ids = [module.security.security_group_id]
  
  root_block_device {
    volume_type = "gp3"
    volume_size = 50
    encrypted   = true
  }
  
  user_data = templatefile("${path.module}/templates/user_data.sh.tpl", {
    environment = var.environment
    region      = var.region
    app_port    = var.app_port
  })
  
  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-app-${count.index + 1}"
    }
  )
  
  lifecycle {
    create_before_destroy = true
  }
}

// Load balancer
resource "aws_lb" "app" {
  name               = "${local.name_prefix}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [module.security.security_group_id]
  subnets            = module.vpc.public_subnets
  
  enable_deletion_protection = var.environment == "production"
  
  tags = local.common_tags
}

// Target group
resource "aws_lb_target_group" "app" {
  name     = "${local.name_prefix}-tg"
  port     = var.app_port
  protocol = "HTTP"
  vpc_id   = module.vpc.vpc_id
  
  health_check {
    path                = "/health"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    interval            = 30
    matcher             = "200"
  }
  
  tags = local.common_tags
}

// Target group attachment
resource "aws_lb_target_group_attachment" "app" {
  count            = var.instance_count
  target_group_arn = aws_lb_target_group.app.arn
  target_id        = aws_instance.app[count.index].id
  port             = var.app_port
}

// Listener
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.app.arn
  port              = 80
  protocol          = "HTTP"
  
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app.arn
  }
  
  tags = local.common_tags
}

// Database
resource "aws_db_instance" "app" {
  count = var.create_database ? 1 : 0
  
  identifier           = "${local.name_prefix}-db"
  engine               = "postgres"
  engine_version       = "13.4"
  instance_class       = var.db_instance_class
  allocated_storage    = var.db_allocated_storage
  max_allocated_storage = var.db_max_allocated_storage
  
  db_name              = var.db_name
  username             = var.db_username
  password             = var.db_password
  
  vpc_security_group_ids = [module.security.security_group_id]
  db_subnet_group_name   = aws_db_subnet_group.app[0].name
  
  backup_retention_period = var.environment == "production" ? 30 : 7
  backup_window           = "03:00-04:00"
  maintenance_window      = "Mon:04:00-Mon:05:00"
  
  skip_final_snapshot     = var.environment != "production"
  final_snapshot_identifier = var.environment == "production" ? "${local.name_prefix}-db-final-snapshot" : null
  
  tags = local.common_tags
}

// DB subnet group
resource "aws_db_subnet_group" "app" {
  count = var.create_database ? 1 : 0
  
  name       = "${local.name_prefix}-db-subnet-group"
  subnet_ids = module.vpc.private_subnets
  
  tags = local.common_tags
}

// Outputs
output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value       = module.vpc.private_subnets
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = module.vpc.public_subnets
}

output "security_group_id" {
  description = "The ID of the security group"
  value       = module.security.security_group_id
}

output "instance_ids" {
  description = "List of IDs of instances"
  value       = aws_instance.app[*].id
}

output "instance_private_ips" {
  description = "List of private IP addresses of instances"
  value       = aws_instance.app[*].private_ip
}

output "lb_dns_name" {
  description = "The DNS name of the load balancer"
  value       = aws_lb.app.dns_name
}

output "db_instance_endpoint" {
  description = "The connection endpoint of the database"
  value       = var.create_database ? aws_db_instance.app[0].endpoint : null
}