package cpc

/**
 * This is a unique-counting sketch that implements the
 * <i>Compressed Probabilistic Counting (CPC, a.k.a FM85)</i> algorithms developed by Kevin Lang in
 * his paper
 * <a href="https://arxiv.org/abs/1708.06839">Back to the Future: an Even More Nearly
 * Optimal Cardinality Estimation Algorithm</a>.
 *
 * <p>This sketch is extremely space-efficient when serialized. In an apples-to-apples empirical
 * comparison against compressed HyperLogLog sketches, this new algorithm simultaneously wins on
 * the two dimensions of the space/accuracy tradeoff and produces sketches that are
 * smaller than the entropy of HLL, so no possible implementation of compressed HLL can match its
 * space efficiency for a given accuracy. As described in the paper this sketch implements a newly
 * developed ICON estimator algorithm that survives unioning operations, another
 * well-known estimator, the
 * <a href="https://arxiv.org/abs/1306.3284">Historical Inverse Probability (HIP)</a> estimator
 * does not.
 * The update speed performance of this sketch is quite fast and is comparable to the speed of HLL.
 * The unioning (merging) capability of this sketch also allows for merging of sketches with
 * different configurations of K.
 *
 * <p>For additional security this sketch can be configured with a user-specified hash seed.
 *
 * @author Lee Rhodes
 * @author Kevin Lang
 */

var kxpByteLookup [256]float64
const defaultLgK uint = 11
const defaultSeed uint64 = 9001

type cpcSketch struct {
	seed uint64
	lgK uint
	numCoupons uint64
	mergeFlag bool
	fiCol uint
	windowOffset uint
	slidingWindow []byte
	pairTable pairTable
	//The following variables are only valid in HIP varients
	kxp float64
	hipEstAccum float64
}

func NewCpcSketch() (cpcSketch, error) {
	return NewCpcSketchWithParams(defaultLgK, defaultSeed)
}

func NewCpcSketchWithLgK(lgK uint) (cpcSketch, error) {
	return NewCpcSketchWithParams(lgK, defaultSeed)
}

func NewCpcSketchWithParams(lgK uint, seed uint64) (cpcSketch, error) {
	if err := checkLgK(lgK); err != nil {
		return cpcSketch{}, nil
	}
	sketch := cpcSketch{lgK: lgK, seed: seed, kxp: float64(uint(1) << lgK) }
	sketch.Reset()
	return sketch, nil
}

func (s *cpcSketch) Reset() {
	s.numCoupons = 0
	s.mergeFlag = false
	s.fiCol = 0

	s.windowOffset = 0
	s.slidingWindow = nil
	s.pairTable = pairTable{}
	s.kxp = float64(uint(1) << s.lgK)
	s.hipEstAccum = 0
}

func (s cpcSketch) Copy() (cpcSketch, error) {
	cpy, err := NewCpcSketchWithParams(s.lgK, s.seed)
	if err != nil {
		return cpcSketch{}, err
	}
	cpy.numCoupons = s.numCoupons
	cpy.mergeFlag = s.mergeFlag
	cpy.fiCol = s.fiCol

	cpy.windowOffset = s.windowOffset
	if s.slidingWindow == nil {
		cpy.slidingWindow = nil
	} else {
		cpy.slidingWindow = make([]byte, len(s.slidingWindow))
		copy(cpy.slidingWindow, s.slidingWindow)
	}

}
