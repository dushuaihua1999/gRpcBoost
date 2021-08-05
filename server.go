package main

import (
	"gRpcBoost/protos"
	"gRpcBoost/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main()  {
	//1.创建grpc服务器对象
	grpcServer := grpc.NewServer()
	//2.服务端注册
	protos.RegisterHelloServiceServer(grpcServer,new(service.HelloServiceImpl))
	//3.监听端口
	listener,err := net.Listen("tcp",":1234")
	if err != nil {
		log.Fatal(err)
	}
	//4.为该端口来的连接提供服务
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
