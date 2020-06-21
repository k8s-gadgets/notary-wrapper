package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func execNotaryCLI(notaryCommand string, requestGUN RequestGun) ([]NotaryList, error) {

	var certDir string
	var certPath string
	var rootCaExists bool
	var cmd *exec.Cmd
	var server string
	var digests []NotaryList

	trustFolder := "/.notary"
	// trustFolder := "/home/daniel/.notary"

	notaryServerSplit := strings.Split(requestGUN.NotaryServer, ".")
	notaryServerName := notaryServerSplit[0]
	server = strings.Join([]string{"https://", requestGUN.NotaryServer}, "")

	certDir = strings.Join([]string{notaryCertPath, "/", notaryServerName}, "")
	certPath = strings.Join([]string{certDir, "/", notaryRootCa}, "")
	rootCaExists = true

	// check if cert exists
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		rootCaExists = false
	}

	if rootCaExists {
		if notaryCommand == "list" {
			cmd = exec.Command(notaryCliPath, "-s", server, "-d", trustFolder, "--tlscacert", certPath, notaryCommand, requestGUN.Gun)
		} else if notaryCommand == "lookup" {
			cmd = exec.Command(notaryCliPath, "-s", server, "-d", trustFolder, "--tlscacert", certPath, notaryCommand, requestGUN.Gun, requestGUN.Tag)
		}
	} else {
		if notaryCommand == "list" {
			cmd = exec.Command(notaryCliPath, "-s", server, "-d", trustFolder, notaryCommand, requestGUN.Gun)
		} else if notaryCommand == "lookup" {
			cmd = exec.Command(notaryCliPath, "-s", server, "-d", trustFolder, notaryCommand, requestGUN.Gun, requestGUN.Tag)
		}
	}

	stdOutStdErr, err := cmd.CombinedOutput()
	stdOutStdErrString := string(stdOutStdErr)
	stdOutStdErrString = strings.TrimSuffix(stdOutStdErrString, "\n")
	stdOutStdErrString = strings.TrimPrefix(stdOutStdErrString, "\n")

	if err != nil {

		if strings.Contains(stdOutStdErrString, "does not have trust data for") {

			//("404 - fatal: <notary-server> does not have trust data for <image>")

			info := strings.SplitN(stdOutStdErrString, " ", 2)[1]
			customError := CustomError{"404", info}
			jsonCustomError, err := json.Marshal(customError)
			err = fmt.Errorf("%s", jsonCustomError)
			logrus.Error(fmt.Sprintf("%s", err))
			return nil, err

		} else {

			//("500 - internal server error")
			customError := CustomError{"500", "internal server error"}
			jsonCustomError, err := json.Marshal(customError)
			err = fmt.Errorf("%s", jsonCustomError)
			logrus.Error(fmt.Sprintf("%s", err))
			logrus.Error(fmt.Sprintf("logtrace: %s", stdOutStdErrString))
			return nil, err
		}
	}

	if strings.Contains(stdOutStdErrString, "x509: certificate signed by unknown authority") {

		//("404 - fatal: x509: certificate signed by unknown authority")

		// info := strings.SplitN(stdOutStdErrString, " ", 2)[1]
		info := strings.Join([]string{"x509: certificate signed by unknown authority", server}, " ")
		customError := CustomError{"404", info}
		jsonCustomError, err := json.Marshal(customError)
		err = fmt.Errorf("%s", jsonCustomError)
		logrus.Error(fmt.Sprintf("%s", err))
		return nil, err

	} else if strings.Contains(stdOutStdErrString, "could not reach") {

		//("404 - fatal: could not reach")

		// info := strings.SplitN(stdOutStdErrString, " ", 2)[1]
		info := strings.Join([]string{"could not reach", server}, " ")
		customError := CustomError{"404", info}
		jsonCustomError, err := json.Marshal(customError)
		err = fmt.Errorf("%s", jsonCustomError)
		logrus.Error(fmt.Sprintf("%s", err))
		return nil, err

	}

	if notaryCommand == "list" {
		strSplit := strings.Split(stdOutStdErrString, "\n")[2:]

		for _, element := range strSplit {

			var digestData []string

			elementSplit := strings.Split(element, " ")
			for _, str := range elementSplit {
				if str != "" {
					digestData = append(digestData, str)
				}
			}
			digests = append(digests, NotaryList{Name: digestData[0], Digest: digestData[1], Size: digestData[2], Role: digestData[3]})
		}

	} else if notaryCommand == "lookup" {
		digestData := strings.Split(stdOutStdErrString, " ")
		digests = append(digests, NotaryList{Name: digestData[0], Digest: digestData[1], Size: digestData[2], Role: ""})
	}

	return digests, nil
}
