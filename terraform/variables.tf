variable "aws_region" {
  description = "Region of the AWS account"
  default     = "us-east-1"
}

variable "db_username" {
  description = "Database username"
}

variable "db_password" {
  description = "Database password"
  sensitive   = true
}