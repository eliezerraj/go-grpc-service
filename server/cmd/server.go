package main

import (
  "context"
  "log"
  "net"
  "time"
  "os"
  "os/signal"
  "syscall"

  "google.golang.org/grpc"
  "google.golang.org/grpc/status"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/keepalive"
  "google.golang.org/protobuf/types/known/timestamppb"
  "google.golang.org/grpc/health/grpc_health_v1"

  "github.com/go-grpc-service/server/healthcheck"

  proto "github.com/go-grpc-service/proto"
)

var (
	HOST = ":50052"
	POD_IP = "X.X.X.X" 
)

type server struct{}

func (s *server) GetBalance(ctx context.Context, in *proto.GetBalanceRequest) (*proto.GetBalanceResponse, error) {
	log.Println("GetBalance in data :", in)

	ts := timestamppb.Now()
	balance := proto.Balance{	Id: "id-01", 
								Account: "acc-01", 
								Amount: 100, 
								DateBalance: ts, 
								Description: "description-01" }

	res := &proto.GetBalanceResponse {
		Balance: &balance,
	}

	return res, nil
}

func (s *server) AddBalance(ctx context.Context, in *proto.AddBalanceRequest) (*proto.AddBalanceResponse, error) {
	log.Println("AddBalance in data :", in)

	res := &proto.AddBalanceResponse {
		Result: true,
	}

	return res, nil
}

func (s *server) ListBalance(ctx context.Context, in *proto.ListBalanceRequest) (*proto.ListBalanceResponse, error) {
	log.Println("ListBalance")

	var array_balance []*proto.Balance
	
	for i:=0; i < 5; i++ {

		ts := timestamppb.Now()
		balance := proto.Balance{	Id: "id-01", 
									Account: "acc-01", 
									Amount: 100 + int32(i), 
									DateBalance: ts, 
									Description: "description-01" }

		array_balance = append(array_balance, &balance)							
	}

	res := &proto.ListBalanceResponse {
		Balance: array_balance,
	}

	return res, nil
}

func (s *server) GetPodInfo(ctx context.Context, in *proto.PodInfoRequest) (*proto.PodInfoResponse, error) {
	log.Println("GetPodInfo")

	podInfo := proto.PodInfo{Ip: POD_IP }

	res := &proto.PodInfoResponse {
		PodInfo: &podInfo,
	}

	return res, nil
}

func middleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error){
	log.Printf("-------------------------------------------- \n")
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "INTERNAL_SERVER_ERROR")
	}
	
	log.Printf("middleware : %v", info.FullMethod)
	log.Printf("req : %v", req)
	log.Printf("headers[client-id] : %v", headers["client-id"])
	log.Printf("headers[authorization] : %v", headers["authorization"])

	if len(headers["authorization"]) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Not Authorized")
	}

	log.Printf("-------------------------------------------- \n")
	return handler(ctx, req) 
}

func init(){
	if os.Getenv("HOST") !=  "" {
		HOST = os.Getenv("HOST")
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Error to get the POD IP address !!!", err)
		os.Exit(3)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				POD_IP = ipnet.IP.String()
			}
		}
	}
}

func main() {
	log.Println("RegisterBalanceServiceServer - Start")

	lis, err := net.Listen("tcp", HOST)
  	if err != nil {
    	log.Fatalf("failed to listen: %v", err)
  	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(middleware))
	opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
												MaxConnectionAge: time.Second * 30,
												MaxConnectionAgeGrace: time.Second * 10,
											}))
 
  	grpcServer := grpc.NewServer(opts...)
  	proto.RegisterBalanceServiceServer(grpcServer, &server{})

	healthService := healthcheck.NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)

	go func(){
		log.Println("Starting server..." + HOST)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to server %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM )
	<-ch

	log.Println("Stopping server")
	grpcServer.Stop()

	log.Println("Stopping listener")
	lis.Close()

	log.Println("Done !!!")
}
