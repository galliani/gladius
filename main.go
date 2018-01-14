package main

import (
    "encoding/json"

    "github.com/apex/go-apex"
)

func main() {
    apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
        // Here we initialize the db and then assign it to a global var of DbCon which is of type *gorm.DB
        // as defined in models.go
        DbCon = InitializeDatabase()
        Run()

        return nil, nil
    })  
}