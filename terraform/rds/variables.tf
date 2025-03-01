variable "db_allocated_storage" {
  description = "Storage size of the database (in GB)"
  default     = 5
}

variable "db_instance_class" {
  description = "Database instance type"
  default     = "db.t3.micro"
}

variable "db_engine" {
  description = "Database engine"
  default     = "postgres"
}

variable "db_engine_version" {
  description = "Database engine version"
  default     = "17.3"
}

variable "db_parameter_group_name" {
  description = "Database parameter group name"
  default     = "default.postgres17"
}

variable "db_skip_final_snapshot" {
  description = "Skip final snapshot"
  default     = true
}

variable "db_publicly_accessible" {
  description = "Publicly accessible"
  default     = true
}

variable "db_name" {
  description = "Database name"
  default     = "smartGo"
}

variable "db_username" {
  description = "Database username"
}

variable "db_password" {
  description = "Database password"
  sensitive   = true
}

variable "subnet_id_main" {
  description = "Subnet ID for the main availability zone"
}

variable "subnet_id_backup" {
  description = "Subnet ID for the backup availability zone"
}

variable "rds_security_group_id" {
  description = "Security group ID for the RDS instance"
}

