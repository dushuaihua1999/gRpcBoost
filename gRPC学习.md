# gRPC学习

## 1.用处

## 2.优势

#### 1.通过protobuf来定义接口--严格约束

#### 2.通过怕protobuf可以将数据序列化为二进制编码--减少数据量

#### 3.

## 3.结构

rpc Add (Request) returns (Request) {}

1. rpc 是一个保留的协议缓冲关键字，表示该函数时一个远程过程调用
2. Add 是函数的名称
3. （Request）表示该函数有一个自定义消息类型的参数Request
4. returns 是一个保留的协议缓冲关键字，表hi函数返回类型的前缀
5. (Request) 表示该函数将返回一个自定义的消息类型,Response

## 4.交互方式分类

#### 1.双向普通函数调用

在消息一样的情况下:

客户端：

```go
resp,err := client.Hello(context.Background(),&protos.String{Value: string(tempM)})
```

服务端:

```go
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
```

#### 2.客户端流式，服务端普通函数

客户端:

```go
//stream,err := client.Channel(context.Background())
//if err != nil {
// log.Fatal(err)
//}
//go func() {
// for{
//    tempM := ""
//    fmt.Scanln(&tempM)
***
//    stream.Send(&protos.String{
***
//       Value: tempM,
//    })
***
//    reply,err := stream.Recv()
***
//    if err != nil {
//       if err == io.EOF {
//          return
//       }
//       log.Fatal(err)
//    }
//    fmt.Println(reply.GetValue())
// }
//}()
```

服务端:

```go
func (p *HelloServiceImpl) ChannelC(stream protos.HelloService_ChannelCServer) error {
   fmt.Println("调用ChannelC方法...")
***
   args,err := stream.Recv()
***
   fmt.Println(args.GetValue())

   if err != nil{
      if err == io.EOF{
         return nil
      }
      return err
   }

   tempM := ""
   fmt.Scanln(&tempM)
   reply := &protos.String{Value: tempM}
    ***
   err = stream.SendAndClose(reply)
    ***
   if err != nil {
      return err
   }
   return err
}
```

#### 3.客户端普通函数，服务端流式

客户端：

```go
//go func() {
// for{
//    tempM := ""
//    fmt.Scanln(&tempM)
***
//    stream,err := client.ChannelS(context.Background(),&protos.String{Value: tempM})
***
//    if err != nil{
//       log.Fatal(err)
//    }
//    for{
***
//       msg,err := stream.Recv()
***
//       if err == io.EOF{
//          break
//       }
//       if err != nil{
//          log.Fatal(err)
//       }
//       fmt.Println(msg)
//    }
// }
//}()
```

服务端:

```go
func (p *HelloServiceImpl) ChannelS(in *protos.String,stream protos.HelloService_ChannelSServer) error {
   fmt.Println("调用ChannelS方法...")
    ***
   fmt.Println(in.GetValue())
	***
   tempM := ""
   fmt.Scanln(&tempM)
    ***
   reply := &protos.String{Value: tempM}
   err := stream.Send(reply)
    ***
   if err != nil {
      return err
   }
   return err
}
```

#### 4.双向流式

客户端:

```go
//stream,err := client.Channel(context.Background())
//if err != nil {
// log.Fatal(err)
//}
//go func() {
// for{
//    tempM := ""
//    fmt.Scanln(&tempM)
***
//    stream.Send(&protos.String{
//       Value: tempM,
//    })
***
****
//    reply,err := stream.Recv()
***
//    if err != nil {
//       if err == io.EOF {
//          return
//       }
//       log.Fatal(err)
//    }
//    fmt.Println(reply.GetValue())
// }
//}()
```

服务端：

```go
func (p *HelloServiceImpl) Channel(stream protos.HelloService_ChannelServer) error {
      fmt.Println("调用Channel方法...")
    ***
      args,err := stream.Recv()
    ***
      fmt.Println(args.GetValue())
      if err != nil{
         if err == io.EOF{
            return nil
         }
         return err
      }
      tempM := ""
      fmt.Scanln(&tempM)
      reply := &protos.String{Value: tempM}
    ***
      err = stream.Send(reply)
    ***
      if err != nil {
         return err
      }
      return err
}
```

## 5.proto文件语法

```protobuf
syntax = "proto3"  ----protobuf版本
package protos;    ----包名
message String {
  string value = 1;   --- 数据类型 变量名 = 标识符
}

service HelloService {
  //普通函数调用
  rpc Hello (String) returns (String);    -- rpc调用  接口名 （请求的数据类型）returns(返回的数据类型)  
  //双向流调用
  rpc Channel (stream String) returns (stream String);
  //客户端普通函数,服务器流式
  rpc ChannelS (String) returns (stream String);
  //客户端流式,服务器普通函数调用
  rpc ChannelC (stream String) returns (String);
}

```

## 6.流式交互的关键点

客户端流与服务端流分析:

1.新流请求的产生,一个新地流调用是下面这个获取strem对象方法的调用也就是想要流不断创建新的，就要反复调用下面ChannelS与ChannelC方法

2.而其对应的Send方法与CloseAndRecv方法是在一个流请求中的，Send与Recv是一一对应的关系。

3.检测本次流是否结束，如果没有自定义检测的话，对于客户端来说就是，传数据过去收不到相应了，也即是返回为空，出现io.EOF
错误。对于服务端来说就是，收不到数据了。

```go
服务端流:
stream,err := client.ChannelS(context.Background(),&protos.String{Value: tempM})
			if err != nil{
				log.Fatal(err)
				return
			}


客户端流:
stream,err := client.ChannelC(context.Background())
if err != nil{
   fmt.Println(err)
   return
}
```

2.接着上面的情况分析一下下面的代码：

#### 客户端：

```go
go func() {
   for{
      stream,err := client.ChannelC(context.Background())
      if err != nil{
         fmt.Println(err)
         return
      }
      //发
      tempM := ""
      fmt.Scanln(&tempM)
      err = stream.Send(&protos.String{
         Value: tempM,
      })
      if err != nil {
         if err == io.EOF{
            fmt.Println("传输完成")
         }
      }
      //收
      mes,err:=stream.CloseAndRecv()
      if err != nil{
         if err != nil {
            if err == io.EOF{
               fmt.Println("本次传输完成")
               continue
            }
         }
      }
      fmt.Println("server: " + mes.GetValue())
   }
}()
```

#### 服务端：

```go
func (p *HelloServiceImpl) ChannelC(stream protos.HelloService_ChannelCServer) error {
   fmt.Println("调用ChannelC方法...")
   for{
      //收
      args,err := stream.Recv()
      fmt.Println("client: ",args.GetValue())
      if err != nil{
         if err == io.EOF{
            fmt.Println("本次传输完成")
            return nil
         }
         return err
      }

      //发
      fmt.Println("开始输入响应信息: ")
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
```

在一次流请求中，一旦流对象获取成功就会开始调用服务端的ChannelC方法，假设服务端发送123，这时服务端如下代码会开始打印123

```go
args,err := stream.Recv()
fmt.Println("client: ",args.GetValue())
```

打印完，进入下面等待“开始输入相应信息：”的操作：

输入之后，因为这时候客户端还在CloseAndRecv那里等待，Close意味着客户端发完这个消息后通道stream要关闭。

```go
mes,err:=stream.CloseAndRecv()
```

这时候，服务端收到扫描的信息后，开始执行SendAndClose的方法，把数据传过去后，服务端这边也关闭了通道的访问。而这时候如果不创建新的流对象的话，客户端也就不知道自己是否应该停止，这时候就会继续调用Send方法，调完之后，依然等待服务端回信息。但是此时通道已经没有了，获取不到返回的信息。就会出现EOF，认为此次文件传输已经结束了。服务端则是，收不到信息，认为此次交互也结束了。所以会出现EOF的错误。这就是前段时间遇到的疑问点。虽然没有了流，但对于Send等操作流的方法来说，依然可以调用，就是是一个空流而已，里面不接受数据，可以认为客户端发的信息丢了。所以服务端那里也收不到消息。

## 7.Send与Recv

1. 首先stream流可以分为两个方向---双流式
   1. 客户端Send-----服务端Recv
   2. 客户端Recv-----服务端Send

2. CS双方的Send与Recv是一一对应关系，一个Send对应一个Recv

3. CloseSend()方法会将该端的Send方向给关闭，如果对方没有捕获该异常，就会异常终止，因为这个方向已经不通了，但是对方还在等待，就会出错。

4. 综上，stream流也是有限资源，在资源有限的情况下，最好的方式是，能用一次Send与Recv解决就一次解决，不能的情况下，要尽量保证几次Send就用几次Recv处理，不要出现某一方有多余的传输与等待存在，比如一个是有限的for循环，一个是无限的for循环，这就会导致一段出现阻塞的情况。

5. io.EOF的错误产生有两种方式：

   1. 用户自定义条件抛出错误
   2. 系统本次调用结束，这种方式判断的权力就交给系统即gRPC框架自己去判断，一般默认就是Recv接受不到数据了就EOF。这适用于双方的判断，客户端如果收不到了消息就报出EOF。

   

