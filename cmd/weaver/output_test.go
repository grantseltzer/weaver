package main

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrintOutput(t *testing.T) {

	testCases := []struct {
		testName       string
		output         output
		expectedError  error
		expectedOutput string
	}{
		{
			testName:       "everything is fine",
			expectedError:  nil,
			expectedOutput: `{"functionName":"main.testFunction","args":[{"type":"INT","value":"420"}],"procInfo":{}}`,
			output: output{
				FunctionName: "main.testFunction",
				Args: []outputArg{
					{
						Type:  "INT",
						Value: "420",
					},
				},
			},
		},
		{
			testName:       "empty args",
			expectedError:  nil,
			expectedOutput: `{"functionName":"main.testFunction","procInfo":{}}`,
			output: output{
				FunctionName: "main.testFunction",
				Args:         nil,
			},
		},
		{
			testName:       "empty output",
			expectedError:  nil,
			expectedOutput: `{"functionName":"","procInfo":{}}`,
			output:         output{},
		},
		{
			testName:       "erroneous functionname?",
			expectedError:  nil,
			expectedOutput: `{"functionName":"\u0000","procInfo":{}}`,
			output: output{
				FunctionName: string(0x0),
			},
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.testName, func(t *testing.T) {

			read, write, err := os.Pipe()
			if err != nil {
				t.Error(err)
			}
			globalOutput = write

			var buf bytes.Buffer
			go func() {
				err = printOutput(testcase.output)
				assert.Equal(t, testcase.expectedError, err)
				io.Copy(&buf, read)
			}()

			time.Sleep(time.Millisecond * 50)
			assert.Equal(t, testcase.expectedOutput, buf.String())

		})
	}
}
