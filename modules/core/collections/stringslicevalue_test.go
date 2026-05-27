package collections_test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSliceLastValue(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		testName       string
		sliceSource    string
		sliceSeperator string
		expectedReturn string
		expectedError  bool
	}{
		{testName: "longSlice", sliceSource: "this/is/a/long/slash/separated/string/success", sliceSeperator: "/", expectedReturn: "success", expectedError: false},
		{testName: "shortendSlice", sliceSource: "this/is/a/long/slash/separated", sliceSeperator: "/", expectedReturn: "separated", expectedError: false},
		{testName: "dashSlice", sliceSource: "this-is-a-long-dash-separated-string-success", sliceSeperator: "-", expectedReturn: "success", expectedError: false},
		{testName: "seperatorNotPresent", sliceSource: "this-is-a-long-dash-separated-string-success", sliceSeperator: "/", expectedReturn: "", expectedError: true},
		{testName: "sourceNoSeperator", sliceSource: "noslicepresent", sliceSeperator: "/", expectedReturn: "", expectedError: true},
		{testName: "emptyStrings", sliceSource: "", sliceSeperator: "", expectedReturn: "", expectedError: true},
	}

	for _, tc := range testCases {
		testFor := tc // necessary range capture

		t.Run(testFor.testName, func(t *testing.T) {
			t.Parallel()

			actualReturn, err := collections.GetSliceLastValueE(testFor.sliceSource, testFor.sliceSeperator)
			switch testFor.expectedError {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}

			assert.Equal(t, testFor.expectedReturn, actualReturn)
		})
	}
}

func TestGetSliceIndexValue(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		expectedReturn string
		sliceIndex     int
		expectedError  bool
	}{
		{expectedReturn: "", sliceIndex: -1, expectedError: true},
		{expectedReturn: "this", sliceIndex: 0, expectedError: false},
		{expectedReturn: "slash", sliceIndex: 4, expectedError: false},
		{expectedReturn: "success", sliceIndex: 7, expectedError: false},
		{expectedReturn: "", sliceIndex: 10, expectedError: true},
	}

	sliceSource := "this/is/a/long/slash/separated/string/success"
	sliceSeperator := "/"

	for _, tc := range testCases {
		testFor := tc // necessary range capture

		t.Run(fmt.Sprintf("Index_%v", testFor.sliceIndex), func(t *testing.T) {
			t.Parallel()

			actualReturn, err := collections.GetSliceIndexValueE(sliceSource, sliceSeperator, testFor.sliceIndex)
			switch testFor.expectedError {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}

			assert.Equal(t, testFor.expectedReturn, actualReturn)
		})
	}
}
