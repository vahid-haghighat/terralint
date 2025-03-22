// This file focuses on complex template directives and interpolation

// Complex template with nested directives
locals {
  complex_template = <<-EOT
    # Configuration File
    
    %{if var.environment == "production"}
    environment = "production"
    log_level = "warn"
    debug = false
    %{else if var.environment == "staging"}
    environment = "staging"
    log_level = "info"
    debug = true
    %{else}
    environment = "development"
    log_level = "debug"
    debug = true
    %{endif}
    
    # Resources section
    [resources]
    %{for name, config in var.resources}
    [[resources.${name}]]
    type = "${config.type}"
    size = ${config.size}
    replicas = ${config.replicas}
    %{if config.high_availability}
    ha_enabled = true
    min_replicas = ${config.min_replicas}
    max_replicas = ${config.max_replicas}
    %{endif}
    
    %{for key, value in config.settings}
    ${key} = ${value}
    %{endfor}
    %{endfor}
    
    # Network configuration
    [network]
    %{for i, subnet in var.subnets}
    [[network.subnet]]
    cidr = "${subnet.cidr}"
    zone = "${subnet.zone}"
    public = ${subnet.public}
    
    %{if subnet.public}
    # Public subnet configuration
    gateway = "${subnet.gateway}"
    %{else}
    # Private subnet configuration
    nat_gateway = "${subnet.nat_gateway}"
    %{endif}
    %{endfor}
    
    # Nested for loops
    [instances]
    %{for type, count in var.instance_counts}
    %{for i in range(count)}
    [[instances.${type}]]
    id = "${type}-${i + 1}"
    %{if i < 3}
    primary = true
    %{else}
    primary = false
    %{endif}
    %{endfor}
    %{endfor}
    
    # Conditional sections with nested interpolation
    [security]
    %{if var.enable_encryption}
    encryption = true
    algorithm = "${var.encryption_algorithm}"
    key_rotation = ${var.key_rotation_days} days
    
    %{if var.encryption_algorithm == "AES256"}
    # AES256 specific settings
    key_size = 256
    %{else if var.encryption_algorithm == "AES128"}
    # AES128 specific settings
    key_size = 128
    %{else}
    # Default settings
    key_size = 192
    %{endif}
    %{else}
    encryption = false
    %{endif}
    
    # Complex nested conditions and loops
    [advanced]
    %{if var.advanced_config}
    enabled = true
    
    %{for feature, settings in var.advanced_features}
    [[advanced.${feature}]]
    %{if settings.enabled}
    enabled = true
    %{for key, value in settings.config}
    ${key} = ${jsonencode(value)}
    %{endfor}
    
    %{if feature == "monitoring"}
    # Monitoring specific settings
    [advanced.${feature}.alerts]
    %{for alert in settings.alerts}
    [[advanced.${feature}.alerts.rule]]
    name = "${alert.name}"
    threshold = ${alert.threshold}
    duration = "${alert.duration}"
    severity = "${alert.severity}"
    %{endfor}
    %{endif}
    %{else}
    enabled = false
    %{endif}
    %{endfor}
    %{else}
    enabled = false
    %{endif}
    
    # Deeply nested conditionals
    [system]
    %{if var.system_config != null}
    %{if var.system_config.enabled}
    enabled = true
    mode = "${var.system_config.mode}"
    
    %{if var.system_config.mode == "high_performance"}
    # High performance settings
    [system.performance]
    cpu_allocation = "maximum"
    memory_overcommit = true
    io_priority = "high"
    %{else if var.system_config.mode == "balanced"}
    # Balanced settings
    [system.performance]
    cpu_allocation = "balanced"
    memory_overcommit = false
    io_priority = "normal"
    %{else}
    # Economy settings
    [system.performance]
    cpu_allocation = "minimum"
    memory_overcommit = false
    io_priority = "low"
    %{endif}
    
    %{if var.system_config.backup != null}
    [system.backup]
    enabled = ${var.system_config.backup.enabled}
    %{if var.system_config.backup.enabled}
    schedule = "${var.system_config.backup.schedule}"
    retention = ${var.system_config.backup.retention} days
    
    %{for target in var.system_config.backup.targets}
    [[system.backup.target]]
    path = "${target.path}"
    priority = ${target.priority}
    %{endfor}
    %{endif}
    %{endif}
    %{else}
    enabled = false
    %{endif}
    %{else}
    # System configuration not provided
    enabled = false
    %{endif}
  EOT
}

// Complex template with strip markers
locals {
  strip_markers_template = <<-EOT
    # Normal line with trailing whitespace    
    
    %{~for item in var.items~}
    ${item}%{if item != var.items[length(var.items) - 1]}, %{endif}
    %{~endfor~}
    
    # The above should render on a single line
    
    %{for user in var.users~}
    username: ${user.name}
      %{~if user.admin~}
      role: administrator
      %{~else~}
      role: user
      %{~endif~}
    %{~endfor}
    
    # Complex stripping with nested directives
    %{~for group in var.groups~}
    [${group.name}]
      %{~for member in group.members~}
      - ${member.name}%{if member.owner} (owner)%{endif}
      %{~endfor~}
    %{~endfor~}
  EOT
}

// Nested template directives in string interpolation
locals {
  nested_template = "The result is: ${
    var.condition ? 
    "Condition is true, items: ${join(", ", [
      for item in var.items : 
      upper(item) if contains(var.allowed_items, item)
    ])}" : 
    "Condition is false, count: ${length(var.items)}"
  }"
}

// Complex string interpolation with functions and conditionals
locals {
  complex_interpolation = {
    // Nested function calls in interpolation
    nested_functions = "Result: ${
      jsonencode(
        merge(
          { 
            name = var.name,
            enabled = var.enabled
          },
          {
            for k, v in var.tags :
            "tag_${k}" => v if v != null
          },
          var.additional_settings != null ? var.additional_settings : {}
        )
      )
    }"
    
    // Conditional with function calls
    conditional = "Status: ${
      var.status == "healthy" ? 
      upper("OK - ${formatdate("YYYY-MM-DD", timestamp())}") :
      var.status == "degraded" ? 
      title("warning - ${var.message}") :
      lower("error - ${var.error_code}: ${var.message}")
    }"
    
    // For expression in interpolation
    for_expression = "Items: ${
      join(", ", [
        for i, item in var.items :
        "${i + 1}. ${item.name} (${item.type})" if item.enabled
      ])
    }"
    
    // Nested conditionals
    nested_conditionals = "Result: ${
      var.level > 5 ? 
      "High (${var.level})" :
      var.level > 2 ?
      "Medium (${var.level})" :
      var.level > 0 ?
      "Low (${var.level})" :
      "None"
    }"
    
    // Complex splat in interpolation
    splat_expression = "IDs: ${
      join(", ", 
        flatten([
          for group in var.groups : [
            for member in group.members : 
            "${group.name}.${member.id}"
          ]
        ])
      )
    }"
  }
}

// Dynamic blocks with complex expressions
resource "aws_security_group" "complex_dynamic" {
  name        = "complex-dynamic-example"
  description = "Example of complex dynamic blocks"
  
  // Dynamic block with nested dynamic blocks
  dynamic "ingress" {
    for_each = {
      for port_obj in var.ingress_ports :
      port_obj.name => port_obj if port_obj.enabled
    }
    
    content {
      description = lookup(ingress.value, "description", "Port ${ingress.value.port}")
      from_port   = ingress.value.port
      to_port     = ingress.value.port
      protocol    = lookup(ingress.value, "protocol", "tcp")
      
      // Nested dynamic block
      dynamic "cidr_blocks" {
        for_each = lookup(ingress.value, "cidr_blocks", [])
        
        content {
          cidr_block = cidr_blocks.value
        }
      }
      
      // Another nested dynamic block
      dynamic "security_groups" {
        for_each = lookup(ingress.value, "security_group_ids", [])
        
        content {
          security_group_id = security_groups.value
        }
      }
      
      // Conditional nested dynamic block
      dynamic "self" {
        for_each = lookup(ingress.value, "self", false) ? [1] : []
        
        content {
          self = true
        }
      }
    }
  }
  
  // Dynamic block with complex iterator and conditional
  dynamic "egress" {
    for_each = {
      for idx, rule in var.egress_rules :
      idx => rule if(
        rule.enabled && 
        (var.environment == "prod" ? contains(rule.environments, "prod") : true) &&
        (length(lookup(rule, "protocols", [])) > 0)
      )
    }
    iterator = egress_rule
    
    content {
      description = lookup(egress_rule.value, "description", "Egress rule ${egress_rule.key}")
      
      // Dynamic block for each protocol
      dynamic "protocol" {
        for_each = egress_rule.value.protocols
        
        content {
          protocol    = protocol.value
          from_port   = lookup(egress_rule.value, "from_port", 0)
          to_port     = lookup(egress_rule.value, "to_port", 0)
          cidr_blocks = lookup(egress_rule.value, "cidr_blocks", ["0.0.0.0/0"])
        }
      }
    }
  }
  
  // Conditional dynamic block
  dynamic "timeouts" {
    for_each = var.custom_timeouts != null ? [var.custom_timeouts] : []
    
    content {
      create = lookup(timeouts.value, "create", null)
      update = lookup(timeouts.value, "update", null)
      delete = lookup(timeouts.value, "delete", null)
    }
  }
  
  // Dynamic block with complex expressions
  dynamic "tags" {
    for_each = merge(
      {
        Name        = "complex-dynamic-example"
        Environment = var.environment
        ManagedBy   = "Terraform"
      },
      var.common_tags != null ? var.common_tags : {},
      {
        for k, v in var.additional_tags :
        k => v if v != null && k != "Name" && k != "Environment"
      }
    )
    
    content {
      key   = tags.key
      value = tags.value
    }
  }
}

// Complex template directives in heredoc
resource "aws_instance" "template_directives" {
  ami           = var.ami_id
  instance_type = var.instance_type
  
  user_data = <<-EOF
    #!/bin/bash
    
    # Set environment variables
    %{for key, value in var.environment_variables}
    export ${key}="${value}"
    %{endfor}
    
    # Configure based on environment
    %{if var.environment == "production"}
    echo "Setting up production environment"
    export LOG_LEVEL=WARN
    export MONITORING=true
    %{else if var.environment == "staging"}
    echo "Setting up staging environment"
    export LOG_LEVEL=INFO
    export MONITORING=true
    %{else}
    echo "Setting up development environment"
    export LOG_LEVEL=DEBUG
    export MONITORING=false
    %{endif}
    
    # Install packages
    %{for package in var.packages}
    echo "Installing ${package}"
    %{if package == "docker"}
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    %{else}
    apt-get update
    apt-get install -y ${package}
    %{endif}
    %{endfor}
    
    # Configure services
    %{for service, config in var.services}
    echo "Configuring ${service}"
    
    cat > /etc/systemd/system/${service}.service <<'SERVICE'
    [Unit]
    Description=${config.description}
    After=network.target
    
    [Service]
    Type=${config.type}
    ExecStart=${config.exec_start}
    %{if config.restart_policy != null}
    Restart=${config.restart_policy}
    %{endif}
    
    [Install]
    WantedBy=multi-user.target
    SERVICE
    
    %{if config.enabled}
    systemctl enable ${service}
    systemctl start ${service}
    %{endif}
    %{endfor}
    
    # Setup users
    %{for user in var.users}
    echo "Creating user ${user.name}"
    useradd -m ${user.name}
    %{if user.groups != null}
    usermod -G ${join(",", user.groups)} ${user.name}
    %{endif}
    
    %{if user.ssh_key != null}
    mkdir -p /home/${user.name}/.ssh
    echo "${user.ssh_key}" > /home/${user.name}/.ssh/authorized_keys
    chmod 600 /home/${user.name}/.ssh/authorized_keys
    chown -R ${user.name}:${user.name} /home/${user.name}/.ssh
    %{endif}
    
    %{if user.sudo}
    echo "${user.name} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${user.name}
    chmod 440 /etc/sudoers.d/${user.name}
    %{endif}
    %{endfor}
    
    # Conditional sections with nested loops
    %{if var.setup_monitoring}
    echo "Setting up monitoring"
    
    %{for agent in var.monitoring_agents}
    echo "Installing ${agent.name} agent"
    
    %{if agent.type == "prometheus"}
    # Prometheus agent setup
    %{for target in agent.targets}
    echo "  - Adding target: ${target.name} (${target.endpoint})"
    %{endfor}
    %{else if agent.type == "datadog"}
    # Datadog agent setup
    export DD_API_KEY="${agent.api_key}"
    %{for tag in agent.tags}
    export DD_TAGS="${tag.key}:${tag.value}"
    %{endfor}
    %{else}
    # Unknown agent type
    echo "Unknown agent type: ${agent.type}"
    %{endif}
    %{endfor}
    %{endif}
    
    echo "Setup complete!"
  EOF
  
  tags = {
    Name = "template-directives-example"
  }
}

// Complex for expressions with template directives
locals {
  // Nested for expressions with template directives
  complex_for_template = [
    for group in var.groups : {
      name = group.name
      id = group.id
      members = [
        for member in group.members : {
          name = member.name
          role = member.role
          permissions = <<-EOT
            # Permissions for ${member.name}
            
            %{if member.role == "admin"}
            # Admin permissions
            read = true
            write = true
            execute = true
            %{else if member.role == "developer"}
            # Developer permissions
            read = true
            write = true
            execute = false
            %{else}
            # Basic permissions
            read = true
            write = false
            execute = false
            %{endif}
            
            %{for resource in member.resources}
            [[resource]]
            name = "${resource.name}"
            type = "${resource.type}"
            %{if resource.type == "database"}
            # Database specific permissions
            %{for action in resource.allowed_actions}
            allow_${action} = true
            %{endfor}
            %{endif}
            %{endfor}
          EOT
        } if member.active
      ]
      config = <<-EOT
        # Group configuration for ${group.name}
        
        [group]
        name = "${group.name}"
        id = "${group.id}"
        
        %{if group.settings != null}
        [group.settings]
        %{for key, value in group.settings}
        ${key} = ${jsonencode(value)}
        %{endfor}
        %{endif}
        
        %{for policy in group.policies}
        [[group.policy]]
        name = "${policy.name}"
        priority = ${policy.priority}
        
        %{for rule in policy.rules}
        [[group.policy.rule]]
        type = "${rule.type}"
        action = "${rule.action}"
        %{if rule.conditions != null}
        %{for condition in rule.conditions}
        condition_${condition.name} = "${condition.value}"
        %{endfor}
        %{endif}
        %{endfor}
        %{endfor}
      EOT
    } if group.enabled
  ]
}

// Complex template with nested interpolation and directives
locals {
  nested_template_directives = templatefile("${path.module}/templates/config.tpl", {
    environment = var.environment
    region = var.region
    services = [
      for service in var.services : {
        name = service.name
        type = service.type
        config = {
          enabled = service.enabled
          port = service.port
          protocol = service.protocol
          settings = merge(
            {
              default_timeout = 30
              max_connections = 100
            },
            service.settings != null ? service.settings : {}
          )
        }
        endpoints = [
          for endpoint in service.endpoints : {
            path = endpoint.path
            method = endpoint.method
            auth_required = endpoint.auth_required
            rate_limit = lookup(endpoint, "rate_limit", 100)
          } if endpoint.enabled
        ]
        dependencies = [
          for dep in service.dependencies : {
            name = dep.name
            required = dep.required
          }
        ]
        template = <<-EOT
          # Service: ${service.name}
          
          %{if service.enabled}
          status = "enabled"
          port = ${service.port}
          protocol = "${service.protocol}"
          
          %{for key, value in service.settings != null ? service.settings : {}}
          ${key} = ${jsonencode(value)}
          %{endfor}
          
          %{for endpoint in service.endpoints}
          %{if endpoint.enabled}
          [[endpoint]]
          path = "${endpoint.path}"
          method = "${endpoint.method}"
          auth_required = ${endpoint.auth_required}
          %{endif}
          %{endfor}
          %{else}
          status = "disabled"
          %{endif}
        EOT
      } if service.type != "external"
    ]
    databases = {
      for db in var.databases : db.name => {
        engine = db.engine
        version = db.version
        size = db.size
        replicas = db.replicas
        config = <<-EOT
          # Database: ${db.name}
          
          engine = "${db.engine}"
          version = "${db.version}"
          
          %{if db.engine == "postgres"}
          # PostgreSQL specific settings
          max_connections = ${db.settings.max_connections}
          shared_buffers = "${db.settings.shared_buffers}"
          %{else if db.engine == "mysql"}
          # MySQL specific settings
          innodb_buffer_pool_size = "${db.settings.innodb_buffer_pool_size}"
          max_connections = ${db.settings.max_connections}
          %{else}
          # Generic database settings
          %{for key, value in db.settings}
          ${key} = ${jsonencode(value)}
          %{endfor}
          %{endif}
          
          %{for user in db.users}
          [[user]]
          name = "${user.name}"
          role = "${user.role}"
          %{if user.role == "admin"}
          privileges = ["ALL"]
          %{else}
          privileges = ${jsonencode(user.privileges)}
          %{endif}
          %{endfor}
        EOT
      }
    }
  })
}