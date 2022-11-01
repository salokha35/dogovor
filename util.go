package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Response ...
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// Send ...
func (res *Response) Send(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	w.WriteHeader(200)

	if res.Payload == nil && res.Code != http.StatusOK {
		res.Payload = struct {
			Error   bool   `json:"error,omitempty"`
			Message string `json:"message,omitempty"`
		}{
			Error:   true,
			Message: strings.ToUpper(string(res.Message[:2])) + string(res.Message[2:]),
		}
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println("ERROR sending response failed:", err)
	}
	return
}

// AddCors ...
func AddCors(router *httprouter.Router) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodDelete, http.MethodGet, http.MethodPost, http.MethodPut},
		AllowedHeaders:   []string{"*"},
		MaxAge:           10,
		AllowCredentials: true,
	}).Handler(router)
}

// Unique ...
func Unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	// turn the map keys into a slice
	uniqSlice := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		replace := strings.ReplaceAll(v, "/#", "")
		uniqSlice = append(uniqSlice, replace)
	}
	return uniqSlice
}
