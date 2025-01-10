terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }

  required_version = ">= 1.3.0"
}

provider "docker" {
}

resource "docker_network" "postgres_network" {
  name = "postgres_network"
}

resource "docker_volume" "postgres_volume" {
  name = "postgres_data"
}

resource "docker_container" "postgres" {
  name  = "postgres"
  image = "postgres:15"

  env = [
    "POSTGRES_USER=${var.username}",
    "POSTGRES_PASSWORD=${var.password}",
    "POSTGRES_DB=${var.database_name}"
  ]

  ports {
    internal = 5432
    external = var.port
  }

  networks_advanced {
    name = docker_network.postgres_network.name
  }

  restart = "always"
}
