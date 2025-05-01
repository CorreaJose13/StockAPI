output "name" {
  description = "Name of lambda function"
  value       = aws_lambda_function.this.function_name
}

output "invoke_arn" {
  description = "Invoke ARN of lambda function"
  value       = aws_lambda_function.this.invoke_arn
}
