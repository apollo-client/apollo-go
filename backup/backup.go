package backup

type Backup interface {
	Write(filePath string, content []byte) error
	Read(filePath string) ([]byte, error)
}

var (
	DefaultBackup = newBackup()
)

func newBackup() Backup {
	return &RawBackup{}
}
