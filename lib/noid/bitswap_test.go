package noid

import (
	"testing"
)

func TestBitSwap(t *testing.T) {
	var value uint64

	// 0101 -> 0101
	value = bitSwap(5, 1, 3)
	assertEqualUint64(5, value, "0101 swapping bits 1 and 3 is still 0101", t)

	// 0110 -> 1100
	value = bitSwap(6, 1, 3)
	assertEqualUint64(12, value, "0110 swapping bits 1 and 3 is 1100", t)

	// 1001001001
	value = 585

	// 1001101000
	value = bitSwap(value, 0, 5)
	assertEqualUint64(616, value, "585 first swap", t)

	// 0001101010
	value = bitSwap(value, 1, 9)
	assertEqualUint64(106, value, "585 second swap", t)

	// 0001101010
	value = bitSwap(value, 2, 7)
	assertEqualUint64(106, value, "585 third swap", t)

	// 0001110010
	value = bitSwap(value, 3, 4)
	assertEqualUint64(114, value, "585 fourth swap", t)

	// 0100110010
	value = bitSwap(value, 6, 8)
	assertEqualUint64(306, value, "585 fifth swap", t)
}
