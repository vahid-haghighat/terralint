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