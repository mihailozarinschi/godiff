package godiff_test

import (
	"github.com/mihailozarinschi/godiff"
	"github.com/stretchr/testify/require"
	"testing"
)

const prime = 7

func TestFingerprint(t *testing.T) {
	tt := []struct {
		data        []byte
		fingerprint int64
	}{
		{
			data:        []byte("abc"),
			fingerprint: 5538,
		},
		{
			data:        []byte("acb"),
			fingerprint: 5544,
		},
		{
			data:        []byte("bac"),
			fingerprint: 5580,
		},
		{
			data:        []byte("bca"),
			fingerprint: 5592,
		},
		{
			data:        []byte("cab"),
			fingerprint: 5628,
		},
		{
			data:        []byte("cba"),
			fingerprint: 5634,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run("", func(t *testing.T) {
			fingerprint := godiff.Fingerprint(tc.data, prime)
			require.Equal(t, tc.fingerprint, fingerprint)
		})
	}
}

func TestSlideFingerprint(t *testing.T) {
	tt := []struct {
		oldFingerprint int64
		byteOut        byte
		byteIn         byte
		windowLen      int
		newFingerprint int64
	}{
		{
			oldFingerprint: 5538, // data: abc
			byteOut:        'a',
			byteIn:         'a',
			windowLen:      3,
			newFingerprint: 5592, // data: bca
		},
		{
			oldFingerprint: 5592, // data: bca
			byteOut:        'b',
			byteIn:         'b',
			windowLen:      3,
			newFingerprint: 5628, // data: cab
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run("", func(t *testing.T) {
			newFingerprint := godiff.SlideFingerprint(tc.oldFingerprint, prime, tc.byteOut, tc.byteIn, tc.windowLen)
			require.Equal(t, tc.newFingerprint, newFingerprint)
		})
	}
}
