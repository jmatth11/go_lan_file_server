package conversion

// IntToBytes takes an int and converts it into a 4 byte array
func IntToBytes(n int) (a []byte) {
	a = make([]byte, 4)
	a[0] = byte(n)
	a[1] = byte(n >> 8)
	a[2] = byte(n >> 16)
	a[3] = byte(n >> 24)
	return
}

// BytesToInt takes 4 bytes and converts them into an int
func BytesToInt(a, b, c, d byte) int {
	return int(a) | (int(b) << 8) | (int(c) << 16) | (int(d) << 24)
}

// Int64ToBytes takes int64 and converts it into an 8 byte array
func Int64ToBytes(n int64) (a []byte) {
	a = make([]byte, 8)
	a[0] = byte(n)
	a[1] = byte(n >> 8)
	a[2] = byte(n >> 16)
	a[3] = byte(n >> 24)
	a[4] = byte(n >> 32)
	a[5] = byte(n >> 40)
	a[6] = byte(n >> 48)
	a[7] = byte(n >> 56)
	return
}

// BytesToInt64 takes 8 bytes and converts them into an int64
func BytesToInt64(a, b, c, d, e, f, g, h byte) int64 {
	return int64(a) | (int64(b) << 8) | (int64(c) << 16) | (int64(d) << 24) | (int64(e) << 32) | (int64(f) << 40) | (int64(g) << 48) | (int64(h) << 56)
}
