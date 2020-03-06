package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpretDataByType(t *testing.T) {

	testcases := []struct {
		testname       string
		inputBytes     []byte
		inputType      goType
		expectedOutput string
	}{
		{
			testname:       "int8",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      INT8,
			expectedOutput: "256",
		},
		{
			testname:       "int16",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      INT16,
			expectedOutput: "256",
		},
		{
			testname:       "int32",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      INT32,
			expectedOutput: "256",
		},
		{
			testname:       "uint8",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      UINT8,
			expectedOutput: "256",
		},
		{
			testname:       "uint16",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      UINT16,
			expectedOutput: "256",
		},
		{
			testname:       "uint32",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      UINT32,
			expectedOutput: "256",
		},
		{
			testname:       "int",
			inputBytes:     []byte{5, 1, 0, 0, 1, 0, 5, 0},
			inputType:      INT,
			expectedOutput: "1407379178520837",
		},
		{
			testname:       "int64",
			inputBytes:     []byte{5, 1, 0, 0, 1, 0, 5, 0},
			inputType:      INT64,
			expectedOutput: "1407379178520837",
		},
		{
			testname:       "uint",
			inputBytes:     []byte{5, 1, 0, 0, 1, 0, 5, 0},
			inputType:      UINT,
			expectedOutput: "1407379178520837",
		},
		{
			testname:       "uint64",
			inputBytes:     []byte{0, 1, 0, 0, 1, 1, 1, 1},
			inputType:      UINT64,
			expectedOutput: "72340172821233920",
		},
		{
			testname:       "float32",
			inputBytes:     []byte{0, 1, 0, 0},
			inputType:      FLOAT32,
			expectedOutput: "3.587324e-43",
		},
		{
			testname:       "float64",
			inputBytes:     []byte{0, 55, 1, 22, 1, 0, 1, 9},
			inputType:      FLOAT64,
			expectedOutput: "2.636108e-265",
		},
		{
			testname:       "bool true",
			inputBytes:     []byte{1, 0, 0, 0},
			inputType:      BOOL,
			expectedOutput: "true",
		},
		{
			testname:       "bool false",
			inputBytes:     []byte{0, 0, 0, 0},
			inputType:      BOOL,
			expectedOutput: "false",
		},
		{
			testname:       "byte",
			inputBytes:     []byte{1, 1, 1, 1},
			inputType:      BYTE,
			expectedOutput: "dec=16843009\tchar='�'",
		},
		{
			testname:       "string",
			inputBytes:     []byte{'h', 'e', 'l', 'l', 'o'},
			inputType:      STRING,
			expectedOutput: "'hello'",
		},
		{
			testname:       "rune",
			inputBytes:     []byte{'h', 0, 0, 0},
			inputType:      RUNE,
			expectedOutput: "h",
		},
		{
			testname:       "struct",
			inputBytes:     []byte{},
			inputType:      STRUCT,
			expectedOutput: "struct interpretation is not yet implemented",
		},
		{
			testname:       "pointer",
			inputBytes:     []byte{},
			inputType:      POINTER,
			expectedOutput: "pointer interpretation is not yet implemented",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.testname, func(t *testing.T) {
			output := interpretDataByType(testcase.inputBytes, testcase.inputType)
			assert.Equal(t, testcase.expectedOutput, output)
		})
	}
}
