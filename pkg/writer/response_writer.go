package writer

import "net/http"

type ResponseWriter interface {
	http.ResponseWriter
	Status()
}

type Writer struct {
	http.ResponseWriter
	status int
}

func (w *Writer) Status() int {
	return w.status
}
