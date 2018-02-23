package proxy

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/examples/proxy/testpb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. testpb/example.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. testpb/example.proto

type EchoService struct{}

func (EchoService) Echo(ctx context.Context, r *testpb.Request) (*testpb.Response, error) {
	return &testpb.Response{
		Data: r.Name,
	}, nil
}

func TestNormal(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := testpb.RegisterExampleServiceHandlerServer(ctx, mux, EchoService{})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	fmt.Println(mux)
	fmt.Println(ts.URL)

	res, err := http.Get(ts.URL + "/v1/echo?name=test")
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	t.Logf("%+v", res.Header)
	t.Logf(string(greeting))
	if len(greeting) == 0 {
		t.Fail()
	}
}
