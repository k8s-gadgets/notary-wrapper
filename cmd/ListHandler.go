package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var requestGun RequestGun
	json.Unmarshal(reqBody, &requestGun)

	logrus.Info(fmt.Sprintf("request: list, content: %s", reqBody))
	digests, err := execNotaryCLI("list", requestGun)

	w.Header().Set("Content-Type", "application/json")

	if err == nil {

		count := 0

		if requestGun.Tag == "" {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(digests)
			logrus.Info(fmt.Sprintf("response: %s", digests))
			count = 1
		} else {

			for _, element := range digests {
				if element.Name == requestGun.Tag {
					w.WriteHeader(200)
					json.NewEncoder(w).Encode(element)
					_, err := json.Marshal((element))
					if err != nil {
						logrus.Error(fmt.Sprintf("logtrace: %s", err.Error()))
					}
					logrus.Info(fmt.Sprintf("response: %s", element))
					count = 1
				}
			}
		}
		if count == 0 {
			w.WriteHeader(404)
			w.Write([]byte("{\"code\":404,\"message\":\"SHA not found\"}"))
			json.NewEncoder(w)
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
