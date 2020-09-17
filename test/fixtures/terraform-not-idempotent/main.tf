resource "null_resource" test {
  triggers = {
    time = timestamp()
  }
}
