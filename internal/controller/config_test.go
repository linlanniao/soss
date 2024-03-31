package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestController_LoadConfig(t *testing.T) {
	//f := "../../config.yaml"
	cfg, err := NewConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.AccessKey)
	assert.NotEmpty(t, cfg.SecretKey)
	assert.NotEmpty(t, cfg.ClientType)
	assert.NotEmpty(t, cfg.Endpoint)
	assert.NotEmpty(t, cfg.Bucket)
}
