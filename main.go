package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

var client http.Client
var orgSrvURL string

func main() {
	orgScheme := os.Getenv("ORIGIN_SCHEME")
	orgHost := os.Getenv("ORIGIN_HOST")
	if orgScheme == "" {
		orgScheme = "https"
	}
	orgSrvURL = orgScheme + "://" + orgHost

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

	img, err := imaging.Decode(orgRes.Body, imaging.AutoOrientation(true))
	if err != nil {
		http.Error(w, "Image transformation failed", http.StatusInternalServerError)
		return
	}

	if err := webp.Encode(w, img, nil); err != nil {
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
