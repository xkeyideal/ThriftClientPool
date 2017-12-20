package client

import (
	"fmt"
	"time"

	"github.com/xkeyideal/ThriftClientPool/thriftPool"
	"github.com/xkeyideal/ThriftClientPool/tutorial"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var GlobalRpcPool *thriftPool.ThriftPool

func Dial(addr, port string, connTimeout time.Duration) (*thriftPool.IdleClient, error) {
	socket, err := thrift.NewTSocketTimeout(fmt.Sprintf("%s:%s", addr, port), connTimeout)
	if err != nil {
		return nil, err
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := tutorial.NewRpcServiceClientFactory(transportFactory.GetTransport(socket), protocolFactory)

	err = client.Transport.Open()
	if err != nil {
		return nil, err
	}

	return &thriftPool.IdleClient{
		Client: client,
		Socket: socket,
	}, nil
}

func Close(c *thriftPool.IdleClient) error {
	err := c.Socket.Close()
	//err = c.Client.(*tutorial.PlusServiceClient).Transport.Close()
	return err
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano()
}

func RpcHelloTest() (msg string, err error) {
	client, err := GlobalRpcPool.Get()
	if err != nil {
		return
	}
	msg, err = client.Client.(*tutorial.RpcServiceClient).Hello()
	err = GlobalRpcPool.Put(client)
	if err != nil {
		msg = "Put Error"
	}
	return
}
