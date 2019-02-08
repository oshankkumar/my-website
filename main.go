package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&struct {
		Health  string `json:"health"`
		LocalIP string `json:"IP"`
	}{"up", getLocalIP()})
}

func main() {
	ipAddr := getLocalIP()
	logrus.WithField("ip", ipAddr).Info("starting server")
	http.HandleFunc("/health", health)
	http.Handle("/", http.FileServer(http.Dir("web")))
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.WithError(err).Error("error in dialing google public dns")
		return ""
	}
	defer conn.Close()
	return conn.LocalAddr().String()
}
