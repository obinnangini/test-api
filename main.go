package main

import (
	"context"
	"fmt"
	"net/http"
	"ngini.com/test-api/internal/api"
	"ngini.com/test-api/internal/dao"
	"os"
	"strconv"
)

func main() {
	var portNum int
	var err error
	var db dao.DAO

	apiPort := os.Getenv("API_PORT")
	if len(apiPort) == 0 {
		panic("API_PORT environment variable not set")
	}
	portNum, err = strconv.Atoi(apiPort)
	if err != nil {
		panic(err)
	}

	db = dao.NewMemoryDAO()
	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) > 0 {
		db = setupDB(context.Background(), dbUrl)
	}
	r := api.SetUpRouter(db)

	address := fmt.Sprintf(":%d", portNum)
	err = http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}

func setupDB(ctx context.Context, dbUrl string) *dao.DBDAO {
	dbDAO := dao.NewDBDAO()
	err := dbDAO.InitConnection(ctx, dbUrl)
	if err != nil {
		panic(err)
	}
	return dbDAO
}
