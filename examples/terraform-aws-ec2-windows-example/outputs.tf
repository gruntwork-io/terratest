output "windows_instance_public_ip" {
  description = "The IPv4 address of the Windows instance. Enter this value into your RDP client when connecting to your instance."
  value       = aws_instance.instance.public_ip
}
