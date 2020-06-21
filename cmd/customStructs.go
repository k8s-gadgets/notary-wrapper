package main

type NotaryList struct {
	Name   string `json:"Name"`
	Digest string `json:"Digest"`
	Size   string `json:"Size"`
	Role   string `json:"Role"`
}

type RequestGun struct {
	NotaryServer string `json:"notaryServer"`
	Gun          string `json:"Gun"`
	Tag          string `json:"Tag"`
}

type VerifySHA struct {
	NotaryServer string `json:"notaryServer"`
	Gun          string `json:"Gun"`
	SHA          string `json:"SHA"`
}

type CustomError struct {
	Code string `json:"code"`
	Info string `json:"info"`
}
