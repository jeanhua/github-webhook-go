package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jeanhua/github-webhook-go/lib"
	"github.com/jeanhua/jokerhttp"
	"github.com/jeanhua/jokerhttp/engine"
	"gopkg.in/yaml.v3"
)

func main() {
	var config Config
	file, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	if config.Port <= 0 || config.Port >= 65535 {
		fmt.Println("port error!")
		os.Exit(1)
	}

	joker := jokerhttp.NewEngine()
	joker.Init()
	joker.SetPort(config.Port)
	joker.Use(func(ctx *engine.JokerContex) {
		sigHeader := ctx.Request.Header.Get("X-Hub-Signature-256")
		if sigHeader == "" {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.ResponseWriter.Write([]byte("Missing signature header"))
			ctx.Abort()
			return
		}
		const prefix = "sha256="
		if !strings.HasPrefix(sigHeader, prefix) || len(sigHeader) <= len(prefix) {
			ctx.ResponseWriter.WriteHeader(401)
			ctx.ResponseWriter.Write([]byte("Invalid signature format"))
			ctx.Abort()
			return
		}
		ctx.Next()
	})

	logfile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	// Set up logging to file
	defer logfile.Close()
	log.SetOutput(logfile)
	log.SetFlags(log.LstdFlags)
	log.Println("Server started on port", config.Port)

	for _, s := range config.Service {
		name := s.Name
		path := s.Path
		secretpath := s.Secret
		script := s.Script
		if name == "" || path == "" || secretpath == "" || script == "" {
			fmt.Println("name or path or secret or script is empty!")
			os.Exit(1)
		}
		secret, err := os.ReadFile(secretpath)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
		joker.MapPost(path, func(request *http.Request, body []byte, params url.Values, setHeader func(key, value string)) (status int, response interface{}) {
			sigHeader := request.Header.Get("X-Hub-Signature-256")
			sig := sigHeader[len("sha256="):]
			delevery := request.Header.Get("X-GitHub-Delivery")
			setHeader("X-GitHub-Delivery", delevery)
			vertify, actualMAC := lib.VerifyHMAC(body, string(secret), sig)
			if vertify {
				log.Println(strings.ReplaceAll(`[Hook]:`, "Hook", s.Name) + s.Path + " @" + delevery)
				go lib.Exec_shell(script, config.Debug)
				return 200, "success"
			} else {
				if config.Debug {
					fmt.Println("Error request: receive="+sig, "expect="+actualMAC)
				}
				return 400, "error key"
			}
		})
	}
	joker.Run()
}

type Config struct {
	Port    int             `yaml:"port"`
	Debug   bool            `yaml:"debug"`
	Service []ServiceConfig `yaml:"service"`
}

type ServiceConfig struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Secret string `yaml:"secret"`
	Script string `yaml:"script"`
}
