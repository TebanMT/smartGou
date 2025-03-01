resource "aws_vpc" "smartGo_vpc" {
  cidr_block = var.cidr_block
  enable_dns_support = true
  enable_dns_hostnames = true
  instance_tenancy = "default"

  tags = {
    Name = "smartGo-vpc"
  }
}

resource "aws_subnet" "smartGo_subnet_main" {
  vpc_id                  = aws_vpc.smartGo_vpc.id
  cidr_block              = var.subnet_cidr_block
  availability_zone       = var.availability_zone_main
  map_public_ip_on_launch = true

  tags = {
    Name = "smartGo-subnet-main"
  }
}

resource "aws_subnet" "smartGo_subnet_backup" {
  vpc_id                  = aws_vpc.smartGo_vpc.id
  cidr_block              = var.subnet_cidr_block_backup
  availability_zone       = var.availability_zone_backup
  map_public_ip_on_launch = true

  tags = {
    Name = "smartGo-subnet-backup"
  }
}

resource "aws_internet_gateway" "smartGo_igw" {
  vpc_id = aws_vpc.smartGo_vpc.id

  tags = {
    Name = "smartGo-igw"
  }
}

resource "aws_route_table" "smartGo_route_table" {
  vpc_id = aws_vpc.smartGo_vpc.id

}

resource "aws_route" "smartGo_route" {
  route_table_id            = aws_route_table.smartGo_route_table.id
  destination_cidr_block    = "0.0.0.0/0"
  gateway_id                = aws_internet_gateway.smartGo_igw.id
}

resource "aws_route_table_association" "smartGo_route_table_association" {
  subnet_id      = aws_subnet.smartGo_subnet_main.id
  route_table_id = aws_route_table.smartGo_route_table.id
}

resource "aws_route_table_association" "smartGo_route_table_association_backup" {
  subnet_id      = aws_subnet.smartGo_subnet_backup.id
  route_table_id = aws_route_table.smartGo_route_table.id
}





