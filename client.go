package main

import (
	"context"
	"fmt"
	"gRpcBoost/protos"
	"google.golang.org/grpc"
	"log"
)

func failOnError(err error, msg string)  {
	if err != nil {
		log.Fatalf("%s: %s",msg,err)
	}
}

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

	//第一种普通函数双向
	//go func() {
	//	for{
	//		tempM := ""
	//		fmt.Scanln(&tempM)
	//		resp,err := client.Hello(context.Background(),&protos.String{Value: string(tempM)})
	//		if err != nil{
	//			log.Fatal(resp)
	//			return
	//		}
	//		fmt.Println(resp)
	//	}
	//}()
	//
	//time.Sleep(time.Second)
/*
==================================分界线====================================
*/

	//第二种:双向流式接受数据--通过
	//stream,err := client.Channel(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//go func() {
	//	for{
	//		//发
	//		tempM := ""
	//		fmt.Scanln(&tempM)
	//		stream.Send(&protos.String{
	//			Value: tempM,
	//		})
	//		//收
	//		reply,err := stream.Recv()
	//		if err != nil {
	//			if err == io.EOF {
	//				return
	//			}
	//			log.Fatal(err)
	//		}
	//		fmt.Println(reply.GetValue())
	//	}
	//}()
/*
   ==================================分界线====================================
*/
	//第三种，服务端流式传输输
	//go func() {
	//	for{
	//		tempM := ""
	//		fmt.Scanln(&tempM)
	//		stream,err := client.ChannelS(context.Background(),&protos.String{Value: tempM})
	//		if err != nil{
	//			log.Fatal(err)
	//			return
	//		}
	//		msg,err := stream.Recv()
	//		if err == io.EOF{
	//			return
	//		}
	//		if err != nil{
	//			log.Fatal(err)
	//			return
	//		}
	//		fmt.Println(msg)
	//	}
	//}()
/*
   ==================================分界线====================================
*/
	//第四种,客户端流式
	stream,err := client.ChannelC(context.Background())
	if err != nil{
		fmt.Println(err)
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
			mes,err:=stream.CloseAndRecv()
			if err != nil{
				log.Fatal(err)
			}
			fmt.Println(mes.GetValue())
		}
	}()
	//阻塞
	<-forever
}