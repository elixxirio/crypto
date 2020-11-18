////////////////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                                       //
//                                                                                        //
// Use of this source code is governed by a license that can be found in the LICENSE file //
////////////////////////////////////////////////////////////////////////////////////////////
package cmix

import (
	"gitlab.com/xx_network/crypto/csprng"
	"math/rand"
	"reflect"
	"testing"
)

//Tests that the ClientGateway key generates consistently and expected
func TestGenerateClientGatewayKey(t *testing.T) {
	grp := grpTest()

	rng := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {

		baseKeyBytes, err := csprng.GenerateInGroup(grp.GetPBytes(), grp.GetP().ByteLen(), rng)
		if err != nil {
			t.Errorf("could not generate base base keys")
		}
		baseKey := grp.NewIntFromBytes(baseKeyBytes)

		clientGatewayKey := GenerateClientGatewayKey(baseKey)
		if !reflect.DeepEqual(clientGatewayKey, precannedClientGatewayKey[i]) {
			t.Errorf("KMAC %v did not match expected:"+
				"\n  Received: %v\n  Expected: %v", i, clientGatewayKey, precannedClientGatewayKey[i])
		}
	}
}

var precannedClientGatewayKey = [][]byte{
	{126, 68, 172, 188, 215, 240, 228, 150, 51, 159, 198, 123, 195, 134, 229, 74, 16, 41, 220, 4, 244, 254, 177, 199, 172, 164, 214, 172, 41, 151, 154, 53},
	{7, 121, 99, 122, 94, 43, 110, 157, 219, 25, 169, 187, 170, 145, 3, 59, 207, 46, 145, 27, 243, 118, 110, 99, 107, 133, 141, 146, 158, 136, 7, 103},
	{150, 37, 45, 211, 117, 154, 137, 60, 177, 183, 170, 195, 98, 123, 14, 228, 235, 90, 30, 91, 123, 52, 19, 69, 34, 93, 47, 60, 49, 183, 24, 200},
	{161, 8, 197, 167, 4, 12, 155, 255, 92, 96, 57, 85, 150, 46, 115, 142, 172, 116, 149, 128, 155, 11, 42, 213, 208, 46, 41, 136, 55, 159, 33, 119},
	{115, 132, 55, 31, 75, 149, 245, 225, 113, 101, 42, 241, 153, 205, 80, 116, 205, 79, 32, 69, 201, 10, 88, 167, 110, 243, 124, 235, 89, 206, 94, 110},
	{229, 210, 144, 228, 246, 41, 189, 83, 157, 184, 182, 184, 81, 146, 117, 124, 23, 162, 252, 14, 191, 205, 240, 123, 3, 233, 63, 206, 118, 207, 215, 133},
	{220, 234, 163, 10, 33, 138, 76, 232, 62, 238, 130, 36, 228, 172, 187, 144, 146, 63, 32, 35, 213, 72, 44, 137, 181, 136, 110, 98, 13, 231, 241, 212},
	{18, 49, 27, 230, 246, 134, 255, 195, 91, 100, 162, 123, 11, 135, 174, 192, 148, 28, 39, 159, 85, 111, 203, 16, 88, 216, 188, 29, 114, 189, 48, 178},
	{213, 140, 117, 100, 32, 148, 48, 133, 163, 171, 115, 255, 227, 89, 11, 42, 146, 193, 181, 11, 136, 52, 228, 32, 239, 121, 170, 82, 61, 101, 143, 109},
	{121, 135, 154, 220, 235, 229, 113, 170, 55, 196, 111, 165, 144, 245, 244, 19, 15, 3, 11, 210, 7, 243, 44, 62, 34, 16, 1, 57, 164, 41, 103, 12},
	{230, 101, 74, 247, 199, 81, 225, 78, 2, 184, 129, 79, 23, 234, 198, 100, 112, 98, 3, 109, 46, 88, 104, 97, 196, 67, 172, 210, 165, 40, 190, 243},
	{190, 139, 47, 69, 155, 88, 98, 61, 113, 253, 55, 36, 140, 204, 115, 85, 30, 43, 179, 40, 5, 111, 36, 200, 118, 116, 222, 48, 12, 112, 99, 90},
	{75, 131, 233, 116, 72, 124, 52, 236, 165, 72, 56, 140, 250, 6, 250, 124, 238, 42, 208, 77, 180, 196, 36, 170, 58, 221, 106, 102, 157, 111, 76, 95},
	{201, 182, 11, 107, 134, 46, 135, 9, 2, 127, 24, 253, 129, 25, 65, 167, 37, 121, 12, 188, 28, 204, 40, 100, 238, 122, 158, 14, 63, 123, 110, 164},
	{249, 205, 175, 11, 155, 242, 210, 204, 15, 26, 190, 203, 229, 107, 27, 77, 107, 150, 189, 27, 235, 126, 229, 110, 8, 200, 170, 253, 19, 131, 228, 183},
	{167, 233, 206, 227, 197, 39, 242, 101, 64, 121, 46, 214, 151, 115, 199, 205, 214, 245, 30, 196, 140, 57, 18, 114, 92, 97, 3, 199, 176, 62, 25, 179},
	{110, 160, 210, 43, 219, 76, 66, 21, 87, 229, 115, 79, 184, 136, 171, 198, 136, 76, 55, 226, 98, 158, 246, 121, 234, 221, 195, 243, 5, 250, 137, 183},
	{203, 115, 87, 194, 120, 77, 192, 225, 186, 227, 67, 22, 68, 113, 220, 50, 116, 122, 201, 198, 78, 41, 54, 98, 241, 230, 134, 129, 40, 105, 172, 143},
	{159, 194, 228, 227, 177, 235, 181, 212, 155, 101, 188, 105, 243, 71, 201, 57, 225, 138, 84, 239, 137, 21, 162, 46, 219, 186, 193, 114, 189, 175, 226, 124},
	{60, 213, 160, 2, 72, 64, 143, 23, 44, 117, 84, 242, 206, 95, 183, 163, 14, 8, 252, 24, 69, 22, 99, 202, 205, 115, 138, 3, 48, 170, 64, 68},
	{224, 244, 227, 152, 214, 61, 214, 162, 213, 151, 163, 142, 123, 11, 150, 186, 192, 166, 184, 149, 228, 85, 22, 231, 128, 239, 92, 226, 244, 96, 2, 127},
	{125, 84, 179, 136, 252, 121, 242, 202, 11, 145, 178, 205, 156, 16, 14, 243, 124, 17, 40, 112, 225, 84, 157, 69, 14, 21, 55, 121, 236, 176, 18, 130},
	{209, 189, 152, 199, 243, 96, 183, 119, 160, 224, 48, 85, 85, 181, 36, 2, 95, 94, 246, 237, 127, 155, 31, 69, 11, 68, 12, 103, 77, 81, 36, 243},
	{157, 196, 32, 137, 217, 51, 110, 80, 51, 171, 229, 219, 191, 199, 70, 26, 137, 138, 0, 208, 33, 113, 50, 45, 217, 230, 62, 80, 106, 67, 200, 211},
	{201, 191, 30, 84, 101, 13, 143, 67, 98, 86, 3, 165, 98, 126, 36, 165, 90, 208, 112, 10, 111, 174, 173, 117, 63, 6, 4, 239, 254, 194, 141, 219},
	{41, 225, 229, 181, 18, 63, 29, 125, 102, 244, 194, 133, 196, 110, 186, 223, 113, 151, 49, 190, 107, 170, 154, 208, 205, 215, 245, 66, 98, 118, 91, 160},
	{170, 218, 57, 224, 53, 34, 240, 91, 216, 146, 217, 208, 89, 215, 7, 22, 147, 128, 226, 169, 155, 159, 86, 206, 128, 206, 198, 253, 218, 182, 26, 76},
	{32, 244, 239, 173, 172, 170, 171, 62, 186, 192, 2, 222, 56, 127, 216, 90, 241, 186, 231, 251, 250, 26, 127, 158, 155, 40, 157, 44, 113, 235, 244, 168},
	{9, 216, 158, 36, 226, 162, 91, 136, 126, 107, 244, 74, 191, 214, 143, 67, 74, 78, 93, 121, 87, 91, 121, 234, 215, 125, 141, 216, 224, 174, 118, 208},
	{57, 179, 231, 10, 107, 248, 151, 43, 149, 23, 155, 161, 123, 244, 224, 65, 157, 189, 249, 83, 214, 253, 53, 92, 38, 136, 128, 161, 0, 146, 211, 137},
	{238, 130, 204, 103, 115, 171, 123, 41, 25, 75, 210, 68, 213, 127, 87, 255, 120, 80, 191, 173, 86, 140, 243, 80, 10, 99, 190, 189, 21, 247, 53, 125},
	{61, 169, 201, 207, 226, 35, 199, 235, 106, 190, 100, 37, 250, 100, 107, 179, 101, 144, 239, 244, 21, 234, 152, 213, 27, 255, 95, 178, 226, 214, 88, 221},
	{26, 132, 227, 100, 63, 51, 199, 218, 148, 141, 226, 206, 168, 95, 98, 86, 102, 136, 254, 143, 82, 52, 254, 37, 167, 122, 231, 243, 113, 229, 191, 30},
	{49, 190, 117, 50, 182, 68, 46, 169, 202, 167, 165, 180, 80, 73, 92, 16, 138, 134, 179, 54, 190, 72, 31, 1, 145, 48, 67, 205, 85, 223, 31, 81},
	{171, 139, 143, 254, 195, 173, 65, 194, 247, 252, 221, 111, 185, 135, 213, 90, 49, 106, 185, 93, 127, 78, 139, 213, 155, 109, 74, 161, 23, 95, 134, 99},
	{217, 30, 14, 118, 61, 173, 44, 157, 230, 149, 34, 239, 4, 122, 204, 160, 26, 223, 107, 221, 114, 56, 148, 66, 121, 13, 179, 66, 84, 91, 252, 54},
	{64, 153, 138, 40, 92, 168, 6, 142, 245, 182, 215, 109, 228, 195, 97, 229, 63, 28, 116, 207, 230, 175, 48, 39, 246, 192, 152, 205, 235, 51, 189, 93},
	{163, 234, 202, 78, 69, 101, 129, 22, 6, 228, 168, 190, 243, 152, 231, 240, 159, 216, 92, 53, 226, 109, 160, 7, 12, 230, 28, 163, 85, 11, 141, 48},
	{158, 101, 157, 1, 17, 175, 139, 248, 4, 121, 193, 93, 180, 128, 162, 213, 115, 120, 11, 163, 33, 142, 64, 66, 253, 170, 55, 92, 52, 229, 52, 113},
	{138, 216, 252, 101, 222, 209, 120, 4, 73, 69, 240, 247, 195, 147, 77, 55, 32, 245, 132, 93, 203, 74, 227, 236, 232, 27, 247, 189, 18, 233, 12, 35},
	{24, 40, 2, 43, 205, 118, 205, 226, 87, 200, 81, 80, 108, 143, 190, 37, 90, 101, 216, 36, 250, 173, 0, 86, 94, 97, 128, 130, 68, 178, 26, 127},
	{128, 119, 111, 189, 1, 228, 20, 28, 17, 25, 93, 6, 169, 98, 220, 95, 121, 152, 215, 158, 59, 86, 15, 86, 203, 42, 149, 17, 21, 2, 49, 109},
	{81, 92, 127, 174, 219, 133, 31, 89, 232, 52, 181, 60, 255, 141, 56, 0, 155, 161, 30, 66, 67, 177, 60, 69, 195, 49, 114, 77, 123, 114, 196, 120},
	{246, 27, 39, 137, 253, 198, 103, 87, 217, 193, 232, 149, 4, 240, 3, 13, 12, 151, 66, 163, 177, 245, 130, 173, 21, 69, 76, 116, 230, 70, 212, 119},
	{102, 173, 61, 20, 76, 27, 174, 185, 78, 224, 120, 209, 97, 234, 96, 250, 214, 241, 1, 225, 242, 61, 154, 231, 8, 162, 183, 127, 130, 99, 1, 218},
	{126, 29, 125, 61, 209, 101, 232, 128, 239, 137, 205, 206, 244, 216, 13, 4, 1, 30, 166, 118, 92, 21, 76, 212, 111, 192, 126, 92, 37, 133, 0, 134},
	{229, 110, 196, 242, 150, 88, 57, 72, 41, 80, 187, 13, 153, 165, 255, 185, 196, 5, 157, 6, 43, 247, 178, 244, 155, 155, 18, 46, 222, 186, 92, 31},
	{136, 81, 105, 176, 113, 61, 136, 0, 127, 179, 158, 201, 126, 218, 254, 111, 28, 181, 207, 172, 25, 59, 181, 115, 243, 229, 195, 121, 169, 179, 219, 42},
	{105, 132, 19, 98, 105, 197, 212, 213, 25, 245, 149, 117, 105, 128, 111, 252, 54, 224, 205, 69, 228, 15, 139, 0, 160, 146, 101, 117, 69, 83, 61, 168},
	{44, 87, 217, 113, 36, 212, 28, 178, 68, 119, 44, 52, 83, 24, 79, 120, 117, 75, 118, 96, 120, 103, 161, 239, 154, 170, 10, 187, 15, 244, 135, 208},
	{67, 231, 201, 55, 231, 181, 112, 185, 29, 221, 242, 90, 58, 105, 252, 145, 22, 20, 159, 99, 33, 94, 165, 179, 171, 246, 201, 212, 83, 216, 104, 76},
	{137, 193, 64, 58, 203, 111, 71, 247, 193, 121, 134, 151, 40, 172, 29, 207, 250, 17, 247, 26, 145, 160, 13, 134, 89, 63, 175, 155, 193, 126, 34, 180},
	{217, 82, 216, 185, 189, 89, 244, 125, 213, 197, 71, 138, 193, 48, 192, 80, 135, 149, 142, 53, 36, 108, 134, 190, 83, 146, 255, 3, 220, 159, 240, 232},
	{69, 60, 193, 160, 214, 121, 84, 151, 56, 124, 110, 54, 243, 205, 39, 128, 239, 112, 204, 191, 156, 78, 0, 211, 111, 8, 169, 227, 11, 135, 128, 103},
	{77, 52, 150, 184, 176, 232, 151, 137, 208, 201, 28, 90, 16, 75, 156, 77, 91, 187, 245, 103, 167, 182, 113, 12, 31, 9, 249, 194, 184, 26, 67, 233},
	{228, 45, 242, 140, 9, 124, 142, 190, 90, 23, 67, 118, 240, 77, 196, 134, 48, 125, 102, 180, 255, 92, 10, 190, 194, 112, 182, 120, 16, 50, 153, 95},
	{193, 196, 29, 155, 209, 178, 140, 157, 89, 242, 123, 85, 162, 109, 93, 23, 74, 227, 58, 91, 119, 227, 150, 39, 192, 91, 238, 162, 161, 194, 101, 131},
	{124, 31, 202, 1, 175, 142, 91, 104, 46, 165, 64, 134, 107, 174, 89, 49, 210, 106, 30, 83, 243, 180, 25, 16, 141, 74, 111, 101, 140, 33, 30, 4},
	{117, 216, 254, 49, 126, 194, 27, 183, 98, 118, 235, 239, 191, 217, 101, 173, 255, 44, 88, 119, 134, 2, 153, 49, 161, 39, 125, 38, 132, 135, 79, 90},
	{100, 199, 230, 126, 144, 100, 181, 100, 86, 135, 110, 69, 227, 204, 1, 48, 11, 182, 233, 216, 92, 30, 121, 214, 169, 118, 179, 238, 118, 163, 71, 175},
	{40, 172, 54, 254, 40, 50, 134, 155, 139, 192, 172, 8, 135, 1, 119, 50, 118, 190, 131, 40, 196, 22, 154, 179, 170, 74, 19, 112, 28, 224, 137, 90},
	{218, 192, 99, 39, 72, 83, 238, 20, 17, 148, 243, 59, 175, 108, 246, 145, 74, 152, 5, 143, 227, 41, 220, 243, 139, 143, 157, 122, 247, 81, 158, 144},
	{164, 178, 198, 149, 12, 54, 211, 126, 17, 133, 254, 47, 168, 225, 98, 152, 31, 55, 63, 250, 178, 157, 165, 129, 106, 139, 150, 102, 183, 76, 179, 0},
	{244, 35, 220, 133, 184, 217, 103, 128, 182, 206, 120, 103, 99, 69, 112, 148, 77, 238, 188, 192, 65, 195, 252, 81, 16, 255, 22, 225, 3, 124, 243, 243},
	{142, 98, 167, 56, 70, 246, 82, 81, 254, 139, 113, 178, 83, 207, 67, 37, 195, 141, 159, 14, 247, 219, 238, 25, 222, 253, 5, 126, 54, 173, 84, 171},
	{35, 118, 3, 132, 24, 192, 101, 230, 101, 190, 28, 16, 26, 133, 232, 31, 82, 160, 247, 123, 83, 124, 107, 47, 187, 216, 140, 173, 3, 149, 110, 204},
	{1, 121, 204, 97, 16, 26, 251, 2, 88, 186, 120, 89, 244, 240, 151, 151, 75, 194, 180, 27, 236, 194, 102, 207, 167, 154, 151, 226, 144, 80, 17, 108},
	{17, 59, 210, 76, 96, 175, 54, 13, 66, 70, 174, 6, 116, 112, 3, 61, 105, 243, 57, 37, 121, 206, 217, 168, 166, 182, 49, 59, 10, 20, 170, 139},
	{98, 26, 201, 16, 86, 228, 172, 140, 204, 52, 222, 217, 195, 75, 8, 35, 155, 157, 164, 233, 138, 127, 65, 51, 210, 175, 210, 229, 190, 118, 154, 201},
	{34, 150, 18, 118, 176, 176, 201, 43, 4, 78, 31, 250, 125, 30, 161, 136, 222, 45, 22, 86, 3, 162, 31, 238, 253, 46, 49, 29, 171, 96, 5, 44},
	{4, 11, 169, 144, 199, 217, 228, 187, 100, 146, 61, 136, 250, 192, 212, 162, 32, 136, 102, 73, 22, 181, 154, 164, 152, 158, 14, 212, 114, 76, 61, 124},
	{128, 45, 97, 183, 234, 126, 157, 3, 92, 74, 25, 243, 243, 36, 84, 224, 78, 117, 83, 147, 248, 114, 80, 144, 44, 121, 5, 75, 240, 58, 94, 125},
	{93, 172, 122, 213, 178, 243, 189, 251, 252, 163, 133, 27, 57, 27, 159, 88, 17, 95, 164, 87, 203, 222, 86, 83, 42, 190, 157, 212, 174, 109, 208, 202},
	{179, 201, 218, 245, 11, 201, 221, 11, 44, 113, 169, 163, 138, 189, 114, 12, 168, 165, 39, 90, 45, 234, 17, 152, 89, 217, 177, 115, 218, 98, 5, 117},
	{30, 91, 112, 204, 100, 213, 199, 2, 55, 43, 160, 21, 167, 184, 222, 109, 156, 56, 7, 146, 58, 13, 53, 166, 63, 66, 15, 221, 46, 18, 232, 72},
	{33, 155, 27, 172, 194, 48, 175, 232, 72, 158, 174, 208, 43, 44, 75, 225, 137, 143, 66, 212, 100, 247, 99, 154, 210, 35, 35, 117, 245, 70, 197, 38},
	{106, 214, 6, 209, 127, 9, 222, 41, 52, 107, 52, 235, 161, 61, 253, 187, 191, 56, 16, 170, 36, 42, 117, 244, 106, 228, 51, 103, 127, 137, 179, 32},
	{40, 119, 232, 199, 15, 110, 57, 254, 63, 215, 173, 229, 231, 104, 252, 155, 59, 135, 5, 1, 21, 71, 222, 118, 184, 245, 221, 154, 135, 196, 130, 81},
	{63, 94, 183, 1, 144, 152, 181, 122, 252, 114, 172, 22, 155, 237, 248, 48, 90, 4, 121, 217, 76, 246, 198, 4, 108, 177, 15, 22, 124, 1, 43, 25},
	{106, 152, 34, 172, 59, 122, 101, 177, 106, 18, 75, 98, 193, 77, 119, 121, 20, 136, 84, 60, 59, 105, 21, 96, 140, 109, 27, 187, 32, 169, 90, 85},
	{85, 97, 97, 221, 67, 106, 135, 223, 14, 248, 47, 54, 152, 32, 54, 243, 115, 83, 98, 240, 229, 143, 192, 35, 184, 165, 235, 97, 109, 68, 107, 33},
	{182, 127, 92, 29, 143, 100, 63, 221, 81, 179, 165, 25, 242, 159, 168, 99, 16, 21, 92, 58, 221, 204, 113, 122, 118, 83, 136, 52, 209, 73, 159, 204},
	{224, 147, 176, 0, 88, 241, 4, 66, 190, 35, 66, 254, 248, 151, 251, 135, 35, 137, 4, 55, 150, 161, 195, 187, 92, 62, 210, 102, 36, 253, 180, 186},
	{245, 107, 114, 96, 225, 246, 148, 144, 17, 76, 72, 249, 117, 91, 17, 27, 226, 98, 52, 143, 239, 242, 211, 178, 11, 227, 115, 100, 121, 8, 45, 95},
	{227, 88, 158, 129, 44, 97, 50, 200, 127, 197, 195, 198, 163, 105, 121, 209, 8, 251, 96, 222, 170, 159, 25, 25, 109, 203, 158, 180, 159, 162, 125, 187},
	{57, 40, 171, 79, 81, 138, 254, 92, 90, 190, 51, 70, 52, 111, 237, 113, 1, 58, 92, 251, 177, 81, 211, 152, 127, 2, 109, 204, 145, 29, 47, 172},
	{226, 61, 29, 12, 210, 238, 58, 130, 30, 152, 142, 231, 38, 73, 170, 75, 141, 203, 155, 197, 91, 13, 61, 23, 71, 7, 211, 139, 241, 109, 139, 165},
	{140, 167, 245, 166, 114, 116, 48, 57, 64, 28, 71, 91, 144, 227, 59, 182, 10, 198, 86, 230, 199, 24, 4, 194, 42, 227, 219, 111, 162, 208, 78, 142},
	{122, 156, 168, 166, 22, 240, 37, 201, 244, 149, 99, 91, 118, 119, 197, 11, 143, 167, 25, 27, 254, 81, 184, 232, 176, 125, 99, 51, 131, 22, 143, 214},
	{168, 32, 176, 196, 48, 146, 165, 53, 6, 175, 231, 238, 118, 36, 34, 249, 192, 53, 1, 169, 126, 194, 228, 20, 165, 103, 74, 120, 36, 195, 5, 229},
	{81, 101, 142, 104, 80, 63, 139, 69, 129, 98, 115, 173, 159, 164, 120, 41, 124, 139, 164, 234, 97, 73, 29, 54, 250, 164, 234, 89, 205, 251, 228, 172},
	{204, 207, 135, 224, 68, 71, 61, 246, 202, 73, 184, 189, 162, 87, 185, 15, 107, 209, 27, 93, 66, 156, 106, 7, 117, 175, 243, 140, 22, 232, 99, 108},
	{164, 81, 209, 99, 189, 251, 46, 186, 48, 222, 217, 137, 115, 117, 79, 97, 8, 230, 124, 22, 109, 134, 144, 103, 252, 86, 113, 184, 254, 66, 30, 62},
	{174, 228, 215, 228, 50, 151, 218, 185, 228, 74, 24, 174, 140, 217, 202, 119, 252, 1, 100, 66, 246, 85, 226, 14, 13, 196, 212, 105, 146, 91, 127, 77},
	{222, 81, 44, 75, 27, 204, 109, 82, 248, 112, 105, 252, 224, 148, 60, 61, 94, 87, 33, 121, 80, 170, 157, 14, 140, 92, 27, 77, 215, 213, 8, 85},
	{243, 84, 29, 120, 121, 254, 206, 23, 197, 192, 117, 38, 25, 56, 115, 208, 80, 152, 143, 85, 155, 54, 9, 70, 58, 4, 10, 124, 191, 23, 106, 84},
	{23, 48, 247, 43, 99, 233, 216, 146, 114, 95, 87, 28, 226, 115, 136, 5, 0, 172, 3, 78, 175, 164, 33, 247, 189, 128, 20, 112, 92, 199, 207, 84},
	{174, 152, 176, 96, 168, 131, 11, 156, 46, 237, 82, 180, 230, 96, 109, 70, 79, 205, 254, 115, 143, 138, 61, 72, 204, 96, 11, 252, 48, 23, 2, 216},
	{19, 212, 75, 163, 189, 179, 208, 114, 85, 44, 1, 219, 110, 50, 251, 30, 41, 4, 249, 0, 186, 168, 123, 161, 40, 119, 14, 77, 137, 55, 104, 112},
	{86, 203, 16, 29, 151, 154, 97, 87, 99, 98, 20, 65, 2, 71, 128, 84, 153, 88, 163, 77, 168, 38, 173, 115, 39, 54, 144, 32, 146, 192, 137, 88},
}
