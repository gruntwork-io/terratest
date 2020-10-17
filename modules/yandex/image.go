package yandex

import (
	"context"
	"github.com/gruntwork-io/terratest/modules/logger"

	compute "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	"testing"
)

func DeleteImage(t *testing.T, imageID string) {
	err := DeleteImageE(t, imageID)
	if err != nil {
		t.Fatal(err)
	}
}

func DeleteImageE(t *testing.T, imageID string) error {
	logger.Logf(t, "Delete Image %s", imageID)

	client, err := NewYandexClientE(t)
	if err != nil {
		return err
	}

	ctx := context.TODO()
	op, err := client.WrapOperation(client.Compute().Image().Delete(ctx, &compute.DeleteImageRequest{
		ImageId: imageID,
	}))

	return op.Wait(ctx)
}
