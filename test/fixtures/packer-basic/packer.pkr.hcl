source "file" "hello" {
  content = "hello world"
  target  = "test.txt"
}

build {
  name = "file_hello"
  sources = [
    "source.file.hello",
  ]
}
