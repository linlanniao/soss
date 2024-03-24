package securer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

type Securer struct {
	secretKey []byte
	block     cipher.Block
	gcm       cipher.AEAD
}

func NewSecurer(password string) (*Securer, error) {
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

	return &Securer{
		secretKey: key,
		block:     block,
		gcm:       gcm,
	}, nil
}

func (s *Securer) initBlock() error {
	if s.block != nil {
		return nil
	}
	k := sha256.Sum256(s.secretKey)
	b, err := aes.NewCipher(k[:])
	if err != nil {
		return err
	}
	s.block = b

	return nil
}

func (s *Securer) initGcm() error {
	if s.gcm != nil {
		return nil
	}
	if err := s.initBlock(); err != nil {
		return err
	}

	g, err := cipher.NewGCM(s.block)
	if err != nil {
		return err
	}
	s.gcm = g

	return nil
}

func (s *Securer) Encrypt(plaintext []byte) ([]byte, error) {
	if err := s.initGcm(); err != nil {
		return nil, err
	}
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return s.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (s *Securer) Decrypt(ciphertext []byte) ([]byte, error) {
	if err := s.initGcm(); err != nil {
		return nil, err
	}

	if len(ciphertext) < s.gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:s.gcm.NonceSize()]
	ciphertext = ciphertext[s.gcm.NonceSize():]

	return s.gcm.Open(nil, nonce, ciphertext, nil)
}
