package middleware

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

func AddLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}

	grpc_zap.ReplaceGrpcLoggerV2(logger)

	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, o...),
	))

	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, o...),
	))
	return opts
}
