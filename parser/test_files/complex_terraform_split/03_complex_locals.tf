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