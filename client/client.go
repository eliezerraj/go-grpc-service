package main

import (
  "context"
  "log"
  "time"

  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "google.golang.org/protobuf/types/known/timestamppb"

  "github.com/golang/protobuf/jsonpb"
  pb "github.com/golang/protobuf/proto"

  proto "github.com/go-grpc-service/proto"
)

var HOST = ":50051" 

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

func ClientGetBalance(ctx context.Context, client proto.BalanceServiceClient){
	log.Println("ClientGetBalance")

	data := &proto.GetBalanceRequest {
		Id: "id-001",
	}

  	response, err := client.GetBalance(ctx, data)
  	if err != nil {
    	log.Fatalf("could not GetBalance: %v", err)
  	}

	result, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
  	}
  	log.Printf("Balance: %s", result)
}

func ClientAddBalance(ctx context.Context, client proto.BalanceServiceClient){
	log.Println("ClientAddBalance")

	ts := timestamppb.Now()
	balance := proto.Balance{	Id: "id", 
								Account: "acc", 
								Amount: 1, 
								DateBalance: ts, 
								Description: "description" }

	data:= &proto.AddBalanceRequest {
		Balance: &balance,
	}

  	response, err := client.AddBalance(ctx, data)
  	if err != nil {
    	log.Fatalf("could not AddBalance: %v", err)
  	}

	result, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
  	}
  	log.Printf("Balance: %s", result)
}

func ClientListBalance(ctx context.Context, client proto.BalanceServiceClient){
	log.Println("ClientListBalance")

	data := &proto.ListBalanceRequest {}

  	response, err := client.ListBalance(ctx, data)
  	if err != nil {
    	log.Fatalf("could not ListBalance: %v", err)
  	}

	result, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
  	}
  	log.Printf("Balance: %s", result)
}

func ClientPodInfo(ctx context.Context, client proto.BalanceServiceClient){
	log.Println("ClientPodInfo")

	data := &proto.PodInfoRequest {}

  	response, err := client.GetPodInfo(ctx, data)
  	if err != nil {
    	log.Fatalf("could not GetPodInfo: %v", err)
  	}

	result, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
  	}
  	log.Printf("Data: %s", result)
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithBlock()) // Wait for ready
	opts = append(opts, grpc.WithInsecure()) // no TLS

	conn, err := grpc.Dial(HOST, opts...)
  	if err != nil {
    	log.Fatalf("did not connect: %v", err)
  	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close gPRC connection: %s", err)
		}
	}()

  	client := proto.NewBalanceServiceClient(conn)

	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
  	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	ctx = metadata.NewOutgoingContext(ctx, header)
  	defer cancel()

	ClientPodInfo(ctx, client)

	ClientGetBalance(ctx, client)

	ClientAddBalance(ctx, client)

	ClientListBalance(ctx, client)
}
