// Variables for the VPC module

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC"
  type        = string
}

variable "availability_zones" {
  description = "A list of availability zones to use"
  type        = list(string)
}

variable "public_subnets" {
  description = "A list of public subnet CIDR blocks"
  type        = list(string)
}

variable "private_subnets" {
  description = "A list of private subnet CIDR blocks"
  type        = list(string)
}

variable "database_subnets" {
  description = "A list of database subnet CIDR blocks"
  type        = list(string)
  default     = []
}

variable "enable_dns_hostnames" {
  description = "Should be true to enable DNS hostnames in the VPC"
  type        = bool
  default     = true
}

variable "enable_dns_support" {
  description = "Should be true to enable DNS support in the VPC"
  type        = bool
  default     = true
}

variable "enable_nat_gateway" {
  description = "Should be true to provision NAT Gateways for each of your private subnets"
  type        = bool
  default     = true
}

variable "single_nat_gateway" {
  description = "Should be true to provision a single shared NAT Gateway across all of your private subnets"
  type        = bool
  default     = false
}

variable "enable_vpn_gateway" {
  description = "Should be true to enable a VPN Gateway in your VPC"
  type        = bool
  default     = false
}

variable "enable_s3_endpoint" {
  description = "Should be true to provision an S3 endpoint to the VPC"
  type        = bool
  default     = false
}

variable "enable_dynamodb_endpoint" {
  description = "Should be true to provision a DynamoDB endpoint to the VPC"
  type        = bool
  default     = false
}

variable "create_network_acls" {
  description = "Should be true to create network ACLs"
  type        = bool
  default     = false
}

variable "enable_flow_logs" {
  description = "Should be true to enable VPC Flow Logs"
  type        = bool
  default     = false
}

variable "flow_logs_destination_arn" {
  description = "The ARN of the S3 bucket where VPC Flow Logs will be pushed"
  type        = string
  default     = null
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}