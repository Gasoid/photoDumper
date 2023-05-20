package main

import (
	"net/http"
	"testing"

	_ "github.com/Gasoid/photoDumper/docs"
)

type testEngine struct{}

func (e *testEngine) Run(addr ...string) error {
	return nil
}

func (e *testEngine) ServeHTTP(http.ResponseWriter, *http.Request) {}

func Test_main(t *testing.T) {
	setupRouterFunc = func() engine { return &testEngine{} }
	openBrowserFunc = func(url string) {}
	tests := []struct {
		name string
	}{
		{name: "main test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_openBrowser(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "fake addr", args: args{url: "http://localhost:1234"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openBrowser(tt.args.url)
		})
	}
}
