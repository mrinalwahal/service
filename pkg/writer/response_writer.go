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

func (w *Writer) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *Writer) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}

func NewWriter(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}
