import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Construct } from 'constructs';
import { HttpApi, CorsHttpMethod, HttpMethod } from 'aws-cdk-lib/aws-apigatewayv2';
import { HttpLambdaIntegration } from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import { env } from './env';
import { Roles } from './roles';

export class SmartGouStack extends cdk.Stack {
  public readonly cognitoRole: Roles;
  // lambda functions
  public readonly requestSingUpFunctionByPhone: lambda.Function;
  public readonly verifyOTPByPhoneFunction: lambda.Function;
  public readonly signUpByEmailFunction: lambda.Function;
  public readonly confirmOtpByEmailFunction: lambda.Function;
  public readonly loginWithEmailFunction: lambda.Function;


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

    // End Auth functions

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
    // End Auth routes


    new cdk.CfnOutput(this, 'HttpApiUrl', {
      value: httpApi.url ?? '',
      description: 'The URL of the API Gateway',
    });

  }

}



/*
Método REST	Endpoint	Método de la Interfaz
POST /auth	Iniciar registro con phone	RequestSignUp
POST /auth/otp	Verificar OTP de registro	VerifySignUpOTP
PATCH /auth	Completar registro con email/password	CompleteSignUp
POST /auth/sessions	Login con email o phone	LoginWithEmail, RequestPhoneLogin
POST /auth/sessions/otp	Verificar OTP para login	VerifyPhoneLoginOTP
PATCH /auth/sessions	Refresh token	RefreshTokens
DELETE /auth/sessions	Logout	Logout
POST /auth/recovery	Solicitar recuperación de contraseña	RequestPasswordRecovery
PATCH /auth/recovery	Resetear contraseña	ResetPassword
*/