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

// to handle line type of widgets
type line struct {
	Id string
	startWidget string
	endWidget string
}

// to handle dfs nodes
type wNode struct {
	nodeId int
	visited bool
	root bool
}

// get all widgets of specific type from specific board - actually, I have to do a part about specific board
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

// parse json response from miro API to slice with line structs
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

 // to find distinct widget ids and write them to a slice
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

// to find (for now) one root, where dfs will start
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

// to get specific widget by its id
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

// to parse text part from a widget
func readWidgetText(widget []byte) string {
	var dat map[string]interface{}
	if err := json.Unmarshal(widget, &dat); err != nil {
		panic(err)
	}
	widgetText := dat["text"].(string)
	widgetTextFinal := widgetText[3:len(widgetText)-4]
	return widgetTextFinal
}

// to create a list of nodes for dfs from widgets we've got from lines starts and ends
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

// find next nodes for current node
func findNextNodes(allLines []line, currId int) []int {
	nextNodes := make([]int, 0)
	for _, line := range allLines {
		startNode, err := strconv.Atoi(line.startWidget)
		if err != nil {
			log.Println("findNextNodes: Atoi error s.w.")
		}
		if startNode == currId {
			nextNode, err := strconv.Atoi(line.endWidget)
			if err != nil {
				log.Println("findNextNodes: Atoi error e.w.")
			}
			nextNodes = append(nextNodes, nextNode)
		}
	}
	return nextNodes
}

// pseudocode to help myself
// DFS(G, u)
// u.visited = true
// for each v ∈ G.Adj[u]
// 	if v.visited == false
// 		DFS(G,v)

// init() may be useful with multiple separated test trees on a map, for future development
// init() {
// For each u ∈ G
// 	u.visited = false
//  For each u ∈ G
//    DFS(G, u)
// }

// func to write to a text file
func writeToFile(nodesList []int) {
	f, err := os.OpenFile("/tmp/testcase", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		f, err := os.Create("/tmp/testcase")
		if err != nil && f == nil {
			panic(err)
		}
	}
	defer f.Close()
	for i, node := range nodesList {
		strNode := strconv.Itoa(node)
		widgetInfo := getWidgetById(strNode)
		widgetText := readWidgetText(widgetInfo)
		if i == 0 {
			log.Println(widgetText)
			n3, err := f.WriteString(widgetText + "\n")
			if err != nil {
				panic(err)
			}
			fmt.Printf("wrote %d bytes\n", n3)
			f.Sync()
		} else {
			log.Printf("%d. %s\n", i, widgetText)
			sI := strconv.Itoa(i)
			n3, err := f.WriteString(sI + ". " + widgetText + "\n")
			if err != nil {
				panic(err)
			}
			fmt.Printf("wrote %d bytes\n", n3)
			f.Sync()
		}
	}
}

// DFS itself with writing to a file (to struct? in future)
func dfs(nextNodes []int, theNode wNode, allNodes []wNode, allLines []line, path []int) {
	path = append(path, theNode.nodeId)
	//log.Printf("dfs starts for %d\n", theNode.nodeId)
	if len(nextNodes) == 0 {
		//log.Printf("last leaf\n")
		//log.Println(path)
		writeToFile(path)
	}
	theNode.visited = true
	for _, next := range nextNodes {
		for _, node := range allNodes {
			if node.nodeId == next && node.visited == false {
				//log.Println(node)
				newNextNodes := findNextNodes(allLines, next)
				dfs(newNextNodes, node, allNodes, allLines, path)
			}
		}
	}
	//log.Printf("dfs ends for %d\n", theNode.nodeId)
}

// this was main
func initEverything() {
	allWidgetsBody := getWidgets("line")
	allLines := parseLines(allWidgetsBody)
	allNodes := findUniqueNodes(allLines)
	//log.Println(allNodes)
	rootWidgetId := findRoot(allLines)
	rootWidgetInfo := getWidgetById(strconv.Itoa(rootWidgetId))
	readWidgetText(rootWidgetInfo)
	newNodes := createNodes(allNodes, rootWidgetId)
	//log.Println(newNodes)
	nextNodes := findNextNodes(allLines, rootWidgetId)
	var rootNode wNode
	for _, node := range newNodes {
		if node.nodeId == rootWidgetId {
			rootNode = node
			break
		}
	}
	path := make([]int, 0)
	dfs(nextNodes, rootNode, newNodes, allLines, path)
}

func main() {
	initEverything()
}


