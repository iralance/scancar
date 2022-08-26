package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/iralance/scancar/server/auth/api/gen/v1"
	todopb "github.com/iralance/scancar/server/todo/api/gen/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
)

func main() {
	// logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true,
				UseProtoNames:  true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true, // If DiscardUnknown is set, unknown fields are ignored.
			},
		},
	))

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:8081",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "todo",
			addr:         "localhost:8082",
			registerFunc: todopb.RegisterTodoServiceHandlerFromEndpoint,
		},
	}

	for _, s := range serverConfig {
		err := s.registerFunc(
			c, mux, s.addr,
			[]grpc.DialOption{grpc.WithInsecure()},
		)
		if err != nil {
			logger.Sugar().Fatalf("cannot register service %s : %v", s.name, err)
		}
	}
	addr := ":8200"
	logger.Sugar().Infof("grpc gateway started at %s", addr)
	logger.Sugar().Fatal(http.ListenAndServe(addr, mux))
}

func demo1() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseEnumNumbers:  true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	))
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(
		c,
		//mux: multiplexer
		mux,
		":8081",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	log.Println(err)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	err = todopb.RegisterTodoServiceHandlerFromEndpoint(
		c,
		//mux: multiplexer
		mux,
		":8082",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	log.Println(err)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	err = http.ListenAndServe(":8200", mux)
	if err != nil {
		log.Fatalf("cannot listen and server: %v", err)
	}
}
