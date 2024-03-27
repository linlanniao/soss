package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

type ContentCipher struct {
	key      string
	bytesKey []byte
	block    cipher.Block
	gcm      cipher.AEAD
}

func NewContentCipher(password string) (*ContentCipher, error) {
	key := []byte(password)
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &ContentCipher{
		key:      password,
		bytesKey: key,
		block:    block,
		gcm:      gcm,
	}, nil
}

func (c *ContentCipher) KeyIsEqual(k string) bool {
	return c.key == k
}

func (c *ContentCipher) initBlock() error {
	if c.block != nil {
		return nil
	}
	k := sha256.Sum256(c.bytesKey)
	b, err := aes.NewCipher(k[:])
	if err != nil {
		return err
	}
	c.block = b

	return nil
}

func (c *ContentCipher) initGcm() error {
	if c.gcm != nil {
		return nil
	}
	if err := c.initBlock(); err != nil {
		return err
	}

	g, err := cipher.NewGCM(c.block)
	if err != nil {
		return err
	}
	c.gcm = g

	return nil
}

func (c *ContentCipher) EncryptBytes(plainBytes []byte) ([]byte, error) {
	if err := c.initGcm(); err != nil {
		return nil, err
	}
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return c.gcm.Seal(nonce, nonce, plainBytes, nil), nil
}

func (c *ContentCipher) Encrypt(writer io.Writer) ([]byte, error) {
	plainBytes := make([]byte, 0)
	if _, err := writer.Write(plainBytes); err != nil {
		return nil, err
	}
	return c.EncryptBytes(plainBytes)
}

func (c *ContentCipher) DecryptBytes(cipherBytes []byte) ([]byte, error) {
	if err := c.initGcm(); err != nil {
		return nil, err
	}

	if len(cipherBytes) < c.gcm.NonceSize() {
		return nil, errors.New("cipherBytes too short")
	}

	nonce := cipherBytes[:c.gcm.NonceSize()]
	cipherBytes = cipherBytes[c.gcm.NonceSize():]

	return c.gcm.Open(nil, nonce, cipherBytes, nil)
}

func (c *ContentCipher) Decrypt(reader io.Reader) ([]byte, error) {
	cipherBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return c.DecryptBytes(cipherBytes)
}
