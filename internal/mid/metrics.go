package mid

import (
	"expvar"
	"net/http"
	"runtime"
	"sales_service/internal/platform/web"
)

var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

func Metrics() web.Middleware {
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			err := before(w, r)

			m.req.Add(1)

			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}
			if err != nil {
				m.err.Add(1)
			}
			return err
		}
		return h
	}
	return f
}
