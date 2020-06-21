package main

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/k8s-gadgets/notary-wrapper/version"
)

type Info struct {
	NotaryWrapperVersion string `json:"NotaryWrapperVersion"`
	GitCommit            string `json:"GitCommit"`
	RuntimeVersion       string `json:"runtimeVersion"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	var info Info

	info.NotaryWrapperVersion = version.NotaryWrapperVersion
	info.GitCommit = version.GitCommit
	info.RuntimeVersion = runtime.Version()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(info)
}
