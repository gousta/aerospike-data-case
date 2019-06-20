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
	key     string
	profile string
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
		log.Fatal("AEROSPIKE | CONNECTION FAILED", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Hello internet",
		})
	})

	r.GET("/profile/:id", func(c *gin.Context) {
		key, _ := aero.NewKey("mobucks", "userProfiles", c.Param("id"))

		testObjToInsert := &userProfileStruct{
			key:     "test",
			profile: "JSON[{\"interestIds\":[1, 2], \"groupId\":1}, {\"interestIds\":[3], \"groupId\":2}]",
		}

		err := client.PutObject(nil, key, testObjToInsert)
		if err != nil {
			log.Fatal("AEROSPIKE | WRITING OBJECT FAILED: ", err)
		}

		userProfileData := &userProfileStruct{}
		err = client.GetObject(nil, key, userProfileData)
		if err != nil {
			log.Fatal("AEROSPIKE | READING OBJECT FAILED: ", err)
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   userProfileData,
		})
	})

	r.GET("/test", func(c *gin.Context) {
		type TestStruct struct {
			mykey int `as:"my_key"` // alias the field to a
		}

		key, _ := aero.NewKey("mobucks", "userProfiles", "test")

		err := client.PutObject(nil, key, TestStruct{mykey: 15})
		if err != nil {
			log.Fatal("AEROSPIKE | WRITING OBJECT FAILED: ", err)
		}

		rObj := &TestStruct{}
		err = client.GetObject(nil, key, rObj)
		if err != nil {
			log.Fatal("AEROSPIKE | READING OBJECT FAILED: ", err)
		}

		c.JSON(200, gin.H{
			"status": "ok",
			"data":   rObj,
		})
	})

	r.Run(fmt.Sprintf(":%d", cfg.ApplicationPort)) // listen and serve on 0.0.0.0:8080
}
