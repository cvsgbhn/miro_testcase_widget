package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"os"
	"github.com/joho/godotenv"
)

type line struct {
	Id string
	startWidget string
	endWidget string
}

func getWidgets(typeName string) []byte {
	url := fmt.Sprintf("https://api.miro.com/v1/boards/o9J_lRePMUc=/widgets/?widgetType=%s", typeName)

	// Create a Bearer string by appending string access token
	if err := godotenv.Load(".env"); err != nil {
        log.Print("No .env file found")
    }
	token, exist := os.LookupEnv("BEARER")
	log.Println("token " + token)
	if exist == false {
		log.Println("missing token")
		return nil
	}
    var bearer = "Bearer " + token

    // Create a new request using http
    req, err := http.NewRequest("GET", url, nil)

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERROR] -", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error while reading the response bytes:", err)
    }

	return []byte(body)
}

func parseLines(linesBytes []byte) []line {
	var dat map[string]interface{}
 	if err := json.Unmarshal(linesBytes, &dat); err != nil {
 		panic(err)
 	}

	allLines := []line{}

	for _,item:=range dat["data"].([]interface{}) {
		theLine := line{
			Id: item.(map[string]interface{})["id"].(string),
			startWidget: (item.(map[string]interface{})["startWidget"]).(map[string]interface {})["id"].(string),
			endWidget: item.(map[string]interface{})["endWidget"].(map[string]interface {})["id"].(string)}
		allLines = append(allLines, theLine)
	}
 	return allLines
 }

 func findUniqueNodes(lines []line) []int {
	var nodes []int
	for _, l := range lines {
		oneNode, err := strconv.Atoi(l.startWidget)
		if err != nil {
			log.Println("Atoi error")
		}
		nodes = append(nodes, oneNode)
		oneNode, err = strconv.Atoi(l.endWidget)
		if err != nil {
			log.Println("Atoi error")
		}
		nodes = append(nodes, oneNode)
	}
	finArr := make([]int, 0, len(nodes))
	mappedArr := make(map[int]bool)
	for _, val := range nodes {
		if _, ok := mappedArr[val]; !ok {
			mappedArr[val] = true
			finArr = append(finArr, val)
		}
	}
	return finArr	
}

func findRoot(lines []line) int {
	var startNodes []int
	var endNodes []int
	for _, l := range lines {
		oneNode, err := strconv.Atoi(l.startWidget)
		if err != nil {
			log.Println("Atoi error")
		}
		startNodes = append(startNodes, oneNode)
		oneNode, err = strconv.Atoi(l.endWidget)
		if err != nil {
			log.Println("Atoi error")
		}
		endNodes = append(endNodes, oneNode)
	}
	for _, sd := range startNodes {
		for _, ed := range endNodes {
			if sd != ed {
				return sd
			}
		}
	}
	return 0
}

func getWidgetById(wId string) []byte {
	url := fmt.Sprintf("https://api.miro.com/v1/boards/o9J_lRePMUc=/widgets/%s", wId)

	// Create a Bearer string by appending string access token
    var bearer = "Bearer " + "EJrhSsVARKZUzy1ipm1-P9L2d2c"

    // Create a new request using http
    req, err := http.NewRequest("GET", url, nil)

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)

    // Send req using http Client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error on response.\n[ERROR] -", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error while reading the response bytes:", err)
    }

	return []byte(body)
}

func readWidgetText(widget []byte) string {
	var dat map[string]interface{}
	if err := json.Unmarshal(widget, &dat); err != nil {
		panic(err)
	}
	widgetText := dat["text"].(string)
	return widgetText
}

///////////////////////////////////////////////////////////////
//////////////////////// DFS //////////////////////////////////
 
// 1. Unique nodes to structs { nodeID int, visited bool, root bool} (array of structs)
// a. Create struct wNode{}
type wNode struct {
	nodeId int
	visited bool
	root bool
}
// b. Create func (uniqueNodes []int, rootId int) []wNode
func createNodes(uniqueNodes []int, rootId int) []wNode {
	nodes := []wNode{}
	for _, n := range uniqueNodes {
		theNode := wNode{}
		if n == rootId {
			theNode = wNode{
				nodeId: n,
				visited: false,
				root: true}
		} else {
			theNode = wNode{
				nodeId: n,
				visited: false,
				root: false}
		}
		nodes = append(nodes, theNode)
	}
	return nodes
}

// 2. Simple dfs
func dfs(lines []line, nodes []wNode) {
}
// 3. Upgrade dfs to write paths



///////////////////////////////////////////////////////////////

func main() {
	allWidgetsBody := getWidgets("line")
	allLines := parseLines(allWidgetsBody)
	allNodes := findUniqueNodes(allLines)
	log.Println(allNodes)
	rootWidgetId := findRoot(allLines)
	rootWidgetInfo := getWidgetById(strconv.Itoa(rootWidgetId))
	readWidgetText(rootWidgetInfo)
	newNodes := createNodes(allNodes, rootWidgetId)
	log.Println(newNodes)
}


