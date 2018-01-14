package main

import (
    "encoding/json"
    "./models"
    "github.com/apex/go-apex"
)

func main() {
    apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
        models.DbCon = models.InitializeDatabase()
        Run()

        return nil, nil
    })  
}