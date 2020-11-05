package main

/********************/
/* SINGLE PARAMETER */
/********************/

//go:noinline
func test_single_byte(x byte) {}

//go:noinline
func test_a_different_single_byte(x byte) {}

//go:noinline
func test_single_rune(x rune) {}

//go:noinline
func test_single_string(x string) {}

//go:noinline
func test_single_bool(x bool) {}

//go:noinline
func test_single_int(x int) {}

//go:noinline
func test_single_int8(x int8) {}

//go:noinline
func test_single_int16(x int16) {}

//go:noinline
func test_single_int32(x int32) {}

//go:noinline
func test_single_int64(x int64) {}

//go:noinline
func test_single_uint(x uint) {}

//go:noinline
func test_single_uint8(x uint8) {}

//go:noinline
func test_single_uint16(x uint16) {}

//go:noinline
func test_single_uint32(x uint32) {}

//go:noinline
func test_single_uint64(x uint64) {}

//go:noinline
func test_single_float32(x float32) {}

//go:noinline
func test_single_float64(x float64) {}

/***********************/
/* Multiple Parameters */
/***********************/

//go:noinline
func test_combined_byte(w byte, x byte, y float32) {}

//go:noinline
func test_combined_rune(w byte, x rune, y float32) {}

//go:noinline
func test_combined_string(w byte, x string, y float32) {}

//go:noinline
func test_combined_bool(w byte, x bool, y float32) {}

//go:noinline
func test_combined_int(w byte, x int, y float32) {}

//go:noinline
func test_combined_int8(w byte, x int8, y float32) {}

//go:noinline
func test_combined_int16(w byte, x int16, y float32) {}

//go:noinline
func test_combined_int32(w byte, x int32, y float32) {}

//go:noinline
func test_combined_int64(w byte, x int64, y float32) {}

//go:noinline
func test_combined_uint(w byte, x uint, y float32) {}

//go:noinline
func test_combined_uint8(w byte, x uint8, y float32) {}

//go:noinline
func test_combined_uint16(w byte, x uint16, y float32) {}

//go:noinline
func test_combined_uint32(w byte, x uint32, y float32) {}

//go:noinline
func test_combined_uint64(w byte, x uint64, y float32) {}

/********************/
/* ARRAY PARAMETERs */
/********************/

//go:noinline
func test_byte_array(x [2]byte) {}

//go:noinline
func test_rune_array(x [2]rune) {}

//go:noinline
func test_string_array(x [2]string) {}

//go:noinline
func test_bool_array(x [2]bool) {}

//go:noinline
func test_int_array(x [2]int) {}

//go:noinline
func test_int8_array(x [2]int8) {}

//go:noinline
func test_int16_array(x [2]int16) {}

//go:noinline
func test_int32_array(x [2]int32) {}

//go:noinline
func test_int64_array(x [2]int64) {}

//go:noinline
func test_uint_array(x [2]uint) {}

//go:noinline
func test_uint8_array(x [2]uint8) {}

//go:noinline
func test_uint16_array(x [2]uint16) {}

//go:noinline
func test_uint32_array(x [2]uint32) {}

//go:noinline
func test_uint64_array(x [2]uint64) {}

func main() {

	//TODO: Read the values to pass into these functions from a file

	test_single_byte('a')
	test_a_different_single_byte('b')
	test_single_rune('a')
	test_single_string("a")
	test_single_bool(true)
	test_single_int(1)
	test_single_int8(1)
	test_single_int16(1)
	test_single_int32(1)
	test_single_int64(1)
	test_single_uint(1)
	test_single_uint8(1)
	test_single_uint16(1)
	test_single_uint32(1)
	test_single_uint64(1)
	test_single_float32(1.1)
	test_single_float64(1.1)
	test_combined_byte(2, 3, 3.0)
	test_combined_rune(2, 'b', 3.0)
	test_combined_string(2, "boo", 3.0)
	test_combined_bool(2, true, 3.0)
	test_combined_int(2, 3, 3.0)
	test_combined_int8(2, 38, 3.0)
	test_combined_int16(2, 316, 3.0)
	test_combined_int32(2, 332, 3.0)
	test_combined_int64(2, 364, 3.0)
	test_combined_uint(2, 12, 3.0)
	test_combined_uint8(2, 128, 3.0)
	test_combined_uint16(2, 1216, 3.0)
	test_combined_uint32(2, 1232, 3.0)
	test_combined_uint64(2, 1264, 3.0)
	test_byte_array([2]byte{3, 9})
	test_rune_array([2]rune{'a', 'b'})
	test_string_array([2]string{"boo", "bah"})
	test_bool_array([2]bool{true, true})
	test_int_array([2]int{5, 7})
	test_int8_array([2]int8{1, 5})
	test_int16_array([2]int16{166, 934})
	test_int32_array([2]int32{1325, 512})
	test_int64_array([2]int64{51, 77})
	test_uint_array([2]uint{1, 55})
	test_uint8_array([2]uint8{15, 55})
	test_uint16_array([2]uint16{951, 123})
	test_uint32_array([2]uint32{5135, 512321})
	test_uint64_array([2]uint64{412412456, 1234134})
}
