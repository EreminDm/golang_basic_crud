package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	//ConnectionURI address for monog db localDB connected
	//if you use another DB address please share it by env
	ConnectionURI = "192.168.99.100:27017"
	//Database name if you use another DB name please share it by env
	Database string
	//Collection name if you use another DB name please share it by env
	Collection string

	//Mongo Client to mongo db connections
	Mongo *mongo.Client
)

func init() {
	err := MongodbURIConnection()
	if err != nil {
		log.Fatal(`Couldn't connect to DB`)
	}
	Database = `personal_data`
	Collection = `information`
}

func main() {
	// Creates a gin router with default middleware
	router := gin.Default()
	//Creates routing group /person for work with PersonalData
	r := router.Group("/persons")
	{
		r.GET("/list", getAllPersonalDataList)
		r.GET("/list/:id", getPersonalDatabyID) //url example: http://localhost:port/person/list/{id}
		r.POST("/add", insertPersonalData)
		r.PUT("/update", updatePersonalData)
		r.DELETE("/remove", deletePersonalData)
	}

	// PORT environment define to 8000
	router.Run(":8000")

}
