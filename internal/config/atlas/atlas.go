package main

import (
	"ariga.io/atlas-provider-gorm/gormschema"
	"fmt"
	"io"
	"itv-movie/internal/models"
	"log"
	"os"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&models.Country{}, &models.Genre{}, &models.Language{}, &models.Movie{}, &models.Session{}, &models.User{})
	if err != nil {
		msg := fmt.Sprintf("failed to load gorm schema: %v\n", err)
		log.Print(msg)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
