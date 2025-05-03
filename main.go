package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var port int
	debug := false
	for idx, param := range os.Args {
		if param == "-p" || param == "-port" {
			if idx == len(os.Args)-1 {
				fmt.Println("error port")
				os.Exit(1)
			}
			var err error
			port, err = strconv.Atoi(os.Args[idx+1])
			if err != nil {
				fmt.Println("error: ", err)
				os.Exit(1)
			}

		}
		if param == "debug=true" {
			debug = true
		}
	}
	_, err := os.Stat(".secret")
	if os.IsNotExist(err) {
		fmt.Println("you have not sava the secret key,enter please:")
		var input string
		fmt.Scanln(&input)
		if len(input) < 10 {
			fmt.Println("length < 10!")
			os.Exit(1)
		}
		err := os.WriteFile(".secret", []byte(input), 0600)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
	}
	key, err := os.ReadFile(".secret")
	keyText := string(key)
	keyText = strings.ReplaceAll(keyText, "\n", "")
	if debug {
		fmt.Println("use key: " + keyText)
	}
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	http.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		sigHeader := r.Header.Get("X-Hub-Signature-256")
		if sigHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing signature header"))
			return
		}
		const prefix = "sha256="
		if !strings.HasPrefix(sigHeader, prefix) || len(sigHeader) <= len(prefix) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid signature format"))
			return
		}
		sig := sigHeader[len(prefix):]
		delevery := r.Header.Get("X-GitHub-Delivery")
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		vertify, actualMAC := verifyHMAC(body, keyText, sig)
		w.Header().Set("X-GitHub-Delivery", delevery)
		if vertify {
			log.Println(exec_shell("./action.sh"))
			w.WriteHeader(200)
			w.Write([]byte("success"))
		} else {
			if debug {
				log.Println("Error request: receive="+sig, "expect="+actualMAC)
				os.WriteFile("request.log", body, 0755)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error key"))
		}
	})
	if port != 0 {
		http.ListenAndServe(":"+strconv.Itoa(port), nil)
	} else {
		http.ListenAndServe(":5599", nil)
	}
}

func computeHMAC(message []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

func verifyHMAC(message []byte, secret, expectedMAC string) (bool, string) {
	actualMAC := computeHMAC(message, secret)
	return hmac.Equal([]byte(actualMAC), []byte(expectedMAC)), actualMAC
}

func exec_shell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	return out.String(), err
}
