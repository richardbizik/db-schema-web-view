package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/richardbizik/db-schema-web-view/database/postgres"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const port = "9090"

var addr = flag.String("addr", ":"+port, "http service address")
var wait = flag.Duration("graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/", serveHome)
	r.Path("/database").Methods("GET").HandlerFunc(getDatabases)
	r.Path("/database/{name}/{schema}").Methods("GET").HandlerFunc(getSchema)
	http.Handle("/", r)

	flag.Parse()
	srv := &http.Server{
		Handler:      r,
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Print("Server running at port: " + port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func getDatabases(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("connections.yaml")
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Println("Failed to read schema")
	}
	var response []GetDatabasesResp
	for _, database := range config.Databases {
		response = append(response, GetDatabasesResp{
			Name:   database.Name,
			Schema: database.Schema,
		})
	}
	marshal, _ := json.Marshal(response)
	w.Write(marshal)
}

func getSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dbName := vars["name"]
	dbSchema := vars["schema"]
	file, err := ioutil.ReadFile("connections.yaml")
	var config Config
	var foundConfig DatabaseConfiguration
	err = yaml.Unmarshal(file, &config)
	for i := range config.Databases {
		if config.Databases[i].Name == dbName {
			foundConfig = config.Databases[i]
			break
		}
	}
	schema, err := postgres.GetDBSchema(dbName, foundConfig.Host, foundConfig.Port, foundConfig.User, foundConfig.Password, dbSchema)
	if err != nil {
		log.Println("Failed to read schema")
	}
	marshal, err := json.Marshal(schema)
	w.Write(marshal)
}

type GetDatabasesResp struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type Config struct {
	Databases []DatabaseConfiguration `yaml:"databases"`
}

type DatabaseConfiguration struct {
	Name     string `yaml:"name"`
	Schema   string `yaml:"schema"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}
