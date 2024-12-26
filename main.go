package main

import (
	"context"
	"fmt"
	"net/http"
	"ngini.com/test-api/internal/api"
	"ngini.com/test-api/internal/dao"
	"os"
	"strconv"

	"go.uber.org/zap"
)

type DBStructure string

const (
	UseMemoryMap  DBStructure = "map"
	UseMemoryList DBStructure = "list"
	UseDB         DBStructure = "db"
)

func main() {
	var portNum int
	var err error
	var db dao.DAO

	logger, _ := zap.NewProduction()
	apiPort := os.Getenv("API_PORT")
	if len(apiPort) == 0 {
		panic("API_PORT environment variable not set")
	}
	portNum, err = strconv.Atoi(apiPort)
	if err != nil {
		panic(err)
	}

	dBStructure := UseMemoryMap
	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) > 0 {
		dBStructure = UseDB
	}

	switch dBStructure {
	case UseMemoryMap:
		db = dao.NewMemoryMapDAO()
	case UseDB:
		db = setupDB(context.Background(), dbUrl)
	case UseMemoryList:
		db = dao.NewMemoryListDAO()
	default:
		db = dao.NewMemoryListDAO()
	}

	logger.Info("Chosen DB Structure", zap.String("dBStructure", string(dBStructure)))
	//dbStructureString := fmt.Sprintf("Chosen DB structure: %v", dBStructure)
	//fmt.Println(dbStructureString)
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
