package random_test

import (
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	t.Parallel()

	min := 0
	max := 100

	for i := 0; i < 100000; i++ {
		value := random.Random(min, max)
		assert.True(t, value >= min && value <= max)
	}
}

func TestRandomInt(t *testing.T) {
	t.Parallel()

	min := 0
	max := 1000

	list := []int{}
	for i := min; i < max; i++ {
		list = append(list, i)
	}

	for i := 0; i < 100000; i++ {
		value := random.RandomInt(list)
		assert.Contains(t, list, value)
	}
}

func TestRandomString(t *testing.T) {
	t.Parallel()

	min := 0
	max := 1000

	list := []string{}
	for i := min; i < max; i++ {
		list = append(list, strconv.Itoa(i))
	}

	for i := 0; i < 100000; i++ {
		value := random.RandomString(list)
		assert.Contains(t, list, value)
	}
}

func TestUniqueID(t *testing.T) {
	t.Parallel()

	previouslySeen := map[string]bool{}

	for i := 0; i < 100; i++ {
		uniqueID := random.UniqueID()
		assert.Len(t, uniqueID, 6)
		assert.NotContains(t, previouslySeen, uniqueID)

		previouslySeen[uniqueID] = true
	}
}
