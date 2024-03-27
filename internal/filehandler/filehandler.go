package filehandler

import (
	"os"
	"path/filepath"

	"github.com/linlanniao/soss/internal"
	"github.com/linlanniao/soss/pkg/cipher"
)

type fileHandler struct {
	_cipher *cipher.ContentCipher
}

var _ internal.IFileHandler = (*fileHandler)(nil)

func NewFileHandler() internal.IFileHandler {
	return &fileHandler{}
}

func (f *fileHandler) cipher(key string) (*cipher.ContentCipher, error) {
	if f._cipher != nil && f._cipher.KeyIsEqual(key) {
		return f._cipher, nil
	}
	c, err := cipher.NewContentCipher(key)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (f *fileHandler) Encrypt(in *internal.File, encryptKey string) (err error) {
	c, err := f.cipher(encryptKey)
	if err != nil {
		return err
	}

	in.Content, err = c.EncryptBytes(in.Content)
	if err != nil {
		return err
	}

	in.Encrypted = true
	return nil
}

func (f *fileHandler) Decrypt(in *internal.File, decryptKey string) (err error) {
	c, err := f.cipher(decryptKey)
	if err != nil {
		return err
	}

	in.Content, err = c.DecryptBytes(in.Content)
	if err != nil {
		return err
	}

	in.Encrypted = false
	return nil
}

func (f *fileHandler) Read(path string) (*internal.File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &internal.File{
		Path:      path,
		Content:   content,
		Encrypted: false,
	}, nil
}

func (f *fileHandler) Write(file *internal.File) error {
	dir := filepath.Dir(file.Path)

	// if directory no exist, create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(file.Path, file.Content, 0644)
}

func (f *fileHandler) SearchFiles(path string) (files []string, err error) {
	// if given path is a file, return list of this file
	if info, err := os.Stat(path); err != nil && !info.IsDir() {
		return []string{path}, nil
	}

	files = make([]string, 0)
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
