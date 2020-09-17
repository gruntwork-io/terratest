output "instance_id" {
  value = aws_instance.example.id
}

output "public_ip" {
  value = aws_instance.example.public_ip
}

output "instance_url" {
  value = "http://${aws_instance.example.public_ip}:${var.instance_port}"
}

