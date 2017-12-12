package cyclic

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
)

//Create the cyclic.Int type as a wrapper of the big.Int type
type Int big.Int

//BigInt converts the givne cyclic Int to a big Int and returns it
func BigInt(n *Int) *big.Int {
	nint := big.Int(*n)
	return &nint
}

//Int64 converts the cyclic Int to an Int64 if possible and returns nil if not
func (n *Int) Int64() int64 {
	return BigInt(n).Int64()
}

//IsInt64 checks if a cyclic Int can be converted to an Int64
func (n *Int) IsInt64() bool {
	return BigInt(n).IsInt64()
}

//(Private) nilInt returns a cyclic Int which is nil
func nilInt() *Int {
	return nil
}

//NewInt allocates and returns a new Int set to x.
func NewInt(x int64) *Int {

	nint := big.NewInt(x)
	mint := Int(*nint)

	return &mint
}

//SetString makes the Int equal to the number held in the string s, interpreted to have a base of b. Returns the set Int and a boolean describing if the operation was successful.
func (z *Int) SetString(s string, b int) (*Int, bool) {
	err := errors.New("Unimplemented function: Int.SetString recieved " + reflect.TypeOf(z).String() +
		", " + reflect.TypeOf(s).String() + ", " + reflect.TypeOf(b).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil, false
}

//SetBytes interprets buf as the bytes of a big-endian unsigned integer, sets z to that value, and returns z.
func (z *Int) SetBytes(buf []byte) *Int {
	err := errors.New("Unimplemented function: Int.SetBytes recieved " + reflect.TypeOf(buf).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil

}

//Mod sets z to the modulus x%y for y != 0 and returns z. If y == 0, a division-by-zero run-time panic occurs. Mod implements Euclidean modulus (unlike Go); see DivMod for more details.
func (z *Int) Mod(x, y *Int) *Int {
	err := errors.New("Unimplemented function: Int.Mod recieved " + reflect.TypeOf(z).String() + ", " +
		reflect.TypeOf(x).String() + ", " + reflect.TypeOf(y).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil
}

//ModInverse sets z to the multiplicative inverse of g in the ring ℤ/nℤ and returns z. If g and n are not relatively prime, the result is undefined.
func (z *Int) ModInverse(g, n *Int) *Int {
	err := errors.New("Unimplemented function: Int.ModInverse recieved " + reflect.TypeOf(z).String() + ", " +
		reflect.TypeOf(g).String() + ", " + reflect.TypeOf(n).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil
}

//Add sets z to the sum x+y and returns z.
func (z *Int) Add(x, y *Int) *Int {
	nint := BigInt(z).Add(BigInt(x), BigInt(y))
	mint := Int(*nint)
	return &mint
}

//Mul sets z to the product x*y and returns z.
func (z *Int) Mul(x, y *Int) *Int {
	err := errors.New("Unimplemented function: Int.Mul recieved " + reflect.TypeOf(z).String() + ", " +
		reflect.TypeOf(x).String() + ", " + reflect.TypeOf(y).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil
}

//Exp sets z = x*y mod |m| (i.e. the sign of m is ignored), and returns z. If y <= 0, the result is 1 mod
//|m|; if m == nil or m == 0, z = x*y. Modular exponentation of inputs of a particular size is not a
//cryptographically constant-time operation.
func (z *Int) Exp(x, y, m *Int) *Int {
	err := errors.New("Unimplemented function: Int.Exp recieved " + reflect.TypeOf(z).String() + ", " +
		reflect.TypeOf(x).String() + ", " + reflect.TypeOf(y).String() + ", " +
		reflect.TypeOf(m).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil

}

//Bytes returns the absolute value of x as a big-endian byte slice.
func (x *Int) Bytes() []byte {
	err := errors.New("Unimplemented function: Int.Bytes recieved " + reflect.TypeOf(x).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return nil
}

//BitLen returns the length of the absolute value of x in bits. The bit length of 0 is 0.
func (x *Int) BitLen() int {
	err := errors.New("Unimplemented function: Int.BitLen recieved " + reflect.TypeOf(x).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return -1
}

//Cmp compares x and y and returns:
//	-1 if x < y
//	 0 if x == y
//	+1 if x > y
func (x *Int) Cmp(y *Int) (r int) {
	return BigInt(x).Cmp(BigInt(y))
}

//Text returns the string representation of x in the given base. Base must be between 2 and 36, inclusive.
//The result uses the lower-case letters 'a' to 'z' for digit values >= 10. No base prefix (such as "0x")
//is added to the string.
func (x *Int) Text(base int) string {
	err := errors.New("Unimplemented function: Int.Text recieved " + reflect.TypeOf(x).String() + ", " +
		reflect.TypeOf(base).String() + "\n")

	if err != nil {
		fmt.Print(err)
	}

	return ""
}
