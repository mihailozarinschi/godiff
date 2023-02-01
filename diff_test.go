package godiff_test

import (
	"crypto/sha1"
	"github.com/mihailozarinschi/godiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"hash"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCalcDiffs(t *testing.T) {
	tt := []struct {
		name         string
		original     func(t *testing.T) godiff.ReaderAt
		updated      func(t *testing.T) godiff.ReaderAt
		hashFn       func() hash.Hash
		minChunkSize int64
		divisor      int64
		prime        int64
		diffs        []*godiff.Diff
	}{
		{
			name: "lorem ipsum (strings.Reader)",
			original: func(_ *testing.T) godiff.ReaderAt {
				return godiff.NewReaderAt(strings.NewReader("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."))
			},
			updated: func(_ *testing.T) godiff.ReaderAt {
				return godiff.NewReaderAt(strings.NewReader("Lorem ipsum dolor sit amet, xxxxxxxxxxx adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniamexercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."))
			},
			hashFn:       sha1.New,
			minChunkSize: 4,
			divisor:      16,
			prime:        7,
			diffs: []*godiff.Diff{
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 147, DataLen: 15, Hash: "7aa83cbbc6a76004f1f1e72644434e26ff635c2c"},
						Type:     godiff.DeltaTypeRemove,
						Position: 6,
					},
					Data: []byte(", quis nostrud "),
				},
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 40, DataLen: 31, Hash: "d78f5778d5b6a19fe3e8273e9575e781796ba8f2"},
						Type:     godiff.DeltaTypeRemove,
						Position: 1,
					},
					Data: []byte("adipiscing elit, sed do eiusmod"),
				},
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 0, DataLen: 40, Hash: "922474181e529d97307d8df727fc5cd18d7e3508"},
						Type:     godiff.DeltaTypeRemove,
						Position: 0,
					},
					Data: []byte("Lorem ipsum dolor sit amet, consectetur "),
				},
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 0, DataLen: 71, Hash: "1eb611d7c6d236c622273d0c6d02d148fd70f7fb"},
						Type:     godiff.DeltaTypeAdd,
						Position: 0,
					},
					Data: []byte("Lorem ipsum dolor sit amet, xxxxxxxxxxx adipiscing elit, sed do eiusmod"),
				},
			},
		},
		{
			// Exact same data as above, but loaded from files
			name: "lorem ipsum (file)",
			original: func(t *testing.T) godiff.ReaderAt {
				f, err := os.Open("testdata/original.txt")
				require.NoError(t, err)
				return f
			},
			updated: func(t *testing.T) godiff.ReaderAt {
				f, err := os.Open("testdata/updated.txt")
				require.NoError(t, err)
				return f
			},
			hashFn:       sha1.New,
			minChunkSize: 4,
			divisor:      16,
			prime:        7,
			diffs: []*godiff.Diff{
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 147, DataLen: 15, Hash: "7aa83cbbc6a76004f1f1e72644434e26ff635c2c"},
						Type:     godiff.DeltaTypeRemove,
						Position: 6,
					},
					Data: []byte(", quis nostrud "),
				},
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 40, DataLen: 31, Hash: "d78f5778d5b6a19fe3e8273e9575e781796ba8f2"},
						Type:     godiff.DeltaTypeRemove,
						Position: 1,
					},
					Data: []byte("adipiscing elit, sed do eiusmod"),
				},
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 0, DataLen: 40, Hash: "922474181e529d97307d8df727fc5cd18d7e3508"},
						Type:     godiff.DeltaTypeRemove,
						Position: 0,
					},
					Data: []byte("Lorem ipsum dolor sit amet, consectetur "),
				},
				{
					ChunkDelta: &godiff.ChunkDelta{
						Chunk:    &godiff.Chunk{DataOffset: 0, DataLen: 71, Hash: "1eb611d7c6d236c622273d0c6d02d148fd70f7fb"},
						Type:     godiff.DeltaTypeAdd,
						Position: 0,
					},
					Data: []byte("Lorem ipsum dolor sit amet, xxxxxxxxxxx adipiscing elit, sed do eiusmod"),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			o := tc.original(t)
			if c, ok := o.(io.Closer); ok {
				defer c.Close()
			}

			u := tc.updated(t)
			if c, ok := u.(io.Closer); ok {
				defer c.Close()
			}

			diffs, err := godiff.CalcDiffs(o, u, tc.hashFn, tc.minChunkSize, tc.divisor, tc.prime)
			require.NoError(t, err)
			assert.Equal(t, len(tc.diffs), len(diffs))
			assert.Equal(t, tc.diffs, diffs)
		})
	}
}
