package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"spotify_challenge.com/src/adapter"
	"spotify_challenge.com/src/application"
	"spotify_challenge.com/src/connector"
	"spotify_challenge.com/src/handler"
	"spotify_challenge.com/src/models"
	"spotify_challenge.com/src/service"
)

func main() {
	client := http.Client{}
	config := models.NewConfig()
	adapter := adapter.NewAdapter(&client, config.SpotifyData)

	sqlObj, err := sql.Open(config.DbProvider, config.ConnectionString)
	if err != nil {
		panic(err)
	}

	conn := connector.NewConnector(sqlObj)
	service := service.NewService(conn)
	app := application.NewApplication(adapter, service)
	handler := handler.NewHandler(app)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/Write/{ISRC}", handler.Write).Methods(http.MethodPost)
	router.HandleFunc("/ReadByISRC/{ISRC}", handler.ReadByISRC)
	router.HandleFunc("/ReadByArtist/{Name}", handler.ReadByArtist)

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalln("ListenAndServer Error", err)
	}
}
