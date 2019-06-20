package main

import (
	"fmt"
	"log"

	aero "github.com/aerospike/aerospike-client-go"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type configSpecification struct {
	ApplicationPort int    `envconfig:"APP_PORT" default:"8089"`
	AerospikeHost   string `envconfig:"AEROSPIKE_HOST" default:"127.0.0.1"`
	AerospikePort   int    `envconfig:"AEROSPIKE_PORT" default:"3000"`
}

type userProfile struct {
	name string
}

func main() {
	// CONFIGURATION
	var cfg configSpecification
	err := envconfig.Process("go-aerospike", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// AEROSPIKE
	client, err := aero.NewClient(cfg.AerospikeHost, cfg.AerospikePort)
	if err != nil {
		log.Fatal("Error connecting to Aerospike Database Server", err)
	}
	defer client.Close()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Hello internet",
		})
	})

	r.GET("/profile/:id", func(c *gin.Context) {
		id := c.Param("id")

		key, err := aero.NewKey("mobucks", "userProfiles", id)
		if err != nil {
			log.Print(err)
		}

		testObj := userProfile{name: "Stratos"}
		client.PutObject(nil, key, &testObj)

		var userProfileData userProfile
		client.GetObject(nil, key, &userProfileData)

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   userProfileData,
		})
	})

	r.Run(fmt.Sprintf(":%d", cfg.ApplicationPort)) // listen and serve on 0.0.0.0:8080
}
