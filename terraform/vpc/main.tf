resource "aws_vpc" "smartGou_vpc" {
  cidr_block = var.cidr_block
  enable_dns_support = true
  enable_dns_hostnames = true
  instance_tenancy = "default"

  tags = {
    Name = "smartGou-vpc"
  }
}

resource "aws_subnet" "smartGou_subnet_main" {
  vpc_id                  = aws_vpc.smartGou_vpc.id
  cidr_block              = var.subnet_cidr_block
  availability_zone       = var.availability_zone_main
  map_public_ip_on_launch = true

  tags = {
    Name = "smartGou-subnet-main"
  }
}

resource "aws_subnet" "smartGou_subnet_backup" {
  vpc_id                  = aws_vpc.smartGou_vpc.id
  cidr_block              = var.subnet_cidr_block_backup
  availability_zone       = var.availability_zone_backup
  map_public_ip_on_launch = true

  tags = {
    Name = "smartGou-subnet-backup"
  }
}

resource "aws_internet_gateway" "smartGou_igw" {
  vpc_id = aws_vpc.smartGou_vpc.id

  tags = {
    Name = "smartGou-igw"
  }
}

resource "aws_route_table" "smartGou_route_table" {
  vpc_id = aws_vpc.smartGou_vpc.id

}

resource "aws_route" "smartGou_route" {
  route_table_id            = aws_route_table.smartGou_route_table.id
  destination_cidr_block    = "0.0.0.0/0"
  gateway_id                = aws_internet_gateway.smartGou_igw.id
}

resource "aws_route_table_association" "smartGou_route_table_association" {
  subnet_id      = aws_subnet.smartGou_subnet_main.id
  route_table_id = aws_route_table.smartGou_route_table.id
}

resource "aws_route_table_association" "smartGou_route_table_association_backup" {
  subnet_id      = aws_subnet.smartGou_subnet_backup.id
  route_table_id = aws_route_table.smartGou_route_table.id
}





