package main

/******************/
/* STRUCT TESTING */
/******************/

type testStruct struct {
}

//go:noinline
func (t *testStruct) testSingleByte(x byte) {}
