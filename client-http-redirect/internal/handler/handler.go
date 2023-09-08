package handler

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-grpc-service/client-http-redirect/internal/service"
	"github.com/go-grpc-service/client-http-redirect/internal/domain"
	"github.com/go-grpc-service/client-http-redirect/internal/erro"
)

var (
	jwtKey = []byte("my_secret_key")
)

type HttpAdapter struct {
	service	service.WorkerService
}

func NewHttpAdapter( service service.WorkerService) *HttpAdapter {
	return &HttpAdapter{
		service: service,
	}
}

func (h *HttpAdapter) MiddleWareHandlerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-------------------------------------------- \n")
		log.Println("MiddleWareHandlerToken (INICIO)")

		token := r.Header.Get("Authorization")
		tokenSlice := strings.Split(token, " ")
		var bearerToken string

		if len(tokenSlice) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(erro.ErrStatusUnauthorized.Error())
			return
		}
		
		bearerToken = tokenSlice[len(tokenSlice)-1]
		res, err := ScopeValidation(bearerToken,"","")
		if err != nil || res != true {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		log.Println("MiddleWareHandlerToken (FIM)")
		log.Printf("-------------------------------------------- \n")
		
		next.ServeHTTP(w, r)
	})
}

func ScopeValidation(token string, path string, method string) (bool, error){
	fmt.Println("==> ScopeValidation :" , token )

	claims := &domain.JwtData{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, erro.ErrStatusUnauthorized
		}
		return false, erro.ErrTokenExpired
	}

	if !tkn.Valid {
		return false, erro.ErrStatusUnauthorized
	}

	return true, nil
}

func (h *HttpAdapter) Authorization(w http.ResponseWriter, r *http.Request) {
	log.Println("==> check => " + r.Method + " => path:  " + r.URL.Path)
	json.NewEncoder(w).Encode(domain.ResponseMessage{Message: "authorization true"})
	return
}

func (h *HttpAdapter) Health(w http.ResponseWriter, r *http.Request) {
	log.Println("==> check => " + r.Method + " => path:  " + r.URL.Path)
	json.NewEncoder(w).Encode(domain.ResponseMessage{Message: "health true"})
	return
}

func (h *HttpAdapter) Live(w http.ResponseWriter, r *http.Request) {
	log.Println("==> live => " + r.Method + " => path:  " + r.URL.Path)
	json.NewEncoder(w).Encode(domain.ResponseMessage{Message: "live true"})
	return
}

func (h *HttpAdapter) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> Index => " + r.Method + " => path:  " + r.URL.Path)
	json.NewEncoder(w).Encode(domain.ResponseMessage{Message: "hello world web"})
	return
}

func (h *HttpAdapter) Version(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> version => " + r.Method + " => path:  " + r.URL.Path)

	res, err := h.service.Version()
	if err != nil{
		log.Printf("ERRO => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}

func (h *HttpAdapter) Info(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> Info => " + r.Method + " => path:  " + r.URL.Path)

	res, err := h.service.Info()
	if err != nil{
		log.Printf("ERRO => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}

func (h *HttpAdapter) InfoPodGrpc(w http.ResponseWriter, r *http.Request){
	fmt.Println("==> InfoPodGrpc => " + r.Method + " => path:  " + r.URL.Path)

	res, err := h.service.InfoPodGrpc()
	if err != nil{
		log.Printf("ERRO => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}

func (h *HttpAdapter) Header(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> header => " + r.Method + " => path:  " + r.URL.Path)

	for k, v := range r.Header { //Iterate over all header fields
		fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
	}
	return
}

func (h *HttpAdapter) Add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> Add => " + r.Method + " => path:  " + r.URL.Path)

	res, err := h.service.Add()
	if err != nil{
		log.Printf("ERRO Add => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}

func (h *HttpAdapter) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> Get => " + r.Method + " => path:  " + r.URL.Path)

	res, err := h.service.Get()
	if err != nil{
		log.Printf("ERRO Get => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}

func (h *HttpAdapter) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==> List => " + r.Method + " => path:  " + r.URL.Path)

	res, err := h.service.List()
	if err != nil{
		log.Printf("ERRO List => %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(res)
	return
}