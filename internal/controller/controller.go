package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/linlanniao/soss/internal"
	"github.com/lmittmann/tint"
)

type Controller struct {
	s3Client    internal.IS3Client
	fileHandler internal.IFileHandler
	config      *Config
	logger      *slog.Logger
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

func NewController(config *Config, client internal.IS3Client, fileHandler internal.IFileHandler, opts ...Option) *Controller {
	c := &Controller{s3Client: client, fileHandler: fileHandler, config: config}

	for _, opt := range opts {
		opt(c)
	}

	if c.logger == nil {
		WithDefaultLogger()(c)
	}

	return c
}

type ListOptions struct {
	Bucket string
	Prefix string
}

func (c *Controller) List(opts ListOptions) error {
	if opts.Bucket == "" {
		opts.Bucket = c.config.Bucket
	}
	objs, err := c.s3Client.List(opts.Bucket, opts.Prefix)

	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	for _, obj := range objs {
		//c.logger.Info(obj.Key)
		fmt.Println(obj.Key)
	}
	return nil
}

type UploadOptions struct {
	Endpoint   string
	Bucket     string
	Prefix     string
	EncryptKey string
	Paths      []string
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

func (c *Controller) uploadSingleFile(bucket, prefix, path, encryptKey string) error {
	file, err := c.fileHandler.Read(path)
	if err != nil {
		c.logger.Error("upload failed", "err", err.Error())
		return err
	}

	// compress file content
	if err := c.fileHandler.Compress(file); err != nil {
		c.logger.Error("compress file failed", "err", err.Error())
		return err
	}

	// encrypt file content
	if err := c.fileHandler.Encrypt(file, encryptKey); err != nil {
		c.logger.Error("encrypt file failed", "err", err.Error())
		return err
	}

	obj, err := c.s3Client.Upload(bucket, prefix, file)
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

func (c *Controller) UploadDirectoryOrFile(bucket, prefix, path, encryptKey string) error {
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
		limiter := make(chan struct{}, uploadParallelism)

		for _, file := range files {
			wg.Add(1)
			limiter <- struct{}{} // Take up a concurrent signal

			go func(file string) {
				defer func() {
					<-limiter // Release a concurrent signal
					wg.Done()
				}()

				subPrefix := filepath.Join(prefix, filepath.Dir(trimDirectory(path, file)))
				if err := c.uploadSingleFile(bucket, subPrefix, file, encryptKey); err != nil {
					c.logger.Error("error uploading file", "file", file, "error", err.Error())
				}
			}(file)
		}

		wg.Wait()
	} else {
		if err := c.uploadSingleFile(bucket, prefix, path, encryptKey); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) Upload(opts UploadOptions) error {
	if opts.Endpoint == "" {
		opts.Endpoint = c.config.Endpoint
	}
	if opts.Bucket == "" {
		opts.Bucket = c.config.Bucket
	}

	if len(opts.Paths) == 0 {
		err := errors.New("no files to upload")
		c.logger.Error("upload failed", "err", err.Error())
		return err
	}

	if opts.Endpoint != "" && opts.Endpoint != c.config.Endpoint {
		if err := c.s3Client.SetEndpoint(opts.Endpoint); err != nil {
			c.logger.Error("upload failed", "err", err.Error())
			return err
		}
		c.config.Endpoint = opts.Endpoint
	}

	if opts.Bucket != "" && opts.Bucket != c.config.Bucket {
		if err := c.s3Client.SetBucket(opts.Bucket); err != nil {
			c.logger.Error("upload failed", "err", err.Error())
			return err
		}
		c.config.Bucket = opts.Bucket
	}

	for _, path := range opts.Paths {
		if err := c.UploadDirectoryOrFile(opts.Bucket, opts.Prefix, path, opts.EncryptKey); err != nil {
			c.logger.Error("upload failed", "err", err.Error())
			return err
		}
	}

	return nil
}

type DownloadOptions struct {
	Endpoint   string
	Bucket     string
	OutputDir  string
	DecryptKey string
	S3keys     []string
}

func (c *Controller) downloadSingleFile(bucket, s3key, outputDir, decryptKey string) error {
	file, err := c.s3Client.Download(
		&internal.S3Object{
			Bucket: bucket,
			Key:    s3key,
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

	// decompress file content
	if err := c.fileHandler.Decompress(file); err != nil {
		c.logger.Error("decompress file failed", "key", s3key, "err", err.Error())
		return err
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

func (c *Controller) downloadDirectoryOrFile(bucket, s3key, outputDir, decryptKey string) error {
	objs, err := c.s3Client.List(bucket, s3key)
	if err != nil {
		c.logger.Error("download directory or file failed", "key", s3key, "err", err.Error())
		return err
	}
	if len(objs) == 0 {
		err = errors.New("directory or file not found")
		c.logger.Error("download directory or file failed", "key", s3key, "err", err.Error())
	}
	var wg sync.WaitGroup
	limiter := make(chan struct{}, downloadParallelism)

	for _, obj := range objs {
		wg.Add(1)
		limiter <- struct{}{} // Take up a concurrent signal

		go func(obj *internal.S3Object) {
			defer func() {
				<-limiter // Release a concurrent signal
				wg.Done()
			}()
			if err := c.downloadSingleFile(bucket, obj.Key, outputDir, decryptKey); err != nil {
				c.logger.Error("download directory or file failed", "key", obj.Key, "err", err.Error())
			}
		}(obj)
	}
	wg.Wait()
	return nil
}

func (c *Controller) Download(opts DownloadOptions) error {
	if opts.Endpoint == "" {
		opts.Endpoint = c.config.Endpoint
	}
	if opts.Bucket == "" {
		opts.Bucket = c.config.Bucket
	}

	if len(opts.S3keys) == 0 {
		err := errors.New("no files to download")
		c.logger.Error(err.Error())
		return err
	}

	if opts.Endpoint != "" && opts.Endpoint != c.config.Endpoint {
		if err := c.s3Client.SetEndpoint(opts.Endpoint); err != nil {
			c.logger.Error("download failed", "err", err.Error())
			return err
		}
		c.config.Endpoint = opts.Endpoint
	}

	if opts.Bucket != "" && opts.Bucket != c.config.Bucket {
		if err := c.s3Client.SetBucket(opts.Bucket); err != nil {
			c.logger.Error("download failed", "err", err.Error())
			return err
		}
		c.config.Bucket = opts.Bucket
	}

	for _, s3key := range opts.S3keys {
		if err := c.downloadDirectoryOrFile(opts.Bucket, s3key, opts.OutputDir, opts.DecryptKey); err != nil {
			return err
		}
	}
	return nil
}
