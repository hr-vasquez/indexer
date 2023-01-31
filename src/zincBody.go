package main

type BodyContent struct {
	Index   string              `json:"index"`
	Records []map[string]string `json:"records"`
}
