package controller_test

import (
	"testing"

	"github.com/linlanniao/soss/internal/controller"
	"github.com/stretchr/testify/assert"
)

func TestController_LoadConfig(t *testing.T) {
	filePath := "../../config.yaml"
	c := controller.Controller{}
	err := c.LoadConfig(filePath)
	assert.NoError(t, err)
}
