// This file contains edge cases and unusual syntax patterns to test the parser

// Empty resource block
resource "null_resource" "empty" {
}

// Resource with only comments
resource "null_resource" "comments_only" {
  # This is a comment
  // This is another comment
  /* This is a block comment */
}

// Extremely long line
locals {
  extremely_long_line = "This is an extremely long line that exceeds the typical line length limit. It's designed to test how the parser handles very long lines. The line continues for a while to ensure it's long enough to potentially cause issues with buffers or other limitations in the parser. It includes some special characters like quotes (\"), backslashes (\\), and other potentially problematic characters: !@#$%^&*()_+-=[]{}|;:,.<>?/"
}

// Unusual whitespace
resource    "aws_instance"     "unusual_whitespace"    {
  ami           =     "ami-12345678"
  instance_type=    "t2.micro"
  
  tags={
    Name    =    "unusual-whitespace"
    Environment="test"
  }
}

// Nested conditionals
locals {
  nested_conditionals = true ? (
    false ? "a" : (
      true ? "b" : (
        false ? "c" : (
          true ? "d" : "e"
        )
      )
    )
  ) : "f"
}

// Multiple blocks with the same name
variable "duplicate" {
  type = string
  default = "first"
}

variable "duplicate" {
  type = number
  default = 123
}

// Unicode characters
locals {
  unicode = "こんにちは世界 • Hello, World! • Привет, мир! • مرحبا بالعالم • 你好，世界！"
}

// Escaped sequences
locals {
  escaped = "Line 1\nLine 2\tTabbed\r\nWindows line ending\\\\ Double backslash \" Quote"
}

// Empty strings and special values
locals {
  empty_string = ""
  single_space = " "
  just_newline = "\n"
  null_value   = null
}

// Unusual numbers
locals {
  zero = 0
  negative = -42
  decimal = 3.14159265359
  scientific = 1.23e45
  hex = 0xDEADBEEF
  octal = 0o755
  binary = 0b10101010
}

// Special characters in identifiers
resource "aws_instance" "special-chars_in.identifier" {
  ami = "ami-12345678"
  instance_type = "t2.micro"
}

// Provider with unusual configuration
provider "aws" {
  region = "us-west-2"
  alias = "unusual"
  
  assume_role {
    role_arn = "arn:aws:iam::123456789012:role/unusual-role"
    session_name = "unusual-session"
  }
  
  default_tags {
    tags = {
      ManagedBy = "Terraform"
      Environment = "Test"
    }
  }
  
  ignore_tags {
    keys = ["IgnoreMe", "AlsoIgnoreMe"]
    key_prefixes = ["temp-", "tmp-"]
  }
}

// Unusual block nesting
resource "aws_security_group" "nested_blocks" {
  name = "nested-blocks"
  
  ingress {
    description = "First level"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  dynamic "egress" {
    for_each = ["one", "two", "three"]
    content {
      description = "Dynamic ${egress.value}"
      from_port = 0
      to_port = 0
      protocol = "-1"
      
      cidr_blocks = ["0.0.0.0/0"]
    }
  }
  
  lifecycle {
    create_before_destroy = true
    
    precondition {
      condition = length(var.allowed_cidrs) > 0
      error_message = "At least one CIDR block must be allowed."
    }
    
    postcondition {
      condition = self.name != ""
      error_message = "Name cannot be empty."
    }
  }
}

// Unusual function calls
locals {
  unusual_functions = {
    nested = merge(
      {
        a = "value"
      },
      {
        b = lookup(
          {
            x = "x-value"
            y = "y-value"
          },
          "z",
          "default"
        )
      }
    )
    
    chained = join(",", concat(split(",", "a,b,c"), ["d", "e", "f"]))
    
    complex = formatlist(
      "%s = %s",
      keys({
        key1 = "value1"
        key2 = "value2"
      }),
      values({
        key1 = "value1"
        key2 = "value2"
      })
    )
  }
}

// Comments in unusual places
locals /* comment */ {
  # Comment
  value1 /* inline comment */ = "test" # end of line comment
  value2 = /* another inline comment */ "test2"
}

// Unusual heredoc
locals {
  unusual_heredoc = <<UNUSUAL
This is a heredoc with unusual content.
It has multiple lines.
It includes special characters: !@#$%^&*()_+-=[]{}|;:,.<>?/
It includes quotes: " and '
It includes backslashes: \\ and \n and \t
It includes unicode: こんにちは世界
UNUSUAL

  indented_heredoc = <<-INDENTED
    This is an indented heredoc.
    The leading whitespace will be trimmed.
      This line has extra indentation.
    Back to normal indentation.
  INDENTED
}

// Unusual for expressions
locals {
  unusual_for = [
    for i, v in ["a", "b", "c"] : {
      index = i
      value = v
      upper = upper(v)
    }
    if v != "b"
  ]
  
  nested_for = {
    for k1, v1 in {
      a = ["x", "y", "z"]
      b = ["p", "q", "r"]
    } : k1 => {
      for i, v2 in v1 : i => "${k1}-${v2}"
    }
  }
}

// Unusual splat expressions
locals {
  unusual_splat = [
    {
      name = "a"
      value = 1
    },
    {
      name = "b"
      value = 2
    },
    {
      name = "c"
      value = 3
    }
  ][*].name
  
  nested_splat = [
    {
      name = "a"
      values = [1, 2, 3]
    },
    {
      name = "b"
      values = [4, 5, 6]
    }
  ][*].values[*]
}
