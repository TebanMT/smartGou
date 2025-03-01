resource "aws_cognito_user_pool" "smartGo" {
  name = var.cognito_user_pool_name

  # authentication methods
  auto_verified_attributes = var.auto_verified_attributes
  alias_attributes = var.alias_attributes
  
  # password policy
  password_policy {
    minimum_length     = var.password_policy.min_length
    require_numbers    = var.password_policy.require_digits
    require_lowercase  = var.password_policy.require_lowercase
    require_uppercase  = var.password_policy.require_uppercase
    require_symbols    = var.password_policy.require_symbols
  }

  # mfa configuration
  mfa_configuration = var.mfa_configuration.type
  
  software_token_mfa_configuration {
    enabled = var.mfa_configuration.software_token_mfa
  }

  
}

resource "aws_cognito_user_pool_client" "smartGo" {
  name = var.cognito_user_pool_client_name
  user_pool_id = aws_cognito_user_pool.smartGo.id

  # authentication methods
  explicit_auth_flows = var.explicit_auth_flows

  # generate secret
  generate_secret = var.generate_secret
  
  
}

