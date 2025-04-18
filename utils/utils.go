package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	log.Print("\n", msg, string(pretty))
	fmt.Print("\n", msg, string(pretty))
}

// ToJSONString beautiful display json
func ToJSONString(v interface{}) string {
	pretty, _ := json.MarshalIndent(v, "", "  ")

	return string(pretty)
}

// DebugRequest ...
func DebugRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("[WARN] Error getting request's dump: %s\n", err)
	}

	log.Printf("[DEBUG] %s\n", string(requestDump))
}

// DebugResponse ...
func DebugResponse(res *http.Response) {
	requestDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Printf("[WARN] Error getting response's dump: %s\n", err)
	}

	log.Printf("[DEBUG] %s\n", string(requestDump))
}

// ConvertMapString from interface to string
func ConvertMapString(o map[string]interface{}) map[string]string {
	converted := make(map[string]string)
	for k, v := range o {
		converted[k] = fmt.Sprint(v)
	}
	return converted
}
