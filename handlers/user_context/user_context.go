// pain: why go no cyclic import

package user_context

import (
	"context"

	"codeberg.org/vnpower/pixivfe/v2/utils"
	"github.com/openzipkin/zipkin-go/model"
)

type UserContext struct {
	Parent model.SpanContext
	Err error
	ErrorStatusCodeOverride int
}

type UserContextKeyType struct{}

var UserContextKey = UserContextKeyType{}

func GetUserContext(context context.Context) *UserContext {
	return context.Value(UserContextKey).(*UserContext)
}

func WithContext(ctx context.Context) context.Context {
	traceId, ctx := utils.Tracer.StartSpanFromContext(ctx, "")
	return context.WithValue(ctx, UserContextKey, &UserContext{
		Parent: traceId.Context(),
	})
}
