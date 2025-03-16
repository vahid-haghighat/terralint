package internal

//switch nodeType {
////case "end": // ts_builtin_sym_end
////case "=": // anon_sym_EQ
////case "{": // anon_sym_LBRACE
////case "}": // anon_sym_RBRACE
//case "identifier": // sym_identifier
//section.Body = node
////case "(": // anon_sym_LPAREN
////case ")": // anon_sym_RPAREN
////case "numeric_lit_token1": // aux_sym_numeric_lit_token1
////case "numeric_lit_token2": // aux_sym_numeric_lit_token2
////case "true": // anon_sym_true
////case "false": // anon_sym_false
////case "null_lit": // sym_null_lit
////case ",": // anon_sym_COMMA
////case "[": // anon_sym_LBRACK
////case "]": // anon_sym_RBRACK
////case ":": // anon_sym_COLON
////case ".": // anon_sym_DOT
////case "legacy_index_token1": // aux_sym_legacy_index_token1
////case ".*": // anon_sym_DOT_STAR
////case "[*]": // anon_sym_LBRACK_STAR_RBRACK
////case "=>": // anon_sym_EQ_GT
////case "for": // anon_sym_for
////case "in": // anon_sym_in
////case "if": // anon_sym_if
////case "ellipsis": // sym_ellipsis
////case "?": // anon_sym_QMARK
////case "-": // anon_sym_DASH
////case "!": // anon_sym_BANG
////case "*": // anon_sym_STAR
////case "/": // anon_sym_SLASH
////case "%": // anon_sym_PERCENT
////case "+": // anon_sym_PLUS
////case ">": // anon_sym_GT
////case ">=": // anon_sym_GT_EQ
////case "<": // anon_sym_LT
////case "<=": // anon_sym_LT_EQ
////case "==": // anon_sym_EQ_EQ
////case "!=": // anon_sym_BANG_EQ
////case "&&": // anon_sym_AMP_AMP
////case "||": // anon_sym_PIPE_PIPE
////case "<<": // anon_sym_LT_LT
////case "<<-": // anon_sym_LT_LT_DASH
////case "strip_marker": // sym_strip_marker
////case "endfor": // anon_sym_endfor
////case "else": // anon_sym_else
////case "endif": // anon_sym_endif
////case "comment": // sym_comment
////case "_whitespace": // sym__whitespace
////case "quoted_template_start": // sym_quoted_template_start
////case "quoted_template_end": // sym_quoted_template_end
////case "_template_literal_chunk": // sym__template_literal_chunk
////case "template_interpolation_start": // sym_template_interpolation_start
////case "template_interpolation_end": // sym_template_interpolation_end
////case "template_directive_start": // sym_template_directive_start
////case "template_directive_end": // sym_template_directive_end
////case "heredoc_identifier": // sym_heredoc_identifier
//case "config_file": // sym_config_file
//section.Children = walkChildNodes(src, node)
//case "body": // sym_body
//section.Children = walkChildNodes(src, node)
//case "attribute": // sym_attribute
//section.Body = node
//section.Children = walkChildNodes(src, node)
////a := Attribute{}
////for i := 0; i < int(node.NamedChildCount()); i++ {
////	child := node.NamedChild(i)
////	childText := string(src[child.StartByte():child.EndByte()])
////	switch child.Type() {
////	case "identifier":
////		a.Name = childText
////	case "expression":
////		a.Value = child.Range()
////		walkChildNodes(src, child)
////	}
////}
//case "block": // sym_block
////b := Block{}
////for i := 0; i < int(node.NamedChildCount()); i++ {
////	child := node.NamedChild(i)
////	childText := string(src[child.StartByte():child.EndByte()])
////	switch child.Type() {
////	case "identifier":
////		b.Type = childText
////	case "string_lit":
////		b.Labels = append(b.Labels, childText)
////	case "body":
////		b.Range = child.Range()
////		walkChildNodes(src, child)
////	}
////}
////section.Body = node
//section.Children = walkChildNodes(src, node)
//case "block_start": // sym_block_start
//case "block_end": // sym_block_end
//case "expression": // sym_expression
//section.Body = node
//section.Children = walkChildNodes(src, node)
////case "_expr_term": // sym__expr_term
//case "literal_value": // sym_literal_value
////case "numeric_lit": // sym_numeric_lit
////case "bool_lit": // sym_bool_lit
//case "string_lit": // sym_string_lit
//section.Body = node
//case "collection_value": // sym_collection_value
//section.Body = node
//section.Children = walkChildNodes(src, node)
////case "_comma": // sym__comma
//case "tuple": // sym_tuple
//section.Body = node
//section.Children = walkChildNodes(src, node)
//case "tuple_start": // sym_tuple_start
//case "tuple_end": // sym_tuple_end
////case "_tuple_elems": // sym__tuple_elems
//case "object": // sym_object
//section.Body = node
//section.Children = walkChildNodes(src, node)
//case "object_start": // sym_object_start
//case "object_end": // sym_object_end
////case "_object_elems": // sym__object_elems
//case "object_elem": // sym_object_elem
//section.Body = node
//section.Children = walkChildNodes(src, node)
////case "index": // sym_index
////case "new_index": // sym_new_index
////case "legacy_index": // sym_legacy_index
//case "get_attr": // sym_get_attr
////case "splat": // sym_splat
////case "attr_splat": // sym_attr_splat
////case "full_splat": // sym_full_splat
//case "for_expr": // sym_for_expr
//section.Body = node
//section.Children = walkChildNodes(src, node)
//case "for_tuple_expr": // sym_for_tuple_expr
//section.Body = node
//section.Children = walkChildNodes(src, node)
//case "for_object_expr": // sym_for_object_expr
//section.Body = node
//section.Children = walkChildNodes(src, node)
//case "for_intro": // sym_for_intro
//section.Body = node
//section.Children = walkChildNodes(src, node)
////case "for_cond": // sym_for_cond
//case "variable_expr": // sym_variable_expr
//case "function_call": // sym_function_call
//section.Body = node
//section.Children = walkChildNodes(src, node)
////case "_function_call_start": // sym__function_call_start
////case "_function_call_end": // sym__function_call_end
//case "function_arguments": // sym_function_arguments
////case "conditional": // sym_conditional
//case "operation": // sym_operation
////case "unary_operation": // sym_unary_operation
////case "binary_operation": // sym_binary_operation
////case "template_expr": // sym_template_expr
////case "quoted_template": // sym_quoted_template
////case "heredoc_template": // sym_heredoc_template
////case "heredoc_start": // sym_heredoc_start
////case "_template": // aux_sym__template
////case "template_literal": // sym_template_literal
////case "template_interpolation": // sym_template_interpolation
////case "template_directive": // sym_template_directive
////case "template_for": // sym_template_for
////case "template_for_start": // sym_template_for_start
////case "template_for_end": // sym_template_for_end
////case "template_if": // sym_template_if
////case "template_if_intro": // sym_template_if_intro
////case "template_else_intro": // sym_template_else_intro
////case "template_if_end": // sym_template_if_end
////case "body_repeat1": // aux_sym_body_repeat1
////case "block_repeat1": // aux_sym_block_repeat1
////case "_tuple_elems_repeat1": // aux_sym__tuple_elems_repeat1
////case "_object_elems_repeat1": // aux_sym__object_elems_repeat1
////case "attr_splat_repeat1": // aux_sym_attr_splat_repeat1
////case "template_literal_repeat1": // aux_sym_template_literal_repeat1
//default:
//fmt.Printf("%s: %s\n", nodeType, nodeText)
////return hclwrite.TokensForIdentifier(string(src[node.StartByte():node.EndByte()]))
//}
