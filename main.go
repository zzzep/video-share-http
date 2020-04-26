package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type Video struct {
	fullPath string
	name     string
	fileType string
}

var videoTypes = map[string]string{
	"default": "video/webm",
	".mp4": "video/mp4",
	".mkv": "video/webm",
	".m4v": "video/webm",
	".ogg": "video/ogg",
}

func main() {
	http.HandleFunc("/home", videoPage)
	http.HandleFunc("/", fileH)
	http.ListenAndServe(":8080", nil)
}

func fileH(w http.ResponseWriter, r *http.Request) {
	f := r.URL.Path
	f = string(bytes.TrimLeft([]byte(f), "/"))
	f = string(bytes.TrimLeft([]byte(f), "\\"))
	http.ServeFile(w, r, f)
}

func videoPage(w http.ResponseWriter, r *http.Request) {
	p := getPath(r)

	videos := getVideos(p)

	_, _ = fmt.Fprintf(w, "<!DOCTYPE html><html>")
	_, _ = fmt.Fprintf(w, "<head></head>")
	_, _ = fmt.Fprintf(w, "<body><center>")
	for _, v := range videos {
		_, _ = fmt.Fprintln(w, v.name)
		_, _ = fmt.Fprintln(w, "<br>")
		//_, _ = fmt.Fprintf(w, "<embed src=\"%s\" width=\"300\" autostart=\"0\">", v.fullPath)
		_, _ = fmt.Fprintf(w, "<video width=\"300px\" controls><source src=\"%s\" type=\"%s\"></video>", v.fullPath, v.fileType)
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

func getVideos(path string) []Video {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var videos []Video

	for _, f := range files {
		var v Video
		v.name = f.Name()
		v.fullPath = path + "/" + v.name
		v.fileType = getVideoType(v)
		videos = append(videos, v)
	}

	return videos
}

func getVideoType(v Video) string {
	ext := filepath.Ext(v.name)

	t, exists := videoTypes[ext]
	if exists == true {
		return t
	}

	return videoTypes["default"]
}
