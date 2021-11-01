package main

// Version 1.0 of a middleware Rest call for Prisma Cloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sshintaku/web_requests"

	"github.com/sshintaku/cloud_types"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

var token, computeBaseUrl string

func main() {
	router := gin.Default()
	username := os.Getenv("APIKEY")
	password := os.Getenv("PASSWORD")
	authResponse, _ := web_requests.GetJWTToken("https://api2.prismacloud.io/login", username, password)
	token = authResponse.Token
	computeUrl, baseUrlError := web_requests.GetComputeBaseUrl(token)
	computeBaseUrl = computeUrl
	if baseUrlError != nil {
		log.Fatal(baseUrlError)
	}
	router.GET("/clouddiscovery", getCloudDiscovery)
	router.Static("/assets", "./assets")
	router.GET("clouddiscovery/bytype", getCloudDiscoveryByType)
	router.GET("/", loadDefaultPage)

	router.Run("localhost:8080")
}

func getCloudDiscoveryByType(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	cloudType := c.Query("type")
	deployed, deployedError := strconv.ParseBool(c.Query("deployed"))
	if deployedError != nil {
		log.Fatal("Error converting query parameter to a bool value.")
	}
	url := computeBaseUrl + "/api/v1/cloud/discovery?project=Central+Console"
	result, resultError := web_requests.GetMethod(url, token)
	if resultError != nil {
		log.Fatal(resultError)
	}
	var cloudResults []cloud_types.DiscoveryResult
	var returnResult []cloud_types.CloudTypeResult
	json.Unmarshal(result, &cloudResults)
	for _, ct := range cloudResults {
		if ct.ServiceType == cloudType {
			var cr cloud_types.CloudTypeResult
			cr.Region = ct.Region
			for _, status := range ct.Entities {
				if status.Defended == deployed {
					cr.Status = "<i><b>Resource group: </i></b>" + status.ResourceGroup + "<br><i><b>Resource Name: </i></b>" + status.Name + "\n"
					returnResult = append(returnResult, cr)
				}
			}
		}
	}
	c.IndentedJSON(http.StatusOK, returnResult)
}

func loadDefaultPage(c *gin.Context) {
	body, err := ioutil.ReadFile("assets/index2.html")
	if err != nil {
		log.Fatal(err)
	}
	markdown := blackfriday.MarkdownCommon(body)
	c.Data(http.StatusOK, "text/html; charset=utf-8", markdown)
}

func getCloudDiscovery(c *gin.Context) {
	url := computeBaseUrl + "/api/v1/cloud/discovery?project=Central+Console"
	result, resultError := web_requests.GetMethod(url, token)

	fmt.Println(string(result))
	if resultError != nil {
		log.Fatal(resultError)
	}
	var cloudResults []cloud_types.DiscoveryResult
	json.Unmarshal(result, &cloudResults)
	c.IndentedJSON(http.StatusOK, cloudResults)
}
