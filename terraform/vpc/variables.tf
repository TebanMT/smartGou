variable "cidr_block" {
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}

variable "subnet_cidr_block" {
  description = "CIDR block for the subnet"
  default     = "10.0.1.0/24"
}

variable "subnet_cidr_block_backup" {
  description = "CIDR block for the backup subnet"
  default     = "10.0.2.0/24"
}

variable "availability_zone_main" {
  description = "Availability zone for the subnet"
  default     = "us-east-1a"
}

variable "availability_zone_backup" {
  description = "Availability zone for the subnet"
  default     = "us-east-1b"
}
