package parser

import "github.com/vahid-haghighat/terralint/parser/types"

// createSimpleTerraformExpected creates the expected structure for simple_test.tf
func createSimpleTerraformExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:   "resource",
				Labels: []string{"aws_instance", "example"},
				Children: []types.Body{
					&types.Attribute{
						Name: "ami",
						Value: &types.LiteralValue{
							Value:     "ami-12345678",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "instance_type",
						Value: &types.LiteralValue{
							Value:     "t2.micro",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Name"},
									},
									Value: &types.LiteralValue{
										Value:     "example-instance",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Environment"},
									},
									Value: &types.LiteralValue{
										Value:     "test",
										ValueType: "string",
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "variable",
				Labels:       []string{"region"},
				BlockComment: "// Variable block",
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "AWS region",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "type",
						Value: &types.LiteralValue{
							Value:     "string",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "default",
						Value: &types.LiteralValue{
							Value:     "us-west-2",
							ValueType: "string",
						},
					},
				},
			},
			&types.Block{
				Type:         "output",
				Labels:       []string{"instance_id"},
				BlockComment: "// Output block",
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "ID of the EC2 instance",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_instance", "example", "id"},
						},
					},
				},
			},
			&types.Block{
				Type:         "data",
				Labels:       []string{"aws_ami", "ubuntu"},
				BlockComment: "// Data source block",
				Children: []types.Body{
					&types.Attribute{
						Name: "most_recent",
						Value: &types.LiteralValue{
							Value:     true,
							ValueType: "bool",
						},
					},
					&types.Block{
						Type: "filter",
						Children: []types.Body{
							&types.Attribute{
								Name: "name",
								Value: &types.LiteralValue{
									Value:     "name",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "values",
								Value: &types.ReferenceExpr{
									Parts: []string{"[\"ubuntu/images/hvm-ssd/ubuntu-focal-20", "04-amd64-server-*\"]"},
								},
							},
						},
					},
					&types.Block{
						Type: "filter",
						Children: []types.Body{
							&types.Attribute{
								Name: "name",
								Value: &types.LiteralValue{
									Value:     "virtualization-type",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "values",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
										&types.LiteralValue{
											Value:     "hvm",
											ValueType: "string",
										},
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "owners",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
								&types.LiteralValue{
									Value:     "099720109477",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "provider",
				Labels:       []string{"aws"},
				BlockComment: "// Provider block",
				Children: []types.Body{
					&types.Attribute{
						Name: "region",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "region"},
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Locals block",
				Children: []types.Body{
					&types.Attribute{
						Name: "common_tags",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Project"},
									},
									Value: &types.LiteralValue{
										Value:     "Test",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Owner"},
									},
									Value: &types.LiteralValue{
										Value:     "Terraform",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Environment"},
									},
									Value: &types.LiteralValue{
										Value:     "Test",
										ValueType: "string",
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "module",
				Labels:       []string{"vpc"},
				BlockComment: "// Module block",
				Children: []types.Body{
					&types.Attribute{
						Name: "source",
						Value: &types.LiteralValue{
							Value:     "terraform-aws-modules/vpc/aws",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "name",
						Value: &types.LiteralValue{
							Value:     "my-vpc",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "cidr",
						Value: &types.LiteralValue{
							Value:     "10.0.0.0/16",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "azs",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
								&types.LiteralValue{
									Value:     "us-west-2a",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "us-west-2b",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "us-west-2c",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
							},
						},
					},
					&types.Attribute{
						Name: "private_subnets",
						Value: &types.ReferenceExpr{
							Parts: []string{"[\"10", "0", "1", "0/24\", \"10", "0", "2", "0/24\", \"10", "0", "3", "0/24\"]"},
						},
					},
					&types.Attribute{
						Name: "public_subnets",
						Value: &types.ReferenceExpr{
							Parts: []string{"[\"10", "0", "101", "0/24\", \"10", "0", "102", "0/24\", \"10", "0", "103", "0/24\"]"},
						},
					},
					&types.Attribute{
						Name: "enable_nat_gateway",
						Value: &types.LiteralValue{
							Value:     true,
							ValueType: "bool",
						},
					},
					&types.Attribute{
						Name: "enable_vpn_gateway",
						Value: &types.LiteralValue{
							Value:     true,
							ValueType: "bool",
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
		},
	}
}

// createComplexTerraformExpected creates the expected structure for complex_terraform_test.tf
func createComplexTerraformExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:   "module",
				Labels: []string{"complex_module"},
				Children: []types.Body{
					&types.Attribute{
						Name: "source",
						Value: &types.LiteralValue{
							Value:     "terraform-aws-modules/vpc/aws",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "version",
						Value: &types.LiteralValue{
							Value:     "~> 3.0",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "vpc_config",
						Value: &types.ForExpr{
							ValueVar: "item",
							ThenKeyExpr: &types.ReferenceExpr{
								Parts: []string{"(var", "enable_dns != null ? var", "enable_dns : false)\n    \n    // Complex object with nested expressions\n    tags = merge(\n      var", "common_tags,\n      {\n        Name        = \"complex-${var", "environment}-vpc\"\n        Environment = var", "environment\n        ManagedBy   = \"terraform\"\n        Complex     = jsonencode({\n          nested = \"value\"\n          list   = [1, 2, 3, 4]\n          map    = {\n            key1 = \"value1\"\n            key2 = 42\n          }\n        })\n      },\n      var", "additional_tags != null ? var", "additional_tags : {}\n    )"},
							},
						},
						BlockComment: "// Complex map with nested objects, expressions, and functions",
					},
					&types.Attribute{
						Name: "subnet_cidrs",
						Value: &types.ForExpr{
							ValueVar: "subnet",
							KeyVar:   "i",
							Collection: &types.ReferenceExpr{
								Parts: []string{"var", "subnets"},
							},
							ThenKeyExpr: &types.FunctionCallExpr{
								Name: "cidrsubnet",
								Args: []types.Expression{},
							},
						},
						BlockComment: "// Complex for expression with filtering and transformation",
					},
					&types.Attribute{
						Name: "subnet_configs",
						Value: &types.ForExpr{
							ValueVar: "zone",
							KeyVar:   "zone_key",
							Collection: &types.ReferenceExpr{
								Parts: []string{"var", "availability_zones"},
							},
							ThenKeyExpr: &types.ReferenceExpr{
								Parts: []string{"zone_key"},
							},
							ThenValueExpr: &types.ObjectExpr{
								Items: []types.ObjectItem{
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"subnet_key"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"> {"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"cidr"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"cidrsubnet("},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"az"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"zone"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"tags"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"merge("},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"Name"},
										},
										Value: &types.LiteralValue{
											Value:     "${var.environment}-${subnet}-${zone}",
											ValueType: "string",
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"Type"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet"},
										},
									},
								},
							},
						},
						BlockComment: "// Nested for expressions with conditional",
					},
					&types.Attribute{
						Name: "all_subnet_ids",
						Value: &types.ReferenceExpr{
							Parts: []string{"flatten([\n    for zone_key, zone in aws_subnet", "main : [\n      for subnet_key, subnet in zone : subnet", "id\n    ]\n  ])"},
						},
						BlockComment: "// Complex splat expression",
					},
					&types.Attribute{
						Name: "user_data",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "<<-EOT\n    #!/bin/bash\n    echo \"Environment: ${var.environment}\"\n    echo \"Region: ${data.aws_region.current.name}\"\n    \n    # Complex interpolation\n    ${join(\"\\n\", [\n      for script in var.bootstrap_scripts :\n      \"source ${script}\"\n    ])}\n    \n    # Conditional section\n    ${var.install_monitoring ? \"setup_monitoring ${var.monitoring_endpoint}\" : \"echo 'Monitoring disabled'\"}\n    \n    # For loop in heredoc\n    %{for pkg in var.packages~}\n    yum install -y ${pkg}\n    %{endfor~}\n    \n    # If directive in heredoc\n    %{if var.environment == \"prod\"~}\n    echo \"Production environment detected, applying strict security\"\n    %{else~}\n    echo \"Non-production environment, using standard security\"\n    %{endif~}\n  EOT",
									ValueType: "string",
								},
							},
						},
						BlockComment: "// Heredoc with interpolation",
					},
					&types.Attribute{
						Name: "timeout",
						Value: &types.ReferenceExpr{
							Parts: []string{"(\n    var", "environment == \"prod\" ? 300 :\n    var", "environment == \"staging\" ? 180 :\n    var", "environment == \"dev\" ? 60 : 30\n  ) + (var", "additional_timeout != null ? var", "additional_timeout : 0)"},
						},
						BlockComment: "// Complex binary expressions with nested conditionals",
					},
					&types.Attribute{
						Name: "complex_calculation",
						Value: &types.ReferenceExpr{
							Parts: []string{"(\n    (var", "base_value * (1 + var", "multiplier)) / \n    (var", "divisor > 0 ? var", "divisor : 1)\n  ) + (\n    var", "environment == \"prod\" ? \n    (var", "prod_adjustment * (1 + var", "prod_factor)) : \n    (var", "non_prod_adjustment * (1 - var", "non_prod_factor))\n  )"},
						},
						BlockComment: "// Nested parentheses and operators",
					},
					&types.Attribute{
						Name: "security_groups",
						Value: &types.ReferenceExpr{
							Parts: []string{"compact(concat(\n    [\n      var", "create_default_security_group ? aws_security_group", "default[0]", "id : \"\",\n      var", "create_bastion_security_group ? aws_security_group", "bastion[0]", "id : \"\"\n    ],\n    var", "additional_security_group_ids\n  ))"},
						},
						BlockComment: "// Complex function calls with nested expressions",
					},
					&types.Attribute{
						Name: "custom_template",
						Value: &types.FunctionCallExpr{
							Name: "templatefile",
							Args: []types.Expression{
								&types.TemplateExpr{
									Parts: []types.Expression{
										&types.LiteralValue{
											Value:     "\"${path.module}/templates/config.tpl\"",
											ValueType: "string",
										},
									},
								},
								&types.ForExpr{
									ValueVar: "subnet",
									Collection: &types.ReferenceExpr{
										Parts: []string{"aws_subnet", "main"},
									},
									ThenKeyExpr: &types.ReferenceExpr{
										Parts: []string{"{\n        id   = subnet", "id\n        cidr = subnet", "cidr_block\n        az   = subnet", "availability_zone\n      }\n    ]\n    features = {\n      monitoring = var", "enable_monitoring\n      logging    = var", "enable_logging\n      encryption = var", "environment == \"prod\" ? true : var", "enable_encryption\n    }"},
									},
								},
							},
						},
						BlockComment: "// Template with directives",
					},
					&types.Block{
						Type:         "dynamic",
						Labels:       []string{"ingress"},
						BlockComment: "// Dynamic blocks with complex expressions",
						Children: []types.Body{
							&types.Attribute{
								Name: "for_each",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "ingress_rules"},
								},
							},
							&types.Attribute{
								Name: "iterator",
								Value: &types.ReferenceExpr{
									Parts: []string{"rule"},
								},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "description",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{
													Parts: []string{"rule"},
												},
												&types.LiteralValue{
													Value:     "description",
													ValueType: "string",
												},
												&types.TemplateExpr{
													Parts: []types.Expression{
														&types.LiteralValue{
															Value:     "\"Ingress Rule ${rule.key}\"",
															ValueType: "string",
														},
													},
												},
											},
										},
									},
									&types.Attribute{
										Name: "from_port",
										Value: &types.ReferenceExpr{
											Parts: []string{"rule", "value", "from_port"},
										},
									},
									&types.Attribute{
										Name: "to_port",
										Value: &types.ReferenceExpr{
											Parts: []string{"rule", "value", "to_port"},
										},
									},
									&types.Attribute{
										Name: "protocol",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(rule", "value, \"protocol\", \"tcp\")"},
										},
									},
									&types.Attribute{
										Name: "cidr_blocks",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(rule", "value, \"cidr_blocks\", null) != null ? rule", "value", "cidr_blocks : [\n        for cidr in var", "default_cidrs :\n        cidr if !contains(var", "excluded_cidrs, cidr)\n      ]"},
										},
									},
									&types.Attribute{
										Name: "security_groups",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(rule", "value, \"security_group_ids\", [])"},
										},
									},
									&types.Attribute{
										Name: "self",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(rule", "value, \"self\", false)"},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "validation",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"condition"},
									},
									Value: &types.ReferenceExpr{
										Parts: []string{"can(regex(\"^(dev|staging|prod)$\", var", "environment))"},
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"error_message"},
									},
									Value: &types.LiteralValue{
										Value:     "Environment must be one of: dev, staging, prod.",
										ValueType: "string",
									},
								},
							},
						},
						BlockComment: "// Complex type constraints",
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_security_group", "complex"},
				BlockComment: "// Resource with complex dynamic blocks and for_each",
				Children: []types.Body{
					&types.Attribute{
						Name: "for_each",
						Value: &types.ForExpr{
							ValueVar: "sg",
							Collection: &types.ReferenceExpr{
								Parts: []string{"var", "security_groups"},
							},
							ThenKeyExpr: &types.ReferenceExpr{
								Parts: []string{"sg", "name"},
							},
							ThenValueExpr: &types.ReferenceExpr{
								Parts: []string{"sg if sg", "create"},
							},
						},
					},
					&types.Attribute{
						Name: "name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${var.prefix}-${each.key}\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "description",
						Value: &types.ReferenceExpr{
							Parts: []string{"each", "value", "description"},
						},
					},
					&types.Attribute{
						Name: "vpc_id",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "vpc_id"},
						},
					},
					&types.Block{
						Type:         "dynamic",
						Labels:       []string{"ingress"},
						BlockComment: "// Dynamic blocks with nested expressions",
						Children: []types.Body{
							&types.Attribute{
								Name: "for_each",
								Value: &types.ReferenceExpr{
									Parts: []string{"each", "value", "ingress_rules"},
								},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "description",
										Value: &types.ReferenceExpr{
											Parts: []string{"ingress", "value", "description"},
										},
									},
									&types.Attribute{
										Name: "from_port",
										Value: &types.ReferenceExpr{
											Parts: []string{"ingress", "value", "from_port"},
										},
									},
									&types.Attribute{
										Name: "to_port",
										Value: &types.ReferenceExpr{
											Parts: []string{"ingress", "value", "to_port"},
										},
									},
									&types.Attribute{
										Name: "protocol",
										Value: &types.ReferenceExpr{
											Parts: []string{"ingress", "value", "protocol"},
										},
									},
									&types.Attribute{
										Name: "cidr_blocks",
										Value: &types.ReferenceExpr{
											Parts: []string{"ingress", "value", "cidr_blocks"},
										},
									},
									&types.Attribute{
										Name: "ipv6_cidr_blocks",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(ingress", "value, \"ipv6_cidr_blocks\", [])"},
										},
									},
									&types.Attribute{
										Name: "prefix_list_ids",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(ingress", "value, \"prefix_list_ids\", [])"},
										},
									},
									&types.Attribute{
										Name: "security_groups",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(ingress", "value, \"security_groups\", [])"},
										},
									},
									&types.Attribute{
										Name: "self",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(ingress", "value, \"self\", false)"},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:   "dynamic",
						Labels: []string{"egress"},
						Children: []types.Body{
							&types.Attribute{
								Name: "for_each",
								Value: &types.ReferenceExpr{
									Parts: []string{"each", "value", "egress_rules"},
								},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "description",
										Value: &types.ReferenceExpr{
											Parts: []string{"egress", "value", "description"},
										},
									},
									&types.Attribute{
										Name: "from_port",
										Value: &types.ReferenceExpr{
											Parts: []string{"egress", "value", "from_port"},
										},
									},
									&types.Attribute{
										Name: "to_port",
										Value: &types.ReferenceExpr{
											Parts: []string{"egress", "value", "to_port"},
										},
									},
									&types.Attribute{
										Name: "protocol",
										Value: &types.ReferenceExpr{
											Parts: []string{"egress", "value", "protocol"},
										},
									},
									&types.Attribute{
										Name: "cidr_blocks",
										Value: &types.ReferenceExpr{
											Parts: []string{"egress", "value", "cidr_blocks"},
										},
									},
									&types.Attribute{
										Name: "ipv6_cidr_blocks",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(egress", "value, \"ipv6_cidr_blocks\", [])"},
										},
									},
									&types.Attribute{
										Name: "prefix_list_ids",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(egress", "value, \"prefix_list_ids\", [])"},
										},
									},
									&types.Attribute{
										Name: "security_groups",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(egress", "value, \"security_groups\", [])"},
										},
									},
									&types.Attribute{
										Name: "self",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(egress", "value, \"self\", false)"},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.FunctionCallExpr{
							Name: "merge",
							Args: []types.Expression{
								&types.ReferenceExpr{
									Parts: []string{"var"},
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"Name"},
											},
											Value: &types.TemplateExpr{
												Parts: []types.Expression{
													&types.LiteralValue{
														Value:     "\"${var.prefix}-${each.key}\"",
														ValueType: "string",
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"Type"},
											},
											Value: &types.LiteralValue{
												Value:     "SecurityGroup",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"Rules"},
											},
											Value: &types.FunctionCallExpr{
												Name: "jsonencode",
												Args: []types.Expression{
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"ingress"},
																},
																Value: &types.ReferenceExpr{
																	Parts: []string{"length(each", "value", "ingress_rules)"},
																},
															},
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"egress"},
																},
																Value: &types.ReferenceExpr{
																	Parts: []string{"length(each", "value", "egress_rules)"},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								&types.LiteralValue{
									Value:     "unknown_arg",
									ValueType: "string",
								},
							},
						},
						BlockComment: "// Complex tags with expressions and functions",
					},
					&types.Block{
						Type: "lifecycle",
						Children: []types.Body{
							&types.Attribute{
								Name: "create_before_destroy",
								Value: &types.LiteralValue{
									Value:     true,
									ValueType: "bool",
								},
							},
							&types.Attribute{
								Name: "prevent_destroy",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "environment == \"prod\" ? true : false"},
								},
							},
							&types.Attribute{
								Name: "ignore_changes",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
										&types.ReferenceExpr{
											Parts: []string{"tags"},
										},
										&types.ReferenceExpr{
											Parts: []string{"tags"},
										},
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Complex locals with nested expressions",
				Children: []types.Body{
					&types.Attribute{
						Name: "subnet_map",
						Value: &types.ForExpr{
							ValueVar: "subnet",
							Collection: &types.ReferenceExpr{
								Parts: []string{"var", "subnets"},
							},
							ThenKeyExpr: &types.ReferenceExpr{
								Parts: []string{"subnet", "name"},
							},
							ThenValueExpr: &types.ObjectExpr{
								Items: []types.ObjectItem{
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"id"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet", "id"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"cidr"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet", "cidr"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"az"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet", "availability_zone"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"public"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet", "public"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"nat_gw"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet", "public ? true : false"},
										},
									},
									{
										Key: &types.ReferenceExpr{
											Parts: []string{"depends_on"},
										},
										Value: &types.ReferenceExpr{
											Parts: []string{"subnet", "public ? [] : [for s in var", "subnets : s", "id if s", "public]"},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "filtered_instances",
						Value: &types.ForExpr{
							ValueVar: "server",
							Collection: &types.ReferenceExpr{
								Parts: []string{"var", "servers"},
							},
							ThenKeyExpr: &types.FunctionCallExpr{
								Name: "{\n      id          = server.id\n      name        = server.name\n      environment = server.environment\n      type        = server.type\n      subnet_id   = server.subnet_id\n      private_ip  = server.private_ip\n      public_ip   = server.public_ip\n      tags        = server.tags\n    }\n    if server.environment == var.environment &&\n    contains",
								Args: []types.Expression{},
							},
						},
						BlockComment: "// Nested for expressions with filtering",
					},
					&types.Attribute{
						Name: "backup_config",
						Value: &types.ConditionalExpr{
							Condition: &types.ReferenceExpr{
								Parts: []string{"var"},
							},
							TrueExpr: &types.ForExpr{
								ValueVar: "target",
								Collection: &types.ReferenceExpr{
									Parts: []string{"var", "backup_targets"},
								},
								ThenKeyExpr: &types.FunctionCallExpr{
									Name: "\"0 1 * * *\"\n    retention =",
									Args: []types.Expression{},
								},
							},
							FalseExpr: &types.LiteralValue{ValueType: "null"},
						},
						BlockComment: "// Complex conditional with multiple nested expressions",
					},
					&types.Attribute{
						Name: "naming_convention",
						Value: &types.ReferenceExpr{
							Parts: []string{"join(\"-\", compact([\n    var", "prefix,\n    var", "environment,\n    var", "region_code,\n    var", "name,\n    var", "suffix != \"\" ? var", "suffix : null\n  ]))"},
						},
						BlockComment: "// Complex string interpolation with functions",
					},
					&types.Attribute{
						Name: "timeout_seconds",
						Value: &types.ReferenceExpr{
							Parts: []string{"(\n    var", "custom_timeout != null ? var", "custom_timeout :\n    var", "environment == \"prod\" ? (\n      var", "high_availability ? 300 : 180\n    ) : (\n      var", "environment == \"staging\" ? 120 : 60\n    )\n  )"},
						},
						BlockComment: "// Nested ternary operators",
					},
					&types.Attribute{
						Name: "schema",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"type"},
									},
									Value: &types.LiteralValue{
										Value:     "object",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"properties"},
									},
									Value: &types.ObjectExpr{
										Items: []types.ObjectItem{
											{
												Key: &types.ReferenceExpr{
													Parts: []string{"id"},
												},
												Value: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key: &types.ReferenceExpr{
																Parts: []string{"type"},
															},
															Value: &types.LiteralValue{
																Value:     "string",
																ValueType: "string",
															},
														},
														{
															Key: &types.ReferenceExpr{
																Parts: []string{"pattern"},
															},
															Value: &types.LiteralValue{
																Value:     "^[a-zA-Z0-9-_]+$",
																ValueType: "string",
															},
														},
													},
												},
											},
											{
												Key: &types.ReferenceExpr{
													Parts: []string{"settings"},
												},
												Value: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key: &types.ReferenceExpr{
																Parts: []string{"type"},
															},
															Value: &types.LiteralValue{
																Value:     "object",
																ValueType: "string",
															},
														},
														{
															Key: &types.ReferenceExpr{
																Parts: []string{"properties"},
															},
															Value: &types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"enabled"},
																		},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"type"},
																					},
																					Value: &types.LiteralValue{
																						Value:     "boolean",
																						ValueType: "string",
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"timeout"},
																		},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"type"},
																					},
																					Value: &types.LiteralValue{
																						Value:     "number",
																						ValueType: "string",
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"retries"},
																		},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"type"},
																					},
																					Value: &types.LiteralValue{
																						Value:     "number",
																						ValueType: "string",
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"options"},
																		},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"type"},
																					},
																					Value: &types.LiteralValue{
																						Value:     "array",
																						ValueType: "string",
																					},
																				},
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"items"},
																					},
																					Value: &types.ObjectExpr{
																						Items: []types.ObjectItem{
																							{
																								Key: &types.ReferenceExpr{
																									Parts: []string{"type"},
																								},
																								Value: &types.LiteralValue{
																									Value:     "object",
																									ValueType: "string",
																								},
																							},
																							{
																								Key: &types.ReferenceExpr{
																									Parts: []string{"properties"},
																								},
																								Value: &types.ObjectExpr{
																									Items: []types.ObjectItem{
																										{
																											Key: &types.ReferenceExpr{
																												Parts: []string{"name"},
																											},
																											Value: &types.ObjectExpr{
																												Items: []types.ObjectItem{
																													{
																														Key: &types.ReferenceExpr{
																															Parts: []string{"type"},
																														},
																														Value: &types.LiteralValue{
																															Value:     "string",
																															ValueType: "string",
																														},
																													},
																												},
																											},
																										},
																										{
																											Key: &types.ReferenceExpr{
																												Parts: []string{"value"},
																											},
																											Value: &types.ObjectExpr{
																												Items: []types.ObjectItem{
																													{
																														Key: &types.ReferenceExpr{
																															Parts: []string{"type"},
																														},
																														Value: &types.LiteralValue{
																															Value:     "string",
																															ValueType: "string",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											{
												Key: &types.ReferenceExpr{
													Parts: []string{"tags"},
												},
												Value: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key: &types.ReferenceExpr{
																Parts: []string{"type"},
															},
															Value: &types.LiteralValue{
																Value:     "object",
																ValueType: "string",
															},
														},
														{
															Key: &types.ReferenceExpr{
																Parts: []string{"additionalProperties"},
															},
															Value: &types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"type"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Complex type expressions",
					},
				},
			},
			&types.Block{
				Type:         "data",
				Labels:       []string{"aws_iam_policy_document", "complex"},
				BlockComment: "// Data source with complex expressions",
				Children: []types.Body{
					&types.Block{
						Type:   "dynamic",
						Labels: []string{"statement"},
						Children: []types.Body{
							&types.Attribute{
								Name: "for_each",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "policy_statements"},
								},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "sid",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{
													Parts: []string{"statement"},
												},
												&types.LiteralValue{
													Value:     "sid",
													ValueType: "string",
												},
												&types.TemplateExpr{
													Parts: []types.Expression{
														&types.LiteralValue{
															Value:     "\"Statement${statement.key}\"",
															ValueType: "string",
														},
													},
												},
											},
										},
									},
									&types.Attribute{
										Name: "effect",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(statement", "value, \"effect\", \"Allow\")"},
										},
									},
									&types.Attribute{
										Name: "actions",
										Value: &types.ReferenceExpr{
											Parts: []string{"statement", "value", "actions"},
										},
									},
									&types.Attribute{
										Name: "resources",
										Value: &types.ReferenceExpr{
											Parts: []string{"statement", "value", "resources"},
										},
									},
									&types.Block{
										Type:   "dynamic",
										Labels: []string{"condition"},
										Children: []types.Body{
											&types.Attribute{
												Name: "for_each",
												Value: &types.ReferenceExpr{
													Parts: []string{"lookup(statement", "value, \"conditions\", [])"},
												},
											},
											&types.Block{
												Type: "content",
												Children: []types.Body{
													&types.Attribute{
														Name: "test",
														Value: &types.ReferenceExpr{
															Parts: []string{"condition", "value", "test"},
														},
													},
													&types.Attribute{
														Name: "variable",
														Value: &types.ReferenceExpr{
															Parts: []string{"condition", "value", "variable"},
														},
													},
													&types.Attribute{
														Name: "values",
														Value: &types.ReferenceExpr{
															Parts: []string{"condition", "value", "values"},
														},
													},
												},
											},
										},
									},
									&types.Block{
										Type:   "dynamic",
										Labels: []string{"principals"},
										Children: []types.Body{
											&types.Attribute{
												Name: "for_each",
												Value: &types.ReferenceExpr{
													Parts: []string{"lookup(statement", "value, \"principals\", [])"},
												},
											},
											&types.Block{
												Type: "content",
												Children: []types.Body{
													&types.Attribute{
														Name: "type",
														Value: &types.ReferenceExpr{
															Parts: []string{"principals", "value", "type"},
														},
													},
													&types.Attribute{
														Name: "identifiers",
														Value: &types.ReferenceExpr{
															Parts: []string{"principals", "value", "identifiers"},
														},
													},
												},
											},
										},
									},
									&types.Attribute{
										Name: "not_actions",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(statement", "value, \"not_actions\", [])"},
										},
									},
									&types.Attribute{
										Name: "not_resources",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(statement", "value, \"not_resources\", [])"},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:         "statement",
						BlockComment: "// Override with inline statement",
						Children: []types.Body{
							&types.Attribute{
								Name: "sid",
								Value: &types.LiteralValue{
									Value:     "ExplicitAllow",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "effect",
								Value: &types.LiteralValue{
									Value:     "Allow",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "actions",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
										&types.LiteralValue{
											Value:     "s3:GetObject",
											ValueType: "string",
										},
										&types.LiteralValue{
											Value:     "s3:ListBucket",
											ValueType: "string",
										},
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
									},
								},
							},
							&types.Attribute{
								Name: "resources",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
										&types.TemplateExpr{
											Parts: []types.Expression{
												&types.LiteralValue{
													Value:     "\"arn:aws:s3:::${var.bucket_name}\"",
													ValueType: "string",
												},
											},
										},
										&types.TemplateExpr{
											Parts: []types.Expression{
												&types.LiteralValue{
													Value:     "\"arn:aws:s3:::${var.bucket_name}/*\"",
													ValueType: "string",
												},
											},
										},
										&types.LiteralValue{
											Value:     "",
											ValueType: "null",
										},
									},
								},
							},
							&types.Block{
								Type: "condition",
								Children: []types.Body{
									&types.Attribute{
										Name: "test",
										Value: &types.LiteralValue{
											Value:     "StringEquals",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "variable",
										Value: &types.LiteralValue{
											Value:     "aws:SourceVpc",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "values",
										Value: &types.ReferenceExpr{
											Parts: []string{"[var", "vpc_id]"},
										},
									},
								},
							},
							&types.Block{
								Type: "condition",
								Children: []types.Body{
									&types.Attribute{
										Name: "test",
										Value: &types.LiteralValue{
											Value:     "StringLike",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "variable",
										Value: &types.LiteralValue{
											Value:     "aws:PrincipalTag/Role",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "values",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
												&types.LiteralValue{
													Value:     "Admin",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "Developer",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
											},
										},
									},
								},
							},
							&types.Block{
								Type: "principals",
								Children: []types.Body{
									&types.Attribute{
										Name: "type",
										Value: &types.LiteralValue{
											Value:     "AWS",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "identifiers",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
												&types.TemplateExpr{
													Parts: []types.Expression{
														&types.LiteralValue{
															Value:     "\"arn:aws:iam::${var.account_id}:root\"",
															ValueType: "string",
														},
													},
												},
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "variable",
				Labels:       []string{"complex_object"},
				BlockComment: "// Variable with complex type constraints and validations",
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "A complex object with nested types and validations",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "type",
						Value: &types.FunctionCallExpr{
							Name: "object",
							Args: []types.Expression{
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"name"},
											},
											Value: &types.LiteralValue{
												Value:     "string",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"environment"},
											},
											Value: &types.LiteralValue{
												Value:     "string",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"enabled"},
											},
											Value: &types.LiteralValue{
												Value:     "bool",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"count"},
											},
											Value: &types.LiteralValue{
												Value:     "number",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"tags"},
											},
											Value: &types.FunctionCallExpr{
												Name: "map",
												Args: []types.Expression{
													&types.LiteralValue{
														Value:     "string",
														ValueType: "string",
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"vpc"},
											},
											Value: &types.FunctionCallExpr{
												Name: "object",
												Args: []types.Expression{
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"id"},
																},
																Value: &types.LiteralValue{
																	Value:     "string",
																	ValueType: "string",
																},
															},
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"cidr"},
																},
																Value: &types.LiteralValue{
																	Value:     "string",
																	ValueType: "string",
																},
															},
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"private_subnets"},
																},
																Value: &types.FunctionCallExpr{
																	Name: "list",
																	Args: []types.Expression{
																		&types.FunctionCallExpr{
																			Name: "object",
																			Args: []types.Expression{
																				&types.ObjectExpr{
																					Items: []types.ObjectItem{
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"id"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"cidr"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"az"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"public_subnets"},
																},
																Value: &types.FunctionCallExpr{
																	Name: "list",
																	Args: []types.Expression{
																		&types.FunctionCallExpr{
																			Name: "object",
																			Args: []types.Expression{
																				&types.ObjectExpr{
																					Items: []types.ObjectItem{
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"id"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"cidr"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"az"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"instances"},
											},
											Value: &types.FunctionCallExpr{
												Name: "list",
												Args: []types.Expression{
													&types.FunctionCallExpr{
														Name: "object",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"id"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"type"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"subnet_id"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"private_ip"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"public_ip"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "string",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"root_volume"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "object",
																			Args: []types.Expression{
																				&types.ObjectExpr{
																					Items: []types.ObjectItem{
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"size"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "number",
																								ValueType: "string",
																							},
																						},
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"type"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "string",
																								ValueType: "string",
																							},
																						},
																						{
																							Key: &types.ReferenceExpr{
																								Parts: []string{"encrypted"},
																							},
																							Value: &types.LiteralValue{
																								Value:     "bool",
																								ValueType: "string",
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"data_volumes"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "list",
																			Args: []types.Expression{
																				&types.FunctionCallExpr{
																					Name: "object",
																					Args: []types.Expression{
																						&types.ObjectExpr{
																							Items: []types.ObjectItem{
																								{
																									Key: &types.ReferenceExpr{
																										Parts: []string{"device_name"},
																									},
																									Value: &types.LiteralValue{
																										Value:     "string",
																										ValueType: "string",
																									},
																								},
																								{
																									Key: &types.ReferenceExpr{
																										Parts: []string{"size"},
																									},
																									Value: &types.LiteralValue{
																										Value:     "number",
																										ValueType: "string",
																									},
																								},
																								{
																									Key: &types.ReferenceExpr{
																										Parts: []string{"type"},
																									},
																									Value: &types.LiteralValue{
																										Value:     "string",
																										ValueType: "string",
																									},
																								},
																								{
																									Key: &types.ReferenceExpr{
																										Parts: []string{"encrypted"},
																									},
																									Value: &types.LiteralValue{
																										Value:     "bool",
																										ValueType: "string",
																									},
																								},
																								{
																									Key: &types.ReferenceExpr{
																										Parts: []string{"iops"},
																									},
																									Value: &types.FunctionCallExpr{
																										Name: "optional",
																										Args: []types.Expression{
																											&types.LiteralValue{
																												Value:     "number",
																												ValueType: "string",
																											},
																										},
																									},
																								},
																								{
																									Key: &types.ReferenceExpr{
																										Parts: []string{"throughput"},
																									},
																									Value: &types.FunctionCallExpr{
																										Name: "optional",
																										Args: []types.Expression{
																											&types.LiteralValue{
																												Value:     "number",
																												ValueType: "string",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"tags"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "map",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "string",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"databases"},
											},
											Value: &types.FunctionCallExpr{
												Name: "map",
												Args: []types.Expression{
													&types.FunctionCallExpr{
														Name: "object",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"engine"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"version"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"instance_class"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"allocated_storage"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "number",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"max_allocated_storage"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "number",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"multi_az"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "bool",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"backup_retention_period"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "number",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"parameters"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "map",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "string",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"subnet_ids"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "list",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "string",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"security_group_ids"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "list",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "string",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"endpoints"},
											},
											Value: &types.FunctionCallExpr{
												Name: "map",
												Args: []types.Expression{
													&types.FunctionCallExpr{
														Name: "object",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"service"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"vpc_endpoint_type"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "string",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"subnet_ids"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.FunctionCallExpr{
																					Name: "list",
																					Args: []types.Expression{
																						&types.LiteralValue{
																							Value:     "string",
																							ValueType: "string",
																						},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"security_group_ids"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.FunctionCallExpr{
																					Name: "list",
																					Args: []types.Expression{
																						&types.LiteralValue{
																							Value:     "string",
																							ValueType: "string",
																						},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"private_dns_enabled"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "bool",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"policy"},
																		},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.LiteralValue{
																					Value:     "string",
																					ValueType: "string",
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "complex_object", "count > 0 && var", "complex_object", "count <= 10"},
								},
							},
							&types.Attribute{
								Name: "error_message",
								Value: &types.LiteralValue{
									Value:     "Count must be between 1 and 10.",
									ValueType: "string",
								},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.ReferenceExpr{
									Parts: []string{"can(regex(\"^(dev|staging|prod)$\", var", "complex_object", "environment))"},
								},
							},
							&types.Attribute{
								Name: "error_message",
								Value: &types.LiteralValue{
									Value:     "Environment must be one of: dev, staging, prod.",
									ValueType: "string",
								},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.ReferenceExpr{
									Parts: []string{"length(var", "complex_object", "vpc", "private_subnets) > 0"},
								},
							},
							&types.Attribute{
								Name: "error_message",
								Value: &types.LiteralValue{
									Value:     "At least one private subnet must be defined.",
									ValueType: "string",
								},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.ReferenceExpr{
									Parts: []string{"alltrue([\n      for instance in var", "complex_object", "instances :\n      instance", "root_volume", "size >= 20 && instance", "root_volume", "encrypted == true\n    ])"},
								},
							},
							&types.Attribute{
								Name: "error_message",
								Value: &types.LiteralValue{
									Value:     "All root volumes must be at least 20GB and encrypted.",
									ValueType: "string",
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "output",
				Labels:       []string{"complex_output"},
				BlockComment: "// Output with complex expressions",
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "Complex output with nested expressions",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ForExpr{
							ValueVar: "sg",
							KeyVar:   "sg_key",
							Collection: &types.ReferenceExpr{
								Parts: []string{"aws_security_group", "complex"},
							},
							ThenKeyExpr: &types.FunctionCallExpr{
								Name: "sg.id\n    ]\n    instance_details = {\n      for instance in local.filtered_instances :\n      instance.id => {\n        name = instance.name\n        private_ip = instance.private_ip\n        public_ip = instance.public_ip\n        subnet = {\n          id = instance.subnet_id\n          details = lookup",
								Args: []types.Expression{},
							},
						},
					},
					&types.Attribute{
						Name: "sensitive",
						Value: &types.LiteralValue{
							Value:     true,
							ValueType: "bool",
						},
					},
					&types.Attribute{
						Name: "depends_on",
						Value: &types.ReferenceExpr{
							Parts: []string{"[\n    module", "complex_module,\n    aws_security_group", "complex,\n    data", "aws_iam_policy_document", "complex\n  ]"},
						},
					},
				},
			},
			&types.Block{
				Type:         "provider",
				Labels:       []string{"aws"},
				BlockComment: "// Provider configuration with complex expressions",
				Children: []types.Body{
					&types.Attribute{
						Name: "region",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "region"},
						},
					},
					&types.Block{
						Type: "assume_role",
						Children: []types.Body{
							&types.Attribute{
								Name: "role_arn",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "environment == \"prod\" ? var", "prod_role_arn : var", "non_prod_role_arn"},
								},
							},
							&types.Attribute{
								Name: "session_name",
								Value: &types.TemplateExpr{
									Parts: []types.Expression{
										&types.LiteralValue{
											Value:     "\"terraform-${var.environment}-${formatdate(\"YYYY-MM-DD-hh-mm-ss\", timestamp())}\"",
											ValueType: "string",
										},
									},
								},
							},
							&types.Attribute{
								Name: "external_id",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "external_id"},
								},
							},
						},
					},
					&types.Block{
						Type: "default_tags",
						Children: []types.Body{
							&types.Attribute{
								Name: "tags",
								Value: &types.FunctionCallExpr{
									Name: "merge",
									Args: []types.Expression{
										&types.ReferenceExpr{
											Parts: []string{"var"},
										},
										&types.ObjectExpr{
											Items: []types.ObjectItem{
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"Environment"},
													},
													Value: &types.ReferenceExpr{
														Parts: []string{"var", "environment"},
													},
												},
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"ManagedBy"},
													},
													Value: &types.LiteralValue{
														Value:     "Terraform",
														ValueType: "string",
													},
												},
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"Project"},
													},
													Value: &types.ReferenceExpr{
														Parts: []string{"var", "project_name"},
													},
												},
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"Owner"},
													},
													Value: &types.ReferenceExpr{
														Parts: []string{"var", "owner"},
													},
												},
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"CreatedAt"},
													},
													Value: &types.FunctionCallExpr{
														Name: "formatdate",
														Args: []types.Expression{
															&types.LiteralValue{
																Value:     "YYYY-MM-DD",
																ValueType: "string",
															},
															&types.FunctionCallExpr{
																Name: "timestamp",
																Args: []types.Expression{},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:   "dynamic",
						Labels: []string{"endpoints"},
						Children: []types.Body{
							&types.Attribute{
								Name: "for_each",
								Value: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left: &types.ReferenceExpr{
											Parts: []string{"var", "custom_endpoints"},
										},
										Operator: "!=",
										Right:    &types.LiteralValue{ValueType: "null"},
									},
									TrueExpr: &types.ReferenceExpr{
										Parts: []string{"var"},
									},
									FalseExpr: &types.ObjectExpr{
										Items: []types.ObjectItem{},
									},
								},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "apigateway",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"apigateway\", null)"},
										},
									},
									&types.Attribute{
										Name: "cloudformation",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"cloudformation\", null)"},
										},
									},
									&types.Attribute{
										Name: "cloudwatch",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"cloudwatch\", null)"},
										},
									},
									&types.Attribute{
										Name: "dynamodb",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"dynamodb\", null)"},
										},
									},
									&types.Attribute{
										Name: "ec2",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"ec2\", null)"},
										},
									},
									&types.Attribute{
										Name: "s3",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"s3\", null)"},
										},
									},
									&types.Attribute{
										Name: "sts",
										Value: &types.ReferenceExpr{
											Parts: []string{"lookup(endpoints", "value, \"sts\", null)"},
										},
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "terraform",
				BlockComment: "// Terraform configuration with complex expressions",
				Children: []types.Body{
					&types.Attribute{
						Name: "required_version",
						Value: &types.LiteralValue{
							Value:     ">= 1.0.0",
							ValueType: "string",
						},
					},
					&types.Block{
						Type: "required_providers",
						Children: []types.Body{
							&types.Attribute{
								Name: "aws",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"source"},
											},
											Value: &types.LiteralValue{
												Value:     "hashicorp/aws",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"version"},
											},
											Value: &types.LiteralValue{
												Value:     ">= 4.0.0, < 5.0.0",
												ValueType: "string",
											},
										},
									},
								},
							},
							&types.Attribute{
								Name: "random",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"source"},
											},
											Value: &types.LiteralValue{
												Value:     "hashicorp/random",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"version"},
											},
											Value: &types.LiteralValue{
												Value:     "~> 3.1",
												ValueType: "string",
											},
										},
									},
								},
							},
							&types.Attribute{
								Name: "null",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"source"},
											},
											Value: &types.LiteralValue{
												Value:     "hashicorp/null",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"version"},
											},
											Value: &types.LiteralValue{
												Value:     "~> 3.1",
												ValueType: "string",
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:   "backend",
						Labels: []string{"s3"},
						Children: []types.Body{
							&types.Attribute{
								Name: "bucket",
								Value: &types.TemplateExpr{
									Parts: []types.Expression{
										&types.LiteralValue{
											Value:     "\"terraform-state-${var.environment}-${var.account_id}\"",
											ValueType: "string",
										},
									},
								},
							},
							&types.Attribute{
								Name: "key",
								Value: &types.LiteralValue{
									Value:     "complex-module/terraform.tfstate",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "region",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "region"},
								},
							},
							&types.Attribute{
								Name: "dynamodb_table",
								Value: &types.LiteralValue{
									Value:     "terraform-locks",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "encrypt",
								Value: &types.LiteralValue{
									Value:     true,
									ValueType: "bool",
								},
							},
							&types.Attribute{
								Name: "kms_key_id",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "kms_key_id"},
								},
							},
							&types.Attribute{
								Name: "role_arn",
								Value: &types.ReferenceExpr{
									Parts: []string{"var", "state_role_arn"},
								},
							},
						},
					},
					&types.Attribute{
						Name: "experiments",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
								&types.ReferenceExpr{
									Parts: []string{"module_variable_optional_attrs"},
								},
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
							},
						},
					},
				},
			},
		},
	}
}

// createModuleExpected creates the expected structure for modules_test/main.tf
func createModuleExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:   "provider",
				Labels: []string{"aws"},
				Children: []types.Body{
					&types.Attribute{
						Name: "region",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "region"},
						},
					},
					&types.Block{
						Type: "default_tags",
						Children: []types.Body{
							&types.Attribute{
								Name: "tags",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"ManagedBy"},
											},
											Value: &types.LiteralValue{
												Value:     "Terraform",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"Project"},
											},
											Value: &types.ReferenceExpr{
												Parts: []string{"var", "project_name"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "terraform",
				BlockComment: "// Terraform configuration",
				Children: []types.Body{
					&types.Attribute{
						Name: "required_version",
						Value: &types.LiteralValue{
							Value:     ">= 1.0.0",
							ValueType: "string",
						},
					},
					&types.Block{
						Type: "required_providers",
						Children: []types.Body{
							&types.Attribute{
								Name: "aws",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"source"},
											},
											Value: &types.LiteralValue{
												Value:     "hashicorp/aws",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"version"},
											},
											Value: &types.LiteralValue{
												Value:     "~> 4.0",
												ValueType: "string",
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:   "backend",
						Labels: []string{"s3"},
						Children: []types.Body{
							&types.Attribute{
								Name: "bucket",
								Value: &types.LiteralValue{
									Value:     "terraform-state-bucket",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "key",
								Value: &types.LiteralValue{
									Value:     "modules-test/terraform.tfstate",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "region",
								Value: &types.LiteralValue{
									Value:     "us-west-2",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "dynamodb_table",
								Value: &types.LiteralValue{
									Value:     "terraform-locks",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "encrypt",
								Value: &types.LiteralValue{
									Value:     true,
									ValueType: "bool",
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Local variables",
				Children: []types.Body{
					&types.Attribute{
						Name: "environment_code",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"development"},
									},
									Value: &types.LiteralValue{
										Value:     "dev",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"staging"},
									},
									Value: &types.LiteralValue{
										Value:     "stg",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"production"},
									},
									Value: &types.LiteralValue{
										Value:     "prod",
										ValueType: "string",
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "env_code",
						Value: &types.ReferenceExpr{
							Parts: []string{"lookup(local", "environment_code, var", "environment, \"dev\")"},
						},
					},
					&types.Attribute{
						Name: "name_prefix",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${var.project_name}-${local.env_code}\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "common_tags",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Environment"},
									},
									Value: &types.ReferenceExpr{
										Parts: []string{"var", "environment"},
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Project"},
									},
									Value: &types.ReferenceExpr{
										Parts: []string{"var", "project_name"},
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"ManagedBy"},
									},
									Value: &types.LiteralValue{
										Value:     "Terraform",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Owner"},
									},
									Value: &types.ReferenceExpr{
										Parts: []string{"var", "owner"},
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "module",
				Labels:       []string{"vpc"},
				BlockComment: "// VPC Module",
				Children: []types.Body{
					&types.Attribute{
						Name: "source",
						Value: &types.LiteralValue{
							Value:     "modules/vpc",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "vpc_name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${local.name_prefix}-vpc\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "vpc_cidr",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "vpc_cidr"},
						},
					},
					&types.Attribute{
						Name: "azs",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "availability_zones"},
						},
					},
					&types.Attribute{
						Name: "private_subnets",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "private_subnet_cidrs"},
						},
					},
					&types.Attribute{
						Name: "public_subnets",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "public_subnet_cidrs"},
						},
					},
					&types.Attribute{
						Name: "enable_nat_gateway",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "enable_nat_gateway"},
						},
					},
					&types.Attribute{
						Name: "single_nat_gateway",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "environment != \"production\""},
						},
					},
					&types.Attribute{
						Name: "enable_vpn_gateway",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "enable_vpn_gateway"},
						},
					},
					&types.Attribute{
						Name: "enable_dns_hostnames",
						Value: &types.LiteralValue{
							Value:     true,
							ValueType: "bool",
						},
					},
					&types.Attribute{
						Name: "enable_dns_support",
						Value: &types.LiteralValue{
							Value:     true,
							ValueType: "bool",
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.FunctionCallExpr{
							Name: "merge",
							Args: []types.Expression{
								&types.ReferenceExpr{
									Parts: []string{"local"},
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"Name"},
											},
											Value: &types.TemplateExpr{
												Parts: []types.Expression{
													&types.LiteralValue{
														Value:     "\"${local.name_prefix}-vpc\"",
														ValueType: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "module",
				Labels:       []string{"security"},
				BlockComment: "// Security group for the application",
				Children: []types.Body{
					&types.Attribute{
						Name: "source",
						Value: &types.LiteralValue{
							Value:     "terraform-aws-modules/security-group/aws",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "version",
						Value: &types.LiteralValue{
							Value:     "~> 4.0",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${local.name_prefix}-sg\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "description",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"Security group for ${var.project_name} application\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "vpc_id",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "vpc_id"},
						},
					},
					&types.Attribute{
						Name: "ingress_with_cidr_blocks",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"from_port"},
											},
											Value: &types.LiteralValue{
												Value:     80,
												ValueType: "number",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"to_port"},
											},
											Value: &types.LiteralValue{
												Value:     80,
												ValueType: "number",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"protocol"},
											},
											Value: &types.LiteralValue{
												Value:     "tcp",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"description"},
											},
											Value: &types.LiteralValue{
												Value:     "HTTP",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"cidr_blocks"},
											},
											Value: &types.LiteralValue{
												Value:     "0.0.0.0/0",
												ValueType: "string",
											},
										},
									},
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"from_port"},
											},
											Value: &types.LiteralValue{
												Value:     443,
												ValueType: "number",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"to_port"},
											},
											Value: &types.LiteralValue{
												Value:     443,
												ValueType: "number",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"protocol"},
											},
											Value: &types.LiteralValue{
												Value:     "tcp",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"description"},
											},
											Value: &types.LiteralValue{
												Value:     "HTTPS",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"cidr_blocks"},
											},
											Value: &types.LiteralValue{
												Value:     "0.0.0.0/0",
												ValueType: "string",
											},
										},
									},
								},
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
							},
						},
					},
					&types.Attribute{
						Name: "egress_with_cidr_blocks",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"from_port"},
											},
											Value: &types.LiteralValue{
												Value:     0,
												ValueType: "number",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"to_port"},
											},
											Value: &types.LiteralValue{
												Value:     0,
												ValueType: "number",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"protocol"},
											},
											Value: &types.LiteralValue{
												Value:     "-1",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"description"},
											},
											Value: &types.LiteralValue{
												Value:     "Allow all outbound traffic",
												ValueType: "string",
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"cidr_blocks"},
											},
											Value: &types.LiteralValue{
												Value:     "0.0.0.0/0",
												ValueType: "string",
											},
										},
									},
								},
								&types.LiteralValue{
									Value:     "",
									ValueType: "null",
								},
							},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_instance", "app"},
				BlockComment: "// EC2 instances",
				Children: []types.Body{
					&types.Attribute{
						Name: "count",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "instance_count"},
						},
					},
					&types.Attribute{
						Name: "ami",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "ami_id"},
						},
					},
					&types.Attribute{
						Name: "instance_type",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "instance_type"},
						},
					},
					&types.Attribute{
						Name: "subnet_id",
						Value: &types.ReferenceExpr{
							Parts: []string{"element(module", "vpc", "private_subnets, count", "index % length(module", "vpc", "private_subnets))"},
						},
					},
					&types.Attribute{
						Name: "vpc_security_group_ids",
						Value: &types.ReferenceExpr{
							Parts: []string{"[module", "security", "security_group_id]"},
						},
					},
					&types.Block{
						Type: "root_block_device",
						Children: []types.Body{
							&types.Attribute{
								Name: "volume_type",
								Value: &types.LiteralValue{
									Value:     "gp3",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "volume_size",
								Value: &types.LiteralValue{
									Value:     50,
									ValueType: "number",
								},
							},
							&types.Attribute{
								Name: "encrypted",
								Value: &types.LiteralValue{
									Value:     true,
									ValueType: "bool",
								},
							},
						},
					},
					&types.Attribute{
						Name: "user_data",
						Value: &types.FunctionCallExpr{
							Name: "templatefile",
							Args: []types.Expression{
								&types.TemplateExpr{
									Parts: []types.Expression{
										&types.LiteralValue{
											Value:     "\"${path.module}/templates/user_data.sh.tpl\"",
											ValueType: "string",
										},
									},
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"environment"},
											},
											Value: &types.ReferenceExpr{
												Parts: []string{"var", "environment"},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"region"},
											},
											Value: &types.ReferenceExpr{
												Parts: []string{"var", "region"},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"app_port"},
											},
											Value: &types.ReferenceExpr{
												Parts: []string{"var", "app_port"},
											},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.FunctionCallExpr{
							Name: "merge",
							Args: []types.Expression{
								&types.ReferenceExpr{
									Parts: []string{"local"},
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"Name"},
											},
											Value: &types.TemplateExpr{
												Parts: []types.Expression{
													&types.LiteralValue{
														Value:     "\"${local.name_prefix}-app-${count.index + 1}\"",
														ValueType: "string",
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type: "lifecycle",
						Children: []types.Body{
							&types.Attribute{
								Name: "create_before_destroy",
								Value: &types.LiteralValue{
									Value:     true,
									ValueType: "bool",
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_lb", "app"},
				BlockComment: "// Load balancer",
				Children: []types.Body{
					&types.Attribute{
						Name: "name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${local.name_prefix}-alb\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "internal",
						Value: &types.LiteralValue{
							Value:     false,
							ValueType: "bool",
						},
					},
					&types.Attribute{
						Name: "load_balancer_type",
						Value: &types.LiteralValue{
							Value:     "application",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "security_groups",
						Value: &types.ReferenceExpr{
							Parts: []string{"[module", "security", "security_group_id]"},
						},
					},
					&types.Attribute{
						Name: "subnets",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "public_subnets"},
						},
					},
					&types.Attribute{
						Name: "enable_deletion_protection",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "environment == \"production\""},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_lb_target_group", "app"},
				BlockComment: "// Target group",
				Children: []types.Body{
					&types.Attribute{
						Name: "name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${local.name_prefix}-tg\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "port",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "app_port"},
						},
					},
					&types.Attribute{
						Name: "protocol",
						Value: &types.LiteralValue{
							Value:     "HTTP",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "vpc_id",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "vpc_id"},
						},
					},
					&types.Block{
						Type: "health_check",
						Children: []types.Body{
							&types.Attribute{
								Name: "path",
								Value: &types.LiteralValue{
									Value:     "/health",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "port",
								Value: &types.LiteralValue{
									Value:     "traffic-port",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "healthy_threshold",
								Value: &types.LiteralValue{
									Value:     3,
									ValueType: "number",
								},
							},
							&types.Attribute{
								Name: "unhealthy_threshold",
								Value: &types.LiteralValue{
									Value:     3,
									ValueType: "number",
								},
							},
							&types.Attribute{
								Name: "timeout",
								Value: &types.LiteralValue{
									Value:     5,
									ValueType: "number",
								},
							},
							&types.Attribute{
								Name: "interval",
								Value: &types.LiteralValue{
									Value:     30,
									ValueType: "number",
								},
							},
							&types.Attribute{
								Name: "matcher",
								Value: &types.LiteralValue{
									Value:     "200",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_lb_target_group_attachment", "app"},
				BlockComment: "// Target group attachment",
				Children: []types.Body{
					&types.Attribute{
						Name: "count",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "instance_count"},
						},
					},
					&types.Attribute{
						Name: "target_group_arn",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_lb_target_group", "app", "arn"},
						},
					},
					&types.Attribute{
						Name: "target_id",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_instance", "app[count", "index]", "id"},
						},
					},
					&types.Attribute{
						Name: "port",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "app_port"},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_lb_listener", "http"},
				BlockComment: "// Listener",
				Children: []types.Body{
					&types.Attribute{
						Name: "load_balancer_arn",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_lb", "app", "arn"},
						},
					},
					&types.Attribute{
						Name: "port",
						Value: &types.LiteralValue{
							Value:     80,
							ValueType: "number",
						},
					},
					&types.Attribute{
						Name: "protocol",
						Value: &types.LiteralValue{
							Value:     "HTTP",
							ValueType: "string",
						},
					},
					&types.Block{
						Type: "default_action",
						Children: []types.Body{
							&types.Attribute{
								Name: "type",
								Value: &types.LiteralValue{
									Value:     "forward",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "target_group_arn",
								Value: &types.ReferenceExpr{
									Parts: []string{"aws_lb_target_group", "app", "arn"},
								},
							},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_db_instance", "app"},
				BlockComment: "// Database",
				Children: []types.Body{
					&types.Attribute{
						Name: "count",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "create_database ? 1 : 0"},
						},
					},
					&types.Attribute{
						Name: "identifier",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${local.name_prefix}-db\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "engine",
						Value: &types.LiteralValue{
							Value:     "postgres",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "engine_version",
						Value: &types.LiteralValue{
							Value:     "13.4",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "instance_class",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "db_instance_class"},
						},
					},
					&types.Attribute{
						Name: "allocated_storage",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "db_allocated_storage"},
						},
					},
					&types.Attribute{
						Name: "max_allocated_storage",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "db_max_allocated_storage"},
						},
					},
					&types.Attribute{
						Name: "db_name",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "db_name"},
						},
					},
					&types.Attribute{
						Name: "username",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "db_username"},
						},
					},
					&types.Attribute{
						Name: "password",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "db_password"},
						},
					},
					&types.Attribute{
						Name: "vpc_security_group_ids",
						Value: &types.ReferenceExpr{
							Parts: []string{"[module", "security", "security_group_id]"},
						},
					},
					&types.Attribute{
						Name: "db_subnet_group_name",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_db_subnet_group", "app[0]", "name"},
						},
					},
					&types.Attribute{
						Name: "backup_retention_period",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "environment == \"production\" ? 30 : 7"},
						},
					},
					&types.Attribute{
						Name: "backup_window",
						Value: &types.LiteralValue{
							Value:     "03:00-04:00",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "maintenance_window",
						Value: &types.LiteralValue{
							Value:     "Mon:04:00-Mon:05:00",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "skip_final_snapshot",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "environment != \"production\""},
						},
					},
					&types.Attribute{
						Name: "final_snapshot_identifier",
						Value: &types.ConditionalExpr{
							Condition: &types.BinaryExpr{
								Left: &types.ReferenceExpr{
									Parts: []string{"var", "environment"},
								},
								Operator: "==",
								Right: &types.LiteralValue{
									Value:     "production",
									ValueType: "string",
								},
							},
							TrueExpr: &types.TemplateExpr{
								Parts: []types.Expression{
									&types.LiteralValue{
										Value:     "\"${local.name_prefix}-db-final-snapshot\"",
										ValueType: "string",
									},
								},
							},
							FalseExpr: &types.LiteralValue{ValueType: "null"},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_db_subnet_group", "app"},
				BlockComment: "// DB subnet group",
				Children: []types.Body{
					&types.Attribute{
						Name: "count",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "create_database ? 1 : 0"},
						},
					},
					&types.Attribute{
						Name: "name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.LiteralValue{
									Value:     "\"${local.name_prefix}-db-subnet-group\"",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "subnet_ids",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "private_subnets"},
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ReferenceExpr{
							Parts: []string{"local", "common_tags"},
						},
					},
				},
			},
			&types.Block{
				Type:         "output",
				Labels:       []string{"vpc_id"},
				BlockComment: "// Outputs",
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "The ID of the VPC",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "vpc_id"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"private_subnets"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "List of IDs of private subnets",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "private_subnets"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"public_subnets"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "List of IDs of public subnets",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "vpc", "public_subnets"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"security_group_id"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "The ID of the security group",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"module", "security", "security_group_id"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"instance_ids"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "List of IDs of instances",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_instance", "app[*]", "id"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"instance_private_ips"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "List of private IP addresses of instances",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_instance", "app[*]", "private_ip"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"lb_dns_name"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "The DNS name of the load balancer",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"aws_lb", "app", "dns_name"},
						},
					},
				},
			},
			&types.Block{
				Type:   "output",
				Labels: []string{"db_instance_endpoint"},
				Children: []types.Body{
					&types.Attribute{
						Name: "description",
						Value: &types.LiteralValue{
							Value:     "The connection endpoint of the database",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ReferenceExpr{
							Parts: []string{"var", "create_database ? aws_db_instance", "app[0]", "endpoint : null"},
						},
					},
				},
			},
		},
	}
}

// createEdgeCasesExpected creates the expected structure for edge_cases_test.tf
func createEdgeCasesExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:   "resource",
				Labels: []string{"null_resource", "empty"},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"null_resource", "comments_only"},
				BlockComment: "// Resource with only comments",
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Extremely long line",
				Children: []types.Body{
					&types.Attribute{
						Name: "extremely_long_line",
						Value: &types.LiteralValue{
							Value:     "This is an extremely long line that exceeds the typical line length limit. It's designed to test how the parser handles very long lines. The line continues for a while to ensure it's long enough to potentially cause issues with buffers or other limitations in the parser. It includes some special characters like quotes (\\\"), backslashes (\\\\), and other potentially problematic characters: !@#$%^&*()_+-=[]{}|;:,.<>?/",
							ValueType: "string",
						},
					},
				},
			},
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_instance", "unusual_whitespace"},
				BlockComment: "// Unusual whitespace",
				Children: []types.Body{
					&types.Attribute{
						Name: "ami",
						Value: &types.LiteralValue{
							Value:     "ami-12345678",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "instance_type",
						Value: &types.LiteralValue{
							Value:     "t2.micro",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "tags",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Name"},
									},
									Value: &types.LiteralValue{
										Value:     "unusual-whitespace",
										ValueType: "string",
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"Environment"},
									},
									Value: &types.LiteralValue{
										Value:     "test",
										ValueType: "string",
									},
								},
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Nested conditionals",
				Children: []types.Body{
					&types.Attribute{
						Name: "nested_conditionals",
						Value: &types.ConditionalExpr{
							Condition: &types.LiteralValue{
								Value:     true,
								ValueType: "bool",
							},
							TrueExpr: &types.ConditionalExpr{
								Condition: &types.LiteralValue{
									Value:     false,
									ValueType: "bool",
								},
								TrueExpr: &types.LiteralValue{
									Value:     "a",
									ValueType: "string",
								},
								FalseExpr: &types.ConditionalExpr{
									Condition: &types.LiteralValue{
										Value:     true,
										ValueType: "bool",
									},
									TrueExpr: &types.LiteralValue{
										Value:     "b",
										ValueType: "string",
									},
									FalseExpr: &types.ConditionalExpr{
										Condition: &types.LiteralValue{
											Value:     false,
											ValueType: "bool",
										},
										TrueExpr: &types.LiteralValue{
											Value:     "c",
											ValueType: "string",
										},
										FalseExpr: &types.ConditionalExpr{
											Condition: &types.LiteralValue{
												Value:     true,
												ValueType: "bool",
											},
											TrueExpr: &types.LiteralValue{
												Value:     "d",
												ValueType: "string",
											},
											FalseExpr: &types.LiteralValue{
												Value:     "e",
												ValueType: "string",
											},
										},
									},
								},
							},
							FalseExpr: &types.LiteralValue{
								Value:     "f",
								ValueType: "string",
							},
						},
					},
				},
			},
			&types.Block{
				Type:         "variable",
				Labels:       []string{"duplicate"},
				BlockComment: "// Multiple blocks with the same name",
				Children: []types.Body{
					&types.Attribute{
						Name: "type",
						Value: &types.LiteralValue{
							Value:     "string",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "default",
						Value: &types.LiteralValue{
							Value:     "first",
							ValueType: "string",
						},
					},
				},
			},
			&types.Block{
				Type:   "variable",
				Labels: []string{"duplicate"},
				Children: []types.Body{
					&types.Attribute{
						Name: "type",
						Value: &types.LiteralValue{
							Value:     "number",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "default",
						Value: &types.LiteralValue{
							Value:     123,
							ValueType: "number",
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Unicode characters",
				Children: []types.Body{
					&types.Attribute{
						Name: "unicode",
						Value: &types.LiteralValue{
							Value:     "こんにちは世界 • Hello, World! • Привет, мир! • مرحبا بالعالم • 你好，世界！",
							ValueType: "string",
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Escaped sequences",
				Children: []types.Body{
					&types.Attribute{
						Name: "escaped",
						Value: &types.LiteralValue{
							Value:     "Line 1\\nLine 2\\tTabbed\\r\\nWindows line ending\\\\\\\\ Double backslash \\\" Quote",
							ValueType: "string",
						},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Empty strings and special values",
				Children: []types.Body{
					&types.Attribute{
						Name: "empty_string",
						Value: &types.LiteralValue{
							Value:     "",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "single_space",
						Value: &types.LiteralValue{
							Value:     " ",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "just_newline",
						Value: &types.LiteralValue{
							Value:     "\\n",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name:  "null_value",
						Value: &types.LiteralValue{ValueType: "null"},
					},
				},
			},
			&types.Block{
				Type:         "locals",
				BlockComment: "// Unusual numbers",
				Children: []types.Body{
					&types.Attribute{
						Name: "zero",
						Value: &types.LiteralValue{
							Value:     0,
							ValueType: "number",
						},
					},
					&types.Attribute{
						Name: "negative",
						Value: &types.UnaryExpr{
							Operator: "-",
							Expr: &types.LiteralValue{
								Value:     42,
								ValueType: "number",
							},
						},
					},
					&types.Attribute{
						Name: "decimal",
						Value: &types.ReferenceExpr{
							Parts: []string{"3", "14159265359"},
						},
					},
					&types.Attribute{
						Name: "scientific",
						Value: &types.ReferenceExpr{
							Parts: []string{"1", "23e45"},
						},
					},
					&types.Attribute{
						Name: "hex",
						Value: &types.LiteralValue{
							Value:     "0xDEADBEEF",
							ValueType: "string",
						},
					},
					&types.Attribute{
						Name: "octal",
						Value: &types.LiteralValue{
							Value:     0,
							ValueType: "number",
						},
					},
					&types.Attribute{
						Name: "o755",
						Value: &types.LiteralValue{
							Value:     0,
							ValueType: "number",
						},
					},
					&types.Block{
						Type:   "b10101010",
						Labels: []string{"aws_instance", "special-chars_in.identifier"},
						Children: []types.Body{
							&types.Attribute{
								Name: "ami",
								Value: &types.LiteralValue{
									Value:     "ami-12345678",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "instance_type",
								Value: &types.LiteralValue{
									Value:     "t2.micro",
									ValueType: "string",
								},
							},
						},
					},
					&types.Block{
						Type:         "provider",
						Labels:       []string{"aws"},
						BlockComment: "// Provider with unusual configuration",
						Children: []types.Body{
							&types.Attribute{
								Name: "region",
								Value: &types.LiteralValue{
									Value:     "us-west-2",
									ValueType: "string",
								},
							},
							&types.Attribute{
								Name: "alias",
								Value: &types.LiteralValue{
									Value:     "unusual",
									ValueType: "string",
								},
							},
							&types.Block{
								Type: "assume_role",
								Children: []types.Body{
									&types.Attribute{
										Name: "role_arn",
										Value: &types.LiteralValue{
											Value:     "arn:aws:iam::123456789012:role/unusual-role",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "session_name",
										Value: &types.LiteralValue{
											Value:     "unusual-session",
											ValueType: "string",
										},
									},
								},
							},
							&types.Block{
								Type: "default_tags",
								Children: []types.Body{
									&types.Attribute{
										Name: "tags",
										Value: &types.ObjectExpr{
											Items: []types.ObjectItem{
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"ManagedBy"},
													},
													Value: &types.LiteralValue{
														Value:     "Terraform",
														ValueType: "string",
													},
												},
												{
													Key: &types.ReferenceExpr{
														Parts: []string{"Environment"},
													},
													Value: &types.LiteralValue{
														Value:     "Test",
														ValueType: "string",
													},
												},
											},
										},
									},
								},
							},
							&types.Block{
								Type: "ignore_tags",
								Children: []types.Body{
									&types.Attribute{
										Name: "keys",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
												&types.LiteralValue{
													Value:     "IgnoreMe",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "AlsoIgnoreMe",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
											},
										},
									},
									&types.Attribute{
										Name: "key_prefixes",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
												&types.LiteralValue{
													Value:     "temp-",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "tmp-",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:         "resource",
						Labels:       []string{"aws_security_group", "nested_blocks"},
						BlockComment: "// Unusual block nesting",
						Children: []types.Body{
							&types.Attribute{
								Name: "name",
								Value: &types.LiteralValue{
									Value:     "nested-blocks",
									ValueType: "string",
								},
							},
							&types.Block{
								Type: "ingress",
								Children: []types.Body{
									&types.Attribute{
										Name: "description",
										Value: &types.LiteralValue{
											Value:     "First level",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "from_port",
										Value: &types.LiteralValue{
											Value:     80,
											ValueType: "number",
										},
									},
									&types.Attribute{
										Name: "to_port",
										Value: &types.LiteralValue{
											Value:     80,
											ValueType: "number",
										},
									},
									&types.Attribute{
										Name: "protocol",
										Value: &types.LiteralValue{
											Value:     "tcp",
											ValueType: "string",
										},
									},
									&types.Attribute{
										Name: "cidr_blocks",
										Value: &types.ReferenceExpr{
											Parts: []string{"[\"0", "0", "0", "0/0\"]"},
										},
									},
								},
							},
							&types.Block{
								Type:   "dynamic",
								Labels: []string{"egress"},
								Children: []types.Body{
									&types.Attribute{
										Name: "for_each",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
												&types.LiteralValue{
													Value:     "one",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "two",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "three",
													ValueType: "string",
												},
												&types.LiteralValue{
													Value:     "",
													ValueType: "null",
												},
											},
										},
									},
									&types.Block{
										Type: "content",
										Children: []types.Body{
											&types.Attribute{
												Name: "description",
												Value: &types.TemplateExpr{
													Parts: []types.Expression{
														&types.LiteralValue{
															Value:     "\"Dynamic ${egress.value}\"",
															ValueType: "string",
														},
													},
												},
											},
											&types.Attribute{
												Name: "from_port",
												Value: &types.LiteralValue{
													Value:     0,
													ValueType: "number",
												},
											},
											&types.Attribute{
												Name: "to_port",
												Value: &types.LiteralValue{
													Value:     0,
													ValueType: "number",
												},
											},
											&types.Attribute{
												Name: "protocol",
												Value: &types.LiteralValue{
													Value:     "-1",
													ValueType: "string",
												},
											},
											&types.Attribute{
												Name: "cidr_blocks",
												Value: &types.ReferenceExpr{
													Parts: []string{"[\"0", "0", "0", "0/0\"]"},
												},
											},
										},
									},
								},
							},
							&types.Block{
								Type: "lifecycle",
								Children: []types.Body{
									&types.Attribute{
										Name: "create_before_destroy",
										Value: &types.LiteralValue{
											Value:     true,
											ValueType: "bool",
										},
									},
									&types.Block{
										Type: "precondition",
										Children: []types.Body{
											&types.Attribute{
												Name: "condition",
												Value: &types.ReferenceExpr{
													Parts: []string{"length(var", "allowed_cidrs) > 0"},
												},
											},
											&types.Attribute{
												Name: "error_message",
												Value: &types.LiteralValue{
													Value:     "At least one CIDR block must be allowed.",
													ValueType: "string",
												},
											},
										},
									},
									&types.Block{
										Type: "postcondition",
										Children: []types.Body{
											&types.Attribute{
												Name: "condition",
												Value: &types.ReferenceExpr{
													Parts: []string{"self", "name != \"\""},
												},
											},
											&types.Attribute{
												Name: "error_message",
												Value: &types.LiteralValue{
													Value:     "Name cannot be empty.",
													ValueType: "string",
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:         "locals",
						BlockComment: "// Unusual function calls",
						Children: []types.Body{
							&types.Attribute{
								Name: "unusual_functions",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"nested"},
											},
											Value: &types.FunctionCallExpr{
												Name: "merge",
												Args: []types.Expression{
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"a"},
																},
																Value: &types.LiteralValue{
																	Value:     "value",
																	ValueType: "string",
																},
															},
														},
													},
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key: &types.ReferenceExpr{
																	Parts: []string{"b"},
																},
																Value: &types.FunctionCallExpr{
																	Name: "lookup",
																	Args: []types.Expression{
																		&types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"x"},
																					},
																					Value: &types.LiteralValue{
																						Value:     "x-value",
																						ValueType: "string",
																					},
																				},
																				{
																					Key: &types.ReferenceExpr{
																						Parts: []string{"y"},
																					},
																					Value: &types.LiteralValue{
																						Value:     "y-value",
																						ValueType: "string",
																					},
																				},
																			},
																		},
																		&types.LiteralValue{
																			Value:     "z",
																			ValueType: "string",
																		},
																		&types.LiteralValue{
																			Value:     "default",
																			ValueType: "string",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"chained"},
											},
											Value: &types.FunctionCallExpr{
												Name: "join",
												Args: []types.Expression{
													&types.LiteralValue{
														Value:     ",",
														ValueType: "string",
													},
													&types.FunctionCallExpr{
														Name: "concat",
														Args: []types.Expression{
															&types.FunctionCallExpr{
																Name: "split",
																Args: []types.Expression{
																	&types.LiteralValue{
																		Value:     ",",
																		ValueType: "string",
																	},
																	&types.LiteralValue{
																		Value:     "a,b,c",
																		ValueType: "string",
																	},
																},
															},
															&types.ArrayExpr{
																Items: []types.Expression{
																	&types.LiteralValue{
																		Value:     "",
																		ValueType: "null",
																	},
																	&types.LiteralValue{
																		Value:     "d",
																		ValueType: "string",
																	},
																	&types.LiteralValue{
																		Value:     "e",
																		ValueType: "string",
																	},
																	&types.LiteralValue{
																		Value:     "f",
																		ValueType: "string",
																	},
																	&types.LiteralValue{
																		Value:     "",
																		ValueType: "null",
																	},
																},
															},
														},
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{
												Parts: []string{"complex"},
											},
											Value: &types.FunctionCallExpr{
												Name: "formatlist",
												Args: []types.Expression{
													&types.LiteralValue{
														Value:     "%s = %s",
														ValueType: "string",
													},
													&types.FunctionCallExpr{
														Name: "keys",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"key1"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "value1",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"key2"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "value2",
																			ValueType: "string",
																		},
																	},
																},
															},
														},
													},
													&types.FunctionCallExpr{
														Name: "values",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"key1"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "value1",
																			ValueType: "string",
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{
																			Parts: []string{"key2"},
																		},
																		Value: &types.LiteralValue{
																			Value:     "value2",
																			ValueType: "string",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Block{
						Type:         "locals",
						BlockComment: "// Comments in unusual places",
						Children: []types.Body{
							&types.Attribute{
								Name: "value1",
								Value: &types.LiteralValue{
									Value:     "test",
									ValueType: "string",
								},
								InlineComment: "# end of line comment",
							},
							&types.Attribute{
								Name: "value2",
								Value: &types.LiteralValue{
									Value:     "test2",
									ValueType: "string",
								},
								BlockComment: "# end of line comment",
							},
						},
					},
					&types.Block{
						Type:         "locals",
						BlockComment: "// Unusual heredoc",
						Children: []types.Body{
							&types.Attribute{
								Name: "unusual_heredoc",
								Value: &types.LiteralValue{
									Value:     " and '\nIt includes backslashes: \\\\ and \\n and \\t\nIt includes unicode: こんにちは世界\nUNUSUAL\n\n  indented_heredoc = <<-INDENTED\n    This is an indented heredoc.\n    The leading whitespace will be trimmed.\n      This line has extra indentation.\n    Back to normal indentation.\n  INDENTED\n}\n\n// Unusual for expressions\nlocals {\n  unusual_for = [\n    for i, v in [",
									ValueType: "string",
								},
							},
						},
					},
				},
			},
		},
	}
}
