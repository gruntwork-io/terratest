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
      even_more = [
        "six"
      ]
    },
    {
      one   = "one"
      two   = 2
      three = 3
      more = [{
        four = 4
        five = "five"
      }]
      even_more = [
        6
      ]
    }
  ]
}

output "not_list_of_maps" {
  value = "Just a string"
}
