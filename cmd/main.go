package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/k8s-gadgets/notary-wrapper/version"
	"github.com/sirupsen/logrus"
)

var (
	notaryCertPath string
	notaryRootCa   string
	notaryCliPath  string
)

func main() {

	// when the wrapper starts print the version for debugging and issue logs later
	logrus.Info(getVersion())

	port := os.Getenv("NOTARY_PORT")
	if port == "" {
		port = "4445"
	}

	notaryCertPath = os.Getenv("NOTARY_CERT_PATH")
	if notaryCertPath == "" {
		notaryCertPath = "/etc/certs/notary"
	}

	notaryRootCa = os.Getenv("NOTARY_ROOT_CA")
	if notaryRootCa == "" {
		notaryRootCa = "root-ca.crt"
	}

	notaryCliPath = os.Getenv("NOTARY_CLI_PATH")
	if notaryCliPath == "" {
		notaryCliPath = "/usr/local/bin/notary"
	}

	info, err := os.Stat("/.notary")
	if err != nil {
		logrus.Info(fmt.Sprintf("err: %s", err))
	}

	logrus.Info(fmt.Sprintf("port: %s", string(port)))
	logrus.Info(fmt.Sprintf("notaryRootCa: %s", string(notaryRootCa)))
	logrus.Info(fmt.Sprintf("notaryCertPath: %s", string(notaryCertPath)))
	logrus.Info(fmt.Sprintf("notaryCliPath: %s", string(notaryCliPath)))
	logrus.Info(fmt.Sprintf("mode: %s", info.Mode()))

	// check if notary binary exists
	if _, err := os.Stat(notaryCliPath); os.IsNotExist(err) {
		logrus.Error(fmt.Sprintf("{\"code\":2,\"message\":\"NOTARY_CLI_PATH: notary cli not found\"}"))
		os.Exit(2)
	}

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// routes
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/list", ListHandler).Methods("POST")
	router.HandleFunc("/lookup", LookupHandler).Methods("POST")
	router.HandleFunc("/verify", VerifyHandler).Methods("POST")
	// router.HandleFunc("/healthz", healthzHandler).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS(notaryCertPath+"/notary-wrapper.crt", notaryCertPath+"/notary-wrapper.key"))
	log.Fatal(srv.ListenAndServe())
}

func getVersion() string {
	return fmt.Sprintf("Version: %s, Git commit: %s, Go version: %s", version.NotaryWrapperVersion, version.GitCommit, runtime.Version())
}
