package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/chord-dht/chord-backend/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Check if a port is available
func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// UpdateConfig updates the config file with the new port
func UpdateConfig(port int, originalPort int) {
	if port != originalPort {
		configPath := "./dist/config.json"
		configData, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}
		var config map[string]interface{}
		if err := json.Unmarshal(configData, &config); err != nil {
			log.Fatalf("Failed to unmarshal config file: %v", err)
		}

		config["CHORD_ADDRESS"] = "http://localhost:" + strconv.Itoa(port)
		newConfigData, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal config file: %v", err)
		}

		if err := os.WriteFile(configPath, newConfigData, os.ModePerm); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
	}
}

func SetupStatic(r *gin.Engine) {
	r.StaticFile("/", "./dist/index.html")
	r.StaticFile("/config.json", "./dist/config.json")
	r.Static("/assets", "./dist/assets")
	r.Static("/icon", "./dist/icon")
}

func main() {
	port := 21776
	originalPort := port
	for !isPortAvailable(port) {
		port++
	}
	// Update config.json only if the port has changed
	UpdateConfig(port, originalPort)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:" + strconv.Itoa(port),
			"http://127.0.0.1:" + strconv.Itoa(port),
		}, // only allow itslef, even not allow other localhost, like request from postman, etc. If you want to allow postman, you could delete the port
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	SetupStatic(r)
	router.SetupAPIRouter("api", r)

	address := "localhost:" + strconv.Itoa(port)
	log.Println("Starting server on", "http://"+address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
