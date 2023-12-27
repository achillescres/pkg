package auth

import (
	"context"
	"github.com/achillescres/pkg/gin/ginmiddleware"
	"github.com/achillescres/pkg/grpc/auth/authProto"
	"google.golang.org/grpc"
)

func NewDefaultTokenGRPCChecker(conn *grpc.ClientConn) (ginmiddleware.TokenChecker[*authProto.UserInfo], error) {
	client := authProto.NewExternalAuthClient(conn)

	return func(ctx context.Context, token string) (*authProto.UserInfo, error) {
		userInfo, err := client.Permissions(ctx, &authProto.CookieAccess{
			Access: token,
		})
		if err != nil {
			return nil, err
		}
		return userInfo, nil
	}, nil
}
