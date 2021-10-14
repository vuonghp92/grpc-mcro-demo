package interceptor

import (
	"context"
	"fmt"

	"github.com/vuonghp92/grpc-mcro-demo/front/support"
	"github.com/vuonghp92/grpc-mcro-demo/shared/md"
	"google.golang.org/grpc"
)

func XTraceID(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	fmt.Println("------xxx", ctx)
	traceID := support.GetTraceIDFromContext(ctx)
	ctx = md.AddTraceIDToContext(ctx, traceID)
	fmt.Println("------xxx111", ctx)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func XUserID(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	user := support.GetUserFromContext(ctx)
	ctx = md.AddUserIDToContext(ctx, user.Id)
	return invoker(ctx, method, req, reply, cc, opts...)
}
