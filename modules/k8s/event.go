package k8s

import (
	"context"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListEvents will retrieve the Events in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
func ListEvents(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Event {
	events, err := ListEventsE(t, options, filters)
	require.NoError(t, err)
	return events
}

// ListEventsE will retrieve the Events that match the given filters and return them.
func ListEventsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Event, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().Events(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
