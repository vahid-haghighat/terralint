package parser

import (
	"github.com/vahid-haghighat/terralint/parser/types"
)

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
								Value: &types.ArrayExpr{
									Items: []types.Expression{
										&types.LiteralValue{
											Value:     "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*",
											ValueType: "string",
										},
									},
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
											Value:     "hvm",
											ValueType: "string",
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
									Value:     "099720109477",
									ValueType: "string",
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
							},
						},
					},
					&types.Attribute{
						Name: "private_subnets",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "10.0.1.0/24",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "10.0.2.0/24",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "10.0.3.0/24",
									ValueType: "string",
								},
							},
						},
					},
					&types.Attribute{
						Name: "public_subnets",
						Value: &types.ArrayExpr{
							Items: []types.Expression{
								&types.LiteralValue{
									Value:     "10.0.101.0/24",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "10.0.102.0/24",
									ValueType: "string",
								},
								&types.LiteralValue{
									Value:     "10.0.103.0/24",
									ValueType: "string",
								},
							},
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
			// Module with nested expressions, conditionals, and for loops
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
						Name: "vpc_config",
						Value: &types.ObjectExpr{
							Items: []types.ObjectItem{
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"name"},
									},
									Value: &types.TemplateExpr{
										Parts: []types.Expression{
											&types.LiteralValue{
												Value:     "complex-",
												ValueType: "string",
											},
											&types.ReferenceExpr{
												Parts: []string{"var", "environment"},
											},
											&types.LiteralValue{
												Value:     "-",
												ValueType: "string",
											},
											&types.ReferenceExpr{
												Parts: []string{"local", "region_code"},
											},
											&types.LiteralValue{
												Value:     "-vpc",
												ValueType: "string",
											},
										},
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"cidr"},
									},
									Value: &types.FunctionCallExpr{
										Name: "cidrsubnet",
										Args: []types.Expression{
											&types.ReferenceExpr{
												Parts: []string{"var", "base_cidr_block"},
											},
											&types.LiteralValue{
												Value:     4,
												ValueType: "number",
											},
											&types.ReferenceExpr{
												Parts: []string{"var", "vpc_index"},
											},
										},
									},
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"enable_dns"},
									},
									Value: &types.ConditionalExpr{
										Condition: &types.BinaryExpr{
											Left: &types.ReferenceExpr{
												Parts: []string{"var", "environment"},
											},
											Operator: "==",
											Right: &types.LiteralValue{
												Value:     "prod",
												ValueType: "string",
											},
										},
										TrueExpr: &types.LiteralValue{
											Value:     true,
											ValueType: "bool",
										},
										FalseExpr: &types.ParenExpr{
											Expression: &types.ConditionalExpr{
												Condition: &types.BinaryExpr{
													Left: &types.ReferenceExpr{
														Parts: []string{"var", "enable_dns"},
													},
													Operator: "!=",
													Right: &types.LiteralValue{
														Value:     nil,
														ValueType: "null",
													},
												},
												TrueExpr: &types.ReferenceExpr{
													Parts: []string{"var", "enable_dns"},
												},
												FalseExpr: &types.LiteralValue{
													Value:     false,
													ValueType: "bool",
												},
											},
										},
									},
									BlockComment: "// Nested conditional expression",
								},
								{
									Key: &types.ReferenceExpr{
										Parts: []string{"tags"},
									},
									Value: &types.FunctionCallExpr{
										Name: "merge",
										Args: []types.Expression{
											&types.ReferenceExpr{
												Parts: []string{"var", "common_tags"},
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
																	Value:     "complex-",
																	ValueType: "string",
																},
																&types.ReferenceExpr{
																	Parts: []string{"var", "environment"},
																},
																&types.LiteralValue{
																	Value:     "-vpc",
																	ValueType: "string",
																},
															},
														},
													},
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
															Value:     "terraform",
															ValueType: "string",
														},
													},
													{
														Key: &types.ReferenceExpr{
															Parts: []string{"Complex"},
														},
														Value: &types.FunctionCallExpr{
															Name: "jsonencode",
															Args: []types.Expression{
																&types.ObjectExpr{
																	Items: []types.ObjectItem{
																		{
																			Key: &types.ReferenceExpr{
																				Parts: []string{"nested"},
																			},
																			Value: &types.LiteralValue{
																				Value:     "value",
																				ValueType: "string",
																			},
																		},
																		{
																			Key: &types.ReferenceExpr{
																				Parts: []string{"list"},
																			},
																			Value: &types.ArrayExpr{
																				Items: []types.Expression{
																					&types.LiteralValue{
																						Value:     1,
																						ValueType: "number",
																					},
																					&types.LiteralValue{
																						Value:     2,
																						ValueType: "number",
																					},
																					&types.LiteralValue{
																						Value:     3,
																						ValueType: "number",
																					},
																					&types.LiteralValue{
																						Value:     4,
																						ValueType: "number",
																					},
																				},
																			},
																		},
																		{
																			Key: &types.ReferenceExpr{
																				Parts: []string{"map"},
																			},
																			Value: &types.ObjectExpr{
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
																							Value:     42,
																							ValueType: "number",
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
											&types.ConditionalExpr{
												Condition: &types.BinaryExpr{
													Left: &types.ReferenceExpr{
														Parts: []string{"var", "additional_tags"},
													},
													Operator: "!=",
													Right: &types.LiteralValue{
														Value:     nil,
														ValueType: "null",
													},
												},
												TrueExpr: &types.ReferenceExpr{
													Parts: []string{"var", "additional_tags"},
												},
												FalseExpr: &types.ObjectExpr{
													Items: []types.ObjectItem{},
												},
											},
										},
									},
									BlockComment: "// Complex object with nested expressions",
								},
							},
						},
						BlockComment: "// Complex map with nested objects, expressions, and functions",
					},
					&types.Attribute{
						Name: "subnet_cidrs",
						Value: &types.ForExpr{
							KeyVar:     "i",
							ValueVar:   "subnet",
							Collection: &types.ReferenceExpr{Parts: []string{"var", "subnets"}},
							ThenKeyExpr: &types.FunctionCallExpr{
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
								Left:     &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"subnet"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "create"}}},
								Operator: "==",
								Right:    &types.LiteralValue{Value: true, ValueType: "bool"},
							},
						},
						BlockComment: "// Complex for expression with filtering and transformation",
					},
					&types.Attribute{
						Name: "subnet_configs",
						Value: &types.ForExpr{
							KeyVar:      "zone_key",
							ValueVar:    "zone",
							Collection:  &types.ReferenceExpr{Parts: []string{"var", "availability_zones"}},
							ThenKeyExpr: &types.ReferenceExpr{Parts: []string{"zone_key"}},
							ThenValueExpr: &types.ForExpr{
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
																		&types.LiteralValue{
																			Value:     "-",
																			ValueType: "string",
																		},
																		&types.ReferenceExpr{Parts: []string{"subnet"}},
																		&types.LiteralValue{
																			Value:     "-",
																			ValueType: "string",
																		},
																		&types.ReferenceExpr{Parts: []string{"zone"}},
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
						BlockComment: "// Nested for expressions with conditional",
					},
					&types.Attribute{
						Name: "all_subnet_ids",
						Value: &types.FunctionCallExpr{
							Name: "flatten",
							Args: []types.Expression{
								&types.ArrayExpr{
									Items: []types.Expression{
										&types.ForExpr{
											KeyVar:     "zone_key",
											ValueVar:   "zone",
											Collection: &types.ReferenceExpr{Parts: []string{"aws_subnet", "main"}},
											ThenValueExpr: &types.ArrayExpr{
												Items: []types.Expression{
													&types.ForExpr{
														KeyVar:      "subnet_key",
														ValueVar:    "subnet",
														Collection:  &types.ReferenceExpr{Parts: []string{"zone"}},
														ThenKeyExpr: &types.ReferenceExpr{Parts: []string{"subnet", "id"}},
													},
												},
											},
										},
									},
								},
							},
						},
						BlockComment: "// Complex splat expression",
					},
					&types.Attribute{
						Name: "user_data",
						Value: &types.HeredocExpr{
							Marker:   "EOT",
							Indented: true,
							Content:  "#!/bin/bash\necho \"Environment: ${var.environment}\"\necho \"Region: ${data.aws_region.current.name}\"\n\n# Complex interpolation\n${join(\"\\n\", [\n  for script in var.bootstrap_scripts :\n  \"source ${script}\"\n])}\n\n# Conditional section\n${var.install_monitoring ? \"setup_monitoring ${var.monitoring_endpoint}\" : \"echo 'Monitoring disabled'\"}\n\n# For loop in heredoc\n%{for pkg in var.packages~}\nyum install -y ${pkg}\n%{endfor~}\n\n# If directive in heredoc\n%{if var.environment == \"prod\"~}\necho \"Production environment detected, applying strict security\"\n%{else~}\necho \"Non-production environment, using standard security\"\n%{endif~}",
						},
						BlockComment: "// Heredoc with interpolation",
					},
					&types.Attribute{
						Name: "timeout",
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
						BlockComment: "// Complex binary expressions with nested conditionals",
					},
					&types.Attribute{
						Name: "complex_calculation",
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
						BlockComment: "// Nested parentheses and operators",
					},
					&types.Attribute{
						Name: "security_groups",
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
						Name: "custom_template",
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
											Key: &types.ReferenceExpr{Parts: []string{"region"}},
											Value: &types.RelativeTraversalExpr{
												Source: &types.RelativeTraversalExpr{
													Source:    &types.ReferenceExpr{Parts: []string{"data", "aws_region", "current"}},
													Traversal: []types.TraversalElem{{Type: "attr", Name: "name"}},
												},
											},
										},
									},
								},
							},
						},
					},
					&types.DynamicBlock{
						ForEach:  &types.ReferenceExpr{Parts: []string{"var", "ingress_rules"}},
						Iterator: "rule", // This matches the file, where iterator is defined as "rule"
						Labels:   []string{"ingress"},
						Content: []types.Body{
							&types.Attribute{
								Name: "description",
								Value: &types.FunctionCallExpr{
									Name: "lookup",
									Args: []types.Expression{
										&types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"rule", "value"}}},
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
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"rule", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "from_port"}}},
							},
							&types.Attribute{
								Name:  "to_port",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"rule", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "to_port"}}},
							},
						},
					},
					&types.Attribute{
						Name: "validation",
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

			// Resource with complex dynamic blocks and for_each
			&types.Block{
				Type:         "resource",
				Labels:       []string{"aws_security_group", "complex"},
				BlockComment: "// Resource with complex dynamic blocks and for_each",
				Children: []types.Body{
					&types.Attribute{
						Name: "for_each",
						Value: &types.ForExpr{
							KeyVar:     "",
							ValueVar:   "sg",
							Collection: &types.ReferenceExpr{Parts: []string{"var", "security_groups"}},
							ThenKeyExpr: &types.RelativeTraversalExpr{
								Source:    &types.ReferenceExpr{Parts: []string{"sg"}},
								Traversal: []types.TraversalElem{{Type: "attr", Name: "name"}},
							},
							ThenValueExpr: &types.ReferenceExpr{Parts: []string{"sg"}},
							Condition: &types.RelativeTraversalExpr{
								Source:    &types.ReferenceExpr{Parts: []string{"sg"}},
								Traversal: []types.TraversalElem{{Type: "attr", Name: "create"}},
							},
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
						Name: "description",
						Value: &types.RelativeTraversalExpr{
							Source:    &types.ReferenceExpr{Parts: []string{"each", "value"}},
							Traversal: []types.TraversalElem{{Type: "attr", Name: "description"}},
						},
					},
					&types.Attribute{
						Name:  "vpc_id",
						Value: &types.ReferenceExpr{Parts: []string{"var", "vpc_id"}},
					},
					&types.DynamicBlock{
						ForEach:  &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"each", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "ingress_rules"}}},
						Iterator: "ingress",
						Labels:   []string{"ingress"},
						Content: []types.Body{
							&types.Attribute{
								Name:  "description",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"ingress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "description"}}},
							},
							&types.Attribute{
								Name:  "from_port",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"ingress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "from_port"}}},
							},
							&types.Attribute{
								Name:  "to_port",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"ingress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "to_port"}}},
							},
							&types.Attribute{
								Name:  "protocol",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"ingress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "protocol"}}},
							},
							&types.Attribute{
								Name:  "cidr_blocks",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"ingress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "cidr_blocks"}}},
							},
						},
					},
					&types.DynamicBlock{
						ForEach:  &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"each", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "egress_rules"}}},
						Iterator: "egress",
						Labels:   []string{"egress"},
						Content: []types.Body{
							&types.Attribute{
								Name:  "description",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"egress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "description"}}},
							},
							&types.Attribute{
								Name:  "from_port",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"egress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "from_port"}}},
							},
							&types.Attribute{
								Name:  "to_port",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"egress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "to_port"}}},
							},
							&types.Attribute{
								Name:  "protocol",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"egress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "protocol"}}},
							},
							&types.Attribute{
								Name:  "cidr_blocks",
								Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"egress", "value"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "cidr_blocks"}}},
							},
						},
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
						},
					},
				},
			},

			// Locals block
			&types.Block{
				Type:         "locals",
				BlockComment: "// Complex locals with nested expressions",
				Children: []types.Body{
					&types.Attribute{
						Name: "subnet_map",
						Value: &types.ForExpr{
							ValueVar:   "subnet",
							Collection: &types.ReferenceExpr{Parts: []string{"var", "subnets"}},
							ThenKeyExpr: &types.RelativeTraversalExpr{
								Source:    &types.ReferenceExpr{Parts: []string{"subnet"}},
								Traversal: []types.TraversalElem{{Type: "attr", Name: "name"}},
							},
							ThenValueExpr: &types.ObjectExpr{
								Items: []types.ObjectItem{
									{
										Key:   &types.ReferenceExpr{Parts: []string{"id"}},
										Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"subnet"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "id"}}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"cidr"}},
										Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"subnet"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "cidr"}}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"az"}},
										Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"subnet"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "availability_zone"}}},
									},
									{
										Key:   &types.ReferenceExpr{Parts: []string{"public"}},
										Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"subnet"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "public"}}},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name: "timeout_seconds",
						Value: &types.ConditionalExpr{
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
								TrueExpr: &types.ConditionalExpr{
									Condition: &types.ReferenceExpr{Parts: []string{"var", "high_availability"}},
									TrueExpr:  &types.LiteralValue{Value: 300, ValueType: "number"},
									FalseExpr: &types.LiteralValue{Value: 180, ValueType: "number"},
								},
								FalseExpr: &types.ConditionalExpr{
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
			},

			// Data source
			&types.Block{
				Type:         "data",
				Labels:       []string{"aws_iam_policy_document", "complex"},
				BlockComment: "// Data source with complex expressions",
				Children: []types.Body{
					// Dynamic statement block removed - the actual Terraform file doesn't have this
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
						},
					},
				},
			},

			// Variable with complex type constraints
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
						Value: &types.TypeExpr{
							TypeName: "object",
							Parameters: []types.Expression{
								&types.ObjectExpr{
									Items: []types.ObjectItem{
										{
											Key:   &types.ReferenceExpr{Parts: []string{"name"}},
											Value: &types.TypeExpr{TypeName: "string"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"environment"}},
											Value: &types.TypeExpr{TypeName: "string"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"enabled"}},
											Value: &types.TypeExpr{TypeName: "bool"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"count"}},
											Value: &types.TypeExpr{TypeName: "number"},
										},
										{
											Key:   &types.ReferenceExpr{Parts: []string{"tags"}},
											Value: &types.TypeExpr{TypeName: "map", Parameters: []types.Expression{&types.TypeExpr{TypeName: "string"}}},
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
									Left: &types.RelativeTraversalExpr{
										Source:    &types.ReferenceExpr{Parts: []string{"var", "complex_object"}},
										Traversal: []types.TraversalElem{{Type: "attr", Name: "count"}},
									},
									Operator: ">",
									Right:    &types.LiteralValue{Value: 0, ValueType: "number"},
								},
							},
							&types.Attribute{
								Name:  "error_message",
								Value: &types.LiteralValue{Value: "Count must be between 1 and 10.", ValueType: "string"},
							},
						},
					},
				},
			},

			// Output with complex expressions
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
									Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"module", "complex_module"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "vpc_id"}}},
								},
								{
									Key:   &types.ReferenceExpr{Parts: []string{"subnet_ids"}},
									Value: &types.RelativeTraversalExpr{Source: &types.ReferenceExpr{Parts: []string{"module", "complex_module"}}, Traversal: []types.TraversalElem{{Type: "attr", Name: "subnet_ids"}}},
								},
								{
									Key: &types.ReferenceExpr{Parts: []string{"security_group_ids"}},
									Value: &types.ForExpr{
										KeyVar:     "sg_key",
										ValueVar:   "sg",
										Collection: &types.ReferenceExpr{Parts: []string{"aws_security_group", "complex"}},
										ThenKeyExpr: &types.RelativeTraversalExpr{
											Source:    &types.ReferenceExpr{Parts: []string{"sg"}},
											Traversal: []types.TraversalElem{{Type: "attr", Name: "id"}},
										},
									},
								},
							},
						},
					},
					&types.Attribute{
						Name:  "sensitive",
						Value: &types.LiteralValue{Value: true, ValueType: "bool"},
					},
				},
			},

			// Provider configuration
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
											Args: []types.Expression{
												&types.LiteralValue{Value: "YYYY-MM-DD-hh-mm-ss", ValueType: "string"},
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

			// Terraform configuration
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
							&types.Block{
								Type: "aws",
								Children: []types.Body{
									&types.Attribute{
										Name:  "source",
										Value: &types.LiteralValue{Value: "hashicorp/aws", ValueType: "string"},
									},
									&types.Attribute{
										Name:  "version",
										Value: &types.LiteralValue{Value: ">= 4.0.0, < 5.0.0", ValueType: "string"},
									},
								},
							},
							&types.Block{
								Type: "random",
								Children: []types.Body{
									&types.Attribute{
										Name:  "source",
										Value: &types.LiteralValue{Value: "hashicorp/random", ValueType: "string"},
									},
									&types.Attribute{
										Name:  "version",
										Value: &types.LiteralValue{Value: "~> 3.1", ValueType: "string"},
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
								Name:  "encrypt",
								Value: &types.LiteralValue{Value: true, ValueType: "bool"},
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
			// Implementation would be here
		},
	}
}

// createEdgeCasesExpected creates the expected structure for edge_cases_test.tf
func createEdgeCasesExpected() types.Body {
	return &types.Root{
		Children: []types.Body{
			// Implementation would be here
		},
	}
}
