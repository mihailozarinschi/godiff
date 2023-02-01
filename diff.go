package godiff

import (
	"fmt"
	"hash"
)

// Diff contains everything to know about a specific change between any 2 given inputs of data
type Diff struct {
	*ChunkDelta
	Data []byte
}

// CalcDiffs provides the differences between any 2 given inputs of data, based on the hashing settings
func CalcDiffs(originalData, updatedData ReaderAt, hashFn func() hash.Hash, minChunkSize, divisor, prime int64) ([]*Diff, error) {

	originalChunks, err := ChunkData(originalData, hashFn(), minChunkSize, divisor, prime)
	if err != nil {
		return nil, fmt.Errorf("error chunking original data: %s", err)
	}

	updatedChunks, err := ChunkData(updatedData, hashFn(), minChunkSize, divisor, prime)
	if err != nil {
		return nil, fmt.Errorf("error chunking updated data: %s", err)
	}

	chunksDeltas, err := GetChunksDeltas(originalChunks, updatedChunks)
	if err != nil {
		return nil, fmt.Errorf("error getting original vs updated chunks deltas: %s", err)
	}

	// Load data target of each diff
	diffs := make([]*Diff, len(chunksDeltas))
	for i, chunkDelta := range chunksDeltas {
		diff := &Diff{ChunkDelta: chunkDelta}
		diffs[i] = diff

		switch chunkDelta.Type {
		case DeltaTypeRemove:
			// NOTE: There's no real need to know the deleted data for deleting it,
			// the position and the length should be enough, this is just for testing purposes.
			diff.Data = make([]byte, chunkDelta.DataLen)
			_, err = originalData.ReadAt(diffs[i].Data, chunkDelta.DataOffset)
			if err != nil {
				return nil, fmt.Errorf("error reading original data diff at %d (len=%d): %s", chunkDelta.DataOffset, chunkDelta.DataLen, err)
			}

		case DeltaTypeAdd:
			diff.Data = make([]byte, chunkDelta.DataLen)
			_, err = updatedData.ReadAt(diffs[i].Data, chunkDelta.DataOffset)
			if err != nil {
				return nil, fmt.Errorf("error reading updated data diff at %d (len=%d): %s", chunkDelta.DataOffset, chunkDelta.DataLen, err)
			}

		}
	}

	return diffs, nil
}
