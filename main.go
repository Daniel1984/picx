package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getImgName(url string) string {
	units := strings.Split(url, "/")
	return units[len(units)-1]
}

func getImgWidthAndHeight(img image.Image) (int, int) {
	ib := img.Bounds()
	iw := ib.Max.X
	ih := ib.Max.Y

	return iw, ih
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	imgURL := r.FormValue("url")

	if imgURL == "" {
		http.Error(w, "Provide image url as query parameter ?url=http://....jpg", http.StatusNotFound)
		return
	}

	resImg, getImgErr := http.Get(imgURL)

	if getImgErr != nil {
		http.Error(w, "Oops, we couldn't find image with this url: "+imgURL, http.StatusNotFound)
	}

	defer resImg.Body.Close()

	imageData, rBuffErr := ioutil.ReadAll(resImg.Body)

	if rBuffErr != nil {
		http.Error(w, "Oops, we're having trouble reading from this url: "+imgURL, http.StatusNotFound)
	}

	iamgeName := getImgName(imgURL)

	// You can now save it to disk or whatever...
	ioutil.WriteFile("./"+iamgeName, imageData, 0666)

	originalImg, _ := imaging.Open("./" + iamgeName)
	iw, ih := getImgWidthAndHeight(originalImg)

	effect1 := imaging.Blur(originalImg, 2)
	effect2 := imaging.Invert(originalImg)

	dst := imaging.New(iw, ih, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, effect1, image.Pt(0, 0))
	dst = imaging.Paste(dst, effect2, image.Pt(256, 0))

	imaging.Save(dst, "./out_example.jpg")

	outImg, _ := ioutil.ReadFile("./out_example.jpg")

	contentType := resImg.Header.Get("Content-Type")
	w.Header().Set("Content-type", contentType)
	w.Write(outImg)
}
