from aws_cdk import (
    core,
    aws_lambda,
    aws_apigateway
)
import os


class InfraStack(core.Stack):

    def __init__(self, scope: core.Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        # Configure path to Dockerfile
        dockerfile = os.path.join(os.getcwd(), "../src")

        _lambda = aws_lambda.DockerImageFunction(
            self, "OlafHandler",
            code=aws_lambda.DockerImageCode.from_image_asset(dockerfile)
        )

        _apigw = aws_apigateway.LambdaRestApi(
            self, "OlafEndpoint",
            handler=_lambda,
        )
