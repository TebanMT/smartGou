provider "aws" {
  region = "us-east-1"
  profile = "duvantis"
}

module "vpc" {
  source = "./vpc"
}

module "security_groups" {
  source = "./security_groups"
  vpc_id = module.vpc.vpc_id
  cidr_blocks = ["0.0.0.0/0"] # TODO: change to the IP of the client
}

module "cognito" {
  source = "./cognito"
}

module "rds" {
  source = "./rds"
  db_username = var.db_username
  db_password = var.db_password
  subnet_id_main = module.vpc.subnet_id_main
  subnet_id_backup = module.vpc.subnet_id_backup
  rds_security_group_id = module.security_groups.rds_security_group_id
}
