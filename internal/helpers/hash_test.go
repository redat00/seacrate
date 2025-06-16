package helpers

import (
	"encoding/hex"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	passw := []byte("my_super_password")
	_, salt, err := GenerateHash(passw, []byte{})
	if err != nil {
		t.Errorf("an error happened during the hash generation : %s", err.Error())
	}
	if len(salt) == 0 {
		t.Errorf("size of the salt is equal to zero")
	}
}

type CompareCase struct {
	Password []byte
	Salt     string
	Hash     string
	Expected bool
}

func TestCompareHash(t *testing.T) {
	cases := []CompareCase{
		{
			Password: []byte("my_super_password"),
			Salt:     "f67ba36ed692c4f41b77d6a88e46a72669e56816b1ba10b7cc0242619a67d6f0",
			Hash:     "b03dd3c4675694224b5f4d4884647a8465631947463b870576566b1ad0d1c5a1",
			Expected: true,
		},
		{
			Password: []byte("wrong!password"),
			Salt:     "f67ba36ed692c4f41b77d6a88e46a72669e56816b1ba10b7cc0242619a67d6f0",
			Hash:     "b03dd3c4675694224b5f4d4884647a8465631947463b870576566b1ad0d1c5a1",
			Expected: false,
		},
	}

	for _, c := range cases {
		decodedSalt, err := hex.DecodeString(c.Salt)
		if err != nil {
			t.Errorf("error during decoding of the salt : %s", err.Error())
		}
		decodedHash, err := hex.DecodeString(c.Hash)
		if err != nil {
			t.Errorf("error during decoding of the hash : %s", err.Error())
		}
		result, err := Compare(c.Password, decodedHash, decodedSalt)
		if result != c.Expected {
			t.Errorf("error during password comparison. Is: %+v. Should be: %+v", result, c.Expected)
		}

	}
}
