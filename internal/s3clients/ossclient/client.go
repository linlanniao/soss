package ossclient

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/linlanniao/soss/internal"
)

type client struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	_cli      *oss.Client
	_bucket   *oss.Bucket
}

var _ internal.IS3Client = (*client)(nil)

func NewClient(endpoint, accessKey, secretKey string) internal.IS3Client {
	c := &client{Endpoint: endpoint, AccessKey: accessKey, SecretKey: secretKey}

	cli, err := oss.New(endpoint, accessKey, secretKey)
	if err != nil {
		panic(err.Error())
	}
	c._cli = cli

	return c
}

func (c *client) SetEndpoint(endpoint string) error {

	// same endpoint, do nothing
	if endpoint == c.Endpoint {
		return nil
	}

	// update endpoint
	c.Endpoint = endpoint

	// create new oss client
	if cli, err := oss.New(endpoint, c.AccessKey, c.SecretKey); err != nil {
		return err
	} else {
		c._cli = cli
		return nil
	}
}

func (c *client) SetBucket(bucket string) error {
	// update bucket
	_, err := c.bucket(bucket)
	return err
}

func (c *client) bucket(bucket string) (*oss.Bucket, error) {
	if c._bucket != nil && c._bucket.BucketName == bucket {
		return c._bucket, nil
	}

	b, err := c._cli.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	c._bucket = b
	return b, nil
}

func (c *client) List(bucket string, prefix string) (objs []*internal.S3Object, err error) {
	b, err := c._cli.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	lsRes, err := b.ListObjects(oss.Prefix(prefix))
	if err != nil {
		return nil, err
	}

	objs = make([]*internal.S3Object, len(lsRes.Objects))
	for i, obj := range lsRes.Objects {
		objs[i] = &internal.S3Object{
			Key:  obj.Key,
			Type: obj.Type,
			Size: obj.Size,
			ETag: obj.ETag,
		}
	}
	return objs, nil
}

func (c *client) Upload(bucket string, prefix string, file *internal.File) (obj *internal.S3Object, err error) {
	if file == nil {
		return nil, errors.New("file is nil")
	}

	b, err := c.bucket(bucket)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(file.Path)
	// TODO: what to deal with the prefix?
	//  1. filepath.Join ?
	//  2. prefix + filepath.Dir ?
	//key := filepath.Join(prefix, fileName)
	key := prefix + fileName
	reader := bytes.NewReader(file.Content)
	err = b.PutObject(key, reader, oss.Prefix(prefix))
	if err != nil {
		return nil, err
	}
	header, err := b.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, err
	}
	cType := header.Get("Content-Type")
	eTag := header.Get("ETag")
	sizeStr := header.Get("Content-Length")

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, errors.New("failed to parse size")
	}

	return &internal.S3Object{
		Bucket: bucket,
		Key:    key,
		Type:   cType,
		Size:   int64(size),
		ETag:   eTag,
	}, nil
}

func (c *client) Download(obj *internal.S3Object, outputDir string) (file *internal.File, err error) {
	if obj == nil {
		return nil, errors.New("obj is nil")
	}

	if obj.Bucket == "" {
		return nil, errors.New("bucket is empty")
	}

	b, err := c.bucket(obj.Bucket)
	if err != nil {
		return nil, err
	}

	reader, err := b.GetObject(obj.Key)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	absDir, _ := filepath.Abs(filepath.Join(outputDir, obj.Key))

	return &internal.File{
		Path:      absDir,
		Content:   content,
		Encrypted: true, // 默认已加密
	}, nil
}
