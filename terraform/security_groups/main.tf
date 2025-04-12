resource "aws_security_group" "smartGou_rds_sg" {
  name        = "smartGou-rds-sg"
  description = "Security group for smartGou RDS DB instance"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = var.cidr_blocks
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"  # allow all outgoing traffic
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "RDS Security Group"
  }
}
