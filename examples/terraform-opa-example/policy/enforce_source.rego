# Enforces that all module blocks source the module from the gruntwork-io github repo on the json representation of the
# terraform source files. A module block in the json representation looks like the following:
#
# {
#   "modules": {
#     "MODULE_LABEL": [{
#       #BLOCK_CONTENT
#     }]
#   }
# }
package enforce_source

allow = true {
    count(violation) == 0
}

violation[module_label] {
    some module_label, i
    startswith(input.module[module_label][i].source, "git::git@github.com:gruntwork-io") == false
}
