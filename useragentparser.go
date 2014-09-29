package main

import (
	"encoding/json"
	"net/http"

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

var parser *uaparser.Parser

func main() {
	isThrottled := true
	th := throttled.RateLimit(throttled.PerHour(3600), &throttled.VaryBy{RemoteAddr: true}, store.NewMemStore(1000))
	regexFile := "regexes.yaml"
	parser = uaparser.New(regexFile)
	r := mux.NewRouter()
	var thandler http.Handler
	thandler = http.HandlerFunc(testHandler())
	if isThrottled {
		thandler = th.Throttle(thandler)
	}
	r.Handle("/api", thandler)
	http.Handle("/", r)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":3000", nil)

}
