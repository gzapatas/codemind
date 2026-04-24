package ingestion

import (
	"fmt"
	"strings"
)

type Chunk struct {
	Kind      string
	Name      string
	StartLine int // 1-based
	EndLine   int // 1-based
	Text      string
}

// ChunkerFunc is a simple function adapter to implement a chunker for a language.
type ChunkerFunc func(content []byte, lang string) ([]Chunk, error)

type registeredChunker struct {
	fn       ChunkerFunc
	priority int // higher value = higher priority
}

var registry = map[string][]registeredChunker{}

// RegisterChunkerWithPriority registers a chunker function for a language with a priority.
// Higher priority chunkers are tried first.
func RegisterChunkerWithPriority(lang string, f ChunkerFunc, priority int) {
	l := strings.ToLower(lang)
	registry[l] = append(registry[l], registeredChunker{fn: f, priority: priority})
	// simple insertion sort to keep slice ordered by descending priority
	arr := registry[l]
	for i := 1; i < len(arr); i++ {
		j := i
		for j > 0 && arr[j].priority > arr[j-1].priority {
			arr[j], arr[j-1] = arr[j-1], arr[j]
			j--
		}
	}
	registry[l] = arr
}

// RegisterChunker is a convenience that registers with priority 0.
func RegisterChunker(lang string, f ChunkerFunc) {
	RegisterChunkerWithPriority(lang, f, 0)
}

// ChunkFile dispatches to the highest-priority registered chunker for the given language.
func ChunkFile(content []byte, lang string) ([]Chunk, error) {
	l := strings.ToLower(lang)
	if arr, ok := registry[l]; ok && len(arr) > 0 {
		for _, rc := range arr {
			chunks, err := rc.fn(content, l)
			if err == nil {
				return chunks, nil
			}
			// if one fails, try the next registered chunker
		}
		return nil, fmt.Errorf("all registered chunkers failed for language: %s", lang)
	}
	return nil, fmt.Errorf("no chunker registered for language: %s", lang)
}
