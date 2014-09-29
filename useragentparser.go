package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/gorilla/mux"
	"github.com/tobie/ua-parser/go/uaparser" // You could change this to a github repo as well
)

func testHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Expose-Headers", "X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		defer r.Body.Close()
		switch r.Method {
		case "GET":
			client := parser.Parse(r.URL.Query().Get("ua"))
			WriteJSON(w, client)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getIpAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIp := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIp == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIp
}

var parser *uaparser.Parser

func main() {
	isThrottled := true
	th := throttled.RateLimit(throttled.PerHour(3600), &throttled.VaryBy{Custom: getIpAddress}, store.NewMemStore(1000))
	regexFile := "regexes.yaml"
	parser = uaparser.New(regexFile)
	r := mux.NewRouter()
	var thandler http.Handler
	thandler = http.HandlerFunc(testHandler())
	if isThrottled {
		thandler = th.Throttle(thandler)
	}
	r.Handle("/api", thandler)
	r.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)

}
