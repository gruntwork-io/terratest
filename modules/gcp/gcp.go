// Package gcp allows interaction with Google Cloud Platform resources.
package gcp

import (
	"google.golang.org/api/option"
)

func withOptions() (opts []option.ClientOption) {
	v, ok := getStaticTokenSource()
	if ok {
		opts = append(opts, option.WithTokenSource(v))
	}

	return
}

// iamPolicyVersionWithConditions is the policy version that carries conditional bindings. Google
// returns a policy at version 1 unless asked otherwise, and a version 1 answer has no room for a
// condition, so every read in this package that fetches an IAM policy asks for this one.
const iamPolicyVersionWithConditions = 3
