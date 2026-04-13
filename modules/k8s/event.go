package k8s

import (
	"context"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListEventsContextE retrieves the Events in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListEventsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Event, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().Events(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListEventsContext retrieves the Events in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListEventsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []corev1.Event {
	t.Helper()
	events, err := ListEventsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return events
}

// ListEvents will retrieve the Events in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
//
// Deprecated: Use [ListEventsContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListEvents(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Event {
	t.Helper()

	return ListEventsContext(t, context.Background(), options, filters)
}

// ListEventsE will retrieve the Events that match the given filters and return them.
//
// Deprecated: Use [ListEventsContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListEventsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Event, error) {
	return ListEventsContextE(t, context.Background(), options, filters)
}
