package log

type Logger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

var defLog Logger

func init() {
	defLog = &Log{}
}

func Init(l Logger) {
	defLog = l
}

func Debugf(format string, params ...interface{}) {
	defLog.Debugf(format, params...)
}
func Infof(format string, params ...interface{}) {
	defLog.Infof(format, params...)
}
func Warnf(format string, params ...interface{}) {
	defLog.Warnf(format, params...)
}
func Errorf(format string, params ...interface{}) {
	defLog.Errorf(format, params...)
}
func Debug(v ...interface{}) {
	defLog.Debug(v...)
}
func Info(v ...interface{}) {
	defLog.Info(v...)
}
func Warn(v ...interface{}) {
	defLog.Warn(v...)
}
func Error(v ...interface{}) {
	defLog.Error(v...)
}

type Log struct{}

func (l *Log) Debugf(format string, params ...interface{}) {
}

func (l *Log) Infof(format string, params ...interface{}) {
}

func (l *Log) Warnf(format string, params ...interface{}) {
}

func (l *Log) Errorf(format string, params ...interface{}) {
}

func (l *Log) Debug(v ...interface{}) {
}

func (l *Log) Info(v ...interface{}) {
}

func (l *Log) Warn(v ...interface{}) {
}

func (l *Log) Error(v ...interface{}) {
}
