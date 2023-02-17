packer {
  required_plugins {
    docker = {
      version = ">=v1.0.1"
      source  = "github.com/hashicorp/docker"
    }
  }
}

source "docker" "ubuntu-docker" {
  changes = ["ENTRYPOINT [\"\"]"]
  commit  = true
  image   = "gruntwork/ubuntu-test:16.04"
}

build {
  sources = ["source.docker.ubuntu-docker"]

  provisioner "shell" {
    inline = ["echo 'Hello, World!' > /test.txt"]
  }

  post-processor "docker-tag" {
    repository = "gruntwork/packer-hello-world-example"
    tag        = ["latest"]
  }
}
