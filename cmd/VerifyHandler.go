package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var verifySHA VerifySHA
	json.Unmarshal(reqBody, &verifySHA)
	requestGun := RequestGun{NotaryServer: verifySHA.NotaryServer, Gun: verifySHA.Gun, Tag: ""}

	logrus.Info(fmt.Sprintf("request: verify, content: %s", reqBody))
	digests, err := execNotaryCLI("list", requestGun)

	w.Header().Set("Content-Type", "application/json")

	if err == nil {

		var matchImages []NotaryList
		count := 0
		for _, element := range digests {
			if element.Digest == verifySHA.SHA {
				matchImages = append(matchImages, element)
				count = 1
			}
		}

		if count == 1 {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(matchImages)
			logrus.Info(fmt.Sprintf("response: %s", matchImages))
		} else if count == 0 {
			w.WriteHeader(404)
			w.Write([]byte("{\"code\":404,\"message\":\"SHA not found\"}"))
			json.NewEncoder(w)
			logrus.Info(fmt.Sprintf("response: %s", matchImages))
		}
	} else {
		var customError CustomError
		json.Unmarshal([]byte(err.Error()), &customError)
		code, newErr := strconv.Atoi(customError.Code)

		if newErr != nil {
			w.WriteHeader(500)
			w.Write([]byte("{\"code\":500,\"message\":\"internal server error\"}"))
			json.NewEncoder(w)
			return
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		json.NewEncoder(w)
	}
}
