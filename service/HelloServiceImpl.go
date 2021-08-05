package service

import (
	"context"
	"fmt"
	"gRpcBoost/protos"
	"io"
	"os"
)

type HelloServiceImpl struct {}

//1.双向流的方式
func (p *HelloServiceImpl) Channel(stream protos.HelloService_ChannelServer) error {
		fmt.Println("调用Channel方法...")
		for {
			args, err := stream.Recv()
			fmt.Println(args.GetValue())
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			tempM := ""
			fmt.Scanln(&tempM)
			reply := &protos.String{Value: tempM}
			err = stream.Send(reply)
			if err != nil {
				return err
			}
		}
		return nil
}

//2.双向普通函数调用
func (p *HelloServiceImpl) Hello(ctx context.Context, in *protos.String) (*protos.String, error) {
	fmt.Println("client: ",in.GetValue())
	tempM := ""
	if in.GetValue() == "bye" || in.GetValue() == "再见"{
		os.Exit(0)
	}

	fmt.Scanln(&tempM)
	if tempM == "bye" || tempM == "再见"{
		resp := &protos.String{Value: tempM}
		return resp,nil
		defer os.Exit(0)
	}
	resp := &protos.String{Value: tempM}
	return resp,nil
}

//3.服务端流式
func (p *HelloServiceImpl) ChannelS(in *protos.String,stream protos.HelloService_ChannelSServer) error {
	fmt.Println("调用ChannelS方法...")
	for{
		fmt.Println(in.GetValue())

		tempM := ""
		fmt.Scanln(&tempM)
		reply := &protos.String{Value: tempM}
		err := stream.Send(reply)
		if err != nil {
			return err
		}
	}
	return nil
}

//4.客户端流式
func (p *HelloServiceImpl) ChannelC(stream protos.HelloService_ChannelCServer) error {
	fmt.Println("调用ChannelC方法...")
	for{
		//收
		args,err := stream.Recv()
		fmt.Println(args.GetValue())
		if err != nil{
			if err == io.EOF{
				return nil
			}
			return err
		}

		//发
		tempM := ""
		fmt.Scanln(&tempM)
		reply := &protos.String{Value: tempM}
		err = stream.SendAndClose(reply)
		if err != nil {
			return err
		}
	}
	return nil
}