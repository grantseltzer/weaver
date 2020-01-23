package main

// Determining stack sizes:
//
// - Look at size of largest data type that's being passed, that sets the window size
// - Each element added is limited by whether or not it will fit into that window
// - If it would go over a limit window then pad until back at 0, add it, then continue
//

func determineStackOffsets(context *traceContext) error {

	var windowSize = 0

	for _, t := range context.Arguments {
		size := goTypeToSizeInBytes[t.goType]
		if size > windowSize {
			windowSize = size
		}
	}

	currentIndex := 8
	bytesInCurrentWindow := 0

	for i := range context.Arguments {
		typeSize := goTypeToSizeInBytes[context.Arguments[i].goType]
		context.Arguments[i].TypeSize = typeSize

		if typeSize+bytesInCurrentWindow > windowSize {
			// Doesn't fit, move index ahead for padding, clear current window
			currentIndex += windowSize - bytesInCurrentWindow
			bytesInCurrentWindow = 0
		}

		context.Arguments[i].StartingOffset = currentIndex

		if context.Arguments[i].ArrayLength > 0 {
			if context.Arguments[i].goType == STRING {
				typeSize = 16
			}
			currentIndex += typeSize * context.Arguments[i].ArrayLength
			bytesInCurrentWindow += (typeSize * context.Arguments[i].ArrayLength) % windowSize
			continue
		}

		currentIndex += typeSize
		bytesInCurrentWindow += typeSize

		//XXX: In go strings take up 16 bytes on the stack, 8 for the pointer and 8 for length
		if context.Arguments[i].goType == STRING {
			currentIndex += 8
		}

	}

	return nil
}
