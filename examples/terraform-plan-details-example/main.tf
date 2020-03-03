resource "null_resource" test {
  for_each = toset(["val1", "val2"])

  triggers = {
    some_attribute   = "attr-${each.value}"
  }
}
