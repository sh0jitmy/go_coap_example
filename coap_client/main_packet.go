package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"bytes"
	"time"

	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v2/dtls"
	"github.com/plgd-dev/go-coap/v2/message"
	//"github.com/plgd-dev/go-coap/v2/udp"
)



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
	path := "/packet"
	packetpath := "packet.bin"
	if len(os.Args) > 1 {
		packetpath = os.Args[1]
	} 

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//reqjson := &ConfigData{ID:"TestID",Value:"TestValue"} 
	//reqdata, err := json.Marshal(reqjson)
	reqdata ,err := os.ReadFile(packetpath)
	if err != nil {
                fmt.Println(err)
        }
        reader := bytes.NewReader(reqdata)
	resp, err := co.Post(ctx, path,message.AppOctets,reader)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	length,_ := resp.BodySize()
	log.Println("Response BodySize:", length)
	readbuf,_ := resp.ReadBody()
	log.Println("Response Body", string(readbuf))
}
