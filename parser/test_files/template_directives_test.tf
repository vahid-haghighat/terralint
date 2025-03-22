// This file contains complex template directives and interpolation to test the parser

locals {
  // Complex template with multiple interpolations
  complex_template = "Hello, ${var.name}! Welcome to ${var.environment} environment. Your IP is ${var.ip_address}."
  
  // Template with nested expressions
  nested_template = "The result is: ${var.enable_feature ? upper(var.feature_name) : "Feature disabled"}"
  
  // Template with function calls
  function_template = "Timestamp: ${timestamp()}, UUID: ${uuid()}, Base64: ${base64encode("Hello")}"
  
  // Template with math expressions
  math_template = "The answer is ${(1 + 2) * 3 / 4}"
  
  // Template with references to other resources
  reference_template = "Instance ID: ${aws_instance.example.id}, Public IP: ${aws_instance.example.public_ip}"
  
  // Template with for expressions
  for_template = "Items: ${join(", ", [for item in var.items : upper(item)])}"
  
  // Template with conditional expressions
  conditional_template = "Status: ${var.status == "active" ? "Active" : "Inactive"}"
  
  // Template with strip markers
  strip_markers_template = <<-EOT
    This is a template with strip markers.
    %{if var.enable_feature~}
    The feature is enabled.
    %{else~}
    The feature is disabled.
    %{endif~}
    
    Available items:
    %{for item in var.items~}
    - ${item}
    %{endfor~}
    
    %{if var.show_timestamp~}
    Timestamp: ${timestamp()}
    %{endif~}
  EOT
  
  // Template with indentation
  indented_template = <<-EOT
    server {
      listen 80;
      server_name ${var.domain_name};
      
      location / {
        proxy_pass http://${var.backend_host}:${var.backend_port};
        %{if var.enable_ssl~}
        proxy_ssl_verify on;
        proxy_ssl_trusted_certificate /etc/ssl/certs/ca-certificates.crt;
        %{endif~}
      }
      
      %{for path, config in var.extra_locations~}
      location ${path} {
        %{for key, value in config~}
        ${key} ${value};
        %{endfor~}
      }
      %{endfor~}
    }
  EOT
  
  // Template with escaping
  escaped_template = "Escaped interpolation: $${var.name}, Escaped directive: %{if true}not processed%{endif}"
  
  // Template with special characters
  special_chars_template = "Special chars: !@#$%^&*()_+-=[]{}|;:'\",.<>?/\\`~"
  
  // Template with unicode
  unicode_template = "Unicode: こんにちは世界 • Hello, World! • Привет, мир! • مرحبا بالعالم • 你好，世界！"
  
  // Template with newlines and tabs
  newlines_template = "Line 1\nLine 2\n\tIndented line\nLine 4"
}

// Resource with template directives
resource "aws_instance" "template_directives" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  // User data with template directives
  user_data = <<-EOF
    #!/bin/bash
    
    # Set variables
    HOSTNAME="${var.hostname}"
    ENVIRONMENT="${var.environment}"
    
    # Update system
    apt-get update
    apt-get upgrade -y
    
    # Install packages
    %{for package in var.packages~}
    apt-get install -y ${package}
    %{endfor~}
    
    # Configure services
    %{if var.enable_nginx~}
    # Nginx configuration
    cat > /etc/nginx/sites-available/default <<'NGINX'
    server {
      listen 80;
      server_name ${var.domain_name};
      
      location / {
        proxy_pass http://localhost:${var.app_port};
      }
    }
    NGINX
    
    systemctl enable nginx
    systemctl start nginx
    %{endif~}
    
    # Set hostname
    hostnamectl set-hostname ${var.hostname}
    
    # Create users
    %{for username, user_config in var.users~}
    # Create user ${username}
    useradd -m -s /bin/bash ${username}
    %{if user_config.sudo~}
    usermod -aG sudo ${username}
    %{endif~}
    %{if user_config.ssh_key != ""~}
    mkdir -p /home/${username}/.ssh
    echo "${user_config.ssh_key}" > /home/${username}/.ssh/authorized_keys
    chmod 600 /home/${username}/.ssh/authorized_keys
    chown -R ${username}:${username} /home/${username}/.ssh
    %{endif~}
    %{endfor~}
    
    # Final message
    echo "Setup completed for ${var.hostname} in ${var.environment} environment"
  EOF
  
  // Tags with template directives
  tags = {
    Name = "${var.name_prefix}-instance"
    Environment = var.environment
    CreatedAt = timestamp()
    CreatedBy = "Terraform"
    %{if var.include_owner~}
    Owner = var.owner
    %{endif~}
  }
}