package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get list of personal data info
func getAllPersonalDataList(c *gin.Context) {
	result, err := SelectAllPersonalData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

//get personal data by id
func getPersonalDatabyID(c *gin.Context) {
	idvalue := c.Param("id")
	result, err := SelectPersonalData("_id", idvalue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func insertPersonalData(c *gin.Context) {
	var p PersonalData
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//p.DocumentID = bson.NewObjectId()
	fmt.Println(p.DocumentID)
	insertResult, err := InsertPersonalData(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(insertResult)
	c.JSON(http.StatusOK, gin.H{"status": "Success"})

}

func updatePersonalData(c *gin.Context) {
	var p PersonalData
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateResult, err := UpdatePersonalDataByID(p.DocumentID, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(`Inserted result: `, updateResult)
	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("Update %v document(s) successfully", updateResult)})
}

func deletePersonalData(c *gin.Context) {
	idvalue := c.Param("id")
	result, err := DeletePersonalData(idvalue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": fmt.Sprintf("Deleted %v document(s) successfully", result)})
}
