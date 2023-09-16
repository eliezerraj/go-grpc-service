package grpc

import (
	"log"

	"google.golang.org/grpc"
	proto "github.com/go-grpc-service/proto"

)

type GrpcClient struct {
	Client 		proto.BalanceServiceClient
	GrcpClient	*grpc.ClientConn
}

func StartGrpcClient(HOST string) (GrpcClient, error){
	log.Print("StartGrpcClient")

	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithBlock()) // Wait for ready
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) // 
	opts = append(opts, grpc.WithInsecure()) // no TLS

	conn, err := grpc.Dial(HOST, opts...)
	if err != nil {
	  log.Fatalf("erro connect to grpc server: %v", err)
	  return GrpcClient{}, err
	}

	client := proto.NewBalanceServiceClient(conn)
	log.Printf("Grpc Client running... %v", client )

	return GrpcClient{
		Client: client,
		GrcpClient : conn,
	}, nil
}

func (s GrpcClient) GetConnection() (proto.BalanceServiceClient) {
	log.Printf("GetConnection") 
	return s.Client
}

func (s GrpcClient) CloseConnection() () {
	log.Printf("CloseConnection") 
	if err := s.GrcpClient.Close(); err != nil {
		log.Printf("Failed to close gPRC connection: %s", err)
	}
}