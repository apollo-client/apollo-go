package backup

import (
	"io/ioutil"
	"os"

	"github.com/apollo-client/apollo-go/log"
)

type RawBackup struct {
}

func (r *RawBackup) Write(filePath string, content []byte) error {
	fd, err := os.Create(filePath)
	if err != nil {
		log.Errorf("write backup err: %v\n", err)
		return err
	}
	_, err = fd.Write(content)
	if err != nil {
		log.Errorf("write backup err: %v\n", err)
	}
	return err
}

func (r *RawBackup) Read(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
