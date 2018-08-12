provider "digitalocean" {
  token = "hello"
}

resource "digitalocean_droplet" "web" {
	count = 2

  image  = "ubuntu-14-04-x64"
  name   = "web-${count.index}"
  region = "nyc2"
  size   = "512mb"
}
