package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	var port = 7125
	var err error
	if len(os.Args) > 1 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("specified port not a number")
		}
	}

	log.Printf("Listening on port %d\n", port)

	uploadHandler := func(w http.ResponseWriter, req *http.Request) {
		suggested := suggestedFilename(req.Header)
		if suggested == "" {
			suggested = "newfile"
		}
		filename := generateAlternateName(suggested)

		log.Printf("downloading %v as %v", suggested, filename)
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = os.WriteFile(filename, body, 0644)
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("wrote %v bytes", len(body))
	}

	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func suggestedFilename(header http.Header) string {
	filename := ""
	re := regexp.MustCompile(`filename\*?=['"]?(?:UTF-\d['"]*)?([^;\r\n"']*)['"]?;?`)
	matches := re.FindStringSubmatch(header.Get("Content-Disposition"))
	if len(matches) == 2 {
		filename = matches[1]
	}
	return filename
}

func generateAlternateName(filename string) string {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return filename
	}

	baseName := filename
	extension := filepath.Ext(filename)
	if extension != "" {
		baseName = filename[:len(filename)-len(extension)]
	}

	i := 1
	for {
		alternate := fmt.Sprintf("%s_%d%s", baseName, i, extension)
		if _, err := os.Stat(alternate); os.IsNotExist(err) {
			return alternate
		}
		i++
	}
}
