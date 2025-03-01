output "rds_instance_address" {
  value = module.rds.db_instance_address
}

output "rds_instance_port" {
  value = module.rds.db_instance_port
}

output "rds_endpoint" {
  value = module.rds.db_endpoint
}

output "cognito_user_pool_id" {
  value = module.cognito.cognito_user_pool_id
}

output "cognito_user_pool_arn" {
  value = module.cognito.cognito_user_pool_arn
}

output "cognito_user_pool_client_id" {
  value = module.cognito.cognito_user_pool_client_id
}


