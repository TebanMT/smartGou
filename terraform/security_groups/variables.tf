variable "vpc_id" {
  description = "VPC ID"
}

variable "cidr_blocks" {
  description = "CIDR blocks"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}
