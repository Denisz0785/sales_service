package web

type Middleware func(Handler) Handler

// wrapMiddleware applies a chain of middleware functions to a handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Iterate over the middleware functions in reverse order.
	for i := len(mw) - 1; i >= 0; i-- {
		// Get the current middleware function.
		h := mw[i]

		// If the middleware function is not nil, apply it to the handler.
		if h != nil {
			handler = h(handler)
		}
	}

	// Return the resulting handler after applying all middleware functions.
	return handler
}
