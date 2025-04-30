output "name" {
  description = "Name of lambda function"
  value       = aws_lambda_function.this.function_name
}
