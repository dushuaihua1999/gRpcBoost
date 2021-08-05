package main

import (
	"context"
	"fmt"
	"gRpcBoost/protos"
	"google.golang.org/grpc"
	"io"
	"log"
)

/*
	第三种方式:服务端流式
*/
func main(){
	//1.客户端发起连接
	conn,err := grpc.Dial("localhost:1234",grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	//2.注册客户端
	client := protos.NewHelloServiceClient(conn)
	//阻塞main主程
	forever := make(chan bool)

	//第三种:双向流式接受数据--通过
	go func() {
		for{
			tempM := ""
			fmt.Scanln(&tempM)
			stream,err := client.ChannelS(context.Background(),&protos.String{Value: tempM})
			if err != nil{
				log.Fatal(err)
				return
			}
			msg,err := stream.Recv()
			if err == io.EOF{
				return
			}
			if err != nil{
				log.Fatal(err)
				return
			}
			fmt.Println(msg)
		}
	}()
	<-forever
}
