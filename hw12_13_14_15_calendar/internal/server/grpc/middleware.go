package internalgrpc

import (
	"context"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(LoggingInterceptor)
}

func LoggingInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	duration := time.Since(start)
	s := []string{
		"[" + time.Now().Format(time.RFC3339) + "]",
		info.FullMethod,
		duration.String(),
		"\n",
	}

	os.Stdout.WriteString(strings.Join(s, " "))

	return h, err
}
