package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sshintaku/web_requests"

	"github.com/sshintaku/cloud_types"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

var token, computeBaseUrl string

func main() {
	router := gin.Default()
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	configFile, _ := ioutil.ReadAll(jsonFile)
	var Auth cloud_types.Authentication
	json.Unmarshal(configFile, &Auth)
	authResponse, _ := web_requests.GetJWTToken("https://api2.prismacloud.io/login", Auth.Username, Auth.Password)
	token = authResponse.Token
	computeUrl, baseUrlError := web_requests.GetComputeBaseUrl(token)
	computeBaseUrl = computeUrl
	if baseUrlError != nil {
		log.Fatal(baseUrlError)
	}
	router.GET("/clouddiscovery", getCloudDiscovery)
	router.Static("/assets", "./assets")
	router.GET("clouddiscovery/:type", getCloudDiscoveryByType)
	router.GET("/", loadDefaultPage)

	router.Run("localhost:8080")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getCloudDiscoveryByType(c *gin.Context) {
	cloudType := c.Param("type")
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
				if status.Defended == false {
					cr.Status = "Resource group: " + status.ResourceGroup + " Resource Name: " + status.Name + "\n"
				}
			}
			if cr.Status != "" {
				returnResult = append(returnResult, cr)
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
