package config

type Context struct {
	Config Application
}

var context *Context

func SetContext(conf Application) {
	context = &Context{
		Config: conf,
	}
}

func GetContext() Context {
	return *context
}
