resource "null_resource" my_null_resource {
  for_each = toset(["val1"])

  triggers = {
    some_attribute   = "attr-${each.value}"
  }
}
