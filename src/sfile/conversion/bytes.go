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
