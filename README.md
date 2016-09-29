#Thrift Pool & Demo

## Introduction

Thrift version is 0.9.3 [Thrift](https://github.com/apache/thrift)

Use gin framework write the http server [Gin](https://github.com/gin-gonic/gin)

Use beego's httplib library write the http client [Beego/httplib](https://github.com/astaxie/beego/tree/master/httplib)


## Hierarchy

src/

    client/

        contains the http/rpc client code and client test code

    server/

        contains the http/rpc server code

    thriftPool/

        contains the Thrift Client Connection Pool code

    tutoral/

        contains the Thrift Gen golang service code by idl file

    tutoral.thrift

        the idl file

    main.go

        the server start code


## Benchmarks

Pool VS Single VS Http

The Sort is server 5000 int32 elements array quicksort and send the array data to client
The Hello is server send "Hello World" string to client

Benchmark name                  | loops     | ns/op
--------------------------------|----------:|----------:
BenchmarkRpcSortPool            |  500      | 3263670
BenchmarkRpcSortSingle          |  500      | 3818633
BenchmarkHttpSort               |  200      | 9296249
BenchmarkRpcHelloPool           |  10000    | 211428
BenchmarkRpcHelloSingle         |   2000    | 571652
BenchmarkHttpHello              |   2000    | 859478


## Code folder

1. Pool code
    ```go
    thriftPool/thrift_pool.go
    ```

2. server code & start:

    ```go
    server/rpcserver.go
    server/httpserver.go

    main.go
    ```

3. client test code

    ```go
    client/client_test.go
    ```

4. thrift idl file
    ```go
    tutorial.thrift
    ```

## Start using the ThriftPool

Your must write the client Dial and Close function, see the client/rpcclient.go

1. Init Pool

```go
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

GlobalRpcPool = thriftPool.NewThriftPool("10.5.20.3", "23455", 100, 32, 600, Dial, Close)
```

2. Get a connection Client from Pool & free it after used

```go
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
```

## Testing

    ```go
        cd ThriftPool/src/client
        go test -v
        go test -test.bench=".*"
    ```

## Other

The ThriftPool supports many servers, please see the thrift_pool.go

```go
type MapPool struct {
    Dial  ThriftDial
    Close ThriftClientClose

    lock *sync.Mutex

    idleTimeout uint32
    connTimeout uint32
    maxConn     uint32

    pools map[string]*ThriftPool
}

func NewMapPool(maxConn, connTimeout, idleTimeout uint32,
    dial ThriftDial, closeFunc ThriftClientClose) *MapPool {

    return &MapPool{
        Dial:        dial,
        Close:       closeFunc,
        maxConn:     maxConn,
        idleTimeout: idleTimeout,
        connTimeout: connTimeout,
        pools:       make(map[string]*ThriftPool),
        lock:        new(sync.Mutex),
    }
}

func (mp *MapPool) getServerPool(ip, port string) (*ThriftPool, error) {
    addr := fmt.Sprintf("%s:%s", ip, port)
    mp.lock.Lock()
    serverPool, ok := mp.pools[addr]
    if !ok {
        mp.lock.Unlock()
        err := errors.New(fmt.Sprintf("Addr:%s thrift pool not exist", addr))
        return nil, err
    }
    mp.lock.Unlock()
    return serverPool, nil
}
```