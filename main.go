package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
)

var client http.Client
var orgSrvURL string
var quality = 90

func main() {
	orgScheme := os.Getenv("ORIGIN_SCHEME")
	orgHost := os.Getenv("ORIGIN_HOST")
	if orgScheme == "" {
		orgScheme = "https"
	}
	orgSrvURL = orgScheme + "://" + orgHost

	if q := os.Getenv("OKOWA_QUALITY"); q != "" {
		quality, _ = strconv.Atoi(q)
	}

	http.HandleFunc("/", webpProxy)
	http.ListenAndServe(":8080", nil)
}

func webpProxy(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	if path == "/" {
		fmt.Fprintln(w, "Okowa lives!")
		return
	}

	orgURL, err := url.Parse(orgSrvURL + path)
	if err != nil {
		http.Error(w, "Invalid origin URL", http.StatusBadRequest)
		return
	}

	orgRes, err := client.Get(orgURL.String())
	if err != nil {
		http.Error(w, "Get origin failed", http.StatusBadGateway)
		return
	}
	defer orgRes.Body.Close()

	if !acceptWepb(r) {
		io.Copy(w, orgRes.Body)
		return
	}

	img, _, err := image.Decode(orgRes.Body)
	if err != nil {
		http.Error(w, "Image transformation failed", http.StatusInternalServerError)
		return
	}

	op := webp.Options{Quality: float32(quality)}

	if err := webp.Encode(w, img, &op); err != nil {
		http.Error(w, "Image transformation failed", http.StatusInternalServerError)
		return
	}
}

func acceptWepb(r *http.Request) bool {
	for _, a := range strings.Split(r.Header.Get("accept"), ",") {
		if a == "image/webp" {
			return true
		}
	}
	return false
}
