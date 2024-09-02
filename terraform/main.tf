provider "aws" {
  region  = "eu-north-1"
  profile = "default"
}

# Variables
variable "route_name" {
  description = "The name of the route to create in API Gateway."
  type        = string
  default     = "parentheses"
}

# Create S3 bucket
resource "aws_s3_bucket" "lambda_bucket" {
  bucket        = "parentheses-lambda-bucket"
  force_destroy = true
}

# Upload Lambda function ZIP to S3 bucket
resource "aws_s3_object" "lambda_zip" {
  bucket = aws_s3_bucket.lambda_bucket.bucket
  key    = "main.zip"
  source = "main.zip"
  acl    = "private"
}

# Create IAM Role for Lambda
resource "aws_iam_role" "lambda_execution_role" {
  name = "parentheses_lambda_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

# Attach IAM policy to the Lambda role
resource "aws_iam_role_policy_attachment" "execution_role_policy_attachment" {
  role       = aws_iam_role.lambda_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Create Lambda function
resource "aws_lambda_function" "parentheses_function" {
  function_name = "parentheses-lambda"
  s3_bucket     = aws_s3_bucket.lambda_bucket.bucket
  s3_key        = aws_s3_object.lambda_zip.key
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_execution_role.arn
  runtime       = "provided.al2"
}

# Create an API Gateway V2
resource "aws_apigatewayv2_api" "http_api_gateway" {
  name          = "lambda-api-gateway"
  protocol_type = "HTTP"
}

# Create integration between API Gateway and Lambda
resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id           = aws_apigatewayv2_api.http_api_gateway.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.parentheses_function.invoke_arn
}

# Create a route in API Gateway
resource "aws_apigatewayv2_route" "parentheses_route" {
  api_id    = aws_apigatewayv2_api.http_api_gateway.id
  route_key = "GET /${var.route_name}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

# Allow API Gateway to invoke the Lambda function
resource "aws_lambda_permission" "allow_api_gateway_invoke" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.parentheses_function.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api_gateway.execution_arn}/*/*"
}

# Create API Gateway stage with a descriptive name
resource "aws_apigatewayv2_stage" "production_stage" {
  api_id      = aws_apigatewayv2_api.http_api_gateway.id
  name        = "v1"
  auto_deploy = true
}

# Output the API Gateway endpoint URL
output "api_gateway_url" {
  description = "The URL of the API Gateway for the Lambda function."
  value       = "${aws_apigatewayv2_api.http_api_gateway.api_endpoint}/${aws_apigatewayv2_stage.production_stage.name}/${var.route_name}"
}
