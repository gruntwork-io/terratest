variable cnt {}

resource "null_resource" test {
  count = var.cnt
}
