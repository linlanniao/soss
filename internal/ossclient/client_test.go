package ossclient_test

import (
	"os"
	"testing"

	"github.com/linlanniao/soss/internal/ossclient"
	"github.com/stretchr/testify/assert"
)

const (
	testEndpoint   = "https://oss-cn-guangzhou.aliyuncs.com"
	testBucket     = "ppops-bucket"
	testEncryptKey = "p@ssw*rd"
)

var (
	testAccessKey = os.Getenv("OSS_ACCESS_KEY_ID")
	testSecretKey = os.Getenv("OSS_ACCESS_KEY_SECRET")
)

func TestClient_List(t *testing.T) {
	client, err := ossclient.NewClient(testEndpoint, testAccessKey, testSecretKey, testEncryptKey)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	err = client.List(testBucket, "")
	assert.NoError(t, err)
}

func TestClient_Upload(t *testing.T) {
	client, err := ossclient.NewClient(testEndpoint, testAccessKey, testSecretKey, testEncryptKey)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	err = client.Upload(testBucket, "", []string{"../../README.md"})
	assert.NoError(t, err)
}
