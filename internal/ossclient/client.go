package ossclient

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/linlanniao/soss/pkg/secure"
)

type Client struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	cli       *oss.Client
	securer   *securer.Securer
}

func NewClient(endpoint, accessKey, secretKey, encryptPassword string) (*Client, error) {
	c := &Client{Endpoint: endpoint, AccessKey: accessKey, SecretKey: secretKey}
	s, err := securer.NewSecurer(encryptPassword)
	if err != nil {
		return nil, err
	}
	c.securer = s

	c.cli, err = oss.New(endpoint, accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) List(bucketName string, prefix string) error {
	b, err := c.cli.Bucket(bucketName)
	if err != nil {
		return err
	}

	lsRes, err := b.ListObjects(oss.Prefix(prefix))
	if err != nil {
		return err
	}

	for _, obj := range lsRes.Objects {
		slog.Info("", "file", obj.Key)
	}
	return nil
}

type File struct {
	path    string
	name    string
	content []byte
}

func (c *Client) readFile(path string) (*File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	encrypted, err := c.securer.Encrypt(raw)
	if err != nil {
		return nil, err
	}
	return &File{
		path:    path,
		name:    filepath.Base(path),
		content: encrypted,
	}, nil
}

func (c *Client) uploadTo(bucket, prefix, name string, content []byte) error {
	b, err := c.cli.Bucket(bucket)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(content)
	err = b.PutObject(name, reader, oss.Prefix(prefix))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Upload(bucket string, prefix string, files []string) error {
	if len(files) == 0 {
		err := fmt.Errorf("no files to upload")
		slog.Error(err.Error(), "bucket", bucket, "prefix", prefix, "files", files)
		return err
	}
	for _, file := range files {
		f, err := c.readFile(file)
		if err != nil {
			slog.Error("upload failed", "file", file, "prefix", prefix, "error", err.Error())
			return err
		}
		err = c.uploadTo(bucket, prefix, f.name, f.content)
		if err != nil {
			slog.Error("upload failed", "file", file, "prefix", prefix, "error", err.Error())
			return err
		}
		slog.Info("upload succeed", "file", f.name)
	}
	return nil
}

func (c *Client) Download(bucket string, files []string, outputDir string) error {
	//TODO implement me
	panic("implement me")
}
