package commons

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	User      = "user"
	UserEmail = "user_email"
)

func GetSingleValue(ctx context.Context, key string) string {
	incomingContext, _ := metadata.FromIncomingContext(ctx)
	values := incomingContext.Get(key)
	if values == nil {
		return ""
	}
	return values[0]
}
