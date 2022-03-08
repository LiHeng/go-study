package resty

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
)

type Library struct {
	Name   string
	Latest string
}

type Libraries struct {
	Results []*Library
}

func TestResty(t *testing.T) {
	client := resty.New()

	resp, err := client.R().Get("https://baidu.com")
	if err != nil {
		t.FailNow()
	}
	fmt.Println("Response Info:")
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status:", resp.Status())
	fmt.Println("Proto:", resp.Proto())
	fmt.Println("Time:", resp.Time())
	fmt.Println("Received At:", resp.ReceivedAt())
	fmt.Println("Size:", resp.Size())
	fmt.Println("Headers:")
	for key, value := range resp.Header() {
		fmt.Println(key, "=", value)
	}
	fmt.Println("Cookies:")
	for i, cookie := range resp.Cookies() {
		fmt.Printf("cookie%d: name:%s value:%s\n", i, cookie.Name, cookie.Value)
	}

	resp, err = client.R().
		SetOutput("test.png").
		Get("http://www.juimg.com/tuku/yulantu/140112/328648-14011213253758.jpg")
	if err != nil {
		t.FailNow()
	}

	libraries := &Libraries{}
	resp, err = client.R().
		SetResult(libraries).
		ForceContentType("application/json").
		Get("https://api.cdnjs.com/libraries")
	if err != nil {
		t.FailNow()
	}
	for _, lib := range libraries.Results {
		fmt.Println("first library:")
		fmt.Printf("name: %s latest: %s\n", lib.Name, lib.Latest)
		break
	}
}
