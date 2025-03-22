// VPC Module - Creates a VPC with public, private, and database subnets

// Create the VPC
resource "aws_vpc" "this" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = var.enable_dns_hostnames
  enable_dns_support   = var.enable_dns_support
  
  tags = merge(
    var.tags,
    {
      Name = var.vpc_name
    }
  )
}

// Create Internet Gateway
resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-igw"
    }
  )
}

// Create public subnets
resource "aws_subnet" "public" {
  count = length(var.public_subnets)
  
  vpc_id                  = aws_vpc.this.id
  cidr_block              = var.public_subnets[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-public-${var.availability_zones[count.index]}"
      Tier = "Public"
    }
  )
}

// Create private subnets
resource "aws_subnet" "private" {
  count = length(var.private_subnets)
  
  vpc_id                  = aws_vpc.this.id
  cidr_block              = var.private_subnets[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = false
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-private-${var.availability_zones[count.index]}"
      Tier = "Private"
    }
  )
}

// Create database subnets
resource "aws_subnet" "database" {
  count = length(var.database_subnets)
  
  vpc_id                  = aws_vpc.this.id
  cidr_block              = var.database_subnets[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = false
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-database-${var.availability_zones[count.index]}"
      Tier = "Database"
    }
  )
}

// Create NAT Gateway(s)
resource "aws_eip" "nat" {
  count = var.enable_nat_gateway ? (var.single_nat_gateway ? 1 : length(var.public_subnets)) : 0
  
  domain = "vpc"
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-nat-eip-${count.index + 1}"
    }
  )
}

resource "aws_nat_gateway" "this" {
  count = var.enable_nat_gateway ? (var.single_nat_gateway ? 1 : length(var.public_subnets)) : 0
  
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-nat-gw-${count.index + 1}"
    }
  )
  
  depends_on = [aws_internet_gateway.this]
}

// Create VPN Gateway
resource "aws_vpn_gateway" "this" {
  count = var.enable_vpn_gateway ? 1 : 0
  
  vpc_id = aws_vpc.this.id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-vpn-gw"
    }
  )
}

// Create route tables
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-public-rt"
      Tier = "Public"
    }
  )
}

resource "aws_route" "public_internet_gateway" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.this.id
  
  timeouts {
    create = "5m"
  }
}

resource "aws_route_table" "private" {
  count = var.enable_nat_gateway ? (var.single_nat_gateway ? 1 : length(var.private_subnets)) : 1
  
  vpc_id = aws_vpc.this.id
  
  tags = merge(
    var.tags,
    {
      Name = var.single_nat_gateway ? "${var.vpc_name}-private-rt" : "${var.vpc_name}-private-rt-${var.availability_zones[count.index]}"
      Tier = "Private"
    }
  )
}

resource "aws_route" "private_nat_gateway" {
  count = var.enable_nat_gateway ? (var.single_nat_gateway ? 1 : length(var.private_subnets)) : 0
  
  route_table_id         = aws_route_table.private[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = var.single_nat_gateway ? aws_nat_gateway.this[0].id : aws_nat_gateway.this[count.index].id
  
  timeouts {
    create = "5m"
  }
}

resource "aws_route_table" "database" {
  count = length(var.database_subnets) > 0 ? (var.single_nat_gateway || !var.enable_nat_gateway ? 1 : length(var.database_subnets)) : 0
  
  vpc_id = aws_vpc.this.id
  
  tags = merge(
    var.tags,
    {
      Name = var.single_nat_gateway || !var.enable_nat_gateway ? "${var.vpc_name}-database-rt" : "${var.vpc_name}-database-rt-${var.availability_zones[count.index]}"
      Tier = "Database"
    }
  )
}

resource "aws_route" "database_nat_gateway" {
  count = var.enable_nat_gateway && length(var.database_subnets) > 0 ? (var.single_nat_gateway ? 1 : length(var.database_subnets)) : 0
  
  route_table_id         = aws_route_table.database[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = var.single_nat_gateway ? aws_nat_gateway.this[0].id : aws_nat_gateway.this[count.index].id
  
  timeouts {
    create = "5m"
  }
}

// Associate route tables with subnets
resource "aws_route_table_association" "public" {
  count = length(var.public_subnets)
  
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count = length(var.private_subnets)
  
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = var.single_nat_gateway || !var.enable_nat_gateway ? aws_route_table.private[0].id : aws_route_table.private[count.index].id
}

resource "aws_route_table_association" "database" {
  count = length(var.database_subnets)
  
  subnet_id      = aws_subnet.database[count.index].id
  route_table_id = var.single_nat_gateway || !var.enable_nat_gateway ? aws_route_table.database[0].id : aws_route_table.database[count.index].id
}

// Create VPC Endpoints for S3 and DynamoDB
resource "aws_vpc_endpoint" "s3" {
  count = var.enable_s3_endpoint ? 1 : 0
  
  vpc_id       = aws_vpc.this.id
  service_name = "com.amazonaws.${data.aws_region.current.name}.s3"
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-s3-endpoint"
    }
  )
}

resource "aws_vpc_endpoint_route_table_association" "s3_public" {
  count = var.enable_s3_endpoint ? 1 : 0
  
  route_table_id  = aws_route_table.public.id
  vpc_endpoint_id = aws_vpc_endpoint.s3[0].id
}

resource "aws_vpc_endpoint_route_table_association" "s3_private" {
  count = var.enable_s3_endpoint ? (var.single_nat_gateway ? 1 : length(var.private_subnets)) : 0
  
  route_table_id  = aws_route_table.private[count.index].id
  vpc_endpoint_id = aws_vpc_endpoint.s3[0].id
}

resource "aws_vpc_endpoint" "dynamodb" {
  count = var.enable_dynamodb_endpoint ? 1 : 0
  
  vpc_id       = aws_vpc.this.id
  service_name = "com.amazonaws.${data.aws_region.current.name}.dynamodb"
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-dynamodb-endpoint"
    }
  )
}

resource "aws_vpc_endpoint_route_table_association" "dynamodb_public" {
  count = var.enable_dynamodb_endpoint ? 1 : 0
  
  route_table_id  = aws_route_table.public.id
  vpc_endpoint_id = aws_vpc_endpoint.dynamodb[0].id
}

resource "aws_vpc_endpoint_route_table_association" "dynamodb_private" {
  count = var.enable_dynamodb_endpoint ? (var.single_nat_gateway ? 1 : length(var.private_subnets)) : 0
  
  route_table_id  = aws_route_table.private[count.index].id
  vpc_endpoint_id = aws_vpc_endpoint.dynamodb[0].id
}

// Create Network ACLs
resource "aws_network_acl" "public" {
  count = var.create_network_acls ? 1 : 0
  
  vpc_id     = aws_vpc.this.id
  subnet_ids = aws_subnet.public[*].id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-public-nacl"
      Tier = "Public"
    }
  )
}

resource "aws_network_acl_rule" "public_ingress" {
  count = var.create_network_acls ? 1 : 0
  
  network_acl_id = aws_network_acl.public[0].id
  rule_number    = 100
  egress         = false
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl_rule" "public_egress" {
  count = var.create_network_acls ? 1 : 0
  
  network_acl_id = aws_network_acl.public[0].id
  rule_number    = 100
  egress         = true
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl" "private" {
  count = var.create_network_acls ? 1 : 0
  
  vpc_id     = aws_vpc.this.id
  subnet_ids = aws_subnet.private[*].id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-private-nacl"
      Tier = "Private"
    }
  )
}

resource "aws_network_acl_rule" "private_ingress" {
  count = var.create_network_acls ? 1 : 0
  
  network_acl_id = aws_network_acl.private[0].id
  rule_number    = 100
  egress         = false
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl_rule" "private_egress" {
  count = var.create_network_acls ? 1 : 0
  
  network_acl_id = aws_network_acl.private[0].id
  rule_number    = 100
  egress         = true
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl" "database" {
  count = var.create_network_acls && length(var.database_subnets) > 0 ? 1 : 0
  
  vpc_id     = aws_vpc.this.id
  subnet_ids = aws_subnet.database[*].id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-database-nacl"
      Tier = "Database"
    }
  )
}

resource "aws_network_acl_rule" "database_ingress" {
  count = var.create_network_acls && length(var.database_subnets) > 0 ? 1 : 0
  
  network_acl_id = aws_network_acl.database[0].id
  rule_number    = 100
  egress         = false
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

resource "aws_network_acl_rule" "database_egress" {
  count = var.create_network_acls && length(var.database_subnets) > 0 ? 1 : 0
  
  network_acl_id = aws_network_acl.database[0].id
  rule_number    = 100
  egress         = true
  protocol       = "-1"
  rule_action    = "allow"
  cidr_block     = "0.0.0.0/0"
}

// Create Flow Logs
resource "aws_flow_log" "this" {
  count = var.enable_flow_logs ? 1 : 0
  
  log_destination      = var.flow_logs_destination_arn
  log_destination_type = "s3"
  traffic_type         = "ALL"
  vpc_id               = aws_vpc.this.id
  
  tags = merge(
    var.tags,
    {
      Name = "${var.vpc_name}-flow-logs"
    }
  )
}

// Data sources
data "aws_region" "current" {}