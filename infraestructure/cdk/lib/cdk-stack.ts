import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Construct } from 'constructs';
import { HttpApi, CorsHttpMethod, HttpMethod } from 'aws-cdk-lib/aws-apigatewayv2';
import { HttpLambdaIntegration } from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import { env } from './env';
import { Roles } from './roles';

export class SmartGouStack extends cdk.Stack {
  public readonly cognitoRole: Roles;
  // lambda functions for security module
  public readonly requestSingUpFunctionByPhone: lambda.Function;
  public readonly verifyOTPByPhoneFunction: lambda.Function;
  public readonly signUpByEmailFunction: lambda.Function;
  public readonly confirmOtpByEmailFunction: lambda.Function;
  public readonly loginWithEmailFunction: lambda.Function;
  public readonly refreshTokenFunction: lambda.Function;
  public readonly logoutFunction: lambda.Function;
  public readonly requestRecoveryPasswordFunction: lambda.Function;
  public readonly resetPasswordFunction: lambda.Function;
  // lambda functions for users module
  public readonly completeOnboardingFunction: lambda.Function;
  public readonly getUserProfileFunction: lambda.Function;
  // lambda functions for reference module
  public readonly getMetaCategoriesFunction: lambda.Function;
  public readonly getCategoriesFunction: lambda.Function;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const envVariables = env(['COGNITO_USER_POOL_ID', 'COGNITO_USER_POOL_CLIENT_ID', 'COGNITO_REGION', 'DATABASE_URL']);

    // Cognito Role
    this.cognitoRole = new Roles(this, 'CognitoRole', 'lambda.amazonaws.com');
    this.cognitoRole.addToPolicy("cognito-idp:*", "*");
    this.cognitoRole.addToPolicy("sns:*", "*");

    // Auth functions
    this.requestSingUpFunctionByPhone = new lambda.Function(this, 'requestSingUpByPhoneLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'requestSingUpByPhoneLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/request_singup_by_phone/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.verifyOTPByPhoneFunction = new lambda.Function(this, 'verifyOTPByPhoneLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'verifyOTPByPhoneLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/verify_otp_from_phone/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.signUpByEmailFunction = new lambda.Function(this, 'signUpByEmailLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'signUpByEmailLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/sign_up_by_email/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.confirmOtpByEmailFunction = new lambda.Function(this, 'confirmOtpByEmailLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'confirmOtpByEmailLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/confirm_otp_by_email/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.loginWithEmailFunction = new lambda.Function(this, 'loginWithEmailLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'loginWithEmailLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/login_with_email/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.refreshTokenFunction = new lambda.Function(this, 'refreshTokenLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'refreshTokenLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/refresh_token/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.logoutFunction = new lambda.Function(this, 'logoutLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'logoutLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/logout/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.requestRecoveryPasswordFunction = new lambda.Function(this, 'requestRecoveryPasswordLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'requestRecoveryPasswordLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/request_recovery_password/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.resetPasswordFunction = new lambda.Function(this, 'resetPasswordLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'resetPasswordLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/reset_password/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    // End Auth functions

    // Users functions
    this.completeOnboardingFunction = new lambda.Function(this, 'completeOnboardingLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'completeOnboardingLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/compleate_onbording/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.getUserProfileFunction = new lambda.Function(this, 'getUserProfileLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'getUserProfileLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/get_user_profile/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });
    // End Users functions

    // Reference functions
    this.getMetaCategoriesFunction = new lambda.Function(this, 'getMetaCategoriesLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'getMetaCategoriesLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/get_meta_categories/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });

    this.getCategoriesFunction = new lambda.Function(this, 'getCategoriesLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      functionName: 'getCategoriesLambda',
      memorySize: 1024,
      code: lambda.Code.fromAsset('../../bin/get_categories/function.zip'),
      handler: 'bootstrap',
      environment: envVariables,
      role: this.cognitoRole.role,
    });
    // End Reference functions

    const httpApi = new HttpApi(this, 'SmartGouApiRest', {
      apiName: 'SmartGou API REST',
      description: 'This API serves enpoints to be consume by SmartGou app.',
      corsPreflight: {
        allowOrigins: [CorsHttpMethod.ANY],
        allowMethods: [CorsHttpMethod.ANY],
        allowCredentials: false,
        allowHeaders: ['Content-Type', 'X-Amz-Date', 'Authorization', 'X-Api-Key'],
      }
    });



    // Auth routes
    httpApi.addRoutes({
      path: '/auth/otp',
      methods: [HttpMethod.POST],
      integration: new HttpLambdaIntegration('RequestSingUpByPhoneIntegration', this.requestSingUpFunctionByPhone),
    });

    httpApi.addRoutes({
      path: '/auth/otp',
      methods: [HttpMethod.PATCH],
      integration: new HttpLambdaIntegration('VerifyOTPByPhoneIntegration', this.verifyOTPByPhoneFunction),
    });

    httpApi.addRoutes({
      path: '/auth',
      methods: [HttpMethod.POST],
      integration: new HttpLambdaIntegration('SignUpByEmailIntegration', this.signUpByEmailFunction),
    });

    httpApi.addRoutes({
      path: '/auth',
      methods: [HttpMethod.PATCH],
      integration: new HttpLambdaIntegration('ConfirmOtpByEmailIntegration', this.confirmOtpByEmailFunction),
    });

    httpApi.addRoutes({
      path: '/auth/sessions',
      methods: [HttpMethod.POST],
      integration: new HttpLambdaIntegration('LoginWithEmailIntegration', this.loginWithEmailFunction),
    });

    httpApi.addRoutes({
      path: '/auth/sessions',
      methods: [HttpMethod.PATCH],
      integration: new HttpLambdaIntegration('RefreshTokenIntegration', this.refreshTokenFunction),
    });

    httpApi.addRoutes({
      path: '/auth/sessions',
      methods: [HttpMethod.DELETE],
      integration: new HttpLambdaIntegration('LogoutIntegration', this.logoutFunction),
    });

    httpApi.addRoutes({
      path: '/auth/recovery-password',
      methods: [HttpMethod.POST],
      integration: new HttpLambdaIntegration('RequestRecoveryPasswordIntegration', this.requestRecoveryPasswordFunction),
    });

    httpApi.addRoutes({
      path: '/auth/recovery-password',
      methods: [HttpMethod.PATCH],
      integration: new HttpLambdaIntegration('ResetPasswordIntegration', this.resetPasswordFunction),
    });
    // End Auth routes

    // Users routes
    httpApi.addRoutes({
      path: '/users/{id-user}/onboarding',
      methods: [HttpMethod.PATCH],
      integration: new HttpLambdaIntegration('CompleteOnboardingIntegration', this.completeOnboardingFunction),
    });

    httpApi.addRoutes({
      path: '/users/{id-user}',
      methods: [HttpMethod.GET],
      integration: new HttpLambdaIntegration('GetUserProfileIntegration', this.getUserProfileFunction),
    });
    // End Users routes

    // Reference routes
    httpApi.addRoutes({
      path: '/references/meta-categories',
      methods: [HttpMethod.GET],
      integration: new HttpLambdaIntegration('GetMetaCategoriesIntegration', this.getMetaCategoriesFunction),
    });

    httpApi.addRoutes({
      path: '/references/categories',
      methods: [HttpMethod.GET],
      integration: new HttpLambdaIntegration('GetCategoriesIntegration', this.getCategoriesFunction),
    });
    // End Reference routes


    new cdk.CfnOutput(this, 'HttpApiUrl', {
      value: httpApi.url ?? '',
      description: 'The URL of the API Gateway',
    });

  }

}
