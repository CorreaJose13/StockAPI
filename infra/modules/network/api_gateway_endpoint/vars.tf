variable "api_gateway_name" {
  description = "Name of the api gateway"
  type        = string
}

variable "path" {
  description = "Path of the api gateway resource"
  type        = string
}

variable "method" {
  description = "HTTP method"
  type        = string
}

variable "stage" {
  description = "Name of the stage"
  type        = string
}

variable "lambda_name" {
  description = "Name of the lambda function"
  type        = string
}
