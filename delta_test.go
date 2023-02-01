package godiff_test

import (
	"github.com/mihailozarinschi/godiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetChunksDeltas(t *testing.T) {
	tt := []struct {
		name     string
		original []*godiff.Chunk
		updated  []*godiff.Chunk
		deltas   []*godiff.ChunkDelta
	}{
		{
			name:     "ABCB-BABC",
			original: []*godiff.Chunk{{Hash: "A"}, {Hash: "B"}, {Hash: "C"}, {Hash: "B"}},
			updated:  []*godiff.Chunk{{Hash: "B"}, {Hash: "A"}, {Hash: "B"}, {Hash: "C"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "B"}, Type: godiff.DeltaTypeRemove, Position: 3},
				{Chunk: &godiff.Chunk{Hash: "B"}, Type: godiff.DeltaTypeAdd, Position: 0},
			},
		},
		{
			name:     "BABC-ABCB",
			original: []*godiff.Chunk{{Hash: "B"}, {Hash: "A"}, {Hash: "B"}, {Hash: "C"}},
			updated:  []*godiff.Chunk{{Hash: "A"}, {Hash: "B"}, {Hash: "C"}, {Hash: "B"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeRemove, Position: 3},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeRemove, Position: 1},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeAdd, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeAdd, Position: 2},
			},
		},
		{
			name:     "BCBCBCA-ABCBCABC",
			original: []*godiff.Chunk{{Hash: "B"}, {Hash: "C"}, {Hash: "B"}, {Hash: "C"}, {Hash: "B"}, {Hash: "C"}, {Hash: "A"}},
			updated:  []*godiff.Chunk{{Hash: "A"}, {Hash: "B"}, {Hash: "C"}, {Hash: "B"}, {Hash: "C"}, {Hash: "A"}, {Hash: "B"}, {Hash: "C"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeRemove, Position: 6},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeAdd, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeAdd, Position: 5},
			},
		},
		{
			name:     "ABCBCABC-BCBCBCA",
			original: []*godiff.Chunk{{Hash: "A"}, {Hash: "B"}, {Hash: "C"}, {Hash: "B"}, {Hash: "C"}, {Hash: "A"}, {Hash: "B"}, {Hash: "C"}},
			updated:  []*godiff.Chunk{{Hash: "B"}, {Hash: "C"}, {Hash: "B"}, {Hash: "C"}, {Hash: "B"}, {Hash: "C"}, {Hash: "A"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeRemove, Position: 5},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeRemove, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeAdd, Position: 6},
			},
		},
		{
			name:     "ABCDEFK-BHDEFCK",
			original: []*godiff.Chunk{{Hash: "A"}, {Hash: "B"}, {Hash: "C"}, {Hash: "D"}, {Hash: "E"}, {Hash: "F"}, {Hash: "K"}},
			updated:  []*godiff.Chunk{{Hash: "B"}, {Hash: "H"}, {Hash: "D"}, {Hash: "E"}, {Hash: "F"}, {Hash: "C"}, {Hash: "K"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeRemove, Position: 2},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeRemove, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "H"}, Type: godiff.DeltaTypeAdd, Position: 1},
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeAdd, Position: 5},
			},
		},
		{
			name:     "BHDEFCK-ABCDEFK",
			original: []*godiff.Chunk{{Hash: "B"}, {Hash: "H"}, {Hash: "D"}, {Hash: "E"}, {Hash: "F"}, {Hash: "C"}, {Hash: "K"}},
			updated:  []*godiff.Chunk{{Hash: "A"}, {Hash: "B"}, {Hash: "C"}, {Hash: "D"}, {Hash: "E"}, {Hash: "F"}, {Hash: "K"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeRemove, Position: 5},
				{Chunk: &godiff.Chunk{Hash: "H"}, Type: godiff.DeltaTypeRemove, Position: 1},
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeAdd, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeAdd, Position: 2},
			},
		},
		{
			name:     "INTENTION-EXECUTION",
			original: []*godiff.Chunk{{Hash: "I"}, {Hash: "N"}, {Hash: "T"}, {Hash: "E"}, {Hash: "N"}, {Hash: "T"}, {Hash: "I"}, {Hash: "O"}, {Hash: "N"}},
			updated:  []*godiff.Chunk{{Hash: "E"}, {Hash: "X"}, {Hash: "E"}, {Hash: "C"}, {Hash: "U"}, {Hash: "T"}, {Hash: "I"}, {Hash: "O"}, {Hash: "N"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "N"}, Type: godiff.DeltaTypeRemove, Position: 4},
				{Chunk: &godiff.Chunk{Hash: "T"}, Type: godiff.DeltaTypeRemove, Position: 2},
				{Chunk: &godiff.Chunk{Hash: "N"}, Type: godiff.DeltaTypeRemove, Position: 1},
				{Chunk: &godiff.Chunk{Hash: "I"}, Type: godiff.DeltaTypeRemove, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "X"}, Type: godiff.DeltaTypeAdd, Position: 1},
				{Chunk: &godiff.Chunk{Hash: "E"}, Type: godiff.DeltaTypeAdd, Position: 2},
				{Chunk: &godiff.Chunk{Hash: "C"}, Type: godiff.DeltaTypeAdd, Position: 3},
				{Chunk: &godiff.Chunk{Hash: "U"}, Type: godiff.DeltaTypeAdd, Position: 4},
			},
		},
		{
			name:     "BENYAM-EPHREM",
			original: []*godiff.Chunk{{Hash: "B"}, {Hash: "E"}, {Hash: "N"}, {Hash: "Y"}, {Hash: "A"}, {Hash: "M"}},
			updated:  []*godiff.Chunk{{Hash: "E"}, {Hash: "P"}, {Hash: "H"}, {Hash: "R"}, {Hash: "E"}, {Hash: "M"}},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{Hash: "A"}, Type: godiff.DeltaTypeRemove, Position: 4},
				{Chunk: &godiff.Chunk{Hash: "Y"}, Type: godiff.DeltaTypeRemove, Position: 3},
				{Chunk: &godiff.Chunk{Hash: "N"}, Type: godiff.DeltaTypeRemove, Position: 2},
				{Chunk: &godiff.Chunk{Hash: "B"}, Type: godiff.DeltaTypeRemove, Position: 0},
				{Chunk: &godiff.Chunk{Hash: "P"}, Type: godiff.DeltaTypeAdd, Position: 1},
				{Chunk: &godiff.Chunk{Hash: "H"}, Type: godiff.DeltaTypeAdd, Position: 2},
				{Chunk: &godiff.Chunk{Hash: "R"}, Type: godiff.DeltaTypeAdd, Position: 3},
				{Chunk: &godiff.Chunk{Hash: "E"}, Type: godiff.DeltaTypeAdd, Position: 4},
			},
		},
		{
			// Following chunks are from chunk_test.go "lorem ipsum sit amet..."
			name: "Lorem impsum",
			original: []*godiff.Chunk{
				{DataOffset: 0, DataLen: 40, Hash: "922474181e529d97307d8df727fc5cd18d7e3508"},  // -
				{DataOffset: 40, DataLen: 31, Hash: "d78f5778d5b6a19fe3e8273e9575e781796ba8f2"}, // -
				{DataOffset: 71, DataLen: 6, Hash: "61d16a2d286b0c0c22c84b576142cc7476ecbb3d"},
				{DataOffset: 77, DataLen: 35, Hash: "dfd2022b0b4fd16fdf2844b1f21aa8adecf54875"},
				{DataOffset: 112, DataLen: 27, Hash: "947023b26b7d2b1661e890e2d7ee961967e4f9f1"},
				{DataOffset: 139, DataLen: 8, Hash: "3953040043b41732b67f3bcfe4a76154d3bd52ca"},
				{DataOffset: 147, DataLen: 15, Hash: "7aa83cbbc6a76004f1f1e72644434e26ff635c2c"}, // -
				{DataOffset: 162, DataLen: 39, Hash: "bd63b5215974481b866885e6a4a9e6e5cb27dc10"},
				{DataOffset: 201, DataLen: 57, Hash: "72fac7a44ef8da03e89938adce841ad39d9eeefc"},
				{DataOffset: 258, DataLen: 24, Hash: "6cb46c9f3e3960d6ce6dd669e0d4cf74ddc1f558"},
				{DataOffset: 282, DataLen: 107, Hash: "0332132187f0c5de9984d117738294f494eeb70f"},
				{DataOffset: 389, DataLen: 15, Hash: "d3e60be4f992049e346fc594c8169f5374df37cc"},
				{DataOffset: 404, DataLen: 41, Hash: "310860931148f42e25c0be31ae27ed4e1a9f35c0"},
			},
			updated: []*godiff.Chunk{
				{DataOffset: 0, DataLen: 71, Hash: "1eb611d7c6d236c622273d0c6d02d148fd70f7fb"}, // +
				{DataOffset: 71, DataLen: 6, Hash: "61d16a2d286b0c0c22c84b576142cc7476ecbb3d"},
				{DataOffset: 77, DataLen: 35, Hash: "dfd2022b0b4fd16fdf2844b1f21aa8adecf54875"},
				{DataOffset: 112, DataLen: 27, Hash: "947023b26b7d2b1661e890e2d7ee961967e4f9f1"},
				{DataOffset: 139, DataLen: 8, Hash: "3953040043b41732b67f3bcfe4a76154d3bd52ca"},
				{DataOffset: 147, DataLen: 39, Hash: "bd63b5215974481b866885e6a4a9e6e5cb27dc10"},
				{DataOffset: 186, DataLen: 57, Hash: "72fac7a44ef8da03e89938adce841ad39d9eeefc"},
				{DataOffset: 243, DataLen: 24, Hash: "6cb46c9f3e3960d6ce6dd669e0d4cf74ddc1f558"},
				{DataOffset: 267, DataLen: 107, Hash: "0332132187f0c5de9984d117738294f494eeb70f"},
				{DataOffset: 374, DataLen: 15, Hash: "d3e60be4f992049e346fc594c8169f5374df37cc"},
				{DataOffset: 389, DataLen: 41, Hash: "310860931148f42e25c0be31ae27ed4e1a9f35c0"},
			},
			deltas: []*godiff.ChunkDelta{
				{Chunk: &godiff.Chunk{DataOffset: 147, DataLen: 15, Hash: "7aa83cbbc6a76004f1f1e72644434e26ff635c2c"}, Type: godiff.DeltaTypeRemove, Position: 6},
				{Chunk: &godiff.Chunk{DataOffset: 40, DataLen: 31, Hash: "d78f5778d5b6a19fe3e8273e9575e781796ba8f2"}, Type: godiff.DeltaTypeRemove, Position: 1},
				{Chunk: &godiff.Chunk{DataOffset: 0, DataLen: 40, Hash: "922474181e529d97307d8df727fc5cd18d7e3508"}, Type: godiff.DeltaTypeRemove, Position: 0},
				{Chunk: &godiff.Chunk{DataOffset: 0, DataLen: 71, Hash: "1eb611d7c6d236c622273d0c6d02d148fd70f7fb"}, Type: godiff.DeltaTypeAdd, Position: 0},
			},
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			deltas, err := godiff.GetChunksDeltas(tc.original, tc.updated)
			require.NoError(t, err)
			assert.Equal(t, tc.deltas, deltas)

			for _, delta := range deltas {
				t.Logf("%s(%d): %q", delta.Type.String(), delta.Position, delta.Hash)
			}
		})
	}
}
