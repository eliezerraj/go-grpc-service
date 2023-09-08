package server

import (
	"time"
	"log"
	"net/http"
	"os"
	"context"
	"os/signal"
	"syscall"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-grpc-service/client-http-redirect/internal/handler"
	"github.com/go-grpc-service/client-http-redirect/internal/domain"
)

type HttpServer struct {
	start 	time.Time
	infoPod domain.InfoPod
}

func NewHttpServer(	start time.Time, 
					infoPod domain.InfoPod ) *HttpServer {
	return &HttpServer{	
		start: start,
		infoPod: infoPod,
	 }
}

func (s HttpServer) StartHttpServer(handler *handler.HttpAdapter) {
	log.Print("StartHttpServer")

	myRouter := mux.NewRouter().StrictSlash(true)

	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", handler.Health) 

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", handler.Live) 

	index := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    index.HandleFunc("/", handler.Index) 

	version := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    version.HandleFunc("/version", handler.Version) 

	info := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    info.HandleFunc("/info", handler.Info)

	infoPodGrpc := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    infoPodGrpc.HandleFunc("/infoPodGrpc", handler.InfoPodGrpc)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", handler.Header) 

	authorization := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    authorization.HandleFunc("/authorization", handler.Authorization) 
	authorization.Use(handler.MiddleWareHandlerToken)

	add := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
    add.HandleFunc("/add", handler.Add) 

	get := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    get.HandleFunc("/get", handler.Get) 

	list := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    list.HandleFunc("/list", handler.List) 

	log.Print("StartHttpServer port: ", s.infoPod.Port)
	http_srv := http.Server{
		Addr:		":" +  strconv.Itoa(s.infoPod.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(5) * time.Second,   
		WriteTimeout: time.Duration(5) * time.Second,  
		IdleTimeout:  time.Duration(5) * time.Second, 
	}

	go func() {
		err := http_srv.ListenAndServe()
		if err != nil {
			log.Print("message ==> ", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	log.Printf("Stopping Server")
	ctx , cancel := context.WithTimeout(context.Background(), time.Duration(5) * time.Second)
	defer cancel()
	
	if err := http_srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Print("WARNING Dirty Shutdown ", err)
		return
	}

	log.Printf("Stop Done !")
}