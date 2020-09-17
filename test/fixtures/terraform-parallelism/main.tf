# This resource just waits for 5 seconds. If we run it with enough parallelism, the whole module should apply in about
# 5 seconds. If we set parallelism to 1, it should take at least 25 seconds.
resource "null_resource" "wait" {
  count = 5

  triggers = {
    run_always = timestamp()
  }

  provisioner "local-exec" {
    command = "sleep 5"
  }
}