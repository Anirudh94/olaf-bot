from aws_cdk import core

from lib.infra_stack import InfraStack


app = core.App()
InfraStack(app, "OlafBotInfra", env={'region': 'us-east-1'})

app.synth()
