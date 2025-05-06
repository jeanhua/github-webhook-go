package lib

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os/exec"
)

func ComputeHMAC(message []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(message []byte, secret, expectedMAC string) (bool, string) {
	actualMAC := ComputeHMAC(message, secret)
	return hmac.Equal([]byte(actualMAC), []byte(expectedMAC)), actualMAC
}

func Exec_shell(s string, debug bool) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	if debug {
		log.Println("exec shell: ", s)
		log.Println("exec result: ", out.String())
	}
}
