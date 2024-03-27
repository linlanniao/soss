package cipher_test

import (
	"testing"

	securer "github.com/linlanniao/soss/pkg/cipher"
	"github.com/stretchr/testify/assert"
)

func TestSecurer_Encrypt_Decrypt(t *testing.T) {
	key := ""
	cipher, _ := securer.NewContentCipher(key)

	raw := "hello world"
	encrypted, err := cipher.EncryptBytes([]byte(raw))
	assert.NoError(t, err)
	decrypted, err := cipher.DecryptBytes(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, raw, string(decrypted))
}

func BenchmarkSecurer_Encrypt_Decrypt(b *testing.B) {
	key := "abcdefg"
	cipher, _ := securer.NewContentCipher(key)

	raw := "FRIY4HIZWA3OSFTG4N0DRW0EDA4VSWQI474HVMFKYS5KUV9BQAGM3J8VLSDPX9CN78430XWNKIXCBSODTBW2L9EXK1WLHUOCPPKHQMHNHTU6OUJE9A9ET5W1S5AY81EP"
	for i := 0; i < b.N; i++ {
		encrypted, _ := cipher.EncryptBytes([]byte(raw))
		raw2, _ := cipher.DecryptBytes(encrypted)
		assert.Equal(b, raw, string(raw2))
	}
}
