provider "aws" {
  region = "us-east-1"
}

# lambda definition
resource "aws_lambda_function" "ipquery" {
  function_name = "ipQuery"

  # I'm using S3 to store the compiled golang project
  s3_bucket = "terraform-serverless-votazu-ipquery"
  s3_key    = "v1.0.0/main.zip"

  # "main" is the filename within the zip file (main.zip)
  handler = "main"
  runtime = "go1.x"
  memory_size = 128
  role = "${aws_iam_role.lambda_exec.arn}"

  environment {
    variables = {
      extapikey = data.aws_secretsmanager_secret_version.secret_credentials.secret_string
    }
  }
}

# AWS services the Lambda function can access.
resource "aws_iam_role" "lambda_exec" {
  name = "goserverless_ipquery_lambda"

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

# Logs configuration
data "aws_iam_policy_document" "lambda_exec_role_policy" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }
}

resource "aws_iam_role_policy" "lambda_exec_role" {
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_exec_role_policy.json
}

# Gateway permission
resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ipquery.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.ipquery.execution_arn}/*/*"
}

# Secret manager config for external api key
data "aws_secretsmanager_secret" "secret_name" {
   name = "test/ipquery/extapikey"
}

data "aws_secretsmanager_secret_version" "secret_credentials" {
  secret_id = data.aws_secretsmanager_secret.secret_name.id
}