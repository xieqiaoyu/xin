package hashid

import (
	"reflect"
	"testing"
)

func TestConsistentShuffle(t *testing.T) {
	alphabet := []byte(DefaultAlphabet)
	salt := []byte("abc1")
	consistentShuffle(alphabet, salt)

	if string(alphabet) != "cUpI6isqCa0brWZnJA8wNTzDHEtLXOYgh5fQm2uRj4deM91oB7FkSGKxvyVP3l" {
		t.Fatalf("shuffle error with salt %s", string(salt))
	}
	t.Logf("%v", string(DefaultAlphabet))

}

func TestHashAndUnHash(t *testing.T) {
	alphabet := []byte(DefaultAlphabet)
	hashNumber := 1234567
	hashResult, err := hash(hashNumber, alphabet)
	if err != nil {
		t.Fatalf("hash error : %s", err)
	}
	t.Logf("hash result: %s", string(hashResult))
	unhashNumber, err := unhash(hashResult, alphabet)
	if err != nil {
		t.Fatalf("unhash error: %s", err)
	}
	if unhashNumber != hashNumber {
		t.Fatalf("unhash result:%d is not hash number :%d", unhashNumber, hashNumber)
	}
}

func TestHasher(t *testing.T) {
	hasher, err := NewDefaultHasher([]byte("abc"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	numbers := []int{1234, 567}
	result, err := hasher.Encode(numbers)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s", result)
	checkNumbers, err := hasher.Decode(result)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !reflect.DeepEqual(checkNumbers, numbers) {
		t.Fatalf("Mismatch Decode (%v) and Encode (%v)", checkNumbers, numbers)
	}

}
