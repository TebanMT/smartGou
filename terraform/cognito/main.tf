# Role for cognito
resource "aws_iam_role" "cognito_role" {
  name = "cognito-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = ["cognito-idp.amazonaws.com", "lambda.amazonaws.com"]
        }
      },
    ]
  })
}

# Policy for cognito
resource "aws_iam_policy" "cognito_policy" {
  name = "cognito-policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "cognito-idp:*",
          "cognito-identity:*",
          "cognito-sync:*",
          "cognito-auth:*",
          "cognito-identity:*",
          "cognito-auth:*",
          "cognito-sync:*",
          "cognito-idp:*",
          "cognito-identity:*",
          "cognito-sync:*",
          "cognito-auth:*",
          "sns:*",
          "s3:*",
          "lambda:*",
          "apigateway:*",
          "cloudwatch:*",
          "logs:*",
        ]
        Effect = "Allow"
        Resource = "*"
      }
    ]
  })
}

# Attach policy to role
resource "aws_iam_role_policy_attachment" "cognito_policy_attachment" {
  role = aws_iam_role.cognito_role.name
  policy_arn = aws_iam_policy.cognito_policy.arn
}

resource "aws_lambda_function" "custom_auth_lambda_functions" {
  for_each = var.custom_auth_lambda_functions
  function_name = each.value
  role = aws_iam_role.cognito_role.arn
  handler = "bootstrap"
  runtime = "provided.al2"
  filename = "${path.module}/../../bin/custom_auth/${each.value}/function.zip"
  source_code_hash = filebase64sha256("${path.module}/../../bin/custom_auth/${each.value}/function.zip")
}

resource "aws_cognito_user_pool" "smartGou" {
  name = var.cognito_user_pool_name

  
  schema {
    name = "external_id"
    attribute_data_type = "String"
    mutable = true
    required = false
  }

  schema {
    name = "username"
    attribute_data_type = "String"
    mutable = true
    required = false
  }

  lifecycle {
    ignore_changes = [
      schema
    ]
  }

  lambda_config {
    define_auth_challenge = aws_lambda_function.custom_auth_lambda_functions["define_auth_challenge"].arn
    create_auth_challenge = aws_lambda_function.custom_auth_lambda_functions["create_auth_challenge"].arn
    verify_auth_challenge_response = aws_lambda_function.custom_auth_lambda_functions["verify_auth_challenge_response"].arn
    custom_message = aws_lambda_function.custom_auth_lambda_functions["custom_email_messages"].arn
  }

  sms_configuration {
    external_id = var.sms_configuration.external_id
    sns_caller_arn = aws_iam_role.cognito_role.arn
    sns_region = var.sms_configuration.sns_region
  }

  verification_message_template {
    sms_message = var.sms_configuration.message
    default_email_option = var.email_configuration.default_email_option
    email_message = var.email_configuration.message
    email_subject = var.email_configuration.subject
  }

  // account recovery setting
  account_recovery_setting {
    recovery_mechanism {
      name = "verified_email"
      priority = 1
    }
    recovery_mechanism {
      name = "verified_phone_number"
      priority = 2
    }
    

  }

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

resource "aws_lambda_permission" "custom_auth_lambda_permissions" {
  statement_id = "AllowExecutionFromCognito"
  action = "lambda:InvokeFunction"
  for_each = var.custom_auth_lambda_functions
  function_name = aws_lambda_function.custom_auth_lambda_functions[each.value].arn
  principal = "cognito-idp.amazonaws.com"
  source_arn = aws_cognito_user_pool.smartGou.arn
}


resource "aws_cognito_user_pool_client" "smartGou" {
  name = var.cognito_user_pool_client_name
  user_pool_id = aws_cognito_user_pool.smartGou.id

  # authentication methods
  explicit_auth_flows = var.explicit_auth_flows

  # generate secret
  generate_secret = var.generate_secret
  
  
}

