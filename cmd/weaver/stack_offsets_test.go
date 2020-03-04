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
		{
			testName: "array of ints",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      INT,
						ArrayLength: 5,
					},
				},
			},
			expectedContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:         INT,
						ArrayLength:    5,
						TypeSize:       8,
						StartingOffset: 8,
					},
				},
			},
			expectedError: nil,
		},
		{
			testName: "multiple ints",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      INT,
						ArrayLength: 0,
					},
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
					{
						goType:         INT,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 16,
					},
				},
			},

			expectedError: nil,
		},
		{
			testName: "multiple variables of different type 1",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      INT,
						ArrayLength: 0,
					},
					{
						goType:      INT8,
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
					{
						goType:         INT8,
						ArrayLength:    0,
						TypeSize:       1,
						StartingOffset: 16,
					},
				},
			},
		},
		{
			testName: "multiple variables of different type 2",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      INT8,
						ArrayLength: 0,
					},
					{
						goType:      INT64,
						ArrayLength: 0,
					},
				},
			},
			expectedContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:         INT8,
						ArrayLength:    0,
						TypeSize:       1,
						StartingOffset: 8,
					},
					{
						goType:         INT64,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 16,
					},
				},
			},
		},
		{
			testName: "multiple variables of different type 3",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      RUNE,
						ArrayLength: 0,
					},
					{
						goType:      FLOAT32,
						ArrayLength: 0,
					},
				},
			},
			expectedContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:         RUNE,
						ArrayLength:    0,
						TypeSize:       4,
						StartingOffset: 8,
					},
					{
						goType:         FLOAT32,
						ArrayLength:    0,
						TypeSize:       4,
						StartingOffset: 12,
					},
				},
			},
		},
		{
			testName: "multiple variables of different type 4",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      INT,
						ArrayLength: 55,
					},
					{
						goType:      FLOAT64,
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
						ArrayLength:    55,
						TypeSize:       8,
						StartingOffset: 8,
					},
					{
						goType:         FLOAT64,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 448,
					},
				},
			},
		},
		{
			testName: "string argument",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      STRING,
						ArrayLength: 0,
					},
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
						goType:         STRING,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 8,
					},
					{
						goType:         INT,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 24,
					},
				},
			},
		},
		{
			testName: "string array argument",
			inputContext: functionTraceContext{
				binaryName:   "foobar",
				HasArguments: true,
				Arguments: []argument{
					{
						goType:      STRING,
						ArrayLength: 2,
					},
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
						goType:         STRING,
						ArrayLength:    2,
						TypeSize:       8,
						StartingOffset: 8,
					},
					{
						goType:         INT,
						ArrayLength:    0,
						TypeSize:       8,
						StartingOffset: 40,
					},
				},
			},
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
