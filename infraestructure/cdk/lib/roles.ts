import * as iam from 'aws-cdk-lib/aws-iam';
import { Construct } from 'constructs';
export class Roles {
  public readonly role: iam.Role;

  constructor(scope: Construct, roleName: string, principal: string) {

    this.role = new iam.Role(scope, roleName, {
      roleName: roleName,
      assumedBy: new iam.ServicePrincipal(principal),
    });
    // Add managed policy for Lambda Basic Execution Role
    this.role.addManagedPolicy(iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'));
  }

  public addToPolicy(action: string, resource = '*') {
    this.role.addToPolicy(new iam.PolicyStatement({
      actions: [action],
      resources: [resource],
    }));
  }
}
