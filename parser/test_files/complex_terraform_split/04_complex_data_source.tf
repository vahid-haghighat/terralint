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