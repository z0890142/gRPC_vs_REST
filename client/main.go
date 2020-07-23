package main

import (
	"bytes"
	"context"
	"encoding/json"
	pb "gRPC_test/pd"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var c pb.EchoClient

func main() {
	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
	defer conn.Close()

	c = pb.NewEchoClient(conn)

	r := mux.NewRouter()

	r.HandleFunc("/EchoHttp", EchoHttp).Methods("POST")
	r.HandleFunc("/EchoGRPC", EchoGRPC).Methods("POST")

	err = http.ListenAndServe(":8088", r)
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
}

type Request struct {
	Number1 int32
	Number2 int32
}

func EchoGRPC(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()

	var request pb.EchoRequest
	json.Unmarshal(body, &request)
	resp, err := c.Echo(context.Background(), &request)
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}
	responseByte, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseByte)
}

func EchoHttp(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //io.LimitReader限制大小
	defer r.Body.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:80/EchoHttp", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("http client DO : %v", err)
	}
	respbody, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576)) //io.LimitReader限制大小
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respbody)
}
