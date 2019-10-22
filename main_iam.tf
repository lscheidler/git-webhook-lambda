resource "aws_iam_policy" "git-webhook-lambda_policy" {
  name        = "git-webhook-lambda_policy"
  path        = "/"
  description = "git-webhook policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "sqs:SendMessage"
      ],
      "Resource": "${aws_sqs_queue.git-webhook.arn}"
    }
  ]
}
EOF

}

resource "aws_iam_role" "git-webhook-lambda_role" {
  name = "git-webhook-lambda_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

}

resource "aws_iam_role_policy_attachment" "attach" {
  role       = aws_iam_role.git-webhook-lambda_role.name
  policy_arn = aws_iam_policy.git-webhook-lambda_policy.arn
}

