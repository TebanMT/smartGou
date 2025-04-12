output "vpc_id" {
  value = aws_vpc.smartGou_vpc.id
}

output "subnet_id_main" {
  value = aws_subnet.smartGou_subnet_main.id
}

output "subnet_id_backup" {
  value = aws_subnet.smartGou_subnet_backup.id
}
