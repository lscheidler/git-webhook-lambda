git-webhook-lambda
==================

This is an AWS lambda function, which processes git webhooks, which are send from
a git hosting provider through AWS alb and executes configurable hooks.

Current supported git hosting provider
--------------------------------------

| Provider      | Plugin                                                                                      |
|---------------|---------------------------------------------------------------------------------------------|
| bitbucket.org | [git-webhook-plugin-bitbucket](https://github.com/lscheidler/git-webhook-plugin-bitbucket)  |

Hooks
-----

| Hook    | Description                       |
|---------|-----------------------------------|
| ec2     | start/stop ec2 instances          |
| lambda  | invokes lambda function           |
| rest    | invoke rest endpoint              |
| sqs     | send webhook data to a sqs queue  |


Environment variables
=====================

rules
-----

```
# terraform example
resource "aws_lambda_function" "git-webhook-lambda" {
  filename         = "git-webhook-lambda.zip"
  function_name    = "git-webhook-lambda"
  role             = aws_iam_role.git-webhook-lambda_role.arn
  handler          = "git-webhook-lambda"
  source_code_hash = filebase64sha256("git-webhook-lambda.zip")
  runtime          = "go1.x"
  publish          = true

  environment {
    variables = {
      rules = jsonencode({
        "bitbucket" : {
          "repo:push" : [
            {"type" : "sqs", "method" : "SendMessage", "queueUrl" : "https://sqs.<region>.amazonaws.com/<account-id>/<queue>"},
            {"type": "ec2", "method": "StartInstances", "instanceIds": ["<instance-id>"]},
            {"type": "rest", "method": "POST", "endpoint": "<rest-endpoint>", "data": {}},
            {"type": "lambda", "function": "<lambda-function-name>", "region": "<lambda-function-region>", "qualifier": "<lambda-function-qualifier>", "data":{"httpMethod": "POST", "path": "/", "queryStringParameters": {}, "headers": {"accept": "application/json"}, "body": "<body>"}}
          ]
        }
      })
    }
  }
}
```
