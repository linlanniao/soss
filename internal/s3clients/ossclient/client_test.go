package ossclient

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/linlanniao/soss/internal"
	"github.com/stretchr/testify/assert"
)

const (
	testEndpoint = "https://oss-cn-guangzhou.aliyuncs.com"
	testBucket   = "ppops-bucket"
)

var (
	testAccessKey = os.Getenv("OSS_ACCESS_KEY_ID")
	testSecretKey = os.Getenv("OSS_ACCESS_KEY_SECRET")
)

func newTestClient() *client {
	x := NewClient(testEndpoint, testAccessKey, testSecretKey)
	return x.(*client)
}

func deleteObject(client *client, objectName string) error {
	bucket, err := client.bucket(testBucket)
	if err != nil {
		return err
	}
	return bucket.DeleteObject(objectName)
}

func TestClient_List(t *testing.T) {
	client := newTestClient()
	objs, err := client.List(testBucket, "tester3")
	assert.NoError(t, err)
	assert.NotEmpty(t, objs)
	for _, obj := range objs {
		t.Logf("key: %s, size: %d", obj.Key, obj.Size)
	}
}

func TestClient_Upload(t *testing.T) {
	client := newTestClient()
	file := &internal.File{
		Path:      "xx/bb/cc/TestClient_Upload.txt",
		Content:   []byte("iam test file"),
		Encrypted: true, //fake
	}
	prefix := "tester"
	obj, err := client.Upload(testBucket, prefix, file)
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	t.Logf("obj: %+v", obj)

	defer deleteObject(client, obj.Key) // clean up

	download := func() {
		// using oss client to download test file
		bucket, err := client.bucket(obj.Bucket)
		assert.NoError(t, err)
		reader, err := bucket.GetObject(obj.Key)
		assert.NoError(t, err)
		defer reader.Close()
		content, _ := io.ReadAll(reader)
		assert.Equal(t, file.Content, content)
	}
	download()
}

func TestClient_Download(t *testing.T) {
	client := newTestClient()
	file := &internal.File{
		Path:    "xx/bb/cc/TestClient_Download.txt",
		Content: []byte("iam test file"),
	}
	prefix := "tester/"

	upload := func() *internal.S3Object {
		// using oss client to upload test file
		bucket, err := client.bucket(testBucket)
		assert.NoError(t, err)

		key := prefix + filepath.Base(file.Path)
		reader := bytes.NewReader(file.Content)

		_ = bucket.PutObject(key, reader)

		return &internal.S3Object{
			Bucket: testBucket,
			Key:    key,
			Type:   "",
			Size:   int64(len(file.Content)),
			ETag:   "",
		}
	}

	obj := upload()

	//defer deleteObject(client, obj.Key) // clean up

	file2, err := client.Download(obj, "/tmp")
	assert.NoError(t, err)
	assert.NotNil(t, file2)
	assert.Equal(t, file.Content, file2.Content)
	assert.True(t, file2.Encrypted)
}
