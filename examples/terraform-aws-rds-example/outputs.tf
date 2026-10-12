output "db_instance_id" {
  # Since AWS provider v5, id is the resource ID (db-XXXX), not the identifier the RDS API looks instances up by.
  value = aws_db_instance.example.identifier
}

