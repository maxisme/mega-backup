package mega

import (
	"testing"
)

func TestEncryptCompressDir(t *testing.T) {
	testFile := "file.tar.gz.encr"
	key := "123"
	err := EncryptCompressDir("./test/", testFile, key)
	if err != nil {
		t.Error(err.Error())
	}

	err = DecryptTar(testFile, "unencrypted.tar.gz", key)
	if err != nil {
		t.Error(err.Error())
	}

}
