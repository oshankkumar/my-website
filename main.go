package main

import (
	"encoding/json"
	"flag"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

var (
	port = flag.String("port", os.Getenv("PORT"), "application server port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/health", health)
	http.Handle("/", http.FileServer(http.Dir("web")))
	mux := AddMiddleware(http.DefaultServeMux, Logging)

	logrus.WithField("addr", ":"+*port).Info("starting server")
	logrus.Fatal(http.ListenAndServe(":"+*port, mux))
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.WithError(err).Error("error in dialing google public dns")
		return ""
	}
	defer conn.Close()
	addr := conn.LocalAddr().String()
	return strings.Split(addr, ":")[0]
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&struct {
		Health  string `json:"health"`
		LocalIP string `json:"IP"`
	}{"up", getLocalIP()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Logging(next http.Handler) http.Handler {
	return handlers.LoggingHandler(logrus.NewEntry(logrus.StandardLogger()).Writer(), next)
}

func AddMiddleware(base http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		base = middlewares[i](base)
	}
	return base
}
