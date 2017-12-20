package server

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/xkeyideal/ThriftClientPool/tutorial"
)

const (
	NetworkAddr = "0.0.0.0:23455"
)

func randInt(limit int32) []int32 {
	rand.Seed(time.Now().UnixNano())
	x := []int32{}
	var i int32 = 0
	for ; i < limit; i++ {
		x = append(x, rand.Int31n(20000000))
	}
	return x
}

func less(a, b int32, asc bool) bool {
	if asc {
		return a < b
	}
	return a > b
}

func partition(arr []int32, l, r int, asc bool) int {
	m := (l + r) >> 1
	arr[m], arr[r] = arr[r], arr[m]
	x := arr[r]
	i := l - 1

	for j := l; j < r; j++ {
		if less(arr[j], x, asc) {
			i += 1
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	i += 1
	arr[i], arr[r] = arr[r], arr[i]
	return i
}

func qsort(arr []int32, l, r int, asc bool) {
	if l < r {
		q := partition(arr, l, r, asc)
		qsort(arr, l, q-1, asc)
		qsort(arr, q+1, r, asc)
	}
}

type RpcServiceImpl struct{}

func (p *RpcServiceImpl) Plus(req *tutorial.Node) (r int32, err error) {
	return req.A + req.B, nil
}
func (p *RpcServiceImpl) Hello() (r string, err error) {
	return "Rpc Hello world", nil
}

func (s *RpcServiceImpl) Sort(sd *tutorial.SortDesc) (r []int32, err error) {
	//startTime := currentTimeMillis()
	arr := randInt(sd.Limit)
	qsort(arr, 0, len(arr)-1, sd.Asc)
	r = arr
	//	for i := 0; i < 10; i++ {
	//		r = append(r, arr[i])
	//	}
	//endTime := currentTimeMillis()
	//fmt.Println("Tcp time->", endTime, startTime, (endTime - startTime))
	return
}

func RunServer() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}
	handler := &RpcServiceImpl{}
	processor := tutorial.NewRpcServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("thrift server in", NetworkAddr)
	err = server.Serve()
	if err != nil {
		fmt.Println("Server Run Error: ", err)
		os.Exit(1)
	}
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano()
}
