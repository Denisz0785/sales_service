package mid

import (
	"context"
	"net/http"
	"sales_service/internal/platform/web"

	"github.com/go-faster/errors"
	"go.opencensus.io/trace"
)

func Panics() web.Middleware {

	f := func(after web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			ctx, span := trace.StartSpan(ctx, "internal.mid.Panics")
			defer span.End()

			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %+v", r)
				}
			}()

			return after(ctx, w, r)
		}
		return h
	}

	return f
}
