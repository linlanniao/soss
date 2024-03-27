package internal

type File struct {
	Path      string // absolute path
	Content   []byte // content
	Encrypted bool   // if true, content is encrypted
}

type S3Object struct {
	Bucket string // Object bucket
	Key    string // Object key
	Type   string // Object type
	Size   int64  // Object size
	ETag   string // Object eTag
}
