import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { Construct } from 'constructs';
import { HttpApi, CorsHttpMethod, HttpMethod } from 'aws-cdk-lib/aws-apigatewayv2';
import { HttpLambdaIntegration } from 'aws-cdk-lib/aws-apigatewayv2-integrations';

export class SmartGouStack extends cdk.Stack {

  public readonly createUserFunction: lambda.Function;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);


    this.createUserFunction = new lambda.Function(this, 'createUserLambda', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromAsset('../../bin/create_user/bootstrap.zip'),
      handler: 'bootstrap',
    });

    const httpApi = new HttpApi(this, 'SmartGouApiRest', {
      apiName: 'SmartGou API REST',
      description: 'This API serves enpoints to be consume by SmartGou app.',
      corsPreflight: {
        allowOrigins: [CorsHttpMethod.ANY],
        allowMethods: [CorsHttpMethod.ANY],
        allowCredentials: true,
        allowHeaders: ['Content-Type', 'X-Amz-Date', 'Authorization', 'X-Api-Key'],
      }
    });
    
    httpApi.addRoutes({
      path: '/users',
      methods: [HttpMethod.POST],
      integration: new HttpLambdaIntegration('CreateUserIntegration', this.createUserFunction),
    });

    new cdk.CfnOutput(this, 'HttpApiUrl', {
      value: httpApi.url ?? '',
      description: 'The URL of the API Gateway',
    });

  }

}
