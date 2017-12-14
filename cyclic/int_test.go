package cyclic

import (
	"math/big"
	"reflect"
	"testing"
)

//TestNewInt checks if the NewInt function returns a cyclic Int with
//the same value of the passed int64
func TestNewInt(t *testing.T) {
	expected := big.NewInt(int64(42))

	actual := NewInt(int64(42))

	actualData := bigInt(actual).Int64()
	expectedData := expected.Int64()

	if actualData != expectedData {
		t.Errorf("Test of NewInt failed, expected: '%v', got: '%v'",
			actualData, expectedData)
	}
}

//!!!TestSet!!!

//TestSetString
func TestSetString(t *testing.T) {
	type testStructure struct {
		str  string
		base int
	}

	testStructs := []testStructure{
		{"42", 0},
		{"100000000", 0},
		{"-5", 0},
		{"0", 0},
		{"f", 0},
		{"12", 5},
		{"9000000000000000000000000000000090090909090909090090909090909090", 0},
		{"-1", 2},
	}

	tests := len(testStructs)
	pass := 0

	for i, testi := range testStructs {
		b := big.NewInt(0)
		b, eSuccess := b.SetString(testi.str, testi.base)

		// Test invalid input
		if eSuccess == false {
			actual := NewInt(0)
			actual, aSuccess := actual.SetString(testi.str, testi.base)
			if aSuccess != eSuccess || actual != nil {
				t.Error("Test of SetString() failed at index:", i,
					"Function didn't handle invalid input correctly")
			} else {
				pass += 1
			}
		} else {

			// Test valid input
			expected := cycInt(b)

			actual := NewInt(0)
			actual, aSuccess := actual.SetString(testi.str, testi.base)

			if actual.Cmp(expected) != 0 {
				t.Errorf("Test of SetString() failed at index: %v Expected: %v, %v;",
					" Actual: %v, %v", i, expected, eSuccess, actual, aSuccess)
			} else {
				pass += 1
			}
		}
	}
	println("SetString()", pass, "out of", tests, "tests passed.")
}

//TestSetBytes
func TestSetBytes(t *testing.T) {
	expected := []*Int{
		NewInt(42),
		NewInt(6553522),
		//*NewInt(867530918239450598372829049587), TODO: When text parsing impl
		NewInt(42)}
	testBytes := [][]byte{
		{0x2A},             // 42
		{0x63, 0xFF, 0xB2}, // 6553522
		// { 0xA, 0xF3, 0x24, 0xC1, 0xA0, 0xAD, 0x87, 0x20,
		//   0x57, 0xCE, 0xF4, 0x32, 0xF3 }, //"867530918239450598372829049587",
		{0x2A}} // TODO: Should be <nil>, not 42
	tests := len(expected)
	pass := 0
	actual := NewInt(0)
	for i, testi := range testBytes {
		actual = actual.SetBytes(testi)
		if actual.Cmp(expected[i]) != 0 {
			t.Errorf("Test of SetBytes failed at index %v, expected: '%v', "+
				"actual: %v", i, expected[i].Text(10), actual.Text(10))
		} else {
			pass += 1
		}
	}
	println("SetBytes()", pass, "out of", tests, "tests passed.")
}

//!!!TestInt64!!!

//!!!TestIsInt64!!!

//TestMod checks if the Mod placeholder exists
func TestMod(t *testing.T) {
	var actual, expected int64
	var xint, mint, zint *Int

	zint = NewInt(30)

	//Test where x<m

	expected = 42

	xint = NewInt(42)
	mint = NewInt(69)

	actual = zint.Mod(xint, mint).Int64()

	if actual != expected {
		t.Errorf("Test 'x<m' of Mod failed, expected: '%v', got: '%v'",
			expected, actual)
	}

	//Test where x == m
	expected = 0

	xint = NewInt(42)
	mint = NewInt(42)

	actual = zint.Mod(xint, mint).Int64()

	if actual != expected {
		t.Errorf("Test 'x==m' of Mod failed, expected: '%v', got: '%v'",
			expected, actual)
	}

	//test where x>m

	expected = 27

	xint = NewInt(69)
	mint = NewInt(42)

	actual = zint.Mod(xint, mint).Int64()

	if actual != expected {
		t.Errorf("Test 'x>m' of Mod failed, expected: '%v', got: '%v'",
			expected, actual)
	}

}

//TestModInverse checks if the ModInverse placeholder exists
func TestModInverse(t *testing.T) {

	var expected, actual int64

	expected = 69

	gint := NewInt(42)
	nint := NewInt(27)
	zint := NewInt(30)

	actual = zint.ModInverse(gint, nint).Int64()

	actual = actual * expected

	/*if actual != expected {
		t.Errorf("Test of Mod failed, expected: '%v', got:  '%v'", expected, actual)
	}*/

}

//TestAdd checks if the Add placeholder exists
func TestAdd(t *testing.T) {

	var actual, expected int64
	var xint, yint, zint *Int

	xint = NewInt(42)
	yint = NewInt(69)
	zint = NewInt(30)

	expected = 111
	actual = zint.Add(xint, yint).Int64()

	if actual != expected {
		t.Errorf("Test of Add failed, expected: '%v', got:  '%v'", actual, expected)
	}

}

//TestMul checks if the Mod placeholder exists
/*func TestMul(t *testing.T) {

	expected := nilInt()

	xint := NewInt(42)
	yint := NewInt(69)
	zint := NewInt(30)

	actual := zint.Mul(xint, yint)

	if actual != expected {
		t.Errorf("Test of Mul failed, expected: '%v', got:  '%v'", actual, expected)
	}

}*/

//TestExp checks if the Exp placeholder exists
/*func TestExp(t *testing.T) {

	expected := nilInt()

	xint := NewInt(42)
	yint := NewInt(69)
	zint := NewInt(30)
	mint := NewInt(87)

	actual := zint.Exp(xint, yint, mint)

	actual = actual * expected

	/*if actual != expected {
		t.Errorf("Test of Exp failed, expected: '%v', got:  '%v'", actual, expected)
	}

}*/

//TestBytes checks if the Bytes placeholder exists
func TestBytes(t *testing.T) {
	tests := []int64{
		42,
		6553522,
		-42,
	}

	for i, testi := range tests {
		expected := big.NewInt(testi).Bytes()
		actual := NewInt(testi).Bytes()
		if len(expected) != len(actual) {
			t.Errorf("Case %v of Bytes() failed, Actual: '%v', Expected: '%v'", i, actual, expected)
		}
	}

	// Changed tests to compare output of cyclic Bytes() to big Bytes()
	/*testints := []Int{
		*NewInt(42),
		*NewInt(6553522),
		//*NewInt(867530918239450598372829049587), TODO: When text parsing impl
		*NewInt(-42)}

	expectedbytes := [][]byte{
		{0x2A},             // 42
		{0x63, 0xFF, 0xB2}, // 6553522
		// { 0xA, 0xF3, 0x24, 0xC1, 0xA0, 0xAD, 0x87, 0x20,
		//   0x57, 0xCE, 0xF4, 0x32, 0xF3 }, //"867530918239450598372829049587",
		{0x2A}} // TODO: Should be <nil>, not 42

	for i, tsti := range testints {
		actual := tsti.Bytes()
		fmt.Printf("Big Int: %v Bytes: %v", testints[i].Text(10), actual)
		for j, tstb := range expectedbytes[i] {
			if actual[j] != tstb {
				t.Errorf("Case %v of Bytes() failed, got: '%v', expected: '%v'", i, actual,
					tstb)
			}
		}
	}*/

}

//TestBitLen checks if the BitLen placeholder exists
func TestBitLen(t *testing.T) {
	testints := []Int{
		*NewInt(42),
		*NewInt(6553522),
		//*NewInt(867530918239450598372829049587), TODO: When text parsing impl
		*NewInt(-42)}

	expectedlens := []int{
		6,
		23,
		// ???, TODO: when text parsing implemented
		6}

	for i, tsti := range testints {
		actual := bigInt(&tsti).BitLen()
		if actual != expectedlens[i] {
			t.Errorf("Case %v of BitLen failed, got: '%v', expected: '%v'", i, actual,
				expectedlens[i])
		}
	}
}

//TestCmp checks if the Cmp placeholder exists
func TestCmp(t *testing.T) {

	var expected, actual int
	var xint, yint *Int

	//Tests for case where x < y

	expected = -1

	xint = NewInt(42)
	yint = NewInt(69)

	actual = xint.Cmp(yint)

	if actual != expected {
		t.Errorf("Test 'less than' of Cmp failed, expected: '%v', got:"+
			" '%v'", actual, expected)
	}

	//Tests for case where x==y

	expected = 0

	xint = NewInt(42)
	yint = NewInt(42)

	actual = xint.Cmp(yint)

	if actual != expected {
		t.Errorf("Test 'equals' of Cmp failed, expected: '%v', got: '%v'",
			actual, expected)
	}

	//Test for case where x > y

	expected = 1

	xint = NewInt(69)
	yint = NewInt(42)

	actual = xint.Cmp(yint)

	if actual != expected {
		t.Errorf("Test 'greater than' of Cmp failed, expected: '%v', got:"+
			" '%v'", actual, expected)
	}

}

//TestText checks if the Text placeholder exists
func TestText(t *testing.T) {
	testints := []Int{
		*NewInt(42),
		*NewInt(6553522),
		//*NewInt(867530918239450598372829049587), TODO: When text parsing impl
		*NewInt(-42)}
	expectedstrs := []string{
		"42",
		"6553522",
		//"867530918239450598372829049587",
		"-42"} // TODO: Should be <nil>, not -42

	for i, tsti := range testints {
		actual := tsti.Text(10)
		expected := expectedstrs[i]
		if actual != expected {
			t.Errorf("Test of Text failed, got: '%v', expected: '%v'", actual,
				expected)
		}
	}
}

//TestBigInt checks if the function GetBigInt returns a big.Int
func TestBigInt(t *testing.T) {
	expected := reflect.TypeOf(big.NewInt(int64(42)))

	actual := reflect.TypeOf(bigInt(NewInt(int64(42))))

	if actual != expected {
		t.Errorf("Test of GetBigInt failed, expected: '%v', got: '%v'",
			actual, expected)
	}
}

///!!!TestNilInt!!!
func TestNilInt(t *testing.T) {
	pass, tests := 0, 0
	actual := nilInt()

	// test that value is nil
	tests += 1
	if actual != nil {
		t.Errorf("Test of nilInt() failed. Expected nil value, actual:",
			actual)
	} else {
		pass += 1
	}

	// test that type is *Int
	tests += 1
	c := NewInt(0)
	if reflect.TypeOf(c) != reflect.TypeOf(actual) {
		t.Errorf("Test of nilInt() failed. Expected *Int type, actual:",
			reflect.TypeOf(actual))
	} else {
		pass += 1
	}
	println("nilInt()", pass, "out of", tests, "tests passed.")
}
