package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	imgURL := r.FormValue("url")

	if imgURL == "" {
		http.Error(w, "Provide image url as query parameter ?url=http://....jpg", http.StatusNotFound)
		return
	}

	resImg, error := http.Get(imgURL)

	if error != nil {
		http.Error(w, "Oops, we couldn't find image with this url: "+imgURL, http.StatusNotFound)
	}

	defer resImg.Body.Close()

	imageData, err := ioutil.ReadAll(resImg.Body)

	if err != nil {
		http.Error(w, "Oops, we're having trouble reading from this url: "+imgURL, http.StatusNotFound)
	}

	// You can now save it to disk or whatever...
	// try image.NewNRGBA or similar to avoid saving to disk
	// and operate in memory
	ioutil.WriteFile("./test.jpg", imageData, 0666)

	// https://godoc.org/github.com/disintegration/imaging
	originalImg, err := imaging.Open("./test.jpg")
	effect1 := imaging.Blur(originalImg, 2)
	effect2 := imaging.Invert(originalImg)

	dst := imaging.New(512, 512, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, effect1, image.Pt(0, 0))
	dst = imaging.Paste(dst, effect2, image.Pt(256, 0))

	imaging.Save(dst, "./out_example.jpg")

	outImg, _ := ioutil.ReadFile("./out_example.jpg")

	contentType := resImg.Header.Get("Content-Type")
	w.Header().Set("Content-type", contentType)
	w.Write(outImg)
}
