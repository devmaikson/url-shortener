package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"url-shortener/urlshort"
)

func main() {

	yamlFile := flag.String("f", "", "setting a yamlfile")

	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	//fmt.Println("Starting the server on :8080")
	//http.ListenAndServe(":8080", mapHandler)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	var yamlHandler http.HandlerFunc
	if *yamlFile != "" {
		yamlData, err := loadYamlDataFromFile(*yamlFile)
		if err != nil {
			panic(err)
		}
		yamlHandler, err = urlshort.YAMLHandler(yamlData, mapHandler)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		yamlHandler, err = urlshort.YAMLHandler([]byte(yaml), mapHandler)
		if err != nil {
			panic(err)
		}
	}

	json := `{ "/go": "https://go.dev", "/golem": "https://golem.de"}`

	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorld)
	return mux
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func loadYamlDataFromFile(filename string) ([]byte, error) {
	//fmt.Println("will load yamlfile")

	ymlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//fmt.Println(ymlData)
	return ymlData, nil
}
