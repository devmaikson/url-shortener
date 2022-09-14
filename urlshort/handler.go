package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path,omitempty"`
	URL  string `yaml:"url,omitempty"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []pathUrl

	err := yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		fmt.Println("error in Unmarshalling YAML")
		return nil, err
	}
	// convert pathURLs to map for mapHandler
	data := make(map[string]string)
	for _, pu := range pathURLs {
		// fmt.Println("Data:", pu.Path, pu.URL)
		data[pu.Path] = pu.URL
	}

	return MapHandler(data, fallback), nil
}

func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	data := make(map[string]string)

	err := json.Unmarshal(yml, &data)
	if err != nil {
		fmt.Println("error in Unmarshalling JSON")
		return nil, err
	}
	//fmt.Println("Data: ", data)

	return MapHandler(data, fallback), nil
}
