package e2e

import (
	"crypto/sha256"
	"encoding/hex"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

// Assuming DHKey in 2048bit group
const (
	TEST_DHKEY = "8da81b3784ff5fc3c3bd5a5a2cf94d5315dd56cd585ef9944306627e7ccd3c69" +
		"77fcf2db6ffe95e859136d07e24c7f9b25bf7aed37f7333152a6456c1babf107" +
		"c9bb2de486690fd5e4389690802f77bd68503ff9bd89d038daf0e899db3cc558" +
		"aec65a3cae8eff0af6e0f04fa0f798ff8b2ca6a7307ebfbbc2059e5e7ab552f3" +
		"c988063401b0cd23bb5587d538ccf17353e38972f36d382c58cddb9446abe464" +
		"57ffff98dc24dc323c1ea04b0e17c3608ba5ea254d7e9cd585445018ce0fb43c" +
		"4f7869a4450fbb8451d6b8e9a4509824988bce01d33688c3afdea173b5206f68" +
		"cdc392fed7267caef8398e817512ee46aedf6019b6d82a1d9040204d09873d78"
	TEST_USERID    = 42
	TEST_PARTNERID = 11
	EXPECTED_KEY   = "73612b3df0defe6fa5227dce1180f1b540d50d6647da2a334753d4b316adc1ac" +
		"bc7b2dd89519e04d072eb8fa973e1567625a07e20d6fc4ed4c3146121f43f5a0" +
		"35660fa38995dbe77238dd92b981c4e8a1d351a793b57644afba38272b6c87df" +
		"2ad83c39fa8881ba066860e8fffa9dbb11dc991d8553045cf4c961145e57f4a6" +
		"6664860bdc72491492fb890685d30c7832dc8ac822b62c1b8a69991d3b0e1412" +
		"893d8ce8c18ff7c82332d1cd1a1a207fb3d100eadb0b8de8a8bc9d7d40cc0661" +
		"75eb5d1dea4cd6e93303922ac470a29f09eb841affa1f285282c9c224aa8790c" +
		"c07fc8026ef843c25db983a5bb8944cfa8d8b93a8e04b8e9876b2998c2d8bea8"
)

// Test for consistency with hardcoded values
func TestDeriveSingleKey(t *testing.T) {
	userID := id.NewUserFromUint(TEST_USERID, t)
	key := cyclic.NewIntFromString(TEST_DHKEY, 16)
	data := append([]byte{}, key.Bytes()...)
	data = append(data, userID.Bytes()...)
	result := deriveSingleKey(sha256.New(), &grp, data, 0)
	expected := cyclic.NewIntFromString(EXPECTED_KEY, 16)
	actual := cyclic.NewIntFromBytes(result)
	if actual.Cmp(expected) != 0 {
		t.Errorf("Generated Key %v doesn't match expected %v", actual.Text(16), EXPECTED_KEY)
	}
}

// Test both functions with various arguments
func TestDeriveKeys_DeriveReKeys(t *testing.T) {
	userID := id.NewUserFromUint(TEST_USERID, t)
	partnerID := id.NewUserFromUint(TEST_PARTNERID, t)
	key := cyclic.NewIntFromString(TEST_DHKEY, 16)

	nkeys := []uint{10000, 0}
	total := func(a []uint) (s int) {
		for _, n := range a {
			s += int(n)
		}
		return s
	}(nkeys)

	ids := []*id.User{userID, partnerID}
	emerg := []bool{false, true}
	type testFun func(a *cyclic.Group, b *cyclic.Int, c *id.User, d bool, e uint) []*cyclic.Int
	fut := []testFun{DeriveKeys, DeriveReKeys}

	pass := 0
	tests := len(nkeys) * len(ids) * len(emerg) * len(fut)

	expectedKeys := (tests / len(nkeys)) * (total)
	testmap := make(map[string]bool)
	var genKeys = []*cyclic.Int{}

	for _, n := range nkeys {
		for _, id := range ids {
			for _, e := range emerg {
				for _, f := range fut {
					genKeys = f(&grp, key, id, e, n)

					// Check array of keys and if the size matches with requested
					if genKeys == nil {
						t.Errorf("Generated Array of Keys is nil")
					} else if uint(len(genKeys)) != n {
						t.Errorf("Requested %d keys but got %d instead", n, len(genKeys))
					} else {
						// Check each key
						for _, k := range genKeys {
							if k == nil {
								t.Errorf("One generated Key is nil")
							} else if !grp.Inside(k) {
								t.Errorf("Generated key is not inside the group")
							} else {
								testmap[hex.EncodeToString(k.Bytes())] = true
							}
						}
						pass++
					}
				}
			}
		}
	}

	// Confirm all generated keys are different
	if len(testmap) != expectedKeys {
		t.Errorf("Expected %d different keys, but got %d", expectedKeys, len(testmap))
	}

	println("TestDeriveKeys_DeriveReKeys()", pass, "out of", tests, "tests passed.")
}
