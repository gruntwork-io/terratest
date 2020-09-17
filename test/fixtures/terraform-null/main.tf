variable "foo" {
  type = object({
    nullable_string    = string
    nonnullable_string = string
  })
}

output "foo" {
  value = var.foo
}

output "bar" {
  value = var.foo.nullable_string == null ? "I AM NULL" : null
}
