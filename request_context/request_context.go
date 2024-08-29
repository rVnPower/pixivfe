// this package is separate because Go disallows cyclic import. pain

package request_context

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

type userContextKeyType struct{}

var userContextKey = userContextKeyType{}

type RequestContext struct {
	// for tracing
	RequestId string
	// for error handling. normal error page don't need to set this.
	CaughtError error
	// for Render[T]
	RenderStatusCode int
}

func Make() RequestContext {
	return RequestContext{
		RequestId:        ulid.Make().String(),
		RenderStatusCode: 200,
	}
}

func ProvideWith(ctx context.Context) context.Context {
	uc := Make()
	return context.WithValue(ctx, userContextKey, &uc)
}

func GetFromContext(context context.Context) *RequestContext {
	return context.Value(userContextKey).(*RequestContext)
}

func Get(r *http.Request) *RequestContext {
	return GetFromContext(r.Context())
}
