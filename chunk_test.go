package godiff_test

import (
	"crypto/sha1"
	"github.com/mihailozarinschi/godiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"hash"
	"io"
	"strings"
	"testing"
)

func TestChunkData(t *testing.T) {
	tt := []struct {
		name string
		// input
		data func(t *testing.T) io.ReadSeeker // a helper to load the data from files, if necessary
		// configs
		hashFn       hash.Hash
		minChunkSize int64
		divisor      int64
		prime        int64
		// output
		chunks []*godiff.Chunk
	}{
		{
			name: "lorem ipsum (original)",
			data: func(_ *testing.T) io.ReadSeeker {
				return strings.NewReader("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
			},
			hashFn:       sha1.New(),
			minChunkSize: 4,
			divisor:      16,
			prime:        7,
			chunks: []*godiff.Chunk{
				{DataOffset: 0, DataLen: 40, Hash: "922474181e529d97307d8df727fc5cd18d7e3508"},
				{DataOffset: 40, DataLen: 31, Hash: "d78f5778d5b6a19fe3e8273e9575e781796ba8f2"},
				{DataOffset: 71, DataLen: 6, Hash: "61d16a2d286b0c0c22c84b576142cc7476ecbb3d"},
				{DataOffset: 77, DataLen: 35, Hash: "dfd2022b0b4fd16fdf2844b1f21aa8adecf54875"},
				{DataOffset: 112, DataLen: 27, Hash: "947023b26b7d2b1661e890e2d7ee961967e4f9f1"},
				{DataOffset: 139, DataLen: 8, Hash: "3953040043b41732b67f3bcfe4a76154d3bd52ca"},
				{DataOffset: 147, DataLen: 15, Hash: "7aa83cbbc6a76004f1f1e72644434e26ff635c2c"},
				{DataOffset: 162, DataLen: 39, Hash: "bd63b5215974481b866885e6a4a9e6e5cb27dc10"},
				{DataOffset: 201, DataLen: 57, Hash: "72fac7a44ef8da03e89938adce841ad39d9eeefc"},
				{DataOffset: 258, DataLen: 24, Hash: "6cb46c9f3e3960d6ce6dd669e0d4cf74ddc1f558"},
				{DataOffset: 282, DataLen: 107, Hash: "0332132187f0c5de9984d117738294f494eeb70f"},
				{DataOffset: 389, DataLen: 15, Hash: "d3e60be4f992049e346fc594c8169f5374df37cc"},
				{DataOffset: 404, DataLen: 41, Hash: "310860931148f42e25c0be31ae27ed4e1a9f35c0"},
			},
		},
		{
			name: "lorem ipsum (updated)",
			data: func(_ *testing.T) io.ReadSeeker {
				return strings.NewReader("Lorem ipsum dolor sit amet, xxxxxxxxxxx adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniamexercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
			},
			hashFn:       sha1.New(),
			minChunkSize: 4,
			divisor:      16,
			prime:        7,
			chunks: []*godiff.Chunk{
				{DataOffset: 0, DataLen: 71, Hash: "1eb611d7c6d236c622273d0c6d02d148fd70f7fb"},
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
		},
		// TODO: add more test cases, maybe with files, different hash function, etc
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rs := tc.data(t)
			if c, ok := rs.(io.Closer); ok {
				defer c.Close()
			}

			chunks, err := godiff.ChunkData(rs, tc.hashFn, tc.minChunkSize, tc.divisor, tc.prime)
			require.NoError(t, err)
			assert.Equal(t, len(tc.chunks), len(chunks))
			assert.Equal(t, tc.chunks, chunks)

			// Print all chunks data
			_, err = rs.Seek(0, io.SeekStart)
			require.NoError(t, err)
			for i, chunk := range chunks {
				data := make([]byte, chunk.DataLen)
				dataLen, err := io.ReadFull(rs, data)
				require.NoError(t, err)
				if _, ok := rs.(*strings.Reader); ok {
					t.Logf("Chunk #%d startsAt=%d len=%d data=%q", i, chunk.DataOffset, chunk.DataLen, string(data[:dataLen]))
				} else {
					t.Logf("Chunk #%d startsAt=%d len=%d data=%v", i, chunk.DataOffset, chunk.DataLen, data[:dataLen])
				}
			}
		})
	}
}
