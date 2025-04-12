
resource "aws_db_subnet_group" "smartgou_db_subnet_group" {
  name       = "smartgou-db-subnet-group"
  subnet_ids = [var.subnet_id_main, var.subnet_id_backup]

  tags = {
    Name = "smartgou-db-subnet-group"
  }
}

resource "aws_db_instance" "smartGou" {
  allocated_storage    = var.db_allocated_storage
  db_name              = var.db_name
  engine               = var.db_engine
  engine_version       = var.db_engine_version
  instance_class       = var.db_instance_class
  username             = var.db_username
  password             = var.db_password
  parameter_group_name = var.db_parameter_group_name
  skip_final_snapshot  = var.db_skip_final_snapshot
  publicly_accessible  = var.db_publicly_accessible
  db_subnet_group_name = aws_db_subnet_group.smartgou_db_subnet_group.name
  vpc_security_group_ids = [var.rds_security_group_id]
}
