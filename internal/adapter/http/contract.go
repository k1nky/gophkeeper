package http

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}
