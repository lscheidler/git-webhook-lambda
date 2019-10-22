resource "aws_sqs_queue" "git-webhook" {
  name                      = "git-webhook"
  delay_seconds             = 90
  max_message_size          = 16 * 1024
  message_retention_seconds = 60 * 60 * 12
  receive_wait_time_seconds = 20
}
