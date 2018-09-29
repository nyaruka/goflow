package utils

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// DefaultRand is the default rand for calls to Rand()
var DefaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var currentRand = DefaultRand

// NewSeededRand creates a new seeded rand
func NewSeededRand(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// RandDecimal returns a random decimal in the range [0.0, 1.0)
func RandDecimal() decimal.Decimal {
	return decimal.NewFromFloat(currentRand.Float64())
}

// RandIntN returns a random integer in the range [0, n)
func RandIntN(n int) int {
	return currentRand.Intn(n)
}

// SetRand sets the rand used by Rand()
func SetRand(rnd *rand.Rand) {
	currentRand = rnd
}

// https://github.com/python/cpython/blob/master/Modules/_randommodule.c

const (
	N          = 624
	M          = 397
	MATRIX_A   = uint32(0x9908b0df) /* constant vector a */
	UPPER_MASK = uint32(0x80000000) /* most significant w-r bits */
	LOWER_MASK = uint32(0x7fffffff) /* least significant r bits */
)

type PythonRand struct {
	index int
	state [N]uint32
}

func (r *PythonRand) State() [N]uint32 {
	return r.state
}

/* initializes mt[N] with a seed */
func (r *PythonRand) genrand(s uint32) {
	var mti int

	r.state[0] = s
	for mti = 1; mti < N; mti++ {
		r.state[mti] = (uint32(1812433253)*(r.state[mti-1]^(r.state[mti-1]>>30)) + uint32(mti))
		/* See Knuth TAOCP Vol2. 3rd Ed. P.106 for multiplier. */
		/* In the previous versions, MSBs of the seed affect   */
		/* only MSBs of the array mt[].                                */
		/* 2002/01/09 modified by Makoto Matsumoto                     */
	}
	r.index = mti
}

/* initialize by an array with array-length */
/* init_key is the array for initializing keys */
/* key_length is its length */
func (r *PythonRand) initByArray(initKey []uint32) {
	var i, j, k uint /* was signed in the original code. RDH 12/16/2002 */

	r.genrand(uint32(19650218))
	i = 1
	j = 0

	if N > len(initKey) {
		k = N
	} else {
		k = uint(len(initKey))
	}

	for ; k > 0; k-- {
		r.state[i] = (r.state[i] ^ ((r.state[i-1] ^ (r.state[i-1] >> 30)) * uint32(1664525))) + initKey[j] + uint32(j) /* non linear */
		i++
		j++
		if i >= N {
			r.state[0] = r.state[N-1]
			i = 1
		}
		if j >= uint(len(initKey)) {
			j = 0
		}
	}
	for k = N - 1; k > 0; k-- {
		r.state[i] = (r.state[i] ^ ((r.state[i-1] ^ (r.state[i-1] >> 30)) * uint32(1566083941))) - uint32(i) /* non linear */
		i++
		if i >= N {
			r.state[0] = r.state[N-1]
			i = 1
		}
	}

	r.state[0] = uint32(0x80000000) /* MSB is 1; assuring non-zero initial array */
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
func (r *PythonRand) Seed(s int64) {
	// get the absolute value of our seed
	if s < 0 {
		s = -s
	}

	// spli into 32-bit chunks
	var key []uint32
	if s == int64(uint32(s)) {
		key = []uint32{uint32(s)}
	} else {
		u1 := uint32(s)
		u2 := uint32(s >> 32)
		key = []uint32{u1, u2}
	}

	r.initByArray(key)
}

/* mag01[x] = x * MATRIX_A  for x=0,1 */
var mag01 = []uint32{0, MATRIX_A}

/* generates a random number on [0,0xffffffff]-interval */
func (r *PythonRand) Int32() uint32 {
	var y uint32

	if r.index >= N { /* generate N words at one time */
		var kk int

		for kk = 0; kk < N-M; kk++ {
			y = (r.state[kk] & UPPER_MASK) | (r.state[kk+1] & LOWER_MASK)
			r.state[kk] = r.state[kk+M] ^ (y >> 1) ^ mag01[y&uint32(1)]
		}
		for ; kk < N-1; kk++ {
			y = (r.state[kk] & UPPER_MASK) | (r.state[kk+1] & LOWER_MASK)
			r.state[kk] = r.state[kk+(M-N)] ^ (y >> 1) ^ mag01[y&uint32(1)]
		}
		y = (r.state[N-1] & UPPER_MASK) | (r.state[0] & LOWER_MASK)
		r.state[N-1] = r.state[M-1] ^ (y >> 1) ^ mag01[y&uint32(1)]

		r.index = 0
	}

	y = r.state[r.index]
	r.index++
	y ^= (y >> 11)
	y ^= (y << 7) & uint32(0x9d2c5680)
	y ^= (y << 15) & uint32(0xefc60000)
	y ^= (y >> 18)
	return y
}

/* random_random is the function named genrand_res53 in the original code;
 * generates a random number on [0,1) with 53-bit resolution; note that
 * 9007199254740992 == 2**53; I assume they're spelling "/2**53" as
 * multiply-by-reciprocal in the (likely vain) hope that the compiler will
 * optimize the division away at compile-time.  67108864 is 2**26.  In
 * effect, a contains 27 random bits shifted left 26, and b fills in the
 * lower 26 bits of the 53-bit numerator.
 * The original code credited Isaku Wada for this algorithm, 2002/01/09.
 */
func (r *PythonRand) Random() float64 {
	a := r.Int32() >> 5
	b := r.Int32() >> 6
	return (float64(a)*67108864.0 + float64(b)) * (1.0 / 9007199254740992.0)
}
