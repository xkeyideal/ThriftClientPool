package client

import (
	"errors"
	"fmt"
	"testing"
	"thriftPool"
	"time"
	"tutorial"

	"github.com/astaxie/beego/httplib"
)

func TestTcpConn(t *testing.T) {
	GlobalRpcPool = thriftPool.NewThriftPool("10.5.20.3", "23455", 100, 32, 600, Dial, Close)
	c1, err := GlobalRpcPool.Get()
	if err != nil {
		t.Fatalf("get conn from pool err:%v", err)
	}
	n := GlobalRpcPool.GetConnCount()

	if n != 1 {
		t.Fatalf("conn count:%d is err", n)
	}
	n = GlobalRpcPool.GetIdleCount()

	if n != 0 {
		t.Fatalf("idle count:%d is err", n)
	}

	c2, err := GlobalRpcPool.Get()
	if err != nil {
		t.Fatalf("get conn from pool err:%v", err)
	}
	n = GlobalRpcPool.GetConnCount()

	if n != 2 {
		t.Fatalf("conn count:%d is err", n)
	}
	GlobalRpcPool.Put(c1)

	n = GlobalRpcPool.GetIdleCount()

	if n != 1 {
		t.Fatalf("idle count:%d is err", n)
	}

	GlobalRpcPool.Put(c2)

	n = GlobalRpcPool.GetIdleCount()

	if n != 2 {
		t.Fatalf("idle count:%d is err", n)
	}
	GlobalRpcPool.Release()
}

func TestRpcSortPool(t *testing.T) {
	req := tutorial.NewSortDesc()
	req.Limit = 5000
	GlobalRpcPool = thriftPool.NewThriftPool("10.5.20.3", "23455", 100, 32, 600, Dial, Close)
	N := 100
	var TotalTime int64 = 0
	for i := 0; i < N; i++ {
		startTime := currentTimeMillis()
		client, err := GlobalRpcPool.Get()
		if err != nil {
			t.Fatal(err)
			continue
		}
		client.Client.(*tutorial.RpcServiceClient).Sort(req)
		//b.Log(msg)
		err = GlobalRpcPool.Put(client)
		if err != nil {
			//b.Fatal("Put Error", err)
			continue
		}
		endTime := currentTimeMillis()
		TotalTime += endTime - startTime
	}
	t.Log(TotalTime / int64(N))
	GlobalRpcPool.Release()
}

func BenchmarkRpcSortPool(b *testing.B) {
	req := tutorial.NewSortDesc()
	req.Limit = 5000
	GlobalRpcPool = thriftPool.NewThriftPool("10.5.20.3", "23455", 500, 32, 600, Dial, Close)
	for i := 0; i < b.N; i++ {
		client, err := GlobalRpcPool.Get()
		if err != nil {
			b.Fatal(err)
			continue
		}
		r, err := client.Client.(*tutorial.RpcServiceClient).Sort(req)
		if err != nil {
			b.Fatal("Put Error", err)
			//continue
		}

		if int32(len(r)) != req.Limit {
			b.Fatal(len(r))
		}

		//b.Log(msg)
		err = GlobalRpcPool.Put(client)
		if err != nil {
			//b.Fatal("Put Error", err)
			continue
		}
	}
	GlobalRpcPool.Release()
}

func BenchmarkRpcSortSingle(b *testing.B) {
	req := tutorial.NewSortDesc()
	req.Limit = 5000
	for i := 0; i < b.N; i++ {
		client, err := Dial("10.5.20.3", "23455", time.Duration(32)*time.Second)
		if err != nil {
			b.Fatal(err)
			continue
		}
		r, err := client.Client.(*tutorial.RpcServiceClient).Sort(req)
		if err != nil {
			b.Fatal("Put Error", err)
			//continue
		}

		if int32(len(r)) != req.Limit {
			b.Fatal(len(r))
		}
		Close(client)
	}
}

func BenchmarkHttpSort(b *testing.B) {
	url := "http://10.5.20.3:23456/sortasc?limit=5000"
	for i := 0; i < b.N; i++ {
		req := httplib.Get(url)

		resp, err := req.Response()

		if err != nil {
			b.Fatal(err)
			continue
		}

		if resp.StatusCode == 200 {
			r := []int32{}
			req.ToJSON(r)
			if len(r) != 5000 {
				b.Fatal(len(r))
			}
		} else {
			//req.String()
			err = errors.New(fmt.Sprintf("Http Connect Error, Status:%s", resp.Status))
			b.Fatal(err)
		}
	}
}

func BenchmarkRpcHelloPool(b *testing.B) {
	GlobalRpcPool = thriftPool.NewThriftPool("10.5.20.3", "23455", 100, 32, 600, Dial, Close)
	for i := 0; i < b.N; i++ {
		client, err := GlobalRpcPool.Get()
		if err != nil {
			//b.Fatal(err)

			continue
		}
		client.Client.(*tutorial.RpcServiceClient).Hello()
		//b.Log(msg)
		err = GlobalRpcPool.Put(client)
		if err != nil {
			//b.Fatal("Put Error", err)
			continue
		}
	}
	GlobalRpcPool.Release()
}

func BenchmarkRpcHelloSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client, err := Dial("10.5.20.3", "23455", time.Duration(32)*time.Second)
		if err != nil {
			b.Fatal(err)

			continue
		}
		client.Client.(*tutorial.RpcServiceClient).Hello()
		//b.Log(msg)
		//		err = GlobalRpcPool.Put(client)
		//		if err != nil {
		//			//b.Fatal("Put Error", err)
		//			continue
		//		}
		Close(client)
	}
}

func BenchmarkHttpHello(b *testing.B) {
	url := "http://10.5.20.3:23456/hello"
	for i := 0; i < b.N; i++ {
		req := httplib.Get(url)

		resp, err := req.Response()

		if err != nil {
			//b.Fatal(err)
			continue
		}

		if resp.StatusCode == 200 {
			req.String()
			//b.Log(msg)
		} else {
			req.String()
			//		err = util.NewError("Http Connect Error, Status:%s, %s", resp.Status, msg)
			//b.Fatal(err)
		}
	}
}
