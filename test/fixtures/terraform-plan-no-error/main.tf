resource "null_resource" "test" {
	count = 2

  triggers = {
    abc = "def"
  }
}
