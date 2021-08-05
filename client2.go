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
	第二种方式:双向流式
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

	//第二种:双向流式接受数据--通过
	stream,err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	go func() {
		for{
			//发
			tempM := ""
			fmt.Scanln(&tempM)
			stream.Send(&protos.String{
				Value: tempM,
			})
			//收
			reply,err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Fatal(err)
			}
			fmt.Println(reply.GetValue())
		}
	}()
	<-forever
}