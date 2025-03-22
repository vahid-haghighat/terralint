// This file contains edge cases and unusual syntax patterns to test the parser's robustness

// Empty blocks
resource "null_resource" "empty" {}

// Blocks with only comments
resource "null_resource" "comments_only" {
  # This block has only comments
  // Multiple comment styles
  /* Block comment */
}

// Extremely long line with multiple nested function calls
locals {
  extremely_long_line = join(",", concat(split(",", join(",", formatlist("%s-%s", [for i in range(1, 50) : format("item-%04d", i)], [for j in range(1, 50) : format("subitem-%04d", j)]))), [for k in range(1, 20) : format("extra-%04d", k)], [for l in setproduct(range(1, 5), ["a", "b", "c", "d", "e"]) : format("product-%d-%s", l[0], l[1])], [for m, n in zipmap(range(1, 10), range(10, 20)) : format("map-%d-%d", m, n)]))
}

// Nested conditionals with mixed types
locals {
  nested_conditionals = true ? (
    false ? (
      1 > 0 ? (
        "a" == "a" ? (
          null == null ? "level5-true" : "level5-false"
        ) : "level4-false"
      ) : "level3-false"
    ) : (
      [] == [] ? (
        {} == {} ? "level4-true" : "level4-false"
      ) : "level3-false"
    )
  ) : "level1-false"
}

// Deeply nested objects
locals {
  deep_nesting = {
    level1 = {
      level2 = {
        level3 = {
          level4 = {
            level5 = {
              level6 = {
                level7 = {
                  level8 = {
                    level9 = {
                      level10 = {
                        value = "deeply nested value"
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}

// Complex string with multiple interpolation types
locals {
  complex_string = <<-EOT
    Regular text
    ${1 + 2 + 3}
    ${"string interpolation"}
    ${true ? "conditional true" : "conditional false"}
    ${[for i in range(1, 5) : i * 2]}
    ${jsonencode({
      key1 = "value1"
      key2 = 42
      key3 = true
      key4 = null
      key5 = ["a", "b", "c"]
      key6 = {
        nested = "value"
      }
    })}
    %{for x in ["a", "b", "c"]~}
    - Item: ${x}
    %{endfor~}
    %{if true~}
    Conditional section
    %{else~}
    Alternative section
    %{endif~}
  EOT
}

// Mixed quotes and escaping
locals {
  mixed_quotes = {
    double_quotes = "This is a \"quoted\" string with \"nested\" quotes"
    single_in_double = "This has 'single' quotes inside double quotes"
    escapes = "Escape sequences: \n \t \r \\ \""
    hex_escapes = "\x41\x42\x43"
    unicode_escapes = "\u2665 \u2764 \u1F60D"
  }
}

// Unusual but valid identifiers
resource "aws_instance" "unusual-identifiers" {
  tags = {
    "key-with-hyphens" = "value"
    "key.with.dots"    = "value"
    "key_with_underscores" = "value"
    "123numeric-prefix"    = "value"
    "mixed.case-and_symbols" = "value"
    "unicode-symbols-☺♥★" = "value"
    "extremely-long-key-name-that-goes-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on-and-on" = "value"
  }
}

// Complex heredoc with indentation and escaping
locals {
  complex_heredoc = <<-EOF
    This is a heredoc string
      With preserved indentation
        And multiple levels
          Of nesting
    
    It can contain "quotes" and 'apostrophes'
    As well as symbols like $, {, }, \, etc.
    
    Even things that look like interpolation: ${not_a_variable}
    
    And escape sequences:
    \n \t \r \\ \"
    
    EOF
}

// Indented heredoc
locals {
  indented_heredoc = <<-SCRIPT
    #!/bin/bash
    
    echo "Hello, world!"
    
    if [ "$ENVIRONMENT" == "production" ]; then
      echo "This is production!"
    else
      echo "This is not production."
    fi
    
    for i in {1..10}; do
      echo "Number: $i"
    done
    
    function greet() {
      local name="$1"
      echo "Hello, $name!"
    }
    
    greet "Terraform"
    
    exit 0
  SCRIPT
}

// Non-indented heredoc
locals {
  non_indented_heredoc = <<EOF
#!/bin/bash

echo "Hello, world!"

if [ "$ENVIRONMENT" == "production" ]; then
  echo "This is production!"
else
  echo "This is not production."
fi

for i in {1..10}; do
  echo "Number: $i"
done

function greet() {
  local name="$1"
  echo "Hello, $name!"
}

greet "Terraform"

exit 0
EOF
}

// Complex for expressions with multiple iterators and conditions
locals {
  complex_for = {
    // For expression with multiple iterators
    multi_iterator = {
      for idx, val in ["a", "b", "c"] :
      idx => {
        value = val
        upper = upper(val)
        index = idx
      }
    }
    
    // Nested for expressions
    nested_for = [
      for outer in ["x", "y", "z"] : [
        for inner in [1, 2, 3] : {
          outer = outer
          inner = inner
          combined = "${outer}-${inner}"
        }
      ]
    ]
    
    // For expression with complex condition
    conditional_for = [
      for i in range(1, 20) :
      i if i % 2 == 0 && i % 3 != 0 && (i < 5 || i > 15)
    ]
    
    // For expression with complex transformation
    transform_for = {
      for k, v in {
        a = 1
        b = 2
        c = 3
      } :
      k => {
        original = v
        doubled = v * 2
        squared = v * v
        as_string = tostring(v)
        is_even = v % 2 == 0
      }
    }
  }
}

// Complex splat expressions
locals {
  splat_expressions = {
    // Basic splat
    basic_splat = [
      {id = "a", value = 1},
      {id = "b", value = 2},
      {id = "c", value = 3}
    ][*].id
    
    // Splat with attribute access
    attribute_splat = [
      {
        id = "a",
        nested = {
          value = 1,
          tags = ["tag1", "tag2"]
        }
      },
      {
        id = "b",
        nested = {
          value = 2,
          tags = ["tag3", "tag4"]
        }
      }
    ][*].nested.tags
    
    // Splat with index access
    index_splat = [
      {
        id = "a",
        values = [10, 20, 30]
      },
      {
        id = "b",
        values = [40, 50, 60]
      }
    ][*].values[0]
    
    // Nested splat
    nested_splat = [
      {
        id = "a",
        subitems = [
          {name = "a1", value = 1},
          {name = "a2", value = 2}
        ]
      },
      {
        id = "b",
        subitems = [
          {name = "b1", value = 3},
          {name = "b2", value = 4}
        ]
      }
    ][*].subitems[*].name
  }
}

// Complex type constraints
variable "complex_types" {
  type = object({
    // Simple types
    string_val = string
    number_val = number
    bool_val = bool
    
    // List and set types
    list_strings = list(string)
    set_numbers = set(number)
    
    // Map types
    map_strings = map(string)
    map_objects = map(object({
      id = string
      value = number
    }))
    
    // Tuple type
    tuple_mixed = tuple([string, number, bool, list(string), map(number)])
    
    // Nested object type
    nested_object = object({
      name = string
      settings = object({
        enabled = bool
        timeout = number
        retries = number
        options = list(object({
          name = string
          value = string
        }))
      })
      tags = map(string)
    })
    
    // Optional attributes
    with_optional = object({
      required = string
      optional1 = optional(string)
      optional2 = optional(number, 42)
      optional3 = optional(bool, true)
      optional_obj = optional(object({
        name = string
        value = optional(number, 0)
      }), {
        name = "default"
      })
    })
    
    // Any type
    any_val = any
  })
  
  default = {
    string_val = "default"
    number_val = 42
    bool_val = true
    list_strings = ["a", "b", "c"]
    set_numbers = [1, 2, 3]
    map_strings = {
      key1 = "value1"
      key2 = "value2"
    }
    map_objects = {
      obj1 = {
        id = "id1"
        value = 1
      }
      obj2 = {
        id = "id2"
        value = 2
      }
    }
    tuple_mixed = ["string", 42, true, ["list"], {key = 1}]
    nested_object = {
      name = "nested"
      settings = {
        enabled = true
        timeout = 30
        retries = 3
        options = [
          {
            name = "option1"
            value = "value1"
          },
          {
            name = "option2"
            value = "value2"
          }
        ]
      }
      tags = {
        environment = "test"
      }
    }
    with_optional = {
      required = "required_value"
    }
    any_val = null
  }
}

// Multiple validations
variable "with_validations" {
  type = string
  
  validation {
    condition = length(var.with_validations) >= 8
    error_message = "Must be at least 8 characters long."
  }
  
  validation {
    condition = can(regex("^[a-zA-Z0-9-_]+$", var.with_validations))
    error_message = "Can only contain alphanumeric characters, hyphens, and underscores."
  }
  
  validation {
    condition = can(regex("[A-Z]", var.with_validations)) && can(regex("[0-9]", var.with_validations))
    error_message = "Must contain at least one uppercase letter and one number."
  }
}

// Dynamic blocks with complex expressions
resource "aws_security_group" "dynamic_blocks" {
  name = "dynamic-blocks-test"
  
  // Simple dynamic block
  dynamic "ingress" {
    for_each = [22, 80, 443]
    content {
      from_port = ingress.value
      to_port = ingress.value
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }
  
  // Dynamic block with iterator
  dynamic "egress" {
    for_each = [
      {
        port = 80,
        protocol = "tcp",
        cidr_blocks = ["10.0.0.0/8"]
      },
      {
        port = 443,
        protocol = "tcp",
        cidr_blocks = ["10.0.0.0/8", "172.16.0.0/12"]
      }
    ]
    iterator = rule
    content {
      from_port = rule.value.port
      to_port = rule.value.port
      protocol = rule.value.protocol
      cidr_blocks = rule.value.cidr_blocks
    }
  }
  
  // Nested dynamic blocks
  dynamic "ingress" {
    for_each = {
      http = 80
      https = 443
    }
    content {
      from_port = ingress.value
      to_port = ingress.value
      protocol = "tcp"
      
      dynamic "cidr_blocks" {
        for_each = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
        content {
          cidr_block = cidr_blocks.value
        }
      }
    }
  }
  
  // Conditional dynamic block
  dynamic "ingress" {
    for_each = var.enable_ssh ? [22] : []
    content {
      from_port = ingress.value
      to_port = ingress.value
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
      description = "SSH access"
    }
  }
}

// Complex function calls
locals {
  function_calls = {
    // Nested function calls
    nested = base64encode(jsonencode({
      data = formatlist("%s-%s", ["a", "b", "c"], ["1", "2", "3"])
      hash = md5(join(",", ["value1", "value2", "value3"]))
      timestamp = formatdate("YYYY-MM-DD'T'hh:mm:ss", timestamp())
    }))
    
    // Function with complex arguments
    template = templatefile("${path.module}/template.tpl", {
      instances = [
        for i in range(1, 5) : {
          name = "instance-${i}"
          type = i < 3 ? "small" : "large"
          tags = {
            Name = "instance-${i}"
            Environment = "test"
            Index = i
          }
        }
      ]
      settings = {
        enabled = true
        timeout = 30
        retries = 3
      }
    })
    
    // Mathematical functions
    math = {
      max = max([1, 2, 3, 4, 5])
      min = min([1, 2, 3, 4, 5])
      ceil = ceil(1.5)
      floor = floor(1.5)
      log = log(16, 2)
      pow = pow(2, 8)
      signum = signum(-42)
    }
    
    // String manipulation
    strings = {
      format = format("Hello, %s!", "World")
      join = join(", ", ["a", "b", "c"])
      split = split(",", "a,b,c")
      substr = substr("abcdef", 1, 3)
      replace = replace("Hello, World!", "World", "Terraform")
      regex = regex("^(\\d+)-(\\w+)$", "123-abc")
      regexall = regexall("\\d+", "abc123def456")
      lower = lower("HELLO")
      upper = upper("hello")
      title = title("hello world")
      chomp = chomp("hello\n")
      indent = indent(2, "hello\nworld")
    }
    
    // Type conversion
    conversion = {
      to_string = tostring(42)
      to_number = tonumber("42")
      to_bool = tobool("true")
      to_list = tolist(["a", "b", "c"])
      to_set = toset(["a", "b", "c", "a"])
      to_map = tomap({
        key1 = "value1"
        key2 = "value2"
      })
    }
    
    // Collection functions
    collection = {
      length = length([1, 2, 3, 4, 5])
      element = element(["a", "b", "c"], 1)
      contains = contains(["a", "b", "c"], "b")
      keys = keys({a = 1, b = 2, c = 3})
      values = values({a = 1, b = 2, c = 3})
      lookup = lookup({a = 1, b = 2, c = 3}, "d", 0)
      zipmap = zipmap(["a", "b", "c"], [1, 2, 3])
      flatten = flatten([[1, 2], [3, 4], [5, 6]])
      compact = compact(["a", "", "b", null, "c"])
      distinct = distinct(["a", "b", "a", "c", "b"])
      chunklist = chunklist(["a", "b", "c", "d", "e"], 2)
      merge = merge(
        {a = 1, b = 2},
        {b = 3, c = 4},
        {c = 5, d = 6}
      )
    }
    
    // Encoding functions
    encoding = {
      jsonencode = jsonencode({
        string = "value"
        number = 42
        bool = true
        null_val = null
        list = [1, 2, 3]
        object = {
          nested = "value"
        }
      })
      jsondecode = jsondecode("{\"key\":\"value\"}")
      base64encode = base64encode("Hello, World!")
      base64decode = base64decode("SGVsbG8sIFdvcmxkIQ==")
      urlencode = urlencode("Hello, World!")
      yamlencode = yamlencode({
        string = "value"
        number = 42
        bool = true
        null_val = null
        list = [1, 2, 3]
        object = {
          nested = "value"
        }
      })
      yamldecode = yamldecode("key: value")
    }
    
    // Filesystem functions
    filesystem = {
      file = file("${path.module}/file.txt")
      fileexists = fileexists("${path.module}/file.txt")
      fileset = fileset("${path.module}", "*.txt")
      filebase64 = filebase64("${path.module}/file.txt")
      filebase64sha256 = filebase64sha256("${path.module}/file.txt")
      filemd5 = filemd5("${path.module}/file.txt")
      filesha1 = filesha1("${path.module}/file.txt")
      filesha256 = filesha256("${path.module}/file.txt")
      filesha512 = filesha512("${path.module}/file.txt")
    }
    
    // IP network functions
    ip = {
      cidrsubnet = cidrsubnet("10.0.0.0/16", 8, 1)
      cidrhost = cidrhost("10.0.0.0/24", 5)
      cidrnetmask = cidrnetmask("10.0.0.0/16")
      cidrsubnets = cidrsubnets("10.0.0.0/16", 8, 8, 8, 8)
    }
  }
}

// Comments in unusual places
resource /*comment*/ "aws_instance" /*comment*/ "comments_everywhere" /*comment*/ {
  /*comment*/ ami /*comment*/ = /*comment*/ "ami-12345678" /*comment*/
  /*comment*/ instance_type /*comment*/ = /*comment*/ "t2.micro" /*comment*/
  
  tags /*comment*/ = /*comment*/ {
    /*comment*/ Name /*comment*/ = /*comment*/ "test" /*comment*/
  } /*comment*/
} /*comment*/

// Unicode characters in strings and identifiers
locals {
  unicode = {
    // Unicode in strings
    emoji = "😀 😃 😄 😁 😆 😅 😂 🤣 🥲 ☺️"
    international = "你好 こんにちは 안녕하세요 Привет Γειά σου"
    symbols = "★ ☆ ✮ ✯ ☄ ☾ ☽ ☼ ☀ ☁ ☂ ☃ ☻ ☺ ☹"
    math = "∀ ∁ ∂ ∃ ∄ ∅ ∆ ∇ ∈ ∉ ∊ ∋ ∌ ∍ ∎ ∏ ∐ ∑ − ∓ ∔ ∕ ∖ ∗ ∘ ∙ √ ∛ ∜ ∝ ∞ ∟ ∠ ∡ ∢ ∣ ∤ ∥ ∦ ∧ ∨ ∩ ∪ ∫ ∬ ∭ ∮ ∯ ∰ ∱ ∲ ∳ ∴ ∵ ∶ ∷ ∸ ∹ ∺ ∻ ∼ ∽ ∾ ∿"
    
    // Unicode in keys
    "你好" = "Hello in Chinese"
    "こんにちは" = "Hello in Japanese"
    "안녕하세요" = "Hello in Korean"
    "Привет" = "Hello in Russian"
    "Γειά σου" = "Hello in Greek"
  }
}

// Extreme whitespace variations
resource    "aws_instance"     "whitespace"    {
  
  
  ami           =           "ami-12345678"
  
  
  instance_type=    "t2.micro"
  tags={
    Name="whitespace-test"
  }
  
  
  
}

// Minimal whitespace
locals{a="no whitespace"b=1+2+3+4+5 c=[1,2,3,4,5]d={"key"="value"}e=true?1:0}

// Unusual but valid syntax for blocks
resource "aws_instance" "unusual_block_syntax" {name = "test" ami = "ami-12345678" instance_type = "t2.micro"}

// Multiple blocks on one line
resource "null_resource" "one" {} resource "null_resource" "two" {} resource "null_resource" "three" {}

// Unusual but valid syntax for attributes
locals {
  one_line_object = {key1 = "value1" key2 = "value2" key3 = "value3"}
  one_line_list = ["item1", "item2", "item3", "item4", "item5"]
  attribute_groups = {(1 + 1) = "two" ("key_" + "name") = "key_name" (upper("name")) = "NAME"}
}

// Unusual but valid syntax for expressions
locals {
  unusual_expressions = {
    parentheses_everywhere = (((1))) + (((2))) * (((3)))
    empty_structures = {
      empty_string = ""
      empty_list = []
      empty_map = {}
      empty_tuple = []
    }
    operator_chains = 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10
    mixed_operators = 1 + 2 * 3 - 4 / 5 % 6
    complex_conditions = true && false || true && !false || (true && (false || true))
    ternary_chains = true ? 1 : false ? 2 : null != null ? 3 : 4
  }
}

// Unusual but valid syntax for references
locals {
  references = {
    self_reference = local.references
    module_reference = module.nonexistent_module.output
    data_reference = data.nonexistent_data.nonexistent_resource.attribute
    resource_reference = aws_instance.nonexistent_resource.id
    provider_reference = aws.alternate.region
    terraform_reference = terraform.workspace
    path_reference = path.module
    var_reference = var.nonexistent_variable
  }
}

// Unusual but valid syntax for functions
locals {
  unusual_functions = {
    zero_args = timestamp()
    one_arg = abs(-42)
    many_args = format("%s %s %s %s %s", "one", "two", "three", "four", "five")
    nested = upper(lower(title(upper(lower("MIXED case TEXT")))))
  }
}

// Unusual but valid syntax for interpolation
locals {
  unusual_interpolation = {
    empty = "${"}"
    nested = "${"${true ? "yes" : "no"}"}"
    multiple = "prefix-${1 + 1}-${upper("middle")}-${3 * 3}-suffix"
    complex = "complex-${true ? "yes-${1 + 1}" : "no-${2 + 2}"}-suffix"
  }
}

// Unusual but valid syntax for comments
locals {
  // Comment before attribute
  comment_before = "value"
  
  comment_after = "value" // Comment after attribute
  
  /* Block comment before */
  block_comment_before = "value"
  
  block_comment_after = "value" /* Block comment after */
  
  /* Block comment
     spanning
     multiple
     lines */
  multi_line_comment = "value"
  
  // Comment with unusual characters: !@#$%^&*()_+-=[]{}|;:'",.<>/?`~
  unusual_comment_chars = "value"
  
  // Comment with Unicode: 你好 こんにちは 안녕하세요 Привет Γειά σου
  unicode_comment = "value"
  
  // Comment with emoji: 😀 😃 😄 😁 😆 😅 😂 🤣 🥲 ☺️
  emoji_comment = "value"
}

// Edge case: Empty provider block
provider "aws" {}

// Edge case: Provider with alias
provider "aws" {
  alias = "alternate"
}

// Edge case: Provider with multiple aliases
provider "aws" {
  alias = "us-east-1"
}

provider "aws" {
  alias = "us-west-1"
}

provider "aws" {
  alias = "eu-west-1"
}

// Edge case: Terraform block with multiple required_providers
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">= 3.0.0, < 4.0.0"
    }
    
    azurerm = {
      source = "hashicorp/azurerm"
      version = ">= 2.0.0"
    }
    
    google = {
      source = "hashicorp/google"
      version = ">= 3.0.0"
    }
    
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    
    helm = {
      source = "hashicorp/helm"
      version = ">= 2.0.0"
    }
  }
}

// Edge case: Multiple backends (only one would be valid in real code)
terraform {
  backend "s3" {
    bucket = "terraform-state"
    key    = "edge-cases/terraform.tfstate"
    region = "us-east-1"
  }
  
  backend "azurerm" {
    resource_group_name  = "terraform-state"
    storage_account_name = "terraformstate"
    container_name       = "tfstate"
    key                  = "edge-cases.tfstate"
  }
  
  backend "gcs" {
    bucket = "terraform-state"
    prefix = "edge-cases"
  }
  
  backend "remote" {
    organization = "example-org"
    
    workspaces {
      name = "edge-cases"
    }
  }
}

// Edge case: Multiple lifecycle blocks (only one would be valid in real code)
resource "aws_instance" "multiple_lifecycles" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  lifecycle {
    create_before_destroy = true
  }
  
  lifecycle {
    prevent_destroy = true
  }
  
  lifecycle {
    ignore_changes = [tags]
  }
}

// Edge case: Multiple provisioner blocks
resource "aws_instance" "multiple_provisioners" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  provisioner "local-exec" {
    command = "echo 'Hello, World!'"
  }
  
  provisioner "remote-exec" {
    inline = [
      "echo 'Hello, World!'",
      "echo 'This is a test.'"
    ]
  }
  
  provisioner "file" {
    source      = "local/path"
    destination = "remote/path"
  }
  
  provisioner "local-exec" {
    when    = destroy
    command = "echo 'Goodbye, World!'"
  }
}

// Edge case: Multiple connection blocks (only one would be valid in real code)
resource "aws_instance" "multiple_connections" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  connection {
    type     = "ssh"
    user     = "ubuntu"
    password = "password"
    host     = self.public_ip
  }
  
  connection {
    type     = "winrm"
    user     = "administrator"
    password = "password"
    host     = self.public_ip
  }
}

// Edge case: Depends on with multiple resources
resource "aws_instance" "depends_on_multiple" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  depends_on = [
    aws_vpc.example,
    aws_subnet.example,
    aws_security_group.example,
    aws_key_pair.example,
    aws_iam_role.example
  ]
}
