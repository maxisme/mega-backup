package backup

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Hash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, key string) ([]byte, error) {
	if key == "" {
		return nil, fmt.Errorf("must enter key")
	}
	block, _ := aes.NewCipher([]byte(Hash(key)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(Hash(key)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func Tar(src string, writer io.Writer) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("unable to tar files - %v", err.Error())
	}

	gzw, _ := gzip.NewWriterLevel(writer, gzip.BestSpeed)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		f.Close()

		return nil
	})
}

func WriteFile(bytes []byte, out string) error {
	file, err := os.OpenFile(
		out,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0600,
	)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	return err
}
