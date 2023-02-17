# An example rego policy of how to enforce that all module blocks in terraform json representation source the module
# from the gruntwork-io github repo on the json representation of the terraform source files. A module block in the json
# representation looks like the
# following:
#
# {
#   "module": {
#     "MODULE_LABEL": [{
#       #BLOCK_CONTENT
#     }]
#   }
# }
package enforce_source


# website::tag::1:: Only define the allow variable and set to true if the violation set is empty.
allow = true {
    count(violation) == 0
}

# website::tag::1:: Add modules with module_label to the violation set if the source attribute does not start with a string indicating it came from gruntwork-io GitHub org.
violation[module_label] {
    some module_label, i
    startswith(input.module[module_label][i].source, "git::git@github.com:gruntwork-io") == false
}
