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
					// Rest of the implementation would be here
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
