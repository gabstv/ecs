package ecs

const MaxFlagCapacity = 256

// Flag is a 256 bit binary flag
type Flag [4]uint64

// Clone returns a new flag with identical data
func (f Flag) Clone() Flag {
	return Flag{f[0], f[1], f[2], f[3]}
}

// Equals checs if g contains the same bits
func (f Flag) Equals(g Flag) bool {
	return f[0] == g[0] && f[1] == g[1] && f[2] == g[2] && f[3] == g[3]
}

// Xor bitwise (f ^ g)
func (f Flag) Xor(g Flag) Flag {
	return Flag{f[0] ^ g[0], f[1] ^ g[1], f[2] ^ g[2], f[3] ^ g[3]}
}

// And bitwise (f & g)
func (f Flag) And(g Flag) Flag {
	return Flag{f[0] & g[0], f[1] & g[1], f[2] & g[2], f[3] & g[3]}
}

// Or bitwise (f | g)
func (f Flag) Or(g Flag) Flag {
	return Flag{f[0] | g[0], f[1] | g[1], f[2] | g[2], f[3] | g[3]}
}

// Contains tests if (f & g == g)
func (f Flag) Contains(g Flag) bool {
	return f.And(g).Equals(g)
}

// ContainsAny tests if f contains at least one bit of g
func (f Flag) ContainsAny(g Flag) bool {
	return !f.And(g).IsZero()
}

// IsZero returns true if all bits are zero
func (f Flag) IsZero() bool {
	return f[0] == 0 && f[1] == 0 && f[2] == 0 && f[3] == 0
}

// Lowest bit position (set to 1)
func (f Flag) Lowest() uint8 {
	ff := uint64(0xffffffffffffffff)
	if f[0]&ff > 0 {
		v, _ := ubit(f[0])
		return v
	}
	if f[1]&ff > 0 {
		v, _ := ubit(f[1])
		return v + 64
	}
	if f[2]&ff > 0 {
		v, _ := ubit(f[2])
		return v + 128
	}
	if f[3]&ff > 0 {
		v, _ := ubit(f[3])
		return v + 192
	}
	return 0
}

func ubit(v uint64) (uint8, bool) {
	for i := 0; i < 64; i++ {
		if v&(uint64(1<<i)) == (uint64(1 << i)) {
			return uint8(i), true
		}
	}
	return 0, false
}

// NewFlagRaw creates a new flag
func NewFlagRaw(a, b, c, d uint64) Flag {
	return Flag{a, b, c, d}
}

// NewFlag creates a new flag
func NewFlag(bit uint8) Flag {
	var a, b, c, d uint64
	if bit >= 192 {
		d = 1 << (bit - 192)
	} else if bit >= 128 {
		c = 1 << (bit - 128)
	} else if bit >= 64 {
		b = 1 << (bit - 64)
	} else {
		a = 1 << bit
	}
	return Flag{a, b, c, d}
}
