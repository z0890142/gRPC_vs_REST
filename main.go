package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "gRPC_test/pd"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type EchoServer struct{}

func (e *EchoServer) Echo(ctx context.Context, req *pb.EchoRequest) (resp *pb.EchoReply, err error) {

	result := cal(req.Number1, req.Number2)
	return &pb.EchoReply{
		Result: result,
	}, nil

}

func cal(n1 int32, n2 int32) int32 {
	return n1 * n2
}

func main() {
	go startGRPC()
	r := mux.NewRouter()
	r.HandleFunc("/EchoHttp", EchoGRPC).Methods("POST")
	fmt.Println("server start")

	http.ListenAndServe(":80", r)
}

func startGRPC() {
	apiListener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Println(err)
		return
	}

	// 註冊 grpc
	es := &EchoServer{}

	grpc := grpc.NewServer()
	pb.RegisterEchoServer(grpc, es)

	reflection.Register(grpc)
	if err := grpc.Serve(apiListener); err != nil {
		log.Fatal(" grpc.Serve Error: ", err)
		return
	}
}

type Request struct {
	Number1 int32
	Number2 int32
}

func EchoGRPC(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()

	var resuest Request
	json.Unmarshal(body, &resuest)
	result := cal(resuest.Number1, resuest.Number2)
	response := pb.EchoReply{
		Result: result,
	}

	responseByte, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseByte)
}
