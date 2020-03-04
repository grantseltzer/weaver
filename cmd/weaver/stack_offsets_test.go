package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineStackOffsets(t *testing.T) {

	testCases := []struct {
		testName        string
		inputContext    functionTraceContext
		expectedContext functionTraceContext
		expectedError   error
	}{
		{
			testName: "one int",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      INT,
						ArrayLength: 0,
					},
				},
			},
			expectedContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:         INT,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 8,
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.testName, func(t *testing.T) {
			err := determineStackOffsets(&testcase.inputContext)
			assert.Equal(t, testcase.expectedError, err)
			assert.Equal(t, testcase.expectedContext, testcase.inputContext)
		})
	}
}
