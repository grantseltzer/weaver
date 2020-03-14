# Stack offsets

<i>Disclaimer: this doc is based on observations using the cmd/print-stack tool in <b>go 1.13 on linux x86_64 5.4</b>. This is subject to change with different versions of Go and linux.</i>

- Look at size of largest data type that's being passed, that sets the window size. The maximum value for the window is 8.

- Each element added is limited by whether or not it will fit into that window.

- If it would go over a limit window then pad until back at 0, add it, then continue.

## Sizes of go types on the stack

8 Bytes - int, int64, uint, uint64, float64, pointers of any kind, maps, string*, arrays*

4 bytes - int32, uint32, float32, rune

2 bytes - int16, uint16

1 byte - int8, uint8, bool, byte

## Strings

Strings are loaded on the stack as an eight byte address to a byte array (each byte being a character), followed by an 8 byte string length.

## Arrays

The size of the elements within the array are used to calculate the padding on the stack. Therefore if the elements are 2 bytes each, and so are the rest of the function arguments (if any), the padding will be calculated with 2 byte windows. If there's an 8 byte argument with 2 byte elements in the array, the array is padded with 8 byte windows.

## Slices

Slices take up 24 bytes on the stack. There are 3 elements that each take up 8 bytes. Address of an array of memory where the elements are stored, the length of that array, and the cap of the slice.

## Structs / Interfaces

TODO