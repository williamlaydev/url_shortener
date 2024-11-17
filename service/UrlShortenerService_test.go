package service

import "testing"

// Should be able to generate a random string
func TestGenerateRandomString_Success(t *testing.T) {
	rand, err := generateRandomString()

	if err != nil {
		t.Fatalf("test failed: %v", err)
	}

	if len(rand) != 6 {
		t.Fatalf("test failed: length of random string is not 6")
	}
}
