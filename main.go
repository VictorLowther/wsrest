package main

import (
	"encoding/json"
	"flag"
	"github.com/VictorLowther/soap"
	"github.com/VictorLowther/wsman"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var listen string

func init() {
	flag.StringVar(&listen, "listen", "", "Address:port to listen on.  Address can be left blank.")
}

type Request struct {
	Endpoint    string
	Username    string
	Password    string
	Method      string
	ResourceURI string
	Options     []string
	Selectors   []string
	Parameters  []string
}

func handler(res http.ResponseWriter, req *http.Request) {
	request := Request{}
	if req.Method != "POST" {
		http.Error(res, "We only accept POST", 500)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	client := wsman.NewClient(request.Endpoint, request.Username, request.Password)
	var msg *wsman.Message
	switch request.Method {
	case "Identify":
		reply, err := client.Identify()
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		res.Header().Set("content-type", soap.ContentType)
		res.Write(reply.Bytes())
		return
	case "Enumerate":
		msg = client.Enumerate(request.ResourceURI)
	case "EnumerateEPR":
		msg = client.EnumerateEPR(request.ResourceURI)
	case "Get":
		msg = client.Get(request.ResourceURI)
	case "Invoke":
		splitIdx := strings.LastIndex(request.ResourceURI, "/")
		resource := request.ResourceURI[:splitIdx]
		method := request.ResourceURI[splitIdx+1:]
		msg = client.Invoke(resource, method)
	default:
		msg = client.NewMessage(request.Method)
		if len(request.ResourceURI) > 0 {
			msg.SetHeader(wsman.Resource(request.ResourceURI))
		}
	}
	if len(request.Options) > 0 {
		msg.Options(request.Options...)
	}
	if len(request.Selectors) > 0 {
		msg.Selectors(request.Selectors...)
	}
	if len(request.Parameters) > 0 {
		msg.Parameters(request.Parameters...)
	}
	reply, err := msg.Send()
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	res.Header().Set("content-type", soap.ContentType)
	res.Write(reply.Bytes())
	return
}

func main() {
	flag.Parse()

	if listen == "" {
		flag.Usage()
		os.Exit(1)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(listen, nil))
}
