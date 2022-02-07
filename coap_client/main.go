package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"bytes"
	"time"
	"encoding/json"

	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v2/dtls"
	"github.com/plgd-dev/go-coap/v2/message"
	//"github.com/plgd-dev/go-coap/v2/udp"
)


type ConfigData struct {
	ID string `json:"ID"`
	Value string `json:"Value"`
}

func main() {
	co, err := dtls.Dial("localhost:5688", &piondtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			fmt.Printf("Server's hint: %s \n", hint)
			return []byte{0xAB, 0xC1, 0x23}, nil
		},
		PSKIdentityHint: []byte("Pion DTLS Client"),
		CipherSuites:    []piondtls.CipherSuiteID{piondtls.TLS_PSK_WITH_AES_128_CCM_8},
	})
	//co, err := udp.Dial("localhost:5688")

	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	path := "/status"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	log.Printf("Response payload: %+v", resp)
	//buf:= new(bytes.Buffer)
	//io.Copy(buf,resp.Body)
	//log.Println("Response Body:", string(buf.Bytes()))
	length,_ := resp.BodySize()
	log.Println("Response BodySize:", length)
	readbuf,_ := resp.ReadBody()
	log.Println("Response Body", string(readbuf))

	reqjson := &ConfigData{ID:"TestID",Value:"TestValue"} 
	reqdata, err := json.Marshal(reqjson)
        if err != nil {
                fmt.Println(err)
        }
        reader := bytes.NewReader(reqdata)
	resp, err = co.Put(ctx, "/config",message.AppJSON,reader)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	length,_ = resp.BodySize()
	log.Println("Response BodySize:", length)
	readbuf,_ = resp.ReadBody()
	log.Println("Response Body", string(readbuf))
}
