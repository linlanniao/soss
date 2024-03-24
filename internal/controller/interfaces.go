package controller

type IDownloader interface {
	Download(bucket string, files []string, outputDir string) error
}

type IUploader interface {
	Upload(bucket string, prefix string, files []string) error
}

type ILister interface {
	List(bucket string, prefix string) error
}

type IRemoveStorageClient interface {
	ILister
	IUploader
	IDownloader
}

type ISecurer interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}
