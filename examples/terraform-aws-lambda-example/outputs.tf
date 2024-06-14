output "lambda_function_arn" {
  description = "ARN of function"
  value       = aws_lambda_function.lambda.arn
}

output "lambda_source_code_hash" {
  description = "The SHA256 hash of the function's deployment package"
  value       = base64decode(aws_lambda_function.lambda.source_code_hash)
}