package backup

import (
	"io/ioutil"
	"os"
)

type RawBackup struct {
}

func (r *RawBackup) Write(filePath string, content []byte) error {
	fd, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = fd.Write(content)
	return err
}

func (r *RawBackup) Read(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
