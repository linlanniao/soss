package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/linlanniao/soss/internal"
	"github.com/lmittmann/tint"
)

type Controller struct {
	clients     S3Clients
	fileHandler internal.IFileHandler
	endpoint    string
	bucket      string
	logger      *slog.Logger
	isCompress  bool
}

type Option func(c *Controller)

func WithLogger(logger *slog.Logger) Option {
	return func(c *Controller) {
		c.logger = logger
	}
}

func WithDefaultLogger() Option {
	return func(c *Controller) {
		w := os.Stderr

		// create a new logger
		logger := slog.New(tint.NewHandler(w, nil))

		// set global logger with custom options
		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.DateTime,
			}),
		))
		c.logger = logger
	}
}

func WithEndpoint(endPoint string) Option {
	return func(c *Controller) {
		c.endpoint = endPoint
	}
}

func WithBucket(bucket string) Option {
	return func(c *Controller) {
		c.bucket = bucket
	}
}

func WithCompression() Option {
	return func(c *Controller) {
		c.isCompress = true
	}
}

func WithS3Client(key S3ClientType, client internal.IS3Client) Option {
	return func(c *Controller) {
		if len(c.clients) == 0 {
			c.clients = make(S3Clients)
		}
		c.clients[key] = client
	}
}
func WithFileHandler(handler internal.IFileHandler) Option {
	return func(c *Controller) {
		c.fileHandler = handler
	}
}

type S3ClientType string

const (
	S3ClientTypeOSS S3ClientType = "oss"
)

func (t S3ClientType) Validate() error {
	switch t {
	case S3ClientTypeOSS:
		return nil
	default:
		return errors.New("invalid client type")
	}
}

type S3Clients map[S3ClientType]internal.IS3Client

func (c *Controller) getClient(cType S3ClientType) (internal.IS3Client, error) {
	clients := c.clients
	if len(clients) == 0 {
		return nil, errors.New("no s3 clients configured")
	}

	client, ok := clients[cType]
	if !ok {
		return nil, fmt.Errorf("no client with type %s", cType)
	}
	return client, nil
}

func NewController(opts ...Option) *Controller {
	c := &Controller{}

	for _, opt := range opts {
		opt(c)
	}

	if c.logger == nil {
		WithDefaultLogger()(c)
	}

	if c.fileHandler == nil {
		panic("fileHandler is nil")
	}
	if len(c.clients) == 0 {
		panic("clients is empty")
	}

	return c
}

type ListOptions struct {
	S3ClientType S3ClientType
	Endpoint     string
	Bucket       string
	Prefix       string
}

func (c *Controller) List(opts ListOptions) error {

	if opts.Endpoint != "" {
		c.endpoint = opts.Endpoint
	}
	if opts.Bucket != "" {
		c.bucket = opts.Bucket
	}

	client, err := c.getClient(opts.S3ClientType)
	if err != nil {
		return err
	}

	objs, err := client.List(c.endpoint, c.bucket, opts.Prefix)

	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	for _, obj := range objs {
		fmt.Println(obj.Key)
	}
	return nil
}

func trimDirectory(path, file string) string {
	// Normalize paths to ensure correct path separators
	path = filepath.Clean(path)
	file = filepath.Clean(file)

	// Check if file is within the directory specified by path
	if !strings.HasPrefix(file, path) {
		return file // If not, return file as is
	}

	// Get the length of the path directory
	pathLength := len(path)

	// Trim the path part from file
	trimmedPath := file[pathLength:]

	// Remove any leading path separator that might exist
	trimmedPath = strings.TrimPrefix(trimmedPath, string(filepath.Separator))

	return trimmedPath
}

func (c *Controller) uploadSingleFile(endpoint, bucket, prefix, path, encryptKey string, client internal.IS3Client) error {
	file, err := c.fileHandler.Read(path)
	if err != nil {
		c.logger.Error("upload failed", "err", err.Error())
		return err
	}

	// compress file content
	if c.isCompress {
		if err := c.fileHandler.Compress(file); err != nil {
			c.logger.Error("compress file failed", "err", err.Error())
			return err
		}
	}

	// encrypt file content
	if err := c.fileHandler.Encrypt(file, encryptKey); err != nil {
		c.logger.Error("encrypt file failed", "err", err.Error())
		return err
	}

	obj, err := client.Upload(endpoint, bucket, prefix, file)
	if err != nil {
		c.logger.Error("upload failed", "err", err.Error())
		return err
	}
	//c.logger.Info(fmt.Sprintf("uploading %s to %s", file, ossFileKey))

	//if !filepath.IsAbs(file.Path) {
	//	file.Path = "./" + file.Path
	//}
	c.logger.Info("uploading",
		"from", file.Path,
		"to", obj.Bucket+":"+obj.Key,
		"size(bytes)", obj.Size,
	)
	return nil
}

const (
	uploadParallelism   = 10
	downloadParallelism = 10
)

func (c *Controller) UploadDirectoryOrFile(
	endpoint, bucket, prefix, path, encryptKey string, client internal.IS3Client) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.logger.Error("does not exist", "path", path)
		} else {
			c.logger.Error(err.Error())
		}
		return err
	}

	if fileInfo.IsDir() {
		files, err := c.fileHandler.SearchFiles(path)
		if err != nil {
			c.logger.Error(err.Error())
			return err
		}

		var wg sync.WaitGroup
		//limiter := make(chan struct{}, uploadParallelism)
		limiter := make(chan struct{}, runtime.NumCPU()*2)

		for _, file := range files {
			wg.Add(1)
			limiter <- struct{}{} // Take up a concurrent signal

			go func(file string) {
				defer func() {
					<-limiter // Release a concurrent signal
					wg.Done()
				}()

				subPrefix := filepath.Join(prefix, filepath.Dir(trimDirectory(path, file)))
				if err := c.uploadSingleFile(endpoint, bucket, subPrefix, file, encryptKey, client); err != nil {
					c.logger.Error("error uploading file", "file", file, "error", err.Error())
				}
			}(file)
		}

		wg.Wait()
	} else {
		if err := c.uploadSingleFile(endpoint, bucket, prefix, path, encryptKey, client); err != nil {
			return err
		}
	}

	return nil
}

type UploadOptions struct {
	S3ClientType S3ClientType
	Endpoint     string
	Bucket       string
	Prefix       string
	EncryptKey   string
	Paths        []string
}

func (c *Controller) Upload(opts UploadOptions) error {
	if opts.Endpoint != "" {
		c.endpoint = opts.Endpoint
	}
	if opts.Bucket != "" {
		c.bucket = opts.Bucket
	}

	client, err := c.getClient(opts.S3ClientType)
	if err != nil {
		return err
	}

	if len(opts.Paths) == 0 {
		err := errors.New("no files to upload")
		c.logger.Error("upload failed", "err", err.Error())
		return err
	}

	for _, path := range opts.Paths {
		if err := c.UploadDirectoryOrFile(c.endpoint, c.bucket, opts.Prefix, path, opts.EncryptKey, client); err != nil {
			c.logger.Error("upload failed", "err", err.Error())
			return err
		}
	}

	return nil
}

func (c *Controller) downloadSingleFile(
	endpoint, bucket, s3key, outputDir, decryptKey string, client internal.IS3Client) error {
	file, err := client.Download(
		&internal.S3Object{
			Endpoint: endpoint,
			Bucket:   bucket,
			Key:      s3key,
		},
		outputDir,
	)
	if err != nil {
		c.logger.Error("download failed", "key", s3key, "err", err.Error())
		return err
	}

	// decrypt file content
	if err := c.fileHandler.Decrypt(file, decryptKey); err != nil {
		c.logger.Error("decrypt file failed", "err", err.Error())
		return err
	}

	if c.isCompress {
		// decompress file content
		if err := c.fileHandler.Decompress(file); err != nil {
			c.logger.Error("decompress file failed", "key", s3key, "err", err.Error())
			return err
		}
	}

	if err := c.fileHandler.Write(file); err != nil {
		c.logger.Error("decrypt file failed", "err", err.Error())
		return err
	}

	//if !filepath.IsAbs(file.Path) {
	//	file.Path = "./" + file.Path
	//}

	c.logger.Info("downloading",
		"from", bucket+":"+s3key,
		"saveTo", file.Path,
		"size(bytes)", len(file.Content),
	)
	return nil
}

func (c *Controller) downloadDirectoryOrFile(
	endpoint, bucket, s3key, outputDir, decryptKey string, client internal.IS3Client) error {
	objs, err := client.List(endpoint, bucket, s3key)
	if err != nil {
		c.logger.Error("download directory or file failed", "key", s3key, "err", err.Error())
		return err
	}
	if len(objs) == 0 {
		err = errors.New("directory or file not found")
		c.logger.Error("download directory or file failed", "key", s3key, "err", err.Error())
	}
	var wg sync.WaitGroup
	//limiter := make(chan struct{}, downloadParallelism)
	limiter := make(chan struct{}, runtime.NumCPU()*2)

	for _, obj := range objs {
		wg.Add(1)
		limiter <- struct{}{} // Take up a concurrent signal

		go func(obj *internal.S3Object) {
			defer func() {
				<-limiter // Release a concurrent signal
				wg.Done()
			}()
			if err := c.downloadSingleFile(endpoint, bucket, obj.Key, outputDir, decryptKey, client); err != nil {
				c.logger.Error("download directory or file failed", "key", obj.Key, "err", err.Error())
			}
		}(obj)
	}
	wg.Wait()
	return nil
}

type DownloadOptions struct {
	S3ClientType S3ClientType
	Endpoint     string
	Bucket       string
	OutputDir    string
	DecryptKey   string
	S3keys       []string
}

func (c *Controller) Download(opts DownloadOptions) error {
	if opts.Endpoint != "" {
		c.endpoint = opts.Endpoint
	}
	if opts.Bucket != "" {
		c.bucket = opts.Bucket
	}

	client, err := c.getClient(opts.S3ClientType)
	if err != nil {
		return err
	}

	if len(opts.S3keys) == 0 {
		err := errors.New("no files to download")
		c.logger.Error(err.Error())
		return err
	}

	for _, s3key := range opts.S3keys {
		if err := c.downloadDirectoryOrFile(c.endpoint, c.bucket, s3key, opts.OutputDir, opts.DecryptKey, client); err != nil {
			return err
		}
	}
	return nil
}
