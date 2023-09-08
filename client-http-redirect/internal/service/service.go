package service

import (
	"log"
	"context"
	"time"
	"encoding/json"

	"github.com/mitchellh/mapstructure"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/golang/protobuf/jsonpb"

	pb "github.com/golang/protobuf/proto"

	proto "github.com/go-grpc-service/proto"

	"github.com/go-grpc-service/client-http-redirect/internal/server/grpc"
	"github.com/go-grpc-service/client-http-redirect/internal/domain"
)

type WorkerService struct {
	infoPod 	domain.InfoPod
	GrpcClient 	grpc.GrpcClient
}

func NewWorkerService(	infoPod domain.InfoPod,
						grpcClient grpc.GrpcClient ) *WorkerService{
	log.Print("NewWorkerService")

	return &WorkerService{
		infoPod: infoPod,
		GrpcClient: grpcClient,
	}
}

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

//-------------------
func (s WorkerService) Version() (domain.InfoPod, error){
	log.Print("Version")
	return s.infoPod, nil
}

func (s WorkerService) Info() (domain.InfoPod, error){
	log.Print("Info")
	return s.infoPod, nil
}

func (s WorkerService) InfoPodGrpc() (domain.InfoPod, error){
	log.Print("InfoPodGrpc")

	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, header)

	data := &proto.PodInfoRequest {}
	client := s.GrpcClient.GetConnection()

	response, err := client.GetPodInfo(ctx, data)
	if err != nil {
	  log.Fatalf("could not GetPodInfo: %v", err)
	  return domain.InfoPod{}, err
	}
	response_str, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
		return domain.InfoPod{}, err
  	}

	var result_final map[string]interface{}
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		log.Fatalf("could not Unmarshal: %v", err)
		return domain.InfoPod{}, err
	}

	result_filtered := result_final["podInfo"].(map[string]interface{})
	var podInfo domain.InfoPod
	err = mapstructure.Decode(result_filtered, &podInfo)
	if err != nil {
		log.Fatalf("could not mapstructure: %v", err)
		return domain.InfoPod{}, err
	}

	return podInfo, nil
}

func (s WorkerService) Add() (bool, error){
	log.Print("Add")

	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, header)

	ts := timestamppb.Now()
	balance := proto.Balance{	Id: "id", 
								Account: "acc", 
								Amount: 1, 
								DateBalance: ts, 
								Description: "description" }

	data:= &proto.AddBalanceRequest {
		Balance: &balance,
	}

	client := s.GrpcClient.GetConnection()

  	response, err := client.AddBalance(ctx, data)
  	if err != nil {
    	log.Fatalf("could not AddBalance: %v", err)
		return false, err
  	}

	response_str, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
		return false, err
  	}
  	log.Printf(response_str)

	var result_final bool
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		log.Println(err)
	}
	log.Print(result_final)

	return true, nil
}

func (s WorkerService) Get() (domain.Balance, error){
	log.Print("Get")
	
	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, header)

	data := &proto.GetBalanceRequest {
		Id: "id-001",
	}

	client := s.GrpcClient.GetConnection()
	
 	response, err := client.GetBalance(ctx, data)
  	if err != nil {
    	log.Fatalf("could not GetBalance: %v", err)
		return domain.Balance{}, err
  	}

	response_str, err := ProtoToJSON(response)
	if err != nil {
    	log.Fatalf("could not ProtoToJSON: %v", err)
		return domain.Balance{}, err
  	}

	var result_final map[string]interface{}
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		log.Fatalf("could not Unmarshal: %v", err)
		return domain.Balance{}, err
	}

	result_filtered := result_final["balance"].(map[string]interface{})
	var balance domain.Balance
	err = mapstructure.Decode(result_filtered, &balance)
	if err != nil {
		log.Fatalf("could not mapstructure: %v", err)
		return domain.Balance{}, err
	}

	return balance, nil
}

func (s WorkerService) List() ([]domain.Balance, error){
	log.Print("List")

	header := metadata.New(map[string]string{"client-id": "client-001", "authorization": "Beared cookie"})
	ctx, cancel := context.WithTimeout(context.Background(), 6 * time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, header)

	data := &proto.ListBalanceRequest {}

	client := s.GrpcClient.GetConnection()
	
	response, err := client.ListBalance(ctx, data)
	if err != nil {
	  log.Fatalf("could not ListBalance: %v", err)
	  return []domain.Balance{}, err
	}

	response_str, err := ProtoToJSON(response)
	if err != nil {
	  log.Fatalf("could not ProtoToJSON: %v", err)
	  return []domain.Balance{}, err
	}

	var result_final map[string]interface{}
	err = json.Unmarshal([]byte(response_str), &result_final)
	if err != nil {
		log.Fatalf("could not Unmarshal: %v", err)
		return []domain.Balance{}, err
	}

	result_filtered := result_final["balance"].([]interface{})
	var list_balance []domain.Balance
	for _,item:=range result_filtered {
		//log.Printf("%v", item.(map[string]interface{}))
		balance := domain.Balance{}
		err = mapstructure.Decode(item, &balance)
		if err != nil {
			log.Fatalf("could not mapstructure: %v", err)
			return []domain.Balance{}, err
		}
		list_balance = append(list_balance,balance)
	}

	return list_balance, nil
}