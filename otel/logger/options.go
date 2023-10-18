package logger

import (
	"io"
)

type Options struct {
	// Log level
	level int8
	// Enable stdout will print out log to stdout
	stdout bool
	// Log output writer
	writers []io.Writer
	// Log flush to writer time interval
	interval int
	// Flush to writer buffer size
	bufferSize int
	// Log message prefix
	prefix string

	skip int
}

type Option interface {
	Apply(o *Options)
}

type optionFunc func(o *Options)

func (f optionFunc) Apply(o *Options) {
	f(o)
}


func AddCallerSkip(skip int) Option  {
	return optionFunc(func(o *Options) {
		o.skip = skip
	})
}



func WithLevel(level int8) Option {
	return optionFunc(func(o *Options) {
		o.level = level
	})
}

func WithStdout(stdout bool) Option {
	return optionFunc(func(o *Options) {
		o.stdout = stdout
	})
}

func WithInterval(interval int) Option {
	return optionFunc(func(o *Options) {
		o.interval = interval
	})
}

func WithBufferSize(size int) Option  {
	return optionFunc(func(o *Options) {
		o.bufferSize = size
	})
}

func WithPrefix(prefix string) Option {
	return optionFunc(func(o *Options) {
		o.prefix = prefix
	})
}

func WithWriters(writers ...io.Writer) Option  {
	return optionFunc(func(o *Options) {
		o.writers = writers
	})
}

