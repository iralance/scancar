package main

import (
	"github.com/iralance/scancar/server/shared/auth"
	"github.com/iralance/scancar/server/shared/server"
	todopb "github.com/iralance/scancar/server/todo/api/gen/v1"
	"github.com/iralance/scancar/server/todo/todo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	logger.Sugar().Fatal(
		server.RunGRPCServer(&server.GRPCConfig{
			Name:              "todo",
			Addr:              ":8082",
			AuthPublicKeyFile: "./../shared/auth/public.key",
			Logger:            logger,
			RegisterFunc: func(s *grpc.Server) {
				todopb.RegisterTodoServiceServer(s, &todo.Service{
					Logger: logger,
				})
			},
		}),
	)
}

func demo1() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("cannot create logger:", zap.Error(err))
	}
	//c := context.Background()

	ln, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Fatal("failed to listen:", zap.Error(err))
	}
	in, err := auth.Interceptor("./../shared/auth/public.key")
	if err != nil {
		logger.Fatal("cannot create auth interceptor")
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(in))
	todopb.RegisterTodoServiceServer(s, &todo.Service{
		Logger: logger,
	})
	log.Printf("server listening at %v", ln.Addr())
	err = s.Serve(ln)
	logger.Fatal("cannot server", zap.Error(err))
}
