package tss

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/antik10ud/go-comb/comb"
	"testing"
)

func TestKat(t *testing.T) {
	secret, _ := hex.DecodeString("a217525ab5cab096e455acba00a4032c0cc1a1ef7ccd280642d994cdee7694ca")
	shares := make([][]byte, 2)
	shares[0], _ = hex.DecodeString("016fbb11a9264bdf4188e3911827ea30e27a283cbea177d7c421524fb448bbfedc")
	shares[1], _ = hex.DecodeString("022354d4a788d36e233c22d6e54e3865abe008804ddda2cd9984d4393fb9f740e6")
	testRecover(t, secret, ShareSet{shares[0], shares[1]})
}

func TestCreateSharesErrors(t *testing.T) {
	testCaseCreateExpect(t, 0, 2, 3, ErrSecretRequired)
	testCaseCreateExpect(t, 32, 2, 3, ErrInvalidThreshold)
	testCaseCreateExpect(t, 1, 2, 2, nil)
	//testCaseCreateExpect(t, MinSecretBytes-1, 2, 2, ErrSecretTooShort)
	testCaseCreateExpect(t, MaxSecretBytes+1, 2, 2, ErrSecretTooLarge)
	testCaseCreateExpect(t, 32, MaxShares+1, 3, ErrTooManyShares)
	testCaseCreateExpect(t, 32, MinShares-1, 3, ErrTooFewShares)
}

func TestCaseCreateThreshold1(t *testing.T) {
	secret := randomBytes(32)
	_, err := CreateShares(secret, 3, 1)
	if err != ErrInvalidThreshold {
		failNow(t, err)
	}
}

func TestCaseCreateThresholdToomany(t *testing.T) {
	secret := randomBytes(32)
	_, err := CreateShares(secret, 3, 4)
	if err != ErrInvalidThreshold {
		failNow(t, err)
	}
}

func dump(secret []byte, shares ShareSet) {
	println(hex.EncodeToString(secret))
	for _, x := range shares {
		println(hex.EncodeToString(x))
	}

}

func testCaseCreateExpect(t *testing.T, secretSize int, sharesCount int, threshold int, expect error) {
	secret := randomBytes(secretSize)
	_, err := CreateShares(secret, sharesCount, threshold)
	if err != expect {
		failNow(t, expected(expect, err))
	}
}

func TestRecoverErrors(t *testing.T) {
	secret := randomBytes(32)
	shares, _ := CreateShares(secret, 10, 2)
	testCaseRecoverExpect(t, ShareSet{shares[0]}, ErrTooFewShares)
	testCaseRecoverExpect(t, ShareSet{shares[0], randomBytes(32)}, ErrInvalidShare)
	testCaseRecoverExpect(t, ShareSet{randomBytes(MinShareBytes - 1), randomBytes(MinShareBytes - 1)}, ErrInvalidShare)
	testCaseRecoverExpect(t, ShareSet{randomBytes(MaxShareBytes + 1), randomBytes(MaxShareBytes + 1)}, ErrInvalidShare)
	testCaseRecoverExpect(t, ShareSet{shares[0], randomBytes(MaxSecretBytes + 1)}, ErrInvalidShare)
	ss := make(ShareSet, MaxShares+1)
	for i := 0; i < MaxShares+1; i++ {
		ss[i] = randomBytes(32)
	}
	testCaseRecoverExpect(t, ss, ErrTooManyShares)

}

func testCaseRecoverExpect(t *testing.T, shares ShareSet, expect error) {
	_, err := RecoverSecret(shares)
	if err != expect {
		failNow(t, expected(expect, err))
	}
}

func TestNM(t *testing.T) {
	for i := 2; i < MaxShares; i += 77 {
		for j := 2; j < min(i, 5); j++ {
			t.Run(fmt.Sprintf("%d-%d", i, j), func(t *testing.T) {
				testCase(t, 32, i, j)
			})

		}
	}
}

func testCase(t *testing.T, secretSize int, sharesCount int, threshold int) {
	secret := randomBytes(secretSize)
	shares, err := CreateShares(secret, sharesCount, threshold)
	if err != nil {
		failNow(t, err)
	}
	testRecover(t, secret, shares)
}

func TestSharesComb(t *testing.T) {
	for i := 3; i < MaxShares; i += 97 {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			testSharesComb(t, i, 3)
		})
	}
}

func BenchmarkCreateShares(b *testing.B) {
	secret := randomBytes(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateShares(secret, MaxShares, MaxShares/2)
	}
}
func BenchmarkRecoverSecret(b *testing.B) {
	secret := randomBytes(32)
	shares, err := CreateShares(secret, MaxShares, MaxShares/2)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RecoverSecret(shares)
	}
}

func ExampleReadme() {

	secret, _ := hex.DecodeString("05cd605252528ab7302ca970c56ef99897cb6c4230e1cebf24516b4f7a9248c1")

	sharesCount := 5 // number of shares

	threshold := 3 // number of required shares to recover the secret

	shares, _ := CreateShares(secret, sharesCount, threshold)

	recoveredSecret, _ := RecoverSecret(ShareSet{shares[0], shares[1], shares[4]})

	fmt.Println(hex.EncodeToString(recoveredSecret))

	// Output:
	// 05cd605252528ab7302ca970c56ef99897cb6c4230e1cebf24516b4f7a9248c1
}

func testSharesComb(t *testing.T, sharesCount int, threshold int) {
	secret := randomBytes(32)
	shares, err := CreateShares(secret, sharesCount, threshold)
	if err != nil {
		failNow(t, err)
	}
	cmb, err := comb.NewNoRepLex(sharesCount, threshold)
	if err != nil {
		failNow(t, err)
	}
	for {
		v := cmb.Next()
		if v == nil {
			break
		}
		testRecover(t, secret, subShares(shares, *v))
	}

}

func testRecover(t *testing.T, secret []byte, shares ShareSet) {
	recoveredSecret, err := RecoverSecret(shares)
	if err != nil {
		failNow(t, err)
	}
	if bytes.Compare(recoveredSecret, secret) != 0 {
		failNow(t, fmt.Errorf("secret mismatch %x, want %x", recoveredSecret, secret))
	}
}

func subShares(shares ShareSet, pickList []int) ShareSet {
	newShares := make(ShareSet, len(pickList))
	for i, pick := range pickList {
		newShares[i] = shares[pick]
	}
	return newShares
}

func expected(exp error, actual error) error {
	return fmt.Errorf("err '%s' but expected '%s'", actual, exp)
}

func failNow(t *testing.T, err error) {
	t.Error(err)
	t.FailNow()
}

func min(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}

func randomBytes(i int) []byte {
	v := make([]byte, i)
	rand.Read(v)
	return v
}
