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
      Environment = "staging"
      rules = jsonencode({
        "bitbucket" : {
          "repo:push" : [
            { "type" : "sqs", "method" : "SendMessage", "queueUrl" : aws_sqs_queue.git-webhook.id },
            #{"type": "ec2", "method": "StartInstances", "instanceIds": var.start_instances},
            #{"type": "rest", "method": "POST", "endpoint": "<rest-endpoint>", "data": {}},
            {"type": "lambda", "function": var.lambda_function, "region": var.lambda_function_region, "qualifier": "dev", "data":{"httpMethod": "POST", "path": "/", "queryStringParameters": {}, "headers": {"accept": "application/json"}, "body": jsonencode({"ids":var.start_instances, "type":"state","newState":"running"})}},
          ]
        }
        }
      )
    }
  }
}

resource "aws_lambda_alias" "git-webhook-lambda-dev" {
  name             = "dev"
  description      = "git-webhook-lambda dev"
  function_name    = aws_lambda_function.git-webhook-lambda.arn
  function_version = aws_lambda_function.git-webhook-lambda.version
}

#resource "aws_lambda_permission" "git-webhook-lambda_allow_cloudwatch" {
#  statement_id  = "AlarmDowntimeAllowExecutionFromCloudWatch"
#  action        = "lambda:InvokeFunction"
#  function_name = aws_lambda_function.git-webhook-lambda.function_name
#  principal     = "events.amazonaws.com"
#  source_arn    = aws_cloudwatch_event_rule.git-webhook-lambda.arn
#}
#
