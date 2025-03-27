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