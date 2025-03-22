// This file contains extremely complex Terraform code to test the parser
// It includes nested expressions, complex interpolations, and edge cases

// Complex module with nested expressions, conditionals, and for loops
module "complex_module" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  // Complex map with nested objects, expressions, and functions
  vpc_config = {
    name = "complex-${var.environment}-${local.region_code}-vpc"
    cidr = cidrsubnet(var.base_cidr_block, 4, var.vpc_index)
    
    // Nested conditional expression
    enable_dns = var.environment == "prod" ? true : (var.enable_dns != null ? var.enable_dns : false)
    
    // Complex object with nested expressions
    tags = merge(
      var.common_tags,
      {
        Name        = "complex-${var.environment}-vpc"
        Environment = var.environment
        ManagedBy   = "terraform"
        Complex     = jsonencode({
          nested = "value"
          list   = [1, 2, 3, 4]
          map    = {
            key1 = "value1"
            key2 = 42
          }
        })
      },
      var.additional_tags != null ? var.additional_tags : {}
    )
  }

  // Complex for expression with filtering and transformation
  subnet_cidrs = [
    for i, subnet in var.subnets :
    cidrsubnet(var.base_cidr_block, 8, i + 10) if subnet.create == true
  ]

  // Nested for expressions with conditional
  subnet_configs = {
    for zone_key, zone in var.availability_zones :
    zone_key => {
      for subnet_key, subnet in var.subnet_types :
      subnet_key => {
        cidr = cidrsubnet(
          var.base_cidr_block,
          var.subnet_newbits,
          index(var.availability_zones, zone) * length(var.subnet_types) + index(var.subnet_types, subnet)
        )
        az   = zone
        tags = merge(
          var.common_tags,
          {
            Name = "${var.environment}-${subnet}-${zone}"
            Type = subnet
          }
        )
      } if subnet.enabled
    } if contains(var.enabled_zones, zone)
  }

  // Complex splat expression
  all_subnet_ids = flatten([
    for zone_key, zone in aws_subnet.main : [
      for subnet_key, subnet in zone : subnet.id
    ]
  ])

  // Heredoc with interpolation
  user_data = <<-EOT
    #!/bin/bash
    echo "Environment: ${var.environment}"
    echo "Region: ${data.aws_region.current.name}"
    
    # Complex interpolation
    ${join("\n", [
      for script in var.bootstrap_scripts :
      "source ${script}"
    ])}
    
    # Conditional section
    ${var.install_monitoring ? "setup_monitoring ${var.monitoring_endpoint}" : "echo 'Monitoring disabled'"}
    
    # For loop in heredoc
    %{for pkg in var.packages~}
    yum install -y ${pkg}
    %{endfor~}
    
    # If directive in heredoc
    %{if var.environment == "prod"~}
    echo "Production environment detected, applying strict security"
    %{else~}
    echo "Non-production environment, using standard security"
    %{endif~}
  EOT

  // Complex binary expressions with nested conditionals
  timeout = (
    var.environment == "prod" ? 300 :
    var.environment == "staging" ? 180 :
    var.environment == "dev" ? 60 : 30
  ) + (var.additional_timeout != null ? var.additional_timeout : 0)

  // Nested parentheses and operators
  complex_calculation = (
    (var.base_value * (1 + var.multiplier)) / 
    (var.divisor > 0 ? var.divisor : 1)
  ) + (
    var.environment == "prod" ? 
    (var.prod_adjustment * (1 + var.prod_factor)) : 
    (var.non_prod_adjustment * (1 - var.non_prod_factor))
  )

  // Complex function calls with nested expressions
  security_groups = compact(concat(
    [
      var.create_default_security_group ? aws_security_group.default[0].id : "",
      var.create_bastion_security_group ? aws_security_group.bastion[0].id : ""
    ],
    var.additional_security_group_ids
  ))

  // Template with directives
  custom_template = templatefile("${path.module}/templates/config.tpl", {
    environment = var.environment
    region      = data.aws_region.current.name
    vpc_id      = aws_vpc.main.id
    subnets     = [
      for subnet in aws_subnet.main : {
        id   = subnet.id
        cidr = subnet.cidr_block
        az   = subnet.availability_zone
      }
    ]
    features = {
      monitoring = var.enable_monitoring
      logging    = var.enable_logging
      encryption = var.environment == "prod" ? true : var.enable_encryption
    }
  })

  // Dynamic blocks with complex expressions
  dynamic "ingress" {
    for_each = var.ingress_rules
    iterator = rule
    content {
      description = lookup(rule.value, "description", "Ingress Rule ${rule.key}")
      from_port   = rule.value.from_port
      to_port     = rule.value.to_port
      protocol    = lookup(rule.value, "protocol", "tcp")
      cidr_blocks = lookup(rule.value, "cidr_blocks", null) != null ? rule.value.cidr_blocks : [
        for cidr in var.default_cidrs :
        cidr if !contains(var.excluded_cidrs, cidr)
      ]
      security_groups = lookup(rule.value, "security_group_ids", [])
      self            = lookup(rule.value, "self", false)
    }
  }

  // Complex type constraints
  validation = {
    condition     = can(regex("^(dev|staging|prod)$", var.environment))
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

// Resource with complex dynamic blocks and for_each
resource "aws_security_group" "complex" {
  for_each = {
    for sg in var.security_groups :
    sg.name => sg if sg.create
  }

  name        = "${var.prefix}-${each.key}"
  description = each.value.description
  vpc_id      = var.vpc_id

  // Dynamic blocks with nested expressions
  dynamic "ingress" {
    for_each = each.value.ingress_rules
    content {
      description      = ingress.value.description
      from_port        = ingress.value.from_port
      to_port          = ingress.value.to_port
      protocol         = ingress.value.protocol
      cidr_blocks      = ingress.value.cidr_blocks
      ipv6_cidr_blocks = lookup(ingress.value, "ipv6_cidr_blocks", [])
      prefix_list_ids  = lookup(ingress.value, "prefix_list_ids", [])
      security_groups  = lookup(ingress.value, "security_groups", [])
      self             = lookup(ingress.value, "self", false)
    }
  }

  dynamic "egress" {
    for_each = each.value.egress_rules
    content {
      description      = egress.value.description
      from_port        = egress.value.from_port
      to_port          = egress.value.to_port
      protocol         = egress.value.protocol
      cidr_blocks      = egress.value.cidr_blocks
      ipv6_cidr_blocks = lookup(egress.value, "ipv6_cidr_blocks", [])
      prefix_list_ids  = lookup(egress.value, "prefix_list_ids", [])
      security_groups  = lookup(egress.value, "security_groups", [])
      self             = lookup(egress.value, "self", false)
    }
  }

  // Complex tags with expressions and functions
  tags = merge(
    var.common_tags,
    {
      Name = "${var.prefix}-${each.key}"
      Type = "SecurityGroup"
      Rules = jsonencode({
        ingress = length(each.value.ingress_rules)
        egress  = length(each.value.egress_rules)
      })
    },
    each.value.additional_tags != null ? each.value.additional_tags : {}
  )

  lifecycle {
    create_before_destroy = true
    prevent_destroy       = var.environment == "prod" ? true : false
    ignore_changes        = [
      tags["LastModified"],
      tags["AutoUpdated"],
    ]
  }
}

// Complex locals with nested expressions
locals {
  // Complex map transformation
  subnet_map = {
    for subnet in var.subnets :
    subnet.name => {
      id         = subnet.id
      cidr       = subnet.cidr
      az         = subnet.availability_zone
      public     = subnet.public
      nat_gw     = subnet.public ? true : false
      depends_on = subnet.public ? [] : [for s in var.subnets : s.id if s.public]
    }
  }

  // Nested for expressions with filtering
  filtered_instances = [
    for server in var.servers :
    {
      id          = server.id
      name        = server.name
      environment = server.environment
      type        = server.type
      subnet_id   = server.subnet_id
      private_ip  = server.private_ip
      public_ip   = server.public_ip
      tags        = server.tags
    }
    if server.environment == var.environment &&
    contains(var.allowed_types, server.type) &&
    !contains(var.excluded_ids, server.id)
  ]

  // Complex conditional with multiple nested expressions
  backup_config = var.enable_backups ? {
    schedule = var.backup_schedule != null ? var.backup_schedule : "0 1 * * *"
    retention = (
      var.environment == "prod" ? 30 :
      var.environment == "staging" ? 14 :
      7
    )
    targets = [
      for target in var.backup_targets :
      {
        id       = target.id
        name     = target.name
        priority = lookup(target.tags, "backup-priority", "medium")
      }
      if lookup(target.tags, "backup-enabled", "false") == "true"
    ]
    storage = {
      type = var.backup_storage_type
      path = var.backup_storage_path != null ? var.backup_storage_path : "/backups/${var.environment}"
      settings = merge(
        var.default_storage_settings,
        var.custom_storage_settings != null ? var.custom_storage_settings : {}
      )
    }
  } : null

  // Complex string interpolation with functions
  naming_convention = join("-", compact([
    var.prefix,
    var.environment,
    var.region_code,
    var.name,
    var.suffix != "" ? var.suffix : null
  ]))

  // Nested ternary operators
  timeout_seconds = (
    var.custom_timeout != null ? var.custom_timeout :
    var.environment == "prod" ? (
      var.high_availability ? 300 : 180
    ) : (
      var.environment == "staging" ? 120 : 60
    )
  )

  // Complex type expressions
  schema = {
    type = "object"
    properties = {
      id = {
        type = "string"
        pattern = "^[a-zA-Z0-9-_]+$"
      }
      settings = {
        type = "object"
        properties = {
          enabled = { type = "boolean" }
          timeout = { type = "number" }
          retries = { type = "number" }
          options = {
            type = "array"
            items = {
              type = "object"
              properties = {
                name = { type = "string" }
                value = { type = "string" }
              }
            }
          }
        }
      }
      tags = {
        type = "object"
        additionalProperties = { type = "string" }
      }
    }
  }
}

// Data source with complex expressions
data "aws_iam_policy_document" "complex" {
  // Dynamic statement blocks
  dynamic "statement" {
    for_each = var.policy_statements
    content {
      sid       = lookup(statement.value, "sid", "Statement${statement.key}")
      effect    = lookup(statement.value, "effect", "Allow")
      actions   = statement.value.actions
      resources = statement.value.resources

      dynamic "condition" {
        for_each = lookup(statement.value, "conditions", [])
        content {
          test     = condition.value.test
          variable = condition.value.variable
          values   = condition.value.values
        }
      }

      dynamic "principals" {
        for_each = lookup(statement.value, "principals", [])
        content {
          type        = principals.value.type
          identifiers = principals.value.identifiers
        }
      }

      not_actions   = lookup(statement.value, "not_actions", [])
      not_resources = lookup(statement.value, "not_resources", [])
    }
  }

  // Override with inline statement
  statement {
    sid    = "ExplicitAllow"
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:ListBucket",
    ]
    resources = [
      "arn:aws:s3:::${var.bucket_name}",
      "arn:aws:s3:::${var.bucket_name}/*",
    ]

    condition {
      test     = "StringEquals"
      variable = "aws:SourceVpc"
      values   = [var.vpc_id]
    }

    condition {
      test     = "StringLike"
      variable = "aws:PrincipalTag/Role"
      values   = ["Admin", "Developer"]
    }

    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${var.account_id}:root"]
    }
  }
}

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

// Output with complex expressions
output "complex_output" {
  description = "Complex output with nested expressions"
  value = {
    vpc_id = module.complex_module.vpc_id
    subnet_ids = module.complex_module.subnet_ids
    security_group_ids = [
      for sg_key, sg in aws_security_group.complex :
      sg.id
    ]
    instance_details = {
      for instance in local.filtered_instances :
      instance.id => {
        name = instance.name
        private_ip = instance.private_ip
        public_ip = instance.public_ip
        subnet = {
          id = instance.subnet_id
          details = lookup(local.subnet_map, instance.subnet_id, null)
        }
        environment = instance.environment
        type = instance.type
        tags = instance.tags
      }
    }
    backup_enabled = var.enable_backups
    backup_config = local.backup_config
    naming_convention = local.naming_convention
    complex_calculation = module.complex_module.complex_calculation
    policy_document = data.aws_iam_policy_document.complex.json
  }

  sensitive = true

  depends_on = [
    module.complex_module,
    aws_security_group.complex,
    data.aws_iam_policy_document.complex
  ]
}

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