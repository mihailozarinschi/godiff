package godiff

import (
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
)

// Chunk contains the hash of a specific block of data, the starting offset in the original input and the length
type Chunk struct {
	DataOffset int64
	DataLen    int64
	Hash       string
}

// ChunkData will split any given data into chunks of hashes based on a rolling-hash algorithm,
// breakpoints are content-based, meaning that the same data patterns will produce always the same breakpoints
func ChunkData(r io.Reader, h hash.Hash, minChunkSize, divisor, prime int64) ([]*Chunk, error) {
	var chunks []*Chunk

	// With the TeeReader, everything we read from the reader, will be written to the hash too
	r = io.TeeReader(r, h)

	var EOF bool
	var currentOffset int64
	var dataWindow = make([]byte, minChunkSize) // minChunkSize will also be our fingerprinting window size

	for !EOF {
		// Reset the hash before starting a new chunk
		h.Reset()

		// Read initial data window
		chunkLen, err := io.ReadFull(r, dataWindow)
		EOF = errors.Is(err, io.EOF)
		if err != nil && !EOF {
			return nil, fmt.Errorf("error reading initial data window: %s", err)
		}
		if chunkLen == 0 {
			// Nothing was read, we're done
			continue
		}
		currentOffset += int64(chunkLen)

		// Adjust the dataWindow length in case we have read less than expected (because of an EOF).
		// In any other case the length should stay the same
		dataWindow = dataWindow[:chunkLen]

		// Calculate the initial data window fingerprint
		dataWindowFingerprint := Fingerprint(dataWindow, prime)
		for {
			if EOF || foundBreakpoint(dataWindowFingerprint, divisor) {
				// We're either done reading, either got to a breakpoint.
				// Get the current chunk's hash, and start the next chunk, if any
				chunks = append(chunks, &Chunk{
					DataOffset: currentOffset - int64(chunkLen),
					DataLen:    int64(chunkLen),
					Hash:       hex.EncodeToString(h.Sum(nil)),
				})
				break
			}

			// Constantly slide data window by popping the first byte, and appending the next byte
			firstByte := dataWindow[0]
			dataWindow = append(dataWindow[1:], 0)                 // Add an empty byte, will be read on the next line
			readLen, err := r.Read(dataWindow[len(dataWindow)-1:]) // Read 1 byte into the last dataWindow byte
			EOF = errors.Is(err, io.EOF)
			if err != nil && !EOF {
				return nil, fmt.Errorf("error reading next byte: %s", err)
			}
			if readLen == 0 {
				// Nothing was read, we're done.
				continue
			}

			currentOffset++
			chunkLen++

			// Calculate the new data window fingerprint
			dataWindowFingerprint = SlideFingerprint(dataWindowFingerprint, prime, firstByte, dataWindow[len(dataWindow)-1], len(dataWindow))
		}
	}

	return chunks, nil
}

func foundBreakpoint(fingerprint, divisor int64) bool {
	return fingerprint%divisor == divisor-1
}
