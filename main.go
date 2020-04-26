package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/home", videoPage)
	http.HandleFunc("/",fileH)
	http.ListenAndServe(":8080", nil)
}

func fileH(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "videos/example.mp4")
}

func videoPage(w http.ResponseWriter, r *http.Request) {
	p := getPath(r)

	videos := getVideos(p)

	_, _ = fmt.Fprintf(w, "<html><head></head><body><center>")
	for _, v := range videos {
		_, _ = fmt.Fprintln(w, v)
		_, _ = fmt.Fprintf(w, "<br><embed src=\"%s\" allowfullscreen=\"true\" width=\"300\">", v)
		_, _ = fmt.Fprintln(w, "<br>")
	}
	_, _ = fmt.Fprintf(w, "</center></body></html>")
}

func getPath(r *http.Request) string {
	const d = "videos"

	q := r.URL.Query()

	if len(q) == 0 {
		return d
	}

	paths, isSet := q["p"]
	if isSet == true && len(paths) > 0 {
		return paths[0]
	}

	return d
}

func getVideos(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var videos []string

	for _, f := range files {
		videos = append(videos, path+"/"+f.Name())
	}

	return videos
}
