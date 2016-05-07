package monitor

// SplitByChunk for Message is too long
func SplitByChunk(longString string, chunkSize int) []string {
	slices := []string{}
	lastIndex := 0
	lastI := 0

	for i := range longString {
		if i-lastIndex > chunkSize {
			slices = append(slices, longString[lastIndex:lastI])
			lastIndex = lastI
		}
		lastI = i
	}

	if len(longString)-lastIndex > chunkSize {
		slices = append(slices, longString[lastIndex:lastIndex+chunkSize], longString[lastIndex+chunkSize:])
	} else {
		slices = append(slices, longString[lastIndex:])
	}

	return slices
}
