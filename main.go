package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gihyeonsung/file/internal/application"
	"github.com/gihyeonsung/file/internal/infrastructure/persistence"
	"github.com/gihyeonsung/file/internal/infrastructure/presentation"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "db/file.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fileRepository := persistence.NewSqliteFileRepository(db)
	err = fileRepository.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	fileService := persistence.NewLocalFileService("data")

	fileCreate := application.NewFileCreate(fileRepository)
	fileDelete := application.NewFileDelete(fileRepository, fileService)
	fileFind := application.NewFileFind(fileRepository)
	fileDownload := application.NewFileDownload(fileRepository, fileService)
	fileUpload := application.NewFileUpload(fileRepository, fileService)

	mux := http.NewServeMux()
	presentation.NewFileController(mux, fileCreate, fileDelete, fileFind, fileDownload, fileUpload)

	log.Printf("Start listening on port 8080")
	http.ListenAndServe(":8080", mux)
}
