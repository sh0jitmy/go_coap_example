package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"io"
	"encoding/json"

	piondtls "github.com/pion/dtls/v2"
	coap "github.com/plgd-dev/go-coap/v2"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

type ResData struct {
	ErrCode uint32 `json:"errcode"`
	ErrDetail string `json:"errdetail"`
}


func BodySize(r *mux.Message) (int64,error) {
	if r.Body == nil {
		return 0, nil
	}
	orig, err := r.Body.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	_, err = r.Body.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	size, err := r.Body.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	_, err = r.Body.Seek(orig, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size, nil
}


func ReadBody(r *mux.Message) ([]byte,error) {
	if r.Body == nil {
		return nil, nil
	}
	size, err := BodySize(r)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return nil, nil
	}
	_, err = r.Body.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	payload := make([]byte, 1024)
	if int64(len(payload)) < size {
		payload = make([]byte, size)
	}
	n, err := io.ReadFull(r.Body, payload)
	if (err == io.ErrUnexpectedEOF || err == io.EOF) && int64(n) == size {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return payload[:n], nil
}


func handleStatus(w mux.ResponseWriter, r *mux.Message) {
	log.Printf("got message in handleStatus:  %+v from %v\n", r, w.Client().RemoteAddr())
	resdata,_ := os.ReadFile("./status.json")
	//resdata,_ := ioutil.ReadFile("./status.json") 	
	log.Println("resdata",string(resdata))
	reader := bytes.NewReader(resdata)
	w.SetResponse(codes.GET,message.AppJSON,reader)
}

func handleConfig(w mux.ResponseWriter, r *mux.Message) {
	log.Printf("got message in handleStatus:  %+v from %v\n", r, w.Client().RemoteAddr())
	configbuf,_:=ReadBody(r)
	log.Println("Request Body",string(configbuf))
	resjson := &ResData{ErrCode:0,ErrDetail:"None"}	
	resdata, err := json.Marshal(resjson)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("resdata",string(resdata))
	reader := bytes.NewReader(resdata)
	w.SetResponse(codes.GET,message.AppJSON,reader)
}

func handlePacket(w mux.ResponseWriter, r *mux.Message) {
	log.Printf("got message in handleStatus:  %+v from %v\n", r, w.Client().RemoteAddr())
	payloadlen,_:=BodySize(r)
	log.Println("Request Body Size",payloadlen)
	configbuf,_:=ReadBody(r)
	log.Println("Request Body",string(configbuf))
	resjson := &ResData{ErrCode:0,ErrDetail:"None"}	
	resdata, err := json.Marshal(resjson)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("resdata",string(resdata))
	reader := bytes.NewReader(resdata)
	w.SetResponse(codes.GET,message.AppJSON,reader)
}

func main() {
	m := mux.NewRouter()
	m.Handle("/status", mux.HandlerFunc(handleStatus))
	m.Handle("/packet", mux.HandlerFunc(handlePacket))
	m.Handle("/config", mux.HandlerFunc(handleConfig))
	//log.Fatal(coap.ListenAndServe("udp", ":5688", m))
	log.Fatal(coap.ListenAndServeDTLS("udp", ":5688", &piondtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			fmt.Printf("Client's hint: %s \n", hint)
			return []byte{0xAB, 0xC1, 0x23}, nil
		},
		PSKIdentityHint: []byte("Pion DTLS Client"),
		CipherSuites:    []piondtls.CipherSuiteID{piondtls.TLS_PSK_WITH_AES_128_CCM_8},
	}, m))
}
