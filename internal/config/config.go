package config

import (
	"flag"
	"log"
	"os"
)

type HttpServer struct {
	Address string 
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string  `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server"`
}

func MustLoad(){
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	
	if configPath == "" {
		flags := flag.String("config" , "", "path to the configuration file")
        flag.Parse()

		configPath = *flags

		if(configPath == ""){
			log.Fatal("config path is not set")
		}
	}
	
}