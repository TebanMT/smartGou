variable "cognito_user_pool_name" {
  description = "Cognito user pool name"
  default     = "smartGo-user-pool"
}

variable "cognito_user_pool_client_name" {
  description = "Cognito user pool client name"
  default     = "smartGo-user-pool-client"
}

variable "alias_attributes" {
  description = "Alias attributes. The users can use this attribute to sign in to the application."
  default     = ["phone_number", "email"]
}

variable "auto_verified_attributes" {
  description = "Auto verified attributes. The users must verify their email and phone number to sign in to the application."
  default     = ["email"]
}

variable "password_policy" {
  description = "Password policy. The users must follow this policy to sign in to the application."
  default     = {
    min_length    = 8
    require_digits = true
    require_lowercase = true
    require_uppercase = true
    require_symbols = true
  }
}

variable "mfa_configuration" {
  description = "MFA configuration."
  default     = {
    type = "OPTIONAL"
    software_token_mfa = true
  }
}

variable "explicit_auth_flows" {
  description = "Explicit auth flows."
  default     = ["ALLOW_USER_SRP_AUTH", "ALLOW_USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"]
}

variable "generate_secret" {
  description = "Generate secret."
  default     = false
}


