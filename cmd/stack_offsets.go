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

		if typeSize+bytesInCurrentWindow > windowSize {
			// Doesn't fit, move index ahead for padding, clear current window
			currentIndex += windowSize - bytesInCurrentWindow
			bytesInCurrentWindow = 0
		}

		context.Arguments[i].StartingOffset = currentIndex
		currentIndex += typeSize
		bytesInCurrentWindow += typeSize

	}

	return nil
}
