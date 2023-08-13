package hash

import "math/bits"

const c1 uint64 = 0x87c37b91114253d5
const c2 uint64 = 0x4cf5ad432745937f

type hashState struct {
	h1 uint64
	h2 uint64
}

func (s *hashState) blockMix128(k1, k2 uint64) {
	s.h1 ^= mixK1(k1)
	s.h1 = bits.RotateLeft64(s.h1, 27)
	s.h1 += s.h2
	s.h1 = s.h1*5 + 0x52dce729

	s.h2 ^= mixK2(k2)
	s.h2 = bits.RotateLeft64(s.h2, 31)
	s.h2 += s.h1
	s.h2 = s.h2*5 + 0x38495ab5
}

func mixK1(k1 uint64) uint64 {
	k1 *= c1
	k1 = bits.RotateLeft64(k1, 31)
	k1 *= c2
	return k1
}

func mixK2(k2 uint64) uint64 {
	k2 *= c2
	k2 = bits.RotateLeft64(k2, 33)
	k2 *= c1
	return k2
}

func (s *hashState) finalMix128(k1, k2 uint64, inputLengthBytes uint64) (uint64, uint64) {
	s.h1 ^= mixK1(k1)
	s.h2 ^= mixK2(k2)
	s.h1 ^= inputLengthBytes
	s.h2 ^= inputLengthBytes
	s.h1 += s.h2
	s.h2 += s.h1
	s.h1 = finalMix64(s.h1)
	s.h2 = finalMix64(s.h2)
	s.h1 += s.h2
	s.h2 += s.h1
	return s.h1, s.h2
}

func finalMix64(h uint64) uint64 {
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}
