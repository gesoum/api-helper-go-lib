package utils

// SplitToChunk splits the slice into several slices of size chunkSize
func SplitToChunk[T any](chunkSize int, collection []T) [][]T {
	var chunks [][]T
	for i := 0; i < len(collection); i += chunkSize {
		end := i + chunkSize

		if end > len(collection) {
			end = len(collection)
		}

		chunks = append(chunks, collection[i:end])
	}
	return chunks
}
