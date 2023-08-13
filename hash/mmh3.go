package hash

import (
	"encoding/binary"
	"errors"
	"hash"
)

type murmurHash3 struct {
	state hashState
	seed uint64
	buf [bufSize]byte
	n int
}

const bufSize = 128
const sizeOfUint64 = 8

type Hash128 interface{
	hash.Hash64
	Sum128() (uint64, uint64)
}

func NewMurmurHash3() murmurHash3 {
	return NewMurmurHash3WithSeed(0)
}

func NewMurmurHash3WithSeed(seed uint64) murmurHash3 {
	return murmurHash3{state: hashState{h1: seed, h2: seed}, seed: seed}
}

func (h *murmurHash3) WriteByte(b byte) error {
	if h.n == len(h.buf) {
		h.flush()
	}
	h.buf[h.n] = b
	h.n++
	return nil
}

func (h *murmurHash3) Reset() {
	h.state = hashState{h1: h.seed, h2: h.seed}
	h.n = 0
}

func (h *murmurHash3) flush() {
	//TODO: hash buffer
}

func HashUint64(key, seed uint64) (uint64, uint64, error) {
	h := NewMurmurHash3WithSeed(seed)
	h1, h2 := h.state.finalMix128(key, 0, sizeOfUint64)
	return h1, h2, nil
}

func errorForNilKey() (uint64, uint64, error) {
	return 0, 0, errors.New("empty key provided for hashing")
}

func HashUint64Slice(key []uint64, seed uint64) (uint64, uint64, error) {
	if key == nil {
		return errorForNilKey()
	}
	h := NewMurmurHash3WithSeed(seed)
	length := len(key)
	nBlocks := length >> 1

	for i := 0; i < nBlocks; i++ {
		k1, k2 := key[2*i], key[2*i+1]
		h.state.blockMix128(k1, k2)
	}
	var k1 uint64
	if rem := length - 2*nBlocks; rem == 0 {
		k1 = 0
	} else {
		k1 = key[2*nBlocks + 1]
	}
	h1, h2 := h.state.finalMix128(k1, 0, uint64(length) << 3)
	return h1, h2, nil
}

func HashBytes(key []byte, seed uint64) (uint64, uint64, error) {
	if key == nil {
		return errorForNilKey()
	}
	h := NewMurmurHash3WithSeed(seed)
	length := len(key)
	nBlocks := length >> 4 // length / 16

	for i := 0; i < nBlocks; i++ {
		k1 := binary.LittleEndian.Uint64(key[8*i:8*i+8])
		k2 := binary.LittleEndian.Uint64(key[8*(i+1):8*(i+1)+8])
		h.state.blockMix128(k1, k2)
	}
	var k1, k2 uint64
	pos := nBlocks * 8
	if rem := length - nBlocks * 8; rem > 8 {
		k1 = binary.LittleEndian.Uint64(key[pos:pos + 8])
		pos += 8
		buf := make([]byte, 8)
		copy(buf,key[pos:])
		k2 = binary.LittleEndian.Uint64(buf)
	} else if(rem == 8) {
		k1 = binary.LittleEndian.Uint64(key[pos:pos + 8])
		k2 = 0
	} else {
		buf := make([]byte, 8)
		copy(buf,key[pos:])
		k1 = binary.LittleEndian.Uint64(buf)
		k2 = 0
	}
	h1, h2 := h.state.finalMix128(k1, k2, uint64(length))
	return h1, h2, nil
}