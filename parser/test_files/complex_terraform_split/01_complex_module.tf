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