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

type userProfileStruct struct {
	PK      int
	key     int
	profile string
}

// type campaignStruct struct {
// 	PK      int
// 	key     int
//   profile string
// }

var aeroDB *aero.Client

func main() {
	// CONFIGURATION
	var cfg configSpecification
	err := envconfig.Process("go-aerospike", &cfg)
	if err != nil {
		log.Print(err.Error())
	}

	// AEROSPIKE
	aeroDB, err := aero.NewClient(cfg.AerospikeHost, cfg.AerospikePort)
	if err != nil {
		log.Print("AEROSPIKE | CONNECTION FAILED: ", err)
	}

	// HTTP SERVER SETUP
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Hello internet",
		})
	})

	r.GET("/profile/:id", func(c *gin.Context) {
		log.Print(c.Param("id"))
		key, _ := aero.NewKey("mobucks", "userProfiles", c.Param("id"))

		if exists, _ := aeroDB.Exists(nil, key); exists {
			userProfileData := &userProfileStruct{}
			aeroDB.GetObject(nil, key, userProfileData)

			c.JSON(200, gin.H{
				"status": "ok",
				"data":   userProfileData,
			})
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   nil,
		})
	})

	// HTTP SERVER START
	r.Run(fmt.Sprintf(":%d", cfg.ApplicationPort)) // listen and serve on 0.0.0.0:8080
}
