output "cognito_user_pool_id" {
  value = aws_cognito_user_pool.smartGou.id
}

output "cognito_user_pool_arn" {
  value = aws_cognito_user_pool.smartGou.arn
}

output "cognito_user_pool_client_id" {
  value = aws_cognito_user_pool_client.smartGou.id
}


