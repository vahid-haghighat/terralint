package parser

import (
	"github.com/vahid-haghighat/terralint/parser/types"
)

// createComplexModuleExpected creates the expected structure for 01_complex_module.tf
func createComplexModuleExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "module",
				Labels:       []string{"complex_module"},
				BlockComment: "// Complex module with nested expressions, conditionals, and for loops",
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
						Name:         "vpc_config",
						BlockComment: "// Complex map with nested objects, expressions, and functions",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{Parts: []string{"name"}},
									Value: &types.TemplateExpr{
										Parts: []types.Expression{
											&types.LiteralValue{Value: "complex-", ValueType: "string"},
											&types.ReferenceExpr{Parts: []string{"var", "environment"}},
											&types.LiteralValue{Value: "-", ValueType: "string"},
											&types.ReferenceExpr{Parts: []string{"local", "region_code"}},
											&types.LiteralValue{Value: "-vpc", ValueType: "string"},
										},
									},
								},
								{
									Key: &types.ReferenceExpr{Parts: []string{"cidr"}},
									Value: &types.FunctionCallExpr{
										Name: "cidrsubnet",
										Args: []types.Expression{
											&types.ReferenceExpr{Parts: []string{"var", "base_cidr_block"}},
											&types.LiteralValue{Value: 4, ValueType: "number"},
											&types.ReferenceExpr{Parts: []string{"var", "vpc_index"}},
										},
									},
								},
								{
									Key:          &types.ReferenceExpr{Parts: []string{"enable_dns"}},
									BlockComment: "// Nested conditional expression",
									Value: &types.ConditionalExpr{
										Condition: &types.BinaryExpr{
											Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
											Operator: "==",
											Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
										},
										TrueExpr: &types.LiteralValue{Value: true, ValueType: "bool"},
										FalseExpr: &types.ParenExpr{
											Expression: &types.ConditionalExpr{
												Condition: &types.BinaryExpr{
													Left:     &types.ReferenceExpr{Parts: []string{"var", "enable_dns"}},
													Operator: "!=",
													Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
												},
												TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "enable_dns"}},
												FalseExpr: &types.LiteralValue{Value: false, ValueType: "bool"},
											},
										},
									},
								},
								{
									Key: &types.ReferenceExpr{Parts: []string{"tags"}},
									Value: &types.FunctionCallExpr{
										Name: "merge",
										Args: []types.Expression{
											&types.ReferenceExpr{Parts: []string{"var", "common_tags"}},
											&types.ObjectExpr{
												Items: []types.ObjectItem{
													{
														Key: &types.ReferenceExpr{Parts: []string{"Name"}},
														Value: &types.TemplateExpr{
															Parts: []types.Expression{
																&types.LiteralValue{Value: "complex-", ValueType: "string"},
																&types.ReferenceExpr{Parts: []string{"var", "environment"}},
																&types.LiteralValue{Value: "-vpc", ValueType: "string"},
															},
														},
													},
													{
														Key:   &types.ReferenceExpr{Parts: []string{"Environment"}},
														Value: &types.ReferenceExpr{Parts: []string{"var", "environment"}},
													},
													{
														Key:   &types.ReferenceExpr{Parts: []string{"ManagedBy"}},
														Value: &types.LiteralValue{Value: "terraform", ValueType: "string"},
													},
													{
														Key: &types.ReferenceExpr{Parts: []string{"Complex"}},
														Value: &types.FunctionCallExpr{
															Name: "jsonencode",
															Args: []types.Expression{
																&types.ObjectExpr{
																	Items: []types.ObjectItem{
																		{
																			Key:   &types.ReferenceExpr{Parts: []string{"nested"}},
																			Value: &types.LiteralValue{Value: "value", ValueType: "string"},
																		},
																		{
																			Key: &types.ReferenceExpr{Parts: []string{"list"}},
																			Value: &types.ArrayExpr{
																				Items: []types.Expression{
																					&types.LiteralValue{Value: 1, ValueType: "number"},
																					&types.LiteralValue{Value: 2, ValueType: "number"},
																					&types.LiteralValue{Value: 3, ValueType: "number"},
																					&types.LiteralValue{Value: 4, ValueType: "number"},
																				},
																			},
																		},
																		{
																			Key: &types.ReferenceExpr{Parts: []string{"map"}},
																			Value: &types.ObjectExpr{
																				Items: []types.ObjectItem{
																					{
																						Key:   &types.ReferenceExpr{Parts: []string{"key1"}},
																						Value: &types.LiteralValue{Value: "value1", ValueType: "string"},
																					},
																					{
																						Key:   &types.ReferenceExpr{Parts: []string{"key2"}},
																						Value: &types.LiteralValue{Value: 42, ValueType: "number"},
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
											&types.ConditionalExpr{
												Condition: &types.BinaryExpr{
													Left:     &types.ReferenceExpr{Parts: []string{"var", "additional_tags"}},
													Operator: "!=",
													Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
												},
												TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "additional_tags"}},
												FalseExpr: &types.ObjectExpr{},
											},
										},
									},
									BlockComment: "// Complex object with nested expressions",
								},
							},
						},
					},
					&types.Attribute{
						Name:         "subnet_cidrs",
						BlockComment: "// Complex for expression with filtering and transformation",
						Value: &types.ForArrayExpr{
							KeyVar:     "i",
							ValueVar:   "subnet",
							Collection: &types.ReferenceExpr{Parts: []string{"var", "subnets"}},
							ThenValueExpr: &types.FunctionCallExpr{
								Name: "cidrsubnet",
								Args: []types.Expression{
									&types.ReferenceExpr{Parts: []string{"var", "base_cidr_block"}},
									&types.LiteralValue{Value: 8, ValueType: "number"},
									&types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"i"}},
										Operator: "+",
										Right:    &types.LiteralValue{Value: 10, ValueType: "number"},
									},
								},
							},
							Condition: &types.BinaryExpr{
								Left:     &types.ReferenceExpr{Parts: []string{"subnet", "create"}},
								Operator: "==",
								Right:    &types.LiteralValue{Value: true, ValueType: "bool"},
							},
						},
					},
					&types.Attribute{
						Name:         "subnet_configs",
						BlockComment: "// Nested for expressions with conditional",
						Value: &types.ForMapExpr{
							KeyVar:      "zone_key",
							ValueVar:    "zone",
							Collection:  &types.ReferenceExpr{Parts: []string{"var", "availability_zones"}},
							ThenKeyExpr: &types.ReferenceExpr{Parts: []string{"zone_key"}},
							ThenValueExpr: &types.ForMapExpr{
								KeyVar:      "subnet_key",
								ValueVar:    "subnet",
								Collection:  &types.ReferenceExpr{Parts: []string{"var", "subnet_types"}},
								ThenKeyExpr: &types.ReferenceExpr{Parts: []string{"subnet_key"}},
								ThenValueExpr: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{Parts: []string{"cidr"}},
											Value: &types.FunctionCallExpr{
												Name: "cidrsubnet",
												Args: []types.Expression{
													&types.ReferenceExpr{Parts: []string{"var", "base_cidr_block"}},
													&types.ReferenceExpr{Parts: []string{"var", "subnet_newbits"}},
													&types.BinaryExpr{
														Left: &types.BinaryExpr{
															Left: &types.FunctionCallExpr{
																Name: "index",
																Args: []types.Expression{
																	&types.ReferenceExpr{Parts: []string{"var", "availability_zones"}},
																	&types.ReferenceExpr{Parts: []string{"zone"}},
																},
															},
															Operator: "*",
															Right: &types.FunctionCallExpr{
																Name: "length",
																Args: []types.Expression{
																	&types.ReferenceExpr{Parts: []string{"var", "subnet_types"}},
																},
															},
														},
														Operator: "+",
														Right: &types.FunctionCallExpr{
															Name: "index",
															Args: []types.Expression{
																&types.ReferenceExpr{Parts: []string{"var", "subnet_types"}},
																&types.ReferenceExpr{Parts: []string{"subnet"}},
															},
														},
													},
												},
											},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"az"}},
											Value: &types.ReferenceExpr{Parts: []string{"zone"}},
										},
										{
											Key: &types.ReferenceExpr{Parts: []string{"tags"}},
											Value: &types.FunctionCallExpr{
												Name: "merge",
												Args: []types.Expression{
													&types.ReferenceExpr{Parts: []string{"var", "common_tags"}},
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key: &types.ReferenceExpr{Parts: []string{"Name"}},
																Value: &types.TemplateExpr{
																	Parts: []types.Expression{
																		&types.ReferenceExpr{Parts: []string{"var", "environment"}},
																		&types.LiteralValue{Value: "-", ValueType: "string"},
																		&types.ReferenceExpr{Parts: []string{"subnet"}},
																		&types.LiteralValue{Value: "-", ValueType: "string"},
																		&types.ReferenceExpr{Parts: []string{"zone"}},
																	},
																},
															},
															{
																Key:   &types.ReferenceExpr{Parts: []string{"Type"}},
																Value: &types.ReferenceExpr{Parts: []string{"subnet"}},
															},
														},
													},
												},
											},
										},
									},
								},
								Condition: &types.ReferenceExpr{Parts: []string{"subnet", "enabled"}},
							},
							Condition: &types.FunctionCallExpr{
								Name: "contains",
								Args: []types.Expression{
									&types.ReferenceExpr{Parts: []string{"var", "enabled_zones"}},
									&types.ReferenceExpr{Parts: []string{"zone"}},
								},
							},
						},
					},
					&types.Attribute{
						Name:         "all_subnet_ids",
						BlockComment: "// Complex splat expression",
						Value: &types.FunctionCallExpr{
							Name: "flatten",
							Args: []types.Expression{
								&types.ForArrayExpr{
									KeyVar:     "zone_key",
									ValueVar:   "zone",
									Collection: &types.ReferenceExpr{Parts: []string{"aws_subnet", "main"}},
									ThenValueExpr: &types.ForArrayExpr{
										KeyVar:        "subnet_key",
										ValueVar:      "subnet",
										Collection:    &types.ReferenceExpr{Parts: []string{"zone"}},
										ThenValueExpr: &types.ReferenceExpr{Parts: []string{"subnet", "id"}},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name:         "user_data",
						BlockComment: "// Heredoc with interpolation",
						Value: &types.HeredocExpr{
							Marker:   "EOT",
							Indented: true,
							Content:  "#!/bin/bash\necho \"Environment: ${var.environment}\"\necho \"Region: ${data.aws_region.current.name}\"\n\n# Complex interpolation\n${join(\"\\n\", [\n  for script in var.bootstrap_scripts :\n  \"source ${script}\"\n])}\n\n# Conditional section\n${var.install_monitoring ? \"setup_monitoring ${var.monitoring_endpoint}\" : \"echo 'Monitoring disabled'\"}\n\n# For loop in heredoc\n%{for pkg in var.packages~}\nyum install -y ${pkg}\n%{endfor~}\n\n# If directive in heredoc\n%{if var.environment == \"prod\"~}\necho \"Production environment detected, applying strict security\"\n%{else~}\necho \"Non-production environment, using standard security\"\n%{endif~}",
						},
					},
					&types.Attribute{
						Name:         "timeout",
						BlockComment: "// Complex binary expressions with nested conditionals",
						Value: &types.BinaryExpr{
							Left: &types.ParenExpr{
								Expression: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
										Operator: "==",
										Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
									},
									TrueExpr: &types.LiteralValue{Value: 300, ValueType: "number"},
									FalseExpr: &types.ConditionalExpr{
										Condition: &types.BinaryExpr{
											Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
											Operator: "==",
											Right:    &types.LiteralValue{Value: "staging", ValueType: "string"},
										},
										TrueExpr: &types.LiteralValue{Value: 180, ValueType: "number"},
										FalseExpr: &types.ConditionalExpr{
											Condition: &types.BinaryExpr{
												Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
												Operator: "==",
												Right:    &types.LiteralValue{Value: "dev", ValueType: "string"},
											},
											TrueExpr:  &types.LiteralValue{Value: 60, ValueType: "number"},
											FalseExpr: &types.LiteralValue{Value: 30, ValueType: "number"},
										},
									},
								},
							},
							Operator: "+",
							Right: &types.ParenExpr{
								Expression: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "additional_timeout"}},
										Operator: "!=",
										Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
									},
									TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "additional_timeout"}},
									FalseExpr: &types.LiteralValue{Value: 0, ValueType: "number"},
								},
							},
						},
					},
					&types.Attribute{
						Name:         "complex_calculation",
						BlockComment: "// Nested parentheses and operators",
						Value: &types.BinaryExpr{
							Left: &types.ParenExpr{
								Expression: &types.BinaryExpr{
									Left: &types.ParenExpr{
										Expression: &types.BinaryExpr{
											Left:     &types.ReferenceExpr{Parts: []string{"var", "base_value"}},
											Operator: "*",
											Right: &types.ParenExpr{
												Expression: &types.BinaryExpr{
													Left:     &types.LiteralValue{Value: 1, ValueType: "number"},
													Operator: "+",
													Right:    &types.ReferenceExpr{Parts: []string{"var", "multiplier"}},
												},
											},
										},
									},
									Operator: "/",
									Right: &types.ParenExpr{
										Expression: &types.ConditionalExpr{
											Condition: &types.BinaryExpr{
												Left:     &types.ReferenceExpr{Parts: []string{"var", "divisor"}},
												Operator: ">",
												Right:    &types.LiteralValue{Value: 0, ValueType: "number"},
											},
											TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "divisor"}},
											FalseExpr: &types.LiteralValue{Value: 1, ValueType: "number"},
										},
									},
								},
							},
							Operator: "+",
							Right: &types.ParenExpr{
								Expression: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
										Operator: "==",
										Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
									},
									TrueExpr: &types.ParenExpr{
										Expression: &types.BinaryExpr{
											Left:     &types.ReferenceExpr{Parts: []string{"var", "prod_adjustment"}},
											Operator: "*",
											Right: &types.ParenExpr{
												Expression: &types.BinaryExpr{
													Left:     &types.LiteralValue{Value: 1, ValueType: "number"},
													Operator: "+",
													Right:    &types.ReferenceExpr{Parts: []string{"var", "prod_factor"}},
												},
											},
										},
									},
									FalseExpr: &types.ParenExpr{
										Expression: &types.BinaryExpr{
											Left:     &types.ReferenceExpr{Parts: []string{"var", "non_prod_adjustment"}},
											Operator: "*",
											Right: &types.ParenExpr{
												Expression: &types.BinaryExpr{
													Left:     &types.LiteralValue{Value: 1, ValueType: "number"},
													Operator: "-",
													Right:    &types.ReferenceExpr{Parts: []string{"var", "non_prod_factor"}},
												},
											},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name:         "security_groups",
						BlockComment: "// Complex function calls with nested expressions",
						Value: &types.FunctionCallExpr{
							Name: "compact",
							Args: []types.Expression{
								&types.FunctionCallExpr{
									Name: "concat",
									Args: []types.Expression{
										&types.ArrayExpr{
											Items: []types.Expression{
												&types.ConditionalExpr{
													Condition: &types.ReferenceExpr{Parts: []string{"var", "create_default_security_group"}},
													TrueExpr: &types.RelativeTraversalExpr{
														Source: &types.IndexExpr{
															Collection: &types.ReferenceExpr{Parts: []string{"aws_security_group", "default"}},
															Key:        &types.LiteralValue{Value: 0, ValueType: "number"},
														},
														Traversal: []types.TraversalElem{{Type: "attr", Name: "id"}},
													},
													FalseExpr: &types.LiteralValue{Value: "", ValueType: "string"},
												},
												&types.ConditionalExpr{
													Condition: &types.ReferenceExpr{Parts: []string{"var", "create_bastion_security_group"}},
													TrueExpr: &types.RelativeTraversalExpr{
														Source: &types.IndexExpr{
															Collection: &types.ReferenceExpr{Parts: []string{"aws_security_group", "bastion"}},
															Key:        &types.LiteralValue{Value: 0, ValueType: "number"},
														},
														Traversal: []types.TraversalElem{{Type: "attr", Name: "id"}},
													},
													FalseExpr: &types.LiteralValue{Value: "", ValueType: "string"},
												},
											},
										},
										&types.ReferenceExpr{Parts: []string{"var", "additional_security_group_ids"}},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name:         "custom_template",
						BlockComment: "// Template with directives",
						Value: &types.FunctionCallExpr{
							Name: "templatefile",
							Args: []types.Expression{
								&types.TemplateExpr{
									Parts: []types.Expression{
										&types.ReferenceExpr{Parts: []string{"path", "module"}},
										&types.LiteralValue{Value: "/templates/config.tpl", ValueType: "string"},
									},
								},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key:   &types.ReferenceExpr{Parts: []string{"environment"}},
											Value: &types.ReferenceExpr{Parts: []string{"var", "environment"}},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"region"}},
											Value: &types.ReferenceExpr{Parts: []string{"data", "aws_region", "current", "name"}},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"vpc_id"}},
											Value: &types.ReferenceExpr{Parts: []string{"aws_vpc", "main", "id"}},
										},
										{
											Key: &types.ReferenceExpr{Parts: []string{"subnets"}},
											Value: &types.ForArrayExpr{
												KeyVar:     "subnet",
												Collection: &types.ReferenceExpr{Parts: []string{"aws_subnet", "main"}},
												ThenValueExpr: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key:   &types.ReferenceExpr{Parts: []string{"id"}},
															Value: &types.ReferenceExpr{Parts: []string{"subnet", "id"}},
														},
														{
															Key:   &types.ReferenceExpr{Parts: []string{"cidr"}},
															Value: &types.ReferenceExpr{Parts: []string{"subnet", "cidr_block"}},
														},
														{
															Key:   &types.ReferenceExpr{Parts: []string{"az"}},
															Value: &types.ReferenceExpr{Parts: []string{"subnet", "availability_zone"}},
														},
													},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{Parts: []string{"features"}},
											Value: &types.ObjectExpr{
												Items: []types.ObjectItem{
													{
														Key:   &types.ReferenceExpr{Parts: []string{"monitoring"}},
														Value: &types.ReferenceExpr{Parts: []string{"var", "enable_monitoring"}},
													},
													{
														Key:   &types.ReferenceExpr{Parts: []string{"logging"}},
														Value: &types.ReferenceExpr{Parts: []string{"var", "enable_logging"}},
													},
													{
														Key: &types.ReferenceExpr{Parts: []string{"encryption"}},
														Value: &types.ConditionalExpr{
															Condition: &types.BinaryExpr{
																Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
																Operator: "==",
																Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
															},
															TrueExpr:  &types.LiteralValue{Value: true, ValueType: "bool"},
															FalseExpr: &types.ReferenceExpr{Parts: []string{"var", "enable_encryption"}},
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
						Labels: []string{"ingress"},
						Children: []types.Body{
							&types.Attribute{
								Name:  "for_each",
								Value: &types.ReferenceExpr{Parts: []string{"var", "ingress_rules"}},
							},
							&types.Attribute{
								Name:  "iterator",
								Value: &types.ReferenceExpr{Parts: []string{"rule"}},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "description",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"rule", "value"}},
												&types.LiteralValue{Value: "description", ValueType: "string"},
												&types.TemplateExpr{
													Parts: []types.Expression{
														&types.LiteralValue{Value: "Ingress Rule ", ValueType: "string"},
														&types.ReferenceExpr{Parts: []string{"rule", "key"}},
													},
												},
											},
										},
									},
									&types.Attribute{
										Name:  "from_port",
										Value: &types.ReferenceExpr{Parts: []string{"rule", "value", "from_port"}},
									},
									&types.Attribute{
										Name:  "to_port",
										Value: &types.ReferenceExpr{Parts: []string{"rule", "value", "to_port"}},
									},
									&types.Attribute{
										Name: "protocol",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"rule", "value"}},
												&types.LiteralValue{Value: "protocol", ValueType: "string"},
												&types.LiteralValue{Value: "tcp", ValueType: "string"},
											},
										},
									},
									&types.Attribute{
										Name: "cidr_blocks",
										Value: &types.ConditionalExpr{
											Condition: &types.BinaryExpr{
												Left: &types.FunctionCallExpr{
													Name: "lookup",
													Args: []types.Expression{
														&types.ReferenceExpr{Parts: []string{"rule", "value"}},
														&types.LiteralValue{Value: "cidr_blocks", ValueType: "string"},
														&types.LiteralValue{Value: nil, ValueType: "null"},
													},
												},
												Operator: "!=",
												Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
											},
											TrueExpr: &types.ReferenceExpr{Parts: []string{"rule", "value", "cidr_blocks"}},
											FalseExpr: &types.ForArrayExpr{
												KeyVar:     "cidr",
												Collection: &types.ReferenceExpr{Parts: []string{"var", "default_cidrs"}},
												ThenValueExpr: &types.ReferenceExpr{
													Parts: []string{"cidr"},
												},
												Condition: &types.UnaryExpr{
													Operator: "!",
													Expr: &types.FunctionCallExpr{
														Name: "contains",
														Args: []types.Expression{
															&types.ReferenceExpr{Parts: []string{"var", "excluded_cidrs"}},
															&types.ReferenceExpr{Parts: []string{"cidr"}},
														},
													},
												},
											},
										},
									},
									&types.Attribute{
										Name: "security_groups",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"rule", "value"}},
												&types.LiteralValue{Value: "security_group_ids", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "self",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"rule", "value"}},
												&types.LiteralValue{Value: "self", ValueType: "string"},
												&types.LiteralValue{Value: false, ValueType: "bool"},
											},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name:         "validation",
						BlockComment: "// Complex type constraints",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{Parts: []string{"condition"}},
									Value: &types.FunctionCallExpr{
										Name: "can",
										Args: []types.Expression{
											&types.FunctionCallExpr{
												Name: "regex",
												Args: []types.Expression{
													&types.LiteralValue{Value: "^(dev|staging|prod)$", ValueType: "string"},
													&types.ReferenceExpr{Parts: []string{"var", "environment"}},
												},
											},
										},
									},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"error_message"}},
									Value: &types.LiteralValue{Value: "Environment must be one of: dev, staging, prod.", ValueType: "string"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// createComplexResourceExpected creates the expected structure for 02_complex_resource.tf
func createComplexResourceExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_security_group", "complex"},
				BlockComment: "// Resource with complex dynamic blocks and for_each",
				Children: []types.Body{
					&types.Attribute{
						Name: "for_each",
						Value: &types.ForMapExpr{
							KeyVar:        "sg",
							Collection:    &types.ReferenceExpr{Parts: []string{"var", "security_groups"}},
							ThenKeyExpr:   &types.ReferenceExpr{Parts: []string{"sg", "name"}},
							ThenValueExpr: &types.ReferenceExpr{Parts: []string{"sg"}},
							Condition:     &types.ReferenceExpr{Parts: []string{"sg", "create"}},
						},
					},
					&types.Attribute{
						Name: "name",
						Value: &types.TemplateExpr{
							Parts: []types.Expression{
								&types.ReferenceExpr{Parts: []string{"var", "prefix"}},
								&types.LiteralValue{Value: "-", ValueType: "string"},
								&types.ReferenceExpr{Parts: []string{"each", "key"}},
							},
						},
					},
					&types.Attribute{
						Name:  "description",
						Value: &types.ReferenceExpr{Parts: []string{"each", "value", "description"}},
					},
					&types.Attribute{
						Name:  "vpc_id",
						Value: &types.ReferenceExpr{Parts: []string{"var", "vpc_id"}},
					},
					&types.Block{
						Type:   "dynamic",
						Labels: []string{"ingress"},
						Children: []types.Body{
							&types.Attribute{
								Name:  "for_each",
								Value: &types.ReferenceExpr{Parts: []string{"each", "value", "ingress_rules"}},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name:  "description",
										Value: &types.ReferenceExpr{Parts: []string{"ingress", "value", "description"}},
									},
									&types.Attribute{
										Name:  "from_port",
										Value: &types.ReferenceExpr{Parts: []string{"ingress", "value", "from_port"}},
									},
									&types.Attribute{
										Name:  "to_port",
										Value: &types.ReferenceExpr{Parts: []string{"ingress", "value", "to_port"}},
									},
									&types.Attribute{
										Name:  "protocol",
										Value: &types.ReferenceExpr{Parts: []string{"ingress", "value", "protocol"}},
									},
									&types.Attribute{
										Name:  "cidr_blocks",
										Value: &types.ReferenceExpr{Parts: []string{"ingress", "value", "cidr_blocks"}},
									},
									&types.Attribute{
										Name: "ipv6_cidr_blocks",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"ingress", "value"}},
												&types.LiteralValue{Value: "ipv6_cidr_blocks", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "prefix_list_ids",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"ingress", "value"}},
												&types.LiteralValue{Value: "prefix_list_ids", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "security_groups",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"ingress", "value"}},
												&types.LiteralValue{Value: "security_groups", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "self",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"ingress", "value"}},
												&types.LiteralValue{Value: "self", ValueType: "string"},
												&types.LiteralValue{Value: false, ValueType: "bool"},
											},
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
								Name:  "for_each",
								Value: &types.ReferenceExpr{Parts: []string{"each", "value", "egress_rules"}},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name:  "description",
										Value: &types.ReferenceExpr{Parts: []string{"egress", "value", "description"}},
									},
									&types.Attribute{
										Name:  "from_port",
										Value: &types.ReferenceExpr{Parts: []string{"egress", "value", "from_port"}},
									},
									&types.Attribute{
										Name:  "to_port",
										Value: &types.ReferenceExpr{Parts: []string{"egress", "value", "to_port"}},
									},
									&types.Attribute{
										Name:  "protocol",
										Value: &types.ReferenceExpr{Parts: []string{"egress", "value", "protocol"}},
									},
									&types.Attribute{
										Name:  "cidr_blocks",
										Value: &types.ReferenceExpr{Parts: []string{"egress", "value", "cidr_blocks"}},
									},
									&types.Attribute{
										Name: "ipv6_cidr_blocks",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"egress", "value"}},
												&types.LiteralValue{Value: "ipv6_cidr_blocks", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "prefix_list_ids",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"egress", "value"}},
												&types.LiteralValue{Value: "prefix_list_ids", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "security_groups",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"egress", "value"}},
												&types.LiteralValue{Value: "security_groups", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "self",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"egress", "value"}},
												&types.LiteralValue{Value: "self", ValueType: "string"},
												&types.LiteralValue{Value: false, ValueType: "bool"},
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
								&types.ReferenceExpr{Parts: []string{"var", "common_tags"}},
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key: &types.ReferenceExpr{Parts: []string{"Name"}},
											Value: &types.TemplateExpr{
												Parts: []types.Expression{
													&types.ReferenceExpr{Parts: []string{"var", "prefix"}},
													&types.LiteralValue{Value: "-", ValueType: "string"},
													&types.ReferenceExpr{Parts: []string{"each", "key"}},
												},
											},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"Type"}},
											Value: &types.LiteralValue{Value: "SecurityGroup", ValueType: "string"},
										},
										{
											Key: &types.ReferenceExpr{Parts: []string{"Rules"}},
											Value: &types.FunctionCallExpr{
												Name: "jsonencode",
												Args: []types.Expression{
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key: &types.ReferenceExpr{Parts: []string{"ingress"}},
																Value: &types.FunctionCallExpr{
																	Name: "length",
																	Args: []types.Expression{&types.ReferenceExpr{Parts: []string{"each", "value", "ingress_rules"}}},
																},
															},
															{
																Key: &types.ReferenceExpr{Parts: []string{"egress"}},
																Value: &types.FunctionCallExpr{
																	Name: "length",
																	Args: []types.Expression{&types.ReferenceExpr{Parts: []string{"each", "value", "egress_rules"}}},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								&types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"each", "value", "additional_tags"}},
										Operator: "!=",
										Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
									},
									TrueExpr:  &types.ReferenceExpr{Parts: []string{"each", "value", "additional_tags"}},
									FalseExpr: &types.ObjectExpr{},
								},
							},
						},
						BlockComment: "// Complex tags with expressions and functions",
					},
					&types.Block{
						Type: "lifecycle",
						Children: []types.Body{
							&types.Attribute{
								Name:  "create_before_destroy",
								Value: &types.LiteralValue{Value: true, ValueType: "bool"},
							},
							&types.Attribute{
								Name: "prevent_destroy",
								Value: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
										Operator: "==",
										Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
									},
									TrueExpr:  &types.LiteralValue{Value: true, ValueType: "bool"},
									FalseExpr: &types.LiteralValue{Value: false, ValueType: "bool"},
								},
							},
							&types.Attribute{
								Name: "ignore_changes",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.RelativeTraversalExpr{
											Source: &types.ReferenceExpr{Parts: []string{"tags"}},
											Traversal: []types.TraversalElem{
												{Type: "index", Index: &types.LiteralValue{Value: "LastModified", ValueType: "string"}},
											},
										},
										&types.RelativeTraversalExpr{
											Source: &types.ReferenceExpr{Parts: []string{"tags"}},
											Traversal: []types.TraversalElem{
												{Type: "index", Index: &types.LiteralValue{Value: "AutoUpdated", ValueType: "string"}},
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
	}
}

// createComplexLocalsExpected creates the expected structure for 03_complex_locals.tf
func createComplexLocalsExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "locals",
				BlockComment: "// Complex locals with nested expressions",
				Children: []types.Body{
					&types.Attribute{
						Name: "subnet_map",
						Value: &types.ForMapExpr{
							KeyVar:      "subnet",
							Collection:  &types.ReferenceExpr{Parts: []string{"var", "subnets"}},
							ThenKeyExpr: &types.ReferenceExpr{Parts: []string{"subnet", "name"}},
							ThenValueExpr: &types.ObjectExpr{
								Items: []types.ObjectItem{
									{
										Key:   &types.ReferenceExpr{Parts: []string{"id"}},
										Value: &types.ReferenceExpr{Parts: []string{"subnet", "id"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"cidr"}},
										Value: &types.ReferenceExpr{Parts: []string{"subnet", "cidr_block"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"az"}},
										Value: &types.ReferenceExpr{Parts: []string{"subnet", "availability_zone"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"public"}},
										Value: &types.ReferenceExpr{Parts: []string{"subnet", "public"}},
									},
									{
										Key: &types.ReferenceExpr{Parts: []string{"nat_gw"}},
										Value: &types.ConditionalExpr{
											Condition: &types.ReferenceExpr{Parts: []string{"subnet", "public"}},
											TrueExpr:  &types.LiteralValue{Value: true, ValueType: "bool"},
											FalseExpr: &types.LiteralValue{Value: false, ValueType: "bool"},
										},
									},
									{
										Key: &types.ReferenceExpr{Parts: []string{"depends_on"}},
										Value: &types.ConditionalExpr{
											Condition: &types.ReferenceExpr{Parts: []string{"subnet", "public"}},
											TrueExpr:  &types.ArrayExpr{Items: []types.Expression{}},
											FalseExpr: &types.ForArrayExpr{
												ValueVar:      "s",
												Collection:    &types.ReferenceExpr{Parts: []string{"var", "subnets"}},
												ThenValueExpr: &types.ReferenceExpr{Parts: []string{"s", "id"}},
												Condition:     &types.ReferenceExpr{Parts: []string{"s", "public"}},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Complex map transformation",
					},
					&types.Attribute{
						Name: "filtered_instances",
						Value: &types.ForArrayExpr{
							KeyVar:     "server",
							Collection: &types.ReferenceExpr{Parts: []string{"var", "servers"}},
							ThenValueExpr: &types.ObjectExpr{
								Items: []types.ObjectItem{
									{
										Key:   &types.ReferenceExpr{Parts: []string{"id"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "id"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"name"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "name"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"environment"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "environment"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"type"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "type"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"subnet_id"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "subnet_id"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"private_ip"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "private_ip"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"public_ip"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "public_ip"}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"tags"}},
										Value: &types.ReferenceExpr{Parts: []string{"server", "tags"}},
									},
								},
							},
							Condition: &types.BinaryExpr{
								Left: &types.BinaryExpr{
									Left:     &types.ReferenceExpr{Parts: []string{"server", "environment"}},
									Operator: "==",
									Right:    &types.ReferenceExpr{Parts: []string{"var", "environment"}},
								},
								Operator: "&&",
								Right: &types.BinaryExpr{
									Left: &types.FunctionCallExpr{
										Name: "contains",
										Args: []types.Expression{
											&types.ReferenceExpr{Parts: []string{"var", "allowed_types"}},
											&types.ReferenceExpr{Parts: []string{"server", "type"}},
										},
									},
									Operator: "&&",
									Right: &types.UnaryExpr{
										Operator: "!",
										Expr: &types.FunctionCallExpr{
											Name: "contains",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"var", "excluded_ids"}},
												&types.ReferenceExpr{Parts: []string{"server", "id"}},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Nested for expressions with filtering",
					},
					&types.Attribute{
						Name: "backup_config",
						Value: &types.ConditionalExpr{
							Condition: &types.ReferenceExpr{Parts: []string{"var", "enable_backups"}},
							TrueExpr: &types.ObjectExpr{
								Items: []types.ObjectItem{
									{
										Key: &types.ReferenceExpr{Parts: []string{"schedule"}},
										Value: &types.ConditionalExpr{
											Condition: &types.BinaryExpr{
												Left:     &types.ReferenceExpr{Parts: []string{"var", "backup_schedule"}},
												Operator: "!=",
												Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
											},
											TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "backup_schedule"}},
											FalseExpr: &types.LiteralValue{Value: "0 1 * * *", ValueType: "string"},
										},
									},
									{
										Key: &types.ReferenceExpr{Parts: []string{"retention"}},
										Value: &types.ParenExpr{
											Expression: &types.ConditionalExpr{
												Condition: &types.BinaryExpr{
													Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
													Operator: "==",
													Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
												},
												TrueExpr: &types.LiteralValue{Value: 30, ValueType: "number"},
												FalseExpr: &types.ConditionalExpr{
													Condition: &types.BinaryExpr{
														Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
														Operator: "==",
														Right:    &types.LiteralValue{Value: "staging", ValueType: "string"},
													},
													TrueExpr:  &types.LiteralValue{Value: 14, ValueType: "number"},
													FalseExpr: &types.LiteralValue{Value: 7, ValueType: "number"},
												},
											},
										},
									},
									{
										Key: &types.ReferenceExpr{Parts: []string{"targets"}},
										Value: &types.ForArrayExpr{
											KeyVar:     "target",
											Collection: &types.ReferenceExpr{Parts: []string{"var", "backup_targets"}},
											ThenValueExpr: &types.ObjectExpr{
												Items: []types.ObjectItem{
													{
														Key:   &types.ReferenceExpr{Parts: []string{"id"}},
														Value: &types.ReferenceExpr{Parts: []string{"target", "id"}},
													},
													{
														Key:   &types.ReferenceExpr{Parts: []string{"name"}},
														Value: &types.ReferenceExpr{Parts: []string{"target", "name"}},
													},
													{
														Key: &types.ReferenceExpr{Parts: []string{"priority"}},
														Value: &types.FunctionCallExpr{
															Name: "lookup",
															Args: []types.Expression{
																&types.ReferenceExpr{Parts: []string{"target", "tags"}},
																&types.LiteralValue{Value: "backup-priority", ValueType: "string"},
																&types.LiteralValue{Value: "medium", ValueType: "string"},
															},
														},
													},
												},
											},
											Condition: &types.BinaryExpr{
												Left: &types.FunctionCallExpr{
													Name: "lookup",
													Args: []types.Expression{
														&types.ReferenceExpr{Parts: []string{"target", "tags"}},
														&types.LiteralValue{Value: "backup-enabled", ValueType: "string"},
														&types.LiteralValue{Value: "false", ValueType: "string"},
													},
												},
												Operator: "==",
												Right:    &types.LiteralValue{Value: "true", ValueType: "string"},
											},
										},
									},
									{
										Key: &types.ReferenceExpr{Parts: []string{"storage"}},
										Value: &types.ObjectExpr{
											Items: []types.ObjectItem{
												{
													Key:   &types.ReferenceExpr{Parts: []string{"type"}},
													Value: &types.ReferenceExpr{Parts: []string{"var", "backup_storage_type"}},
												},
												{
													Key: &types.ReferenceExpr{Parts: []string{"path"}},
													Value: &types.ConditionalExpr{
														Condition: &types.BinaryExpr{
															Left:     &types.ReferenceExpr{Parts: []string{"var", "backup_storage_path"}},
															Operator: "!=",
															Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
														},
														TrueExpr: &types.ReferenceExpr{Parts: []string{"var", "backup_storage_path"}},
														FalseExpr: &types.TemplateExpr{
															Parts: []types.Expression{
																&types.LiteralValue{Value: "/backups/", ValueType: "string"},
																&types.ReferenceExpr{Parts: []string{"var", "environment"}},
															},
														},
													},
												},
												{
													Key: &types.ReferenceExpr{Parts: []string{"settings"}},
													Value: &types.FunctionCallExpr{
														Name: "merge",
														Args: []types.Expression{
															&types.ReferenceExpr{Parts: []string{"var", "default_storage_settings"}},
															&types.ConditionalExpr{
																Condition: &types.BinaryExpr{
																	Left:     &types.ReferenceExpr{Parts: []string{"var", "custom_storage_settings"}},
																	Operator: "!=",
																	Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
																},
																TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "custom_storage_settings"}},
																FalseExpr: &types.ObjectExpr{},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							FalseExpr: &types.LiteralValue{Value: nil, ValueType: "null"},
						},
						BlockComment: "// Complex conditional with multiple nested expressions",
					},
					&types.Attribute{
						Name: "naming_convention",
						Value: &types.FunctionCallExpr{
							Name: "join",
							Args: []types.Expression{
								&types.LiteralValue{Value: "-", ValueType: "string"},
								&types.FunctionCallExpr{
									Name: "compact",
									Args: []types.Expression{
										&types.ArrayExpr{
											Items: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"var", "prefix"}},
												&types.ReferenceExpr{Parts: []string{"var", "environment"}},
												&types.ReferenceExpr{Parts: []string{"var", "region_code"}},
												&types.ReferenceExpr{Parts: []string{"var", "name"}},
												&types.ConditionalExpr{
													Condition: &types.BinaryExpr{
														Left:     &types.ReferenceExpr{Parts: []string{"var", "suffix"}},
														Operator: "!=",
														Right:    &types.LiteralValue{Value: "", ValueType: "string"},
													},
													TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "suffix"}},
													FalseExpr: &types.LiteralValue{Value: nil, ValueType: "null"},
												},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Complex string interpolation with functions",
					},
					&types.Attribute{
						Name: "timeout_seconds",
						Value: &types.ParenExpr{
							Expression: &types.ConditionalExpr{
								Condition: &types.BinaryExpr{
									Left:     &types.ReferenceExpr{Parts: []string{"var", "custom_timeout"}},
									Operator: "!=",
									Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
								},
								TrueExpr: &types.ReferenceExpr{Parts: []string{"var", "custom_timeout"}},
								FalseExpr: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
										Operator: "==",
										Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
									},
									TrueExpr: &types.ParenExpr{
										Expression: &types.ConditionalExpr{
											Condition: &types.ReferenceExpr{Parts: []string{"var", "high_availability"}},
											TrueExpr:  &types.LiteralValue{Value: 300, ValueType: "number"},
											FalseExpr: &types.LiteralValue{Value: 180, ValueType: "number"},
										},
									},
									FalseExpr: &types.ParenExpr{
										Expression: &types.ConditionalExpr{
											Condition: &types.BinaryExpr{
												Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
												Operator: "==",
												Right:    &types.LiteralValue{Value: "staging", ValueType: "string"},
											},
											TrueExpr:  &types.LiteralValue{Value: 120, ValueType: "number"},
											FalseExpr: &types.LiteralValue{Value: 60, ValueType: "number"},
										},
									},
								},
							},
						},
						BlockComment: "// Nested ternary operators",
					},
					&types.Attribute{
						Name: "schema",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key:   &types.ReferenceExpr{Parts: []string{"type"}},
									Value: &types.LiteralValue{Value: "object", ValueType: "string"},
								},
								{
									Key: &types.ReferenceExpr{Parts: []string{"properties"}},
									Value: &types.ObjectExpr{
										Items: []types.ObjectItem{
											{
												Key: &types.ReferenceExpr{Parts: []string{"id"}},
												Value: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key:   &types.ReferenceExpr{Parts: []string{"type"}},
															Value: &types.LiteralValue{Value: "string", ValueType: "string"},
														},
														{
															Key:   &types.ReferenceExpr{Parts: []string{"pattern"}},
															Value: &types.LiteralValue{Value: "^[a-zA-Z0-9-_]+$", ValueType: "string"},
														},
													},
												},
											},
											{
												Key: &types.ReferenceExpr{Parts: []string{"settings"}},
												Value: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key:   &types.ReferenceExpr{Parts: []string{"type"}},
															Value: &types.LiteralValue{Value: "object", ValueType: "string"},
														},
														{
															Key: &types.ReferenceExpr{Parts: []string{"properties"}},
															Value: &types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"enabled"}},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																					Value: &types.LiteralValue{Value: "boolean", ValueType: "string"},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"timeout"}},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																					Value: &types.LiteralValue{Value: "number", ValueType: "string"},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"retries"}},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																					Value: &types.LiteralValue{Value: "number", ValueType: "string"},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"options"}},
																		Value: &types.ObjectExpr{
																			Items: []types.ObjectItem{
																				{
																					Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																					Value: &types.LiteralValue{Value: "array", ValueType: "string"},
																				},
																				{
																					Key: &types.ReferenceExpr{Parts: []string{"items"}},
																					Value: &types.ObjectExpr{
																						Items: []types.ObjectItem{
																							{
																								Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																								Value: &types.LiteralValue{Value: "object", ValueType: "string"},
																							},
																							{
																								Key: &types.ReferenceExpr{Parts: []string{"properties"}},
																								Value: &types.ObjectExpr{
																									Items: []types.ObjectItem{
																										{
																											Key: &types.ReferenceExpr{Parts: []string{"name"}},
																											Value: &types.ObjectExpr{
																												Items: []types.ObjectItem{
																													{
																														Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																														Value: &types.LiteralValue{Value: "string", ValueType: "string"},
																													},
																												},
																											},
																										},
																										{
																											Key: &types.ReferenceExpr{Parts: []string{"value"}},
																											Value: &types.ObjectExpr{
																												Items: []types.ObjectItem{
																													{
																														Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																														Value: &types.LiteralValue{Value: "string", ValueType: "string"},
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
												Key: &types.ReferenceExpr{Parts: []string{"tags"}},
												Value: &types.ObjectExpr{
													Items: []types.ObjectItem{
														{
															Key:   &types.ReferenceExpr{Parts: []string{"type"}},
															Value: &types.LiteralValue{Value: "object", ValueType: "string"},
														},
														{
															Key: &types.ReferenceExpr{Parts: []string{"additionalProperties"}},
															Value: &types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																		Value: &types.LiteralValue{Value: "string", ValueType: "string"},
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
		},
	}
}

// createComplexDataSourceExpected creates the expected structure for 04_complex_data_source.tf
func createComplexDataSourceExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
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
												&types.ReferenceExpr{Parts: []string{"statement", "value"}},
												&types.LiteralValue{Value: "sid", ValueType: "string"},
												&types.TemplateExpr{
													Parts: []types.Expression{&types.LiteralValue{Value: "Statement", ValueType: "string"}, &types.ReferenceExpr{Parts: []string{"statement", "key"}}},
												},
											},
										},
									},
									&types.Attribute{
										Name: "effect",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"statement", "value"}},
												&types.LiteralValue{Value: "effect", ValueType: "string"},
												&types.LiteralValue{Value: "Allow", ValueType: "string"},
											},
										},
									},
									&types.Attribute{
										Name:  "actions",
										Value: &types.ReferenceExpr{Parts: []string{"statement", "value", "actions"}},
									},
									&types.Attribute{
										Name:  "resources",
										Value: &types.ReferenceExpr{Parts: []string{"statement", "value", "resources"}},
									},
									&types.Block{
										Type:   "dynamic",
										Labels: []string{"condition"},
										Children: []types.Body{
											&types.Attribute{
												Name: "for_each",
												Value: &types.FunctionCallExpr{
													Name: "lookup",
													Args: []types.Expression{
														&types.ReferenceExpr{Parts: []string{"statement", "value"}},
														&types.LiteralValue{Value: "conditions", ValueType: "string"},
														&types.ArrayExpr{},
													},
												},
											},
											&types.Block{
												Type: "content",
												Children: []types.Body{
													&types.Attribute{
														Name:  "test",
														Value: &types.ReferenceExpr{Parts: []string{"condition", "value", "test"}},
													},
													&types.Attribute{
														Name:  "variable",
														Value: &types.ReferenceExpr{Parts: []string{"condition", "value", "variable"}},
													},
													&types.Attribute{
														Name:  "values",
														Value: &types.ReferenceExpr{Parts: []string{"condition", "value", "values"}},
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
												Value: &types.FunctionCallExpr{
													Name: "lookup",
													Args: []types.Expression{
														&types.ReferenceExpr{Parts: []string{"statement", "value"}},
														&types.LiteralValue{Value: "principals", ValueType: "string"},
														&types.ArrayExpr{},
													},
												},
											},
											&types.Block{
												Type: "content",
												Children: []types.Body{
													&types.Attribute{
														Name:  "type",
														Value: &types.ReferenceExpr{Parts: []string{"principals", "value", "type"}},
													},
													&types.Attribute{
														Name:  "identifiers",
														Value: &types.ReferenceExpr{Parts: []string{"principals", "value", "identifiers"}},
													},
												},
											},
										},
									},
									&types.Attribute{
										Name: "not_actions",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"statement", "value"}},
												&types.LiteralValue{Value: "not_actions", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
									&types.Attribute{
										Name: "not_resources",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"statement", "value"}},
												&types.LiteralValue{Value: "not_resources", ValueType: "string"},
												&types.ArrayExpr{},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Dynamic statement blocks",
					},
					&types.Block{
						Type: "statement",
						Children: []types.Body{
							&types.Attribute{
								Name:  "sid",
								Value: &types.LiteralValue{Value: "ExplicitAllow", ValueType: "string"},
							},
							&types.Attribute{
								Name:  "effect",
								Value: &types.LiteralValue{Value: "Allow", ValueType: "string"},
							},
							&types.Attribute{
								Name: "actions",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.LiteralValue{Value: "s3:GetObject", ValueType: "string"},
										&types.LiteralValue{Value: "s3:ListBucket", ValueType: "string"},
									},
								},
							},
							&types.Attribute{
								Name: "resources",
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.TemplateExpr{
											Parts: []types.Expression{
												&types.LiteralValue{Value: "arn:aws:s3:::", ValueType: "string"},
												&types.ReferenceExpr{Parts: []string{"var", "bucket_name"}},
											},
										},
										&types.TemplateExpr{
											Parts: []types.Expression{
												&types.LiteralValue{Value: "arn:aws:s3:::", ValueType: "string"},
												&types.ReferenceExpr{Parts: []string{"var", "bucket_name"}},
												&types.LiteralValue{Value: "/*", ValueType: "string"},
											},
										},
									},
								},
							},
							&types.Block{
								Type: "condition",
								Children: []types.Body{
									&types.Attribute{
										Name:  "test",
										Value: &types.LiteralValue{Value: "StringEquals", ValueType: "string"},
									},
									&types.Attribute{
										Name:  "variable",
										Value: &types.LiteralValue{Value: "aws:SourceVpc", ValueType: "string"},
									},
									&types.Attribute{
										Name: "values",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"var", "vpc_id"}},
											},
										},
									},
								},
							},
							&types.Block{
								Type: "condition",
								Children: []types.Body{
									&types.Attribute{
										Name:  "test",
										Value: &types.LiteralValue{Value: "StringLike", ValueType: "string"},
									},
									&types.Attribute{
										Name:  "variable",
										Value: &types.LiteralValue{Value: "aws:PrincipalTag/Role", ValueType: "string"},
									},
									&types.Attribute{
										Name: "values",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.LiteralValue{Value: "Admin", ValueType: "string"},
												&types.LiteralValue{Value: "Developer", ValueType: "string"},
											},
										},
									},
								},
							},
							&types.Block{
								Type: "principals",
								Children: []types.Body{
									&types.Attribute{
										Name:  "type",
										Value: &types.LiteralValue{Value: "AWS", ValueType: "string"},
									},
									&types.Attribute{
										Name: "identifiers",
										Value: &types.ArrayExpr{
											Items: []types.Expression{
												&types.TemplateExpr{
													Parts: []types.Expression{
														&types.LiteralValue{Value: "arn:aws:iam::", ValueType: "string"},
														&types.ReferenceExpr{Parts: []string{"var", "account_id"}},
														&types.LiteralValue{Value: ":root", ValueType: "string"},
													},
												},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Override with inline statement",
					},
				},
			},
		},
	}
}

// createComplexVariableExpected creates the expected structure for 05_complex_variable.tf
func createComplexVariableExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "variable",
				Labels:       []string{"complex_object"},
				BlockComment: "// Variable with complex type constraints and validations",
				Children: []types.Body{
					&types.Attribute{
						Name:  "description",
						Value: &types.LiteralValue{Value: "A complex object with nested types and validations", ValueType: "string"},
					},

					&types.Attribute{
						Name: "type",
						Value: &types.FunctionCallExpr{
							Name: "object",
							Args: []types.Expression{
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key:   &types.ReferenceExpr{Parts: []string{"name"}},
											Value: &types.ReferenceExpr{Parts: []string{"string"}},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"environment"}},
											Value: &types.ReferenceExpr{Parts: []string{"string"}},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"enabled"}},
											Value: &types.ReferenceExpr{Parts: []string{"bool"}},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"count"}},
											Value: &types.ReferenceExpr{Parts: []string{"number"}},
										},
										{
											Key: &types.ReferenceExpr{Parts: []string{"tags"}},
											Value: &types.FunctionCallExpr{
												Name: "map",
												Args: []types.Expression{
													&types.ReferenceExpr{Parts: []string{"string"}},
												},
											},
										},
										{
											Key: &types.ReferenceExpr{Parts: []string{"vpc"}},
											Value: &types.FunctionCallExpr{
												Name: "object",
												Args: []types.Expression{
													&types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key:   &types.ReferenceExpr{Parts: []string{"id"}},
																Value: &types.ReferenceExpr{Parts: []string{"string"}},
															},
															{
																Key:   &types.ReferenceExpr{Parts: []string{"cidr"}},
																Value: &types.ReferenceExpr{Parts: []string{"string"}},
															},
															{
																Key: &types.ReferenceExpr{Parts: []string{"private_subnets"}},
																Value: &types.FunctionCallExpr{
																	Name: "list",
																	Args: []types.Expression{
																		&types.FunctionCallExpr{
																			Name: "object",
																			Args: []types.Expression{
																				&types.ObjectExpr{
																					Items: []types.ObjectItem{
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"id"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
																						},
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"cidr"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
																						},
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"az"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
															{
																Key: &types.ReferenceExpr{Parts: []string{"public_subnets"}},
																Value: &types.FunctionCallExpr{
																	Name: "list",
																	Args: []types.Expression{
																		&types.FunctionCallExpr{
																			Name: "object",
																			Args: []types.Expression{
																				&types.ObjectExpr{
																					Items: []types.ObjectItem{
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"id"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
																						},
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"cidr"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
																						},
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"az"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
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
											Key: &types.ReferenceExpr{Parts: []string{"instances"}},
											Value: &types.FunctionCallExpr{
												Name: "list",
												Args: []types.Expression{
													&types.FunctionCallExpr{
														Name: "object",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"id"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"subnet_id"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"private_ip"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"public_ip"}},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"string"}},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"root_volume"}},
																		Value: &types.FunctionCallExpr{
																			Name: "object",
																			Args: []types.Expression{
																				&types.ObjectExpr{
																					Items: []types.ObjectItem{
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"size"}},
																							Value: &types.ReferenceExpr{Parts: []string{"number"}},
																						},
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																							Value: &types.ReferenceExpr{Parts: []string{"string"}},
																						},
																						{
																							Key:   &types.ReferenceExpr{Parts: []string{"encrypted"}},
																							Value: &types.ReferenceExpr{Parts: []string{"bool"}},
																						},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"data_volumes"}},
																		Value: &types.FunctionCallExpr{
																			Name: "list",
																			Args: []types.Expression{
																				&types.FunctionCallExpr{
																					Name: "object",
																					Args: []types.Expression{
																						&types.ObjectExpr{
																							Items: []types.ObjectItem{
																								{
																									Key:   &types.ReferenceExpr{Parts: []string{"device_name"}},
																									Value: &types.ReferenceExpr{Parts: []string{"string"}},
																								},
																								{
																									Key:   &types.ReferenceExpr{Parts: []string{"size"}},
																									Value: &types.ReferenceExpr{Parts: []string{"number"}},
																								},
																								{
																									Key:   &types.ReferenceExpr{Parts: []string{"type"}},
																									Value: &types.ReferenceExpr{Parts: []string{"string"}},
																								},
																								{
																									Key:   &types.ReferenceExpr{Parts: []string{"encrypted"}},
																									Value: &types.ReferenceExpr{Parts: []string{"bool"}},
																								},
																								{
																									Key: &types.ReferenceExpr{Parts: []string{"iops"}},
																									Value: &types.FunctionCallExpr{
																										Name: "optional",
																										Args: []types.Expression{
																											&types.ReferenceExpr{Parts: []string{"number"}},
																										},
																									},
																								},
																								{
																									Key: &types.ReferenceExpr{Parts: []string{"throughput"}},
																									Value: &types.FunctionCallExpr{
																										Name: "optional",
																										Args: []types.Expression{
																											&types.ReferenceExpr{Parts: []string{"number"}},
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
																		Key: &types.ReferenceExpr{Parts: []string{"tags"}},
																		Value: &types.FunctionCallExpr{
																			Name: "map",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"string"}},
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
											Key: &types.ReferenceExpr{Parts: []string{"databases"}},
											Value: &types.FunctionCallExpr{
												Name: "map",
												Args: []types.Expression{
													&types.FunctionCallExpr{
														Name: "object",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"engine"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"version"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"instance_class"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"allocated_storage"}},
																		Value: &types.ReferenceExpr{Parts: []string{"number"}},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"max_allocated_storage"}},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"number"}},
																			},
																		},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"multi_az"}},
																		Value: &types.ReferenceExpr{Parts: []string{"bool"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"backup_retention_period"}},
																		Value: &types.ReferenceExpr{Parts: []string{"number"}},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"Args"}},
																		Value: &types.FunctionCallExpr{
																			Name: "map",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"string"}},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"subnet_ids"}},
																		Value: &types.FunctionCallExpr{
																			Name: "list",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"string"}},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"security_group_ids"}},
																		Value: &types.FunctionCallExpr{
																			Name: "list",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"string"}},
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
											Key: &types.ReferenceExpr{Parts: []string{"endpoints"}},
											Value: &types.FunctionCallExpr{
												Name: "map",
												Args: []types.Expression{
													&types.FunctionCallExpr{
														Name: "object",
														Args: []types.Expression{
															&types.ObjectExpr{
																Items: []types.ObjectItem{
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"service"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key:   &types.ReferenceExpr{Parts: []string{"vpc_endpoint_type"}},
																		Value: &types.ReferenceExpr{Parts: []string{"string"}},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"subnet_ids"}},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.FunctionCallExpr{
																					Name: "list",
																					Args: []types.Expression{
																						&types.ReferenceExpr{Parts: []string{"string"}},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"security_group_ids"}},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.FunctionCallExpr{
																					Name: "list",
																					Args: []types.Expression{
																						&types.ReferenceExpr{Parts: []string{"string"}},
																					},
																				},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"private_dns_enabled"}},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"bool"}},
																			},
																		},
																	},
																	{
																		Key: &types.ReferenceExpr{Parts: []string{"policy"}},
																		Value: &types.FunctionCallExpr{
																			Name: "optional",
																			Args: []types.Expression{
																				&types.ReferenceExpr{Parts: []string{"string"}},
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
								Value: &types.BinaryExpr{
									Left: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "complex_object", "count"}},
										Operator: ">",
										Right:    &types.LiteralValue{Value: 0, ValueType: "number"},
									},
									Operator: "&&",
									Right: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "complex_object", "count"}},
										Operator: "<=",
										Right:    &types.LiteralValue{Value: 10, ValueType: "number"},
									},
								},
							},
							&types.Attribute{
								Name:  "error_message",
								Value: &types.LiteralValue{Value: "Count must be between 1 and 10.", ValueType: "string"},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.FunctionCallExpr{
									Name: "can",
									Args: []types.Expression{
										&types.FunctionCallExpr{
											Name: "regex",
											Args: []types.Expression{
												&types.LiteralValue{Value: "^(dev|staging|prod)$", ValueType: "string"},
												&types.ReferenceExpr{Parts: []string{"var", "complex_object", "environment"}},
											},
										},
									},
								},
							},
							&types.Attribute{
								Name:  "error_message",
								Value: &types.LiteralValue{Value: "Environment must be one of: dev, staging, prod.", ValueType: "string"},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.BinaryExpr{
									Left: &types.FunctionCallExpr{
										Name: "length",
										Args: []types.Expression{
											&types.ReferenceExpr{Parts: []string{"var", "complex_object", "vpc", "private_subnets"}},
										},
									},
									Operator: ">",
									Right:    &types.LiteralValue{Value: 0, ValueType: "number"},
								},
							},
							&types.Attribute{
								Name:  "error_message",
								Value: &types.LiteralValue{Value: "At least one private subnet must be defined.", ValueType: "string"},
							},
						},
					},
					&types.Block{
						Type: "validation",
						Children: []types.Body{
							&types.Attribute{
								Name: "condition",
								Value: &types.FunctionCallExpr{
									Name: "alltrue",
									Args: []types.Expression{
										&types.ForArrayExpr{
											KeyVar:     "instance",
											Collection: &types.ReferenceExpr{Parts: []string{"var", "complex_object", "instances"}},
											ThenValueExpr: &types.BinaryExpr{
												Left: &types.BinaryExpr{
													Left:     &types.ReferenceExpr{Parts: []string{"instance", "root_volume", "size"}},
													Operator: ">=",
													Right:    &types.LiteralValue{Value: 20, ValueType: "number"},
												},
												Operator: "&&",
												Right: &types.BinaryExpr{
													Left:     &types.ReferenceExpr{Parts: []string{"instance", "root_volume", "encrypted"}},
													Operator: "==",
													Right:    &types.LiteralValue{Value: true, ValueType: "bool"},
												},
											},
										},
									},
								},
							},
							&types.Attribute{
								Name:  "error_message",
								Value: &types.LiteralValue{Value: "All root volumes must be at least 20GB and encrypted.", ValueType: "string"},
							},
						},
					},
				},
			},
		},
	}
}

// createComplexOutputExpected creates the expected structure for 06_complex_output.tf
func createComplexOutputExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "output",
				Labels:       []string{"complex_output"},
				BlockComment: "// Output with complex expressions",
				Children: []types.Body{
					&types.Attribute{
						Name:  "description",
						Value: &types.LiteralValue{Value: "Complex output with nested expressions", ValueType: "string"},
					},
					&types.Attribute{
						Name: "value",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key:   &types.ReferenceExpr{Parts: []string{"vpc_id"}},
									Value: &types.ReferenceExpr{Parts: []string{"module", "complex_module", "vpc_id"}},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"subnet_ids"}},
									Value: &types.ReferenceExpr{Parts: []string{"module", "complex_module", "subnet_ids"}},
								},
								{
									Key: &types.ReferenceExpr{Parts: []string{"security_group_ids"}},
									Value: &types.ForArrayExpr{
										ValueVar:      "sg",
										KeyVar:        "sg_key",
										Collection:    &types.ReferenceExpr{Parts: []string{"aws_security_group", "complex"}},
										ThenValueExpr: &types.ReferenceExpr{Parts: []string{"sg", "id"}},
									},
								},
								{
									Key: &types.ReferenceExpr{Parts: []string{"instance_details"}},
									Value: &types.ForMapExpr{
										KeyVar:      "instance",
										Collection:  &types.ReferenceExpr{Parts: []string{"local", "filtered_instances"}},
										ThenKeyExpr: &types.ReferenceExpr{Parts: []string{"instance", "id"}},
										ThenValueExpr: &types.ObjectExpr{
											Items: []types.ObjectItem{
												{
													Key:   &types.ReferenceExpr{Parts: []string{"name"}},
													Value: &types.ReferenceExpr{Parts: []string{"instance", "name"}},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"private_ip"}},
													Value: &types.ReferenceExpr{Parts: []string{"instance", "private_ip"}},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"public_ip"}},
													Value: &types.ReferenceExpr{Parts: []string{"instance", "public_ip"}},
												},
												{
													Key: &types.ReferenceExpr{Parts: []string{"subnet"}},
													Value: &types.ObjectExpr{
														Items: []types.ObjectItem{
															{
																Key:   &types.ReferenceExpr{Parts: []string{"id"}},
																Value: &types.ReferenceExpr{Parts: []string{"instance", "subnet_id"}},
															},
															{
																Key: &types.ReferenceExpr{Parts: []string{"details"}},
																Value: &types.FunctionCallExpr{
																	Name: "lookup",
																	Args: []types.Expression{
																		&types.ReferenceExpr{Parts: []string{"local", "subnet_map"}},
																		&types.ReferenceExpr{Parts: []string{"instance", "subnet_id"}},
																		&types.LiteralValue{Value: nil, ValueType: "null"},
																	},
																},
															},
														},
													},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"environment"}},
													Value: &types.ReferenceExpr{Parts: []string{"instance", "environment"}},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"type"}},
													Value: &types.ReferenceExpr{Parts: []string{"instance", "type"}},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"tags"}},
													Value: &types.ReferenceExpr{Parts: []string{"instance", "tags"}},
												},
											},
										},
									},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"backup_enabled"}},
									Value: &types.ReferenceExpr{Parts: []string{"var", "enable_backups"}},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"backup_config"}},
									Value: &types.ReferenceExpr{Parts: []string{"local", "backup_config"}},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"naming_convention"}},
									Value: &types.ReferenceExpr{Parts: []string{"local", "naming_convention"}},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"complex_calculation"}},
									Value: &types.ReferenceExpr{Parts: []string{"module", "complex_module", "complex_calculation"}},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"policy_document"}},
									Value: &types.ReferenceExpr{Parts: []string{"data", "aws_iam_policy_document", "complex", "json"}},
								},
							},
						},
					},
					&types.Attribute{
						Name:  "sensitive",
						Value: &types.LiteralValue{Value: true, ValueType: "bool"},
					},
					&types.Attribute{
						Name: "depends_on",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.ReferenceExpr{Parts: []string{"module", "complex_module"}},
								&types.ReferenceExpr{Parts: []string{"aws_security_group", "complex"}},
								&types.ReferenceExpr{Parts: []string{"data", "aws_iam_policy_document", "complex"}},
							},
						},
					},
				},
			},
		},
	}
}

// createComplexProviderExpected creates the expected structure for 07_complex_provider.tf
func createComplexProviderExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "provider",
				Labels:       []string{"aws"},
				BlockComment: "// Provider configuration with complex expressions",
				Children: []types.Body{
					&types.Attribute{
						Name:  "region",
						Value: &types.ReferenceExpr{Parts: []string{"var", "region"}},
					},
					&types.Block{
						Type: "assume_role",
						Children: []types.Body{
							&types.Attribute{
								Name: "role_arn",
								Value: &types.ConditionalExpr{
									Condition: &types.BinaryExpr{
										Left:     &types.ReferenceExpr{Parts: []string{"var", "environment"}},
										Operator: "==",
										Right:    &types.LiteralValue{Value: "prod", ValueType: "string"},
									},
									TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "prod_role_arn"}},
									FalseExpr: &types.ReferenceExpr{Parts: []string{"var", "non_prod_role_arn"}},
								},
							},
							&types.Attribute{
								Name: "session_name",
								Value: &types.TemplateExpr{
									Parts: []types.Expression{
										&types.LiteralValue{Value: "terraform-", ValueType: "string"},
										&types.ReferenceExpr{Parts: []string{"var", "environment"}},
										&types.LiteralValue{Value: "-", ValueType: "string"},
										&types.FunctionCallExpr{
											Name: "formatdate",
											Args: []types.Expression{&types.LiteralValue{Value: "YYYYMMDDhhmmss", ValueType: "string"}, &types.FunctionCallExpr{Name: "timestamp"}},
										},
									},
								},
							},
							&types.Attribute{
								Name:  "external_id",
								Value: &types.ReferenceExpr{Parts: []string{"var", "external_id"}},
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
										&types.ReferenceExpr{Parts: []string{"var", "common_tags"}},
										&types.ObjectExpr{
											Items: []types.ObjectItem{
												{
													Key:   &types.ReferenceExpr{Parts: []string{"Environment"}},
													Value: &types.ReferenceExpr{Parts: []string{"var", "environment"}},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"ManagedBy"}},
													Value: &types.LiteralValue{Value: "Terraform", ValueType: "string"},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"Project"}},
													Value: &types.ReferenceExpr{Parts: []string{"var", "project_name"}},
												},
												{
													Key:   &types.ReferenceExpr{Parts: []string{"Owner"}},
													Value: &types.ReferenceExpr{Parts: []string{"var", "owner"}},
												},
												{
													Key: &types.ReferenceExpr{Parts: []string{"CreatedAt"}},
													Value: &types.FunctionCallExpr{
														Name: "formatdate",
														Args: []types.Expression{
															&types.LiteralValue{Value: "YYYY-MM-DD", ValueType: "string"},
															&types.FunctionCallExpr{Name: "timestamp"},
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
										Left:     &types.ReferenceExpr{Parts: []string{"var", "custom_endpoints"}},
										Operator: "!=",
										Right:    &types.LiteralValue{Value: nil, ValueType: "null"},
									},
									TrueExpr:  &types.ReferenceExpr{Parts: []string{"var", "custom_endpoints"}},
									FalseExpr: &types.ObjectExpr{},
								},
							},
							&types.Block{
								Type: "content",
								Children: []types.Body{
									&types.Attribute{
										Name: "apigateway",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "apigateway", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
											},
										},
									},
									&types.Attribute{
										Name: "cloudformation",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "cloudformation", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
											},
										},
									},
									&types.Attribute{
										Name: "cloudwatch",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "cloudwatch", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
											},
										},
									},
									&types.Attribute{
										Name: "dynamodb",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "dynamodb", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
											},
										},
									},
									&types.Attribute{
										Name: "ec2",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "ec2", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
											},
										},
									},
									&types.Attribute{
										Name: "s3",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "s3", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
											},
										},
									},
									&types.Attribute{
										Name: "sts",
										Value: &types.FunctionCallExpr{
											Name: "lookup",
											Args: []types.Expression{
												&types.ReferenceExpr{Parts: []string{"endpoints", "value"}},
												&types.LiteralValue{Value: "sts", ValueType: "string"},
												&types.LiteralValue{Value: nil, ValueType: "null"},
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
	}
}

// createComplexTerraformConfigExpected creates the expected structure for 08_complex_terraform_config.tf
func createComplexTerraformConfigExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			&types.Block{
				Type:         "terraform",
				BlockComment: "// Terraform configuration with complex expressions",
				Children: []types.Body{
					&types.Attribute{
						Name:  "required_version",
						Value: &types.LiteralValue{Value: ">= 1.0.0", ValueType: "string"},
					},
					&types.Block{
						Type: "required_providers",
						Children: []types.Body{
							&types.Attribute{
								Name: "aws",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key:   &types.ReferenceExpr{Parts: []string{"source"}},
											Value: &types.LiteralValue{Value: "hashicorp/aws", ValueType: "string"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"version"}},
											Value: &types.LiteralValue{Value: ">= 4.0.0, < 5.0.0", ValueType: "string"},
										},
									},
								},
							},
							&types.Attribute{
								Name: "random",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key:   &types.ReferenceExpr{Parts: []string{"source"}},
											Value: &types.LiteralValue{Value: "hashicorp/random", ValueType: "string"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"version"}},
											Value: &types.LiteralValue{Value: "~> 3.1", ValueType: "string"},
										},
									},
								},
							},
							&types.Attribute{
								Name: "null",
								Value: &types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key:   &types.ReferenceExpr{Parts: []string{"source"}},
											Value: &types.LiteralValue{Value: "hashicorp/null", ValueType: "string"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"version"}},
											Value: &types.LiteralValue{Value: "~> 3.1", ValueType: "string"},
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
										&types.LiteralValue{Value: "terraform-state-", ValueType: "string"},
										&types.ReferenceExpr{Parts: []string{"var", "environment"}},
										&types.LiteralValue{Value: "-", ValueType: "string"},
										&types.ReferenceExpr{Parts: []string{"var", "account_id"}},
									},
								},
							},
							&types.Attribute{
								Name:  "key",
								Value: &types.LiteralValue{Value: "complex-module/terraform.tfstate", ValueType: "string"},
							},
							&types.Attribute{
								Name:  "region",
								Value: &types.ReferenceExpr{Parts: []string{"var", "region"}},
							},
							&types.Attribute{
								Name:  "dynamodb_table",
								Value: &types.LiteralValue{Value: "terraform-locks", ValueType: "string"},
							},
							&types.Attribute{
								Name:  "encrypt",
								Value: &types.LiteralValue{Value: true, ValueType: "bool"},
							},
							&types.Attribute{
								Name:  "kms_key_id",
								Value: &types.ReferenceExpr{Parts: []string{"var", "kms_key_id"}},
							},
							&types.Attribute{
								Name:  "role_arn",
								Value: &types.ReferenceExpr{Parts: []string{"var", "state_role_arn"}},
							},
						},
					},
					&types.Attribute{
						Name: "experiments",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.ReferenceExpr{Parts: []string{"module_variable_optional_attrs"}},
							},
						},
					},
				},
			},
		},
	}
}
