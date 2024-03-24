package securer_test

import (
	"testing"

	securer "github.com/linlanniao/soss/pkg/secure"
	"github.com/stretchr/testify/assert"
)

func TestSecurer_Encrypt_Decrypt(t *testing.T) {
	key := ""
	s, _ := securer.NewSecurer(key)

	raw := "hello world"
	encrypted, err := s.Encrypt([]byte(raw))
	assert.NoError(t, err)
	decrypted, err := s.Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, raw, string(decrypted))
}

func BenchmarkSecurer_Encrypt_Decrypt(b *testing.B) {
	key := "abcdefg"
	s, _ := securer.NewSecurer(key)

	raw := "FRIY4HIZWA3OSFTG4N0DRW0EDA4VSWQI474HVMFKYS5KUV9BQAGM3J8VLSDPX9CN78430XWNKIXCBSODTBW2L9EXK1WLHUOCPPKHQMHNHTU6OUJE9A9ET5W1S5AY81EP"
	for i := 0; i < b.N; i++ {
		encrypted, _ := s.Encrypt([]byte(raw))
		raw2, _ := s.Decrypt(encrypted)
		assert.Equal(b, raw, string(raw2))
	}
}
