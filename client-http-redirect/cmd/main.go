package main

import (
	"fmt"
	"net"
	"strconv"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"log"
	"os"

	"time"

	"github.com/go-grpc-service/client-http-redirect/internal/handler"
	"github.com/go-grpc-service/client-http-redirect/internal/server"

	"github.com/go-grpc-service/client-http-redirect/internal/server/grpc"
	"github.com/go-grpc-service/client-http-redirect/internal/service"
	"github.com/go-grpc-service/client-http-redirect/internal/domain"

	pb "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
)

var (
	my_secret_loaded_from_volume 	domain.Secret
	my_info_pod						domain.InfoPod
	PORT = 3000
	HOST = ":50051"
	API_VERSION = "no-version"
	POD_NAME = "pod no-name"
	POD_PATH = ""
	JWTKEY = "my_secret_key"
)

func ProtoToJSON(msg pb.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	return marshaler.MarshalToString(msg)
}

func JSONToProto(data string, msg pb.Message) error {
	return jsonpb.UnmarshalString(data, msg)
}

func infoPodGrpc(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> infoPodGrpc => " + r.Method + " => path:  " + r.URL.Path)
	
	json.NewEncoder(w).Encode("result")
	return
}

func addBalance(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> addBalance => " + r.Method + " => path:  " + r.URL.Path)
	
	json.NewEncoder(w).Encode("intVar")
	return
}

func getBalance(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> getBalance => " + r.Method + " => path:  " + r.URL.Path)
	
	json.NewEncoder(w).Encode("intVar")
	return
}

func listBalance(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> listBalance => " + r.Method + " => path:  " + r.URL.Path)
	
	json.NewEncoder(w).Encode("intVar")
	return
}

//------------------------
func setInfoPod(){
	fmt.Println("==> setInfoPod")

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error to get the POD IP address !!!", err)
		os.Exit(3)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				my_info_pod.Ip = ipnet.IP.String()
			}
		}
	}
	my_info_pod.OsPid = strconv.Itoa(os.Getpid())

	if os.Getenv("API_VERSION") !=  "" {
		API_VERSION = os.Getenv("API_VERSION")
	}
	my_info_pod.ApiVersion = API_VERSION

	if os.Getenv("POD_NAME") !=  "" {
		POD_NAME = os.Getenv("POD_NAME")
	}
	my_info_pod.PodName = POD_NAME

	if os.Getenv("PORT") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("PORT"))
		PORT = intVar
	}
	my_info_pod.Port = PORT

	if os.Getenv("POD_PATH") !=  "" {
		POD_PATH = os.Getenv("POD_PATH")
	}
	my_info_pod.PodPath = POD_PATH

	if os.Getenv("JWTKEY") !=  "" {
		JWTKEY = os.Getenv("JWTKEY")
	}
	my_info_pod.JwtKey = []byte(JWTKEY)

	if os.Getenv("HOST") !=  "" {
		HOST = os.Getenv("HOST")
	}
}

func init() {
	file_user, err := ioutil.ReadFile("/var/go-hello-world-web/secret/username")
    if err != nil {
        log.Fatal(err)
		os.Exit(3)
    }
	file_pass, err := ioutil.ReadFile("/var/go-hello-world-web/secret/password")
    if err != nil {
        log.Fatal(err)
		os.Exit(3)
    }
	my_secret_loaded_from_volume.Username = string(file_user)
	my_secret_loaded_from_volume.Password = string(file_pass)

	setInfoPod()
}

func main() {
	log.Printf("Starting !")

	grpcClient, err  := grpc.StartGrpcClient(HOST)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	service		:= service.NewWorkerService(my_info_pod, grpcClient)
	handler		:= handler.NewHttpAdapter(*service)
	httpServer	:= server.NewHttpServer(time.Now(), my_info_pod)
	
	httpServer.StartHttpServer(handler)
}