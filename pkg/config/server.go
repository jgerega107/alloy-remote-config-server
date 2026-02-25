package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func Start() {
	httpPort := 8080
	httpPortEnv := os.Getenv("HTTP_PORT")
	if len(httpPortEnv) > 0 {
		httpPort, _ = strconv.Atoi(httpPortEnv)
	}
	grpcPort := 8888
	grpcPortEnv := os.Getenv("GRPC_PORT")
	if len(grpcPortEnv) > 0 {
		grpcPort, _ = strconv.Atoi(grpcPortEnv)
	}
	listenAddr := "0.0.0.0"
	listenAddrEnv := os.Getenv("LISTEN_ADDR")
	if len(listenAddrEnv) > 0 {
		listenAddr = listenAddrEnv
	}
	configFolder := "conf"
	configFolderEnv := os.Getenv("CONFIG_FOLDER")
	if len(configFolderEnv) > 0 {
		configFolder = configFolderEnv
	}
	err := LoadTemplates(configFolder)
	if err != nil {
		log.Println(fmt.Sprintf("Error loading templates: %v", err))
	}
	go StartConnectGrpcServer(listenAddr, grpcPort)
	StartRestServer(listenAddr, httpPort)
}
