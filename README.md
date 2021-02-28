# olaf-bot
Olaf is a bot for [Discord](https://discord.com/).

Infrastructure is serverless using [AWS API Gateway](https://aws.amazon.com/api-gateway/) and [Lambda](https://aws.amazon.com/lambda) functions. It is defined under the `/cdk` directory and uses Python [CDK](https://aws.amazon.com/api-gateway/).

You may find `/scripts/register_command.py` useful to create/get/delete [slash commands](https://discord.com/developers/docs/interactions/slash-commands) because Discord doesn't support a graphical way to do this yet.

The main code base is in `/src` using the [Go](https://golang.org/) programming language.

## Features
It's gonna have whatever features I personally want.

What's available:

`/image <keywords>`: Send a random photo related to `keywords`

