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
    }
  ]
}

output "not_list_of_maps" {
  value = "Just a string"
}
