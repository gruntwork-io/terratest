resource "null_resource" "fail_on_first_run" {
  provisioner "local-exec" {
    command     = "if [[ -f terraform.tfstate.backup ]]; then echo 'This is not the first run, so exiting successfully' && exit 0; else echo 'This is the first run, exiting with an error' && exit 1; fi"
    interpreter = ["/bin/bash", "-c"]
  }
}