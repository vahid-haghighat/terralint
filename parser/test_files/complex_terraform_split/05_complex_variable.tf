// Variable with complex type constraints and validations
variable "complex_object" {
  description = "A complex object with nested types and validations"
  type = object({
    name = string
    environment = string
    enabled = bool
    count = number
    tags = map(string)
    vpc = object({
      id = string
      cidr = string
      private_subnets = list(object({
        id = string
        cidr = string
        az = string
      }))
      public_subnets = list(object({
        id = string
        cidr = string
        az = string
      }))
    })
    instances = list(object({
      id = string
      type = string
      subnet_id = string
      private_ip = string
      public_ip = optional(string)
      root_volume = object({
        size = number
        type = string
        encrypted = bool
      })
      data_volumes = list(object({
        device_name = string
        size = number
        type = string
        encrypted = bool
        iops = optional(number)
        throughput = optional(number)
      }))
      tags = map(string)
    }))
    databases = map(object({
      engine = string
      version = string
      instance_class = string
      allocated_storage = number
      max_allocated_storage = optional(number)
      multi_az = bool
      backup_retention_period = number
      parameters = map(string)
      subnet_ids = list(string)
      security_group_ids = list(string)
    }))
    endpoints = map(object({
      service = string
      vpc_endpoint_type = string
      subnet_ids = optional(list(string))
      security_group_ids = optional(list(string))
      private_dns_enabled = optional(bool)
      policy = optional(string)
    }))
  })

  validation {
    condition     = var.complex_object.count > 0 && var.complex_object.count <= 10
    error_message = "Count must be between 1 and 10."
  }

  validation {
    condition     = can(regex("^(dev|staging|prod)$", var.complex_object.environment))
    error_message = "Environment must be one of: dev, staging, prod."
  }

  validation {
    condition     = length(var.complex_object.vpc.private_subnets) > 0
    error_message = "At least one private subnet must be defined."
  }

  validation {
    condition     = alltrue([
      for instance in var.complex_object.instances :
      instance.root_volume.size >= 20 && instance.root_volume.encrypted == true
    ])
    error_message = "All root volumes must be at least 20GB and encrypted."
  }
}