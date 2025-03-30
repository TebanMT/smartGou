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
  default     = ["email", "phone_number"]
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
  default     = ["ALLOW_USER_SRP_AUTH", "ALLOW_USER_PASSWORD_AUTH", "ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_CUSTOM_AUTH"]
}

variable "custom_auth_lambda_functions" {
  description = "Custom auth lambda functions."
  default     = {
    define_auth_challenge = "define_auth_challenge"
    create_auth_challenge = "create_auth_challenge"
    verify_auth_challenge_response = "verify_auth_challenge_response"
  }
}

variable "generate_secret" {
  description = "Generate secret."
  default     = false
}

variable "sms_configuration" {
  description = "SMS configuration."
  default     = {
    external_id = "1234567890"
    sns_region = "us-east-1"
    message = "Hello from SmartGou, your verification code is {####}"
  }
}

variable "email_configuration" {
  description = "Email configuration."
  default     = {
    default_email_option = "CONFIRM_WITH_CODE"
    message = "Hello from SmartGou, your verification code is {####}"
    subject = "SmartGou Verification Code"
  }

}

