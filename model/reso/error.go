package reso

type HTTPError struct {
	Code int
	Msg  string
	Data interface{}
}
