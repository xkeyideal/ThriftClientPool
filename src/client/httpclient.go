package client

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/httplib"
)

func httpTest() {
	url := "http://10.14.40.112:23456/sortasc?limit=5000000"

	for i := 0; i < 10; i++ {
		startTime := currentTimeMillis()
		req := httplib.Get(url)

		resp, err := req.Response()

		if err != nil {
			return
		}

		if resp.StatusCode == 200 {
			r := []int32{}
			err = req.ToJSON(&r)
			fmt.Println(r)
		} else {
			msg, _ := req.String()
			err = errors.New(fmt.Sprintf("Http Connect Error, Status:%s, %s", resp.Status, msg))
		}
		endTime := currentTimeMillis()
		fmt.Println("HTTP time->", endTime, startTime, (endTime - startTime))
	}
}

func HttpHelloTest() (msg string, err error) {
	url := "http://10.14.40.112:23456/hello"
	req := httplib.Get(url)

	resp, err := req.Response()

	if err != nil {
		return
	}

	if resp.StatusCode == 200 {
		msg, err = req.String()
	} else {
		msg, _ = req.String()
		err = errors.New(fmt.Sprintf("Http Connect Error, Status:%s, %s", resp.Status, msg))
	}
	return
}
