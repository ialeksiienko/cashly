package test

import "cashly/pkg/slogx"

func newTestLogger() slogx.Logger {
	return slogx.New("local")
}

var key = []byte("examplekey1234567890examplekey12")
