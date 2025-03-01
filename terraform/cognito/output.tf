output "cognito_user_pool_id" {
  value = aws_cognito_user_pool.smartGo.id
}

output "cognito_user_pool_arn" {
  value = aws_cognito_user_pool.smartGo.arn
}

output "cognito_user_pool_client_id" {
  value = aws_cognito_user_pool_client.smartGo.id
}


