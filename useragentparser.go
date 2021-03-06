package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/gorilla/mux"
	"github.com/tobie/ua-parser/go/uaparser" // You could change this to a github repo as well
)

type UserAgentResult struct {
	UAString  string `json:"ua_string"`
	UAVersion string `json:"ua_version"`
	UAMajor   string `json:"ua_major"`
	UAMinor   string `json:"ua_minor"`
	UAPatch   string `json:"ua_patch"`

	OSString     string `json:"os_string"`
	OSVersion    string `json:"os_version"`
	OSMajor      string `json:"os_major"`
	OSMinor      string `json:"os_minor"`
	OSPatch      string `json:"os_patch"`
	OSPatchMinor string `json:"os_patch_minor"`

	DeviceFamily string `json:"device_family"`
}

func apiHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Expose-Headers", "X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		defer r.Body.Close()
		switch r.Method {
		case "GET":
			ua := r.URL.Query().Get("ua")
			ip := getIpAddress(r)
			log.Printf("API request by %v with ua=%v", ip, ua)
			client := parser.Parse(ua)
			res := &UserAgentResult{
				UAString:     client.UserAgent.ToString(),
				UAVersion:    client.UserAgent.ToVersionString(),
				UAMajor:      client.UserAgent.Major,
				UAMinor:      client.UserAgent.Minor,
				UAPatch:      client.UserAgent.Patch,
				OSString:     client.Os.ToString(),
				OSVersion:    client.Os.ToVersionString(),
				OSMajor:      client.Os.Major,
				OSMinor:      client.Os.Minor,
				OSPatch:      client.Os.Patch,
				OSPatchMinor: client.Os.PatchMinor,
				DeviceFamily: client.Device.Family,
			}
			WriteJSON(w, res)
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
	isThrottled, _ := strconv.ParseBool(os.Getenv("USERAGENTPARSER_THROTTLED"))
	serveDocs, _ := strconv.ParseBool(os.Getenv("USERAGENTPARSER_DOCS"))
	log.Printf("isThrottled=%v, serveDocs=%v", isThrottled, serveDocs)
	th := throttled.RateLimit(throttled.PerHour(3600), &throttled.VaryBy{Custom: getIpAddress}, store.NewMemStore(1000))
	regexFile := "regexes.yaml"
	parser = uaparser.New(regexFile)
	r := mux.NewRouter()
	var thandler http.Handler
	thandler = http.HandlerFunc(apiHandler())
	if isThrottled {
		thandler = th.Throttle(thandler)
	}
	r.Handle("/api", thandler)
	r.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)

}
