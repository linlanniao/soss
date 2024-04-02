package secret

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/linlanniao/soss/pkg/utils"
)

type Secret struct {
	key  string
	path string
}

func GenerateSecret() *Secret {
	s := &Secret{
		key:  utils.RandLowerUpperNumStr(64),
		path: savePath,
	}
	return s
}

var (
	home, _  = os.UserHomeDir()
	savePath = filepath.Join(home, ".soss", ".secret")
)

func (secret *Secret) Save() error {
	if secret == nil {
		return errors.New("secret is nil")
	}

	if len(secret.key) == 0 {
		return errors.New("key is empty")
	}

	dir := filepath.Dir(secret.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}

	if utils.IsFile(secret.path) {
		return errors.New("secret file already exists")
	}

	return os.WriteFile(secret.path, []byte(secret.key), 0600)
}

func (secret *Secret) Replace() (newerSecret *Secret, err error) {
	if secret == nil {
		return nil, errors.New("secret is nil")
	}

	if len(secret.key) == 0 {
		return nil, errors.New("key is empty")
	}

	dir := filepath.Dir(secret.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, err
		}
	}

	newer := &Secret{
		key:  secret.key,
		path: secret.path,
	}

	if utils.IsFile(newer.path) {
		now := time.Now()
		t := now.Format("2006-01-02_15-04-05")
		backupFileName := fmt.Sprintf("%s.%s.backup", secret.path, t)
		if err := os.Rename(secret.path, backupFileName); err != nil {
			return nil, err
		}
		secret.path = backupFileName
	}

	err = os.WriteFile(newer.path, []byte(secret.key), 0600)
	if err != nil {
		return nil, err
	}
	return newer, nil
}

func (secret *Secret) Path() string {
	return secret.path
}

func (secret *Secret) Key() string {
	return secret.key
}

func Load() (*Secret, error) {
	if !utils.IsFile(savePath) {
		return nil, errors.New("secret file not found")
	}
	b, err := os.ReadFile(savePath)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, errors.New("secret file is empty")
	}

	return &Secret{key: string(b)}, nil
}
