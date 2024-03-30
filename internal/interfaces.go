package internal

type IDownloader interface {
	Download(obj *S3Object, outputDir string) (file *File, err error)
}

type IUploader interface {
	Upload(bucket string, prefix string, file *File) (obj *S3Object, err error)
}

type ILister interface {
	List(bucket string, prefix string) (objs []*S3Object, err error)
}

type IS3ClientConfigurator interface {
	SetEndpoint(endpoint string) error
	SetBucket(bucket string) error
}

type IS3Client interface {
	ILister
	IUploader
	IDownloader
	IS3ClientConfigurator
}

type IContentCipher interface {
	Encrypt(in *File, encryptKey string) (err error)
	Decrypt(in *File, decryptKey string) (err error)
}

type IContentCompressor interface {
	Compress(in *File) (err error)
	Decompress(in *File) (err error)
}

type IFileReadWriter interface {
	Read(path string) (*File, error)
	Write(file *File) error
}

type IFileScanner interface {
	SearchFiles(path string) (files []string, err error)
}

type IFileHandler interface {
	IContentCipher
	IFileReadWriter
	IFileScanner
	IContentCompressor
}
