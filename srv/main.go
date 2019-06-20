package main

import (
	"log"
	"math/rand"
	"time"

	"context"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/etcdv3"
	demo "github.com/weiwenwang/go-mcro-demo/srv/proto/demo"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *demo.Request, rsp *demo.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name + " abcd " + string(rand.Intn(100))
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	demo.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
