package godiff

import "math"

//TODO: Think of a mod (%) mechanism over a constant to prevent integer overflows

// Fingerprint calculates a "weak hash" of any given input. It's based on the Rabin fingerprint method,
// combined with SlideFingerprint, allows to initiate a fingerprint once, and then slide it by adding
// the next byte and dropping the first one.
func Fingerprint(s []byte, prime int64) int64 {
	var (
		n = len(s)
		f int64
	)

	for i := 0; i < n; i++ {
		pow := n - i - 1
		primePow := math.Pow(float64(prime), float64(pow))
		f += int64(primePow) * int64(s[i])
	}

	return f
}

// SlideFingerprint will recalculate the "weak hash" of a previous input that
// slid right by dropping the first byte and getting a new one at the end.
func SlideFingerprint(prevFingerprint, prime int64, out, in byte, windowLen int) int64 {
	pow := windowLen - 1
	primePow := math.Pow(float64(prime), float64(pow))
	outFingerprint := int64(out) * int64(primePow)

	return (prevFingerprint-outFingerprint)*prime + int64(in)
}
