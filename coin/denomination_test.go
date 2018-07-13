package coin

import (
	"testing"
)

// Tests the function that returns the value of denominations fails with
// denominations greater than 63
func TestDenomination_ValueTooBig(t *testing.T) {
	d := Denomination(64)

	if d.Value() != 0 {
		t.Errorf("Denomination.Value: Did not zero the output with too large"+
			" of values; Expected: 0, Received: %v+", d.Value())
	}
}

// Tests that all denominations return the correct value
func TestDenomination_ValueAllInputs(t *testing.T) {
	expectedOutputs := []uint64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024,
		2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576,
		2097152, 4194304, 8388608, 16777216, 33554432, 67108864, 134217728,
		268435456, 536870912, 1073741824, 2147483648, 4294967296, 8589934592,
		17179869184, 34359738368, 68719476736, 137438953472, 274877906944,
		549755813888, 1099511627776, 2199023255552, 4398046511104,
		8796093022208, 17592186044416, 35184372088832, 70368744177664,
		140737488355328, 281474976710656, 562949953421312, 1125899906842624,
		2251799813685248, 4503599627370496, 9007199254740992,
		18014398509481984, 36028797018963968, 72057594037927936,
		144115188075855872, 288230376151711744, 576460752303423488,
		1152921504606846976, 2305843009213693952, 4611686018427387904,
		9223372036854775808}

	for d := Denomination(0); d < Denomination(64); d++ {
		if d.Value() != expectedOutputs[d] {
			t.Errorf("Denomination.Value: Denomination %v's value does"+
				" not match expected value: Expected: %v, Received: %v",
				d, expectedOutputs[d], d.Value())
		}
	}
}

func TestGenerateDenominationList_TooLong(t *testing.T) {
	value := uint32(1<<30 - 1)

	_, err := GenerateDenominationList(value)

	if err != ErrInvalidValue {
		t.Errorf("GenerateDenominationList: Did not return correct error"+
			" for value of %v", value)
	}
}

func TestGenerateDenominationList_Happy(t *testing.T) {
	value := uint32(12)

	dl, err := GenerateDenominationList(value)

	if err != nil {
		t.Errorf("GenerateDenominationList: Returned error for valid input"+
			" for value of %v", value)
	}

	resultValue := uint32(0)

	for _, d := range dl {
		resultValue += uint32(d.Value())
	}

	if resultValue != value {
		t.Errorf("GenerateDenominationList: Value of denomination list"+
			"not equal to input value. Input: %v  Output: %v", value, resultValue)
	}
}