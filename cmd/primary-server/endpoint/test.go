package endpoint

import (
	"github.com/gorilla/mux"
	"net/http"
	"plans/internal/session"
)

type test struct{}

// RegisterTest registers /test endpoints.
func RegisterTest(r *mux.Router) {
	t := test{}

	s := r.PathPrefix("/test").Subrouter()
	s.HandleFunc("/{message}", t.Hello).Methods(http.MethodGet)
}

func (t test) Hello(w http.ResponseWriter, r *http.Request) {
	l := session.MustLoggerFromContext(r.Context())

	vars := mux.Vars(r)

	_, err := w.Write([]byte(vars["message"]))
	if err != nil {
		l.Warn(err.Error())
	}
}
