package noid

func bitSwap(val uint64, b1 byte, b2 byte) uint64 {
	if b1 == b2 {
		return val
	}

	if b1 > b2 {
		return bitSwap(val, b2, b1)
	}

	if (val&uint64(1<<b1))^((val&uint64(1<<b2))>>(b2-b1)) != 0 {
		return (val ^ uint64(1<<b1)) ^ (1 << b2)
	}
	return val
}
