package main

import (
	"os"
	"fmt"
	"net/http"
	"context"
	"github.com/owasp-amass/asset-db/repository/neo4j"

	"github.com/sirupsen/logrus"
)


func main() {
	mux := http.NewServeMux()

	logger := logrus.New()

	loglevel, ok := os.LookupEnv("LOGLEVEL")
	if !ok {
		loglevel = "INFO"
	}
	
	ll, err := logrus.ParseLevel(loglevel)
	if err != nil {
		ll = logrus.InfoLevel
	}
	
	logger.SetLevel(ll)
	
	store, err := neo4j.New(neo4j.Neo4j, "bolt://neo4j:password@localhost:7687/neo4j")
	if err != nil {
		fmt.Println("Unable to connect to asset store: "+err.Error())
		return		
	}
	
	api := &ApiV1{
		ctx: context.Background(),
		store: store,
		bus: &EventBus{
			subscribers: make(map[chan ServerSentEvent]bool),
		},
		logger: logger,
	}

	mux.HandleFunc("GET /listen", api.ListenEvents)
	
	mux.HandleFunc("POST /emit/entity", api.CreateEntity)
	mux.HandleFunc("DELETE /emit/entity/{id}", api.DeleteEntity)
	mux.HandleFunc("PUT /emit/entity/{id}", api.UpdateEntity)
	
	mux.HandleFunc("POST /emit/edge", api.CreateEdge)
	mux.HandleFunc("DELETE /emit/edge/{id}", api.DeleteEdge)
	mux.HandleFunc("PUT /emit/edge/{id}", api.UpdateEdge)
	
	mux.HandleFunc("POST /emit/entity_tag", api.CreateEntityTag)
	mux.HandleFunc("DELETE /emit/entity_tag/{id}", api.DeleteEntityTag)
	mux.HandleFunc("PUT /emit/entity_tag/{id}", api.UpdateEntityTag)

	mux.HandleFunc("POST /emit/edge_tag", api.CreateEdgeTag)
	mux.HandleFunc("DELETE /emit/edge_tag/{id}", api.DeleteEdgeTag)
	mux.HandleFunc("PUT /emit/edge_tag/{id}", api.UpdateEdgeTag)

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
	}
	
	if err := server.ListenAndServeTLS("tls/cert.pem", "tls/key.pem"); err != nil {
		panic(err)
	}
}
