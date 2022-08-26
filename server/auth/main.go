package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	authpb "github.com/iralance/scancar/server/auth/api/gen/v1"
	"github.com/iralance/scancar/server/auth/auth"
	"github.com/iralance/scancar/server/auth/dao"
	"github.com/iralance/scancar/server/auth/token"
	"github.com/iralance/scancar/server/auth/wechat"
	"github.com/iralance/scancar/server/shared/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

var wechatAppID = "wx2321c27788f305c3"
var wechatAppSecret = "2f005a852e0b20244bbe5b9f8f3edf2e"
var mongoURI = "mongodb://localhost:27016/?readPreference=primary&ssl=false"
var privateKeyFile = "./private.key"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("cannot create logger:", zap.Error(err))
	}
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	pkFile, err := os.Open(privateKeyFile)
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}
	// RunGRPCServer
	logger.Sugar().Fatal(
		server.RunGRPCServer(&server.GRPCConfig{
			Name:   "auth",
			Addr:   ":8081",
			Logger: logger,
			RegisterFunc: func(s *grpc.Server) {
				authpb.RegisterAuthServiceServer(s, &auth.Service{
					OpenIDResolver: &wechat.Service{
						AppID:     wechatAppID,
						AppSecret: wechatAppSecret,
					},
					Mongo:          dao.NewMongo(mongoClient.Database("scancar")),
					Logger:         logger,
					TokenExpire:    2 * time.Hour,
					TokenGenerator: token.NewJWTTokenGen("server/auth", privKey),
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
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	pkFile, err := os.Open(privateKeyFile)
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal("failed to listen:", zap.Error(err))
	}
	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		Logger: logger,
		OpenIDResolver: &wechat.Service{
			AppID:     wechatAppID,
			AppSecret: wechatAppSecret,
		},
		Mongo:          dao.NewMongo(mongoClient.Database("scancar")),
		TokenExpire:    2 * time.Hour,
		TokenGenerator: token.NewJWTTokenGen("scancar/auth", privKey),
	})
	log.Printf("server listening at %v", ln.Addr())
	err = s.Serve(ln)
	logger.Fatal("cannot server", zap.Error(err))
}
