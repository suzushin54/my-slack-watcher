import * as cdk from '@aws-cdk/core';
import { Function, Runtime, Code } from "@aws-cdk/aws-lambda"
import { RestApi, Integration, LambdaIntegration, Resource } from "@aws-cdk/aws-apigateway"
import * as iam from "@aws-cdk/aws-iam"

export class MySlackWatcherStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create Lambda Function
    const lambdaFunction: Function = new Function(this, "MySlackWatcher", {
      functionName: "my-slack-watcher",
      runtime: Runtime.GO_1_X,
      code: Code.asset("./lambdaSource"),
      handler: "main",
      memorySize: 256,
      timeout: cdk.Duration.seconds(10),
      environment: {
        "BOT_TOKEN": "xo...",
        "CHANNEL_ID": "XXXXXXXXX",
        "SIGNING_SECRETS": "aaa...",
      }
    })

    // Add policy to function
    lambdaFunction.addToRolePolicy(new iam.PolicyStatement({
      resources: ["*"],
      actions: ["ec2:DescribeInstances"],
    }))

    // Create API Gateway
    const restApi: RestApi = new RestApi(this, "my-slack-watcher", {
      restApiName: "My-Slack-Watcher",
      description: "Deployed by AWS CDK",
    })

    // Create Integration
    const integration: Integration = new LambdaIntegration(lambdaFunction)

    // Create Resource
    const getResource: Resource = restApi.root.addResource("event")

    // Create Method
    getResource.addMethod("POST", integration)
  }
}
