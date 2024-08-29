// pain: why go no cyclic import

package user_context

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

type UserContextKeyType struct{}

var UserContextKey = UserContextKeyType{}

type UserContext struct {
	RequestId string
	Error error
	ErrorStatusCode int
}

func MakeUserContext() UserContext {
	return UserContext{
		RequestId: ulid.Make().String(),
		ErrorStatusCode: http.StatusInternalServerError,
	}
}

func WithContext(ctx context.Context) context.Context {
	uc := MakeUserContext()
	return context.WithValue(ctx, UserContextKey, &uc)
}

func GetUserContext(context context.Context) *UserContext {
	return context.Value(UserContextKey).(*UserContext)
}
