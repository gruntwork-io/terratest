output "list_of_maps" {
  value = [
    {
      one   = 1
      two   = "two"
      three = "three"
      more = {
        four = 4
        five = "five"
      }
    },
    {
      one   = "one"
      two   = 2
      three = 3
      more = [{
        four = 4
        five = "five"
      }]
    },
    {
      one   = "one"
      two   = 2
      three = 3
      more  = ["one", 2.0, 3.4, ["one", 2.0, 3.4], { "one" : 2.0, "three" : 3.4 }]
    }
  ]
}

output "not_list_of_maps" {
  value = "Just a string"
}
