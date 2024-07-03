package models

type HttpRequest struct {
	Method  string
	Path    string
	Version string
	Headers string
	Body    string
}
