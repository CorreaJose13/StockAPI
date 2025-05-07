module "lambda_role" {
  source = "../../modules/iam/iam_role/"

  name     = var.lambda_role
  services = ["lambda.amazonaws.com"]
}

module "lambda_logs_policy" {
  source = "../../modules/iam/iam_policy_attachment/"

  name        = var.lambda_logs_policy
  description = "IAM policy for Lambda function to write logs to CloudWatch"
  action      = ["logs:CreateLogStream", "logs:PutLogEvents"]
  resource    = "arn:aws:logs:*:*:*"
  role_name   = module.lambda_role.name
}
