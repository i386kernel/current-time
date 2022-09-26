package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var conversionMAP = map[string]string{
	"ASR": "-3h",    // North America Atlantic standard time
	"EST": "-5h",    // North american eastern standard time
	"BST": "+1h",    // British summer time
	"IST": "+5h30m", // Indian standard time
	"HKT": "+8h",    // Hong Kong Time
	"ART": "-3h",    // Argentina Time
}

func getCurrentTimeWithTimeDifference(timeDifference string) (string, error) {
	now := time.Now().UTC()
	diff, err := time.ParseDuration(timeDifference)
	if err != nil {
		return " ", err
	}
	now = now.Add(diff)
	return now.Format("15:04:05"), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	timeZone := r.URL.Query().Get("tz")
	if timeZone == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error 400: 'tz' query parameter not found")
		return
	}
	w.WriteHeader(http.StatusOK)
	timediff, ok := conversionMAP[timeZone]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `Error 404: the tz value "%s" does not match available timezones %s`, timeZone, conversionMAP)
		return
	}
	tf, err := getCurrentTimeWithTimeDifference(timediff)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "Time in timezone %v is %s", timeZone, tf)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Error 404: The requested url does not exist.")
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

// Middleware take http handler, validates something and then re-directs from there.
func loggingMiddleware(handler handlerFunc) handlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s = %s - %s", time.Now().Format("2020-06-28 13:00"), r.Method, r.URL.String())
		handler(w, r)
	}
	return fn
}

func main() {
	http.HandleFunc("/convert", loggingMiddleware(handler))
	http.HandleFunc("/", loggingMiddleware(notFoundHandler))
	log.Printf("%s - Starting server on port: 8080", time.Now().Format("12-12-1983 20:18"))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
