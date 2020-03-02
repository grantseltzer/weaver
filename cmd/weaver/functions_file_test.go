package main

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFunctionsFile(t *testing.T) {

	testCases := []struct {
		testName         string
		fileContent      string
		expectedContexts []functionTraceContext
		expectedError    error
	}{
		{
			testName:         "multiple functions with multiple parameters",
			fileContent:      contentA,
			expectedContexts: expectedValueA,
			expectedError:    nil,
		},
		{
			testName:         "multiple functions with multiple parameters, valid whitespace",
			fileContent:      contentAA,
			expectedContexts: expectedValueAA,
			expectedError:    nil,
		},
		{
			testName:         "multiple functions with multiple parameters, and duplicate",
			fileContent:      contentDuplicate,
			expectedContexts: expectedValueDuplicate,
			expectedError:    nil,
		},
		{
			testName:         "single array argument",
			fileContent:      contentB,
			expectedContexts: expectedValueB,
			expectedError:    nil,
		},
		{
			testName:         "multiple array arguments",
			fileContent:      contentC,
			expectedContexts: expectedValueC,
			expectedError:    nil,
		},
		{
			testName:         "empty line in middle of file",
			fileContent:      contentD,
			expectedContexts: expectedValueD,
			expectedError:    nil,
		},
		{
			testName:         "empty line at end of file",
			fileContent:      contentE,
			expectedContexts: expectedValueE,
			expectedError:    nil,
		},
		{
			testName:         "empty line at begining of file",
			fileContent:      contentF,
			expectedContexts: expectedValueF,
			expectedError:    nil,
		},
		{
			testName:         "multiple empty lines",
			fileContent:      contentG,
			expectedContexts: expectedValueG,
			expectedError:    nil,
		},
		{
			testName:         "one bad function name",
			fileContent:      errorValueA,
			expectedContexts: expectedValueErrorA,
			expectedError:    expectedErrorA,
		},
		{
			testName:         "one good, one bad function name",
			fileContent:      errorValueB,
			expectedContexts: expectedValueErrorB,
			expectedError:    expectedErrorB,
		},
		{
			testName:         "invalid newline in middle of function name",
			fileContent:      errorValueD,
			expectedContexts: expectedValueErrorD,
			expectedError:    expectedErrorD,
		},
		// {
		// 	testName:         "",
		// 	fileContent:      errorValueE,
		// 	expectedContexts: expectedValueErrorE,
		// 	expectedError:    expectedErrorE,
		// },
	}

	for _, testcase := range testCases {
		t.Run(testcase.testName, func(t *testing.T) {
			// CREATE TEMPORARY FILE WITH TEST DATA
			tmpFile, err := ioutil.TempFile(os.TempDir(), "funcFilesTest")
			if err != nil {
				t.Error("Cannot create temporary file", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			if _, err = tmpFile.Write([]byte(testcase.fileContent)); err != nil {
				t.Error("Failed to write to temporary file", err)
			}

			// Call function we're testing
			contexts, err := readFunctionsFile(tmpFile.Name())
			isCorrectError := assert.Equal(t, testcase.expectedError, err)
			if !isCorrectError {
				t.Error("Test failed.")
			}

			// Compare contexts with expectedContexts
			isEqual := assert.Equal(t, testcase.expectedContexts, contexts)
			if !isEqual {
				t.Error("Test failed.")
			}
		})
	}
}

const contentA = `
main.foobar(int, int64)
main.buzbaz(bool, rune)
`

var expectedValueA = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},

	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const contentAA = `
   main.foobar(int,       int64)  
main.buzbaz(bool,rune)
`

var expectedValueAA = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const contentDuplicate = `
main.foobar(int,int64)  
main.buzbaz(bool,rune)
main.buzbaz(bool,rune)
`

var expectedValueDuplicate = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const contentB = `
main.foobar([2]byte)
`

var expectedValueB = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BYTE,
				StartingOffset: 8,
				PrintfFormat:   "%c",
				TypeSize:       1,
				ArrayLength:    2,
				VariableName:   "argument1",
			},
		},
	},
}

const contentC = `
main.foobar([3]int, [2]byte)
`

var expectedValueC = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    3,
				VariableName:   "argument1",
			},
			{
				CType:          "char",
				goType:         BYTE,
				StartingOffset: 32,
				PrintfFormat:   "%c",
				TypeSize:       1,
				ArrayLength:    2,
				VariableName:   "argument2",
			},
		},
	},
}

const contentD = `
main.foobar(int, int64)


main.buzbaz(bool, rune)
`

var expectedValueD = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const contentE = `
main.foobar(int, int64)
main.buzbaz(bool, rune)

`

var expectedValueE = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},

	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const contentF = `

main.foobar(int, int64)
main.buzbaz(bool, rune)
`

var expectedValueF = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},

	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const contentG = `

main.foobar(int, int64)



main.buzbaz(bool, rune)

`

var expectedValueG = []functionTraceContext{
	{
		FunctionName: "main.foobar",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "long",
				goType:         INT,
				StartingOffset: 8,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "long",
				goType:         INT64,
				StartingOffset: 16,
				PrintfFormat:   "%ld",
				TypeSize:       8,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
	{
		FunctionName: "main.buzbaz",
		HasArguments: true,
		Arguments: []argument{
			{
				CType:          "char",
				goType:         BOOL,
				StartingOffset: 8,
				PrintfFormat:   "%t",
				TypeSize:       1,
				ArrayLength:    0,
				VariableName:   "argument1",
			},
			{
				CType:          "int",
				goType:         RUNE,
				StartingOffset: 12,
				PrintfFormat:   "%c",
				TypeSize:       4,
				ArrayLength:    0,
				VariableName:   "argument2",
			},
		},
	},
}

const errorValueA = `
main.foooba^r.baz()
`

var expectedValueErrorA []functionTraceContext = nil
var expectedErrorA = errors.New("could not parse function string 'main.foooba^r.baz()': encountered invalid char: ^")

const errorValueB = `
main.foobar(int, int64)
main.foooba^r.baz()
`

var expectedValueErrorB []functionTraceContext = nil
var expectedErrorB = errors.New("could not parse function string 'main.foooba^r.baz()': encountered invalid char: ^")

const errorValueD = `
main.foobar(int,
	 int64)
`

var expectedValueErrorD []functionTraceContext = nil
var expectedErrorD = errors.New("could not parse function string 'main.foobar(int,': incomplete function signature: main.foobar(int,")

const errorValueE = ``

var expectedValueErrorE []functionTraceContext = nil
var expectedErrorE = errors.New("")
