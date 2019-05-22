package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func downloadConfigFile(file string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if file == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("file required"))
			return
		}

		if _, err := os.Stat(file); err != nil {
			http.ServeFile(w, r, file)
			return
		}

		filename := filepath.Base(file)

		w.Header().Set("Content-Disposition", "attachment; filename="+url.PathEscape(filename))
		w.Header().Set("Content-Description", "File Transfer")
		w.Header().Set("Content-Type", "application/x-apple-aspen-config; chatset=utf-8")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		w.Header().Set("Cache-Control", "must-revalidate")
		w.Header().Set("Pragma", "public")

		log.Printf("download file: %s", filename)

		http.ServeFile(w, r, file)
	}

}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	result := make(url.Values)

	//data, err := ioutil.ReadFile("receive.xml")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	dataStr := string(data)

	start := strings.Index(dataStr, "<?xml")
	end := strings.LastIndex(dataStr, "</plist>")

	dataStr = dataStr[start:end]

	log.Printf("receive data: %s", dataStr)
	decoder := xml.NewDecoder(bytes.NewBufferString(dataStr))

	res := make([]string, 0)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.CharData:
			if strings.TrimSpace(string(token)) != "" {
				res = append(res, string(token))
			}
		}
	}

	for i, info := range res {
		switch info {
		case "IMEI":
			result.Set("IMEI", res[i+1])
		case "PRODUCT":
			result.Set("PRODUCT", res[i+1])
		case "UDID":
			result.Set("UDID", res[i+1])
		case "VERSION":
			result.Set("VERSION", res[i+1])
		}
	}

	//out, err := json.Marshal(result)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//w.WriteHeader(http.StatusOK)
	//w.Write(out)

	log.Printf("out data: %+v", result)

	http.Redirect(w, r, "/udid?"+result.Encode(), http.StatusMovedPermanently)

}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Printf("query string: %+v", r.Form)
		w.Write([]byte("UDID: " + r.Form.Get("UDID")))
	})

	http.HandleFunc("/receive", receiveHandler)

	http.HandleFunc("/download", downloadConfigFile("udid.mobileconfig"))

	http.ListenAndServe(":8080", nil)

}
