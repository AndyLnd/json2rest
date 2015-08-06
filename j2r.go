package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"net/http"
	"strings"
	"errors"
	"flag"
)

var jsonData interface{}

func loadJson(fileName string) interface{} {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var geoData interface{}
	e = json.Unmarshal(file, &geoData)
	if e != nil {
		fmt.Printf("JSON error: %v\n", e)
		os.Exit(2)
	}

	return geoData
}

func saveJson(data interface{}, fileName string) {
	jsonData, e := json.Marshal(data)
	if e != nil {
		fmt.Printf("Error encoding JSON: %v\n", e)
		os.Exit(3)
	}

	e = ioutil.WriteFile(fileName, []byte(jsonData), 0644)
	if e != nil {
		fmt.Printf("Error writing file: %v\n", e)
		os.Exit(4)
	}
}

func main() {
	optFileName := flag.String("i", "", "The input json file")
	optPort := flag.Int("p", 1337, "The listening port.")
	flag.Parse()
	fileName := *optFileName
	port := int64(*optPort)
	if fileName == "" {
		flag.PrintDefaults()
		os.Exit(6)
	}
	portString := strconv.FormatInt(port, 10)
	jsonData = loadJson(fileName)
	http.HandleFunc("/", handleReq)
	fmt.Printf("Starting server on port %s ...\n",portString)
	err := http.ListenAndServe(":"+portString, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func handleReq(w http.ResponseWriter, req *http.Request) {
	path := strings.Split(req.RequestURI, "/")
	method := req.Method
	fmt.Println(method)
	result, e := getData(path)
	if e != nil {
		w.Write([]byte(e.Error()))
	} else {
		jsonString, e := json.Marshal(result)
		if e != nil {
			fmt.Printf("Error encoding JSON: %v\n", e)
			os.Exit(5)
		}
		w.Write([]byte(jsonString))
	}
}

func getData(path []string) (interface{}, error) {
	currRoot := jsonData
	index := 0
	l := len(path)
	for index < l {
		if path[index]=="" {
			index++
			continue
		}
		switch  v := currRoot.(type){
			case map[string]interface{}:
			currRoot = currRoot.(map[string]interface{})[path[index]]
			index++
			break
			case []interface{}:
			i, e := strconv.ParseInt(path[index], 10, 64)
			if e != nil {
				return nil, errors.New("index required")
			}
			arrayCurr := currRoot.([]interface{})
			if (len(arrayCurr)<=int(i)) {
				return nil, errors.New("out of array bounds")
			}
			currRoot = arrayCurr[i]
			index++
			break
			default:
			return nil, errors.New(fmt.Sprintf("unhandled type: %T", v))
			index=l
		}
	}
	if currRoot == nil {
		return nil, errors.New("path not found")
	}
	return currRoot, nil
}
