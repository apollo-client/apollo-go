package apollo

import (
	"time"

	"github.com/xnzone/apollo-go/auth"
	"github.com/xnzone/apollo-go/backup"
	"github.com/xnzone/apollo-go/log"
	"github.com/xnzone/apollo-go/transport"
)

type Option func(o *Options)

type Options struct {
	Auth          auth.Auth           // auth interface
	Transport     transport.Transport // transport interface
	WatchInterval time.Duration       // watch interval
	Backup        backup.Backup       // backup interface
	EnableBackup  bool                // enable backup
	BackupPath    string              // backup path
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		Auth:          auth.DefaultAuth,
		Transport:     transport.DefaultTransport,
		Backup:        backup.DefaultBackup,
		WatchInterval: 5 * time.Second,
		EnableBackup:  false,
		BackupPath:    "./",
	}
	for _, o := range opts {
		o(&opt)
	}
	return &opt
}

// Auth custom auth
func Auth(a auth.Auth) Option {
	return func(o *Options) { o.Auth = a }
}

// Transport custom request transport
func Transport(t transport.Transport) Option {
	return func(o *Options) { o.Transport = t }
}

// Backup custom backup read and write
func Backup(b backup.Backup) Option {
	return func(o *Options) { o.Backup = b }
}

// WatchInterval watch interval
func WatchInterval(t time.Duration) Option {
	return func(o *Options) { o.WatchInterval = t }
}

// EnableBackup enable backup
func EnableBackup(enable bool) Option {
	return func(o *Options) { o.EnableBackup = enable }
}

// BackupPath backup path
func BackupPath(p string) Option {
	return func(o *Options) { o.BackupPath = p }
}

func Logger(l log.Logger) Option {
	return func(o *Options) {
		log.Init(l)
	}
}
