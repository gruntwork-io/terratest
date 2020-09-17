output "map_of_objects" {
  value = {
    somebool  = true
    somefloat = 1.1
    one       = 1
    two       = "two"
    three     = "three"
    nest = {
      four = 4
      five = "five"
    }
    nest_list = [
      {
        six   = 6
        seven = "seven"
      },
    ]
  }
}

output "not_map_of_objects" {
  value = "Just a string"
}
