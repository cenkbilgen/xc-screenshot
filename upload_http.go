package main

import (
  "net/http"
  "io/ioutil"
  "log"
  "regexp"
  "strconv"
)

func main() {
	const port = 7000
 	uploadHandler := func(w http.ResponseWriter, req *http.Request) {
 		suggested := suggestedFilename(req.Header)
 		if suggested == "" {
 			suggested = "newfile"
 		}
 		filename := generateFilename(suggested)
 		
 		log.Printf("downloading %v as %v", suggested, filename)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = ioutil.WriteFile(filename, body, 0644)
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}	
		log.Printf("wrote %v bytes", len(body))
	}

	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
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

func generateFilename(suggested string) string {
	files, err := ioutil.ReadDir(`.`)
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`^` + suggested + `(-(\d+)){0,1}$`)
	nextn := 0
	for _, file := range files {
		matches := re.FindStringSubmatch(file.Name())
		mnum := len(matches)
		if mnum == 0 {
			continue
		} else if mnum == 1 {
			if nextn == 0 {
				nextn = 1
			}
		} else if mnum == 3 {
			filen, _ := strconv.Atoi(matches[2])
			if filen >= nextn {
				nextn += 1
			}
		} else if mnum == 2 || mnum > 3 {
			log.Fatal("internal error. check regexp")
		}
	}
	if nextn == 0 {
		return suggested
	} else {
		return suggested + "-" + strconv.Itoa(nextn)
	}
}