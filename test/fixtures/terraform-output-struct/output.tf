output "object" {
  value = {
    somebool   = true
    somefloat  = 0.1
    someint    = 1
    somestring = "two"
    somemap = {
      three = 3
      four  = "four"
    },
    listmaps = [
      {
        five = 5
        six  = "six"
      },
    ]
    liststrings = [
      "seven",
      "eight",
      "nine",
    ]
  }
}

output "list_of_objects" {
  value = [
    {
      somebool   = true
      somefloat  = 0.1
      someint    = 1
      somestring = "two"
    },
    {
      somebool   = false
      somefloat  = 0.3
      someint    = 4
      somestring = "five"
    }
  ]
}
