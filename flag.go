package ecs

type flag [4]uint64

func (f flag) clone() flag {
	return flag{f[0], f[1], f[2], f[3]}
}

func (f flag) equals(g flag) bool {
	return f[0] == g[0] && f[1] == g[1] && f[2] == g[2] && f[3] == g[3]
}

func (f flag) xor(g flag) flag {
	return flag{f[0] ^ g[0], f[1] ^ g[1], f[2] ^ g[2], f[3] ^ g[3]}
}

func (f flag) and(g flag) flag {
	return flag{f[0] & g[0], f[1] & g[1], f[2] & g[2], f[3] & g[3]}
}

func (f flag) or(g flag) flag {
	return flag{f[0] | g[0], f[1] | g[1], f[2] | g[2], f[3] | g[3]}
}

// contains: f & g == g
func (f flag) contains(g flag) bool {
	return f.and(g).equals(g)
}

func (f flag) iszero() bool {
	return f[0] == 0 && f[1] == 0 && f[2] == 0 && f[3] == 0
}

func newflag(a, b, c, d uint64) flag {
	return flag{a, b, c, d}
}

func newflagbit(bit uint8) flag {
	var a, b, c, d uint64
	if bit > 192 {
		d = 1 << (bit - 192)
	} else if bit > 128 {
		c = 1 << (bit - 128)
	} else if bit > 64 {
		b = 1 << (bit - 64)
	} else {
		a = 1 << bit
	}
	return flag{a, b, c, d}
}
