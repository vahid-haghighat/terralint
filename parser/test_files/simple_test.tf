// Simple Terraform file for testing the parser

// Resource block
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  tags = {
    Name = "example-instance"
    Environment = "test"
  }
}

// Variable block
variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

// Output block
output "instance_id" {
  description = "ID of the EC2 instance"
  value       = aws_instance.example.id
}

// Data source block
data "aws_ami" "ubuntu" {
  most_recent = true
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
  
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  
  owners = ["099720109477"] # Canonical
}

// Provider block
provider "aws" {
  region = var.region
}

// Locals block
locals {
  common_tags = {
    Project     = "Test"
    Owner       = "Terraform"
    Environment = "Test"
  }
}

// Module block
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  
  name = "my-vpc"
  cidr = "10.0.0.0/16"
  
  azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  
  enable_nat_gateway = true
  enable_vpn_gateway = true
  
  tags = local.common_tags
}