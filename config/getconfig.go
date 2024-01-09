package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	Bucketname	string `yaml:"bucketname"`
	Mailhost	string `yaml:"mailhost"`
	Mailport	int `yaml:"mailport"`
	User	string `yaml:"user"`
	Password	string `yaml:"password"`
}

var (
  conf *Configuration
)

func readConf() *Configuration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func GetConfig() Configuration {
	conf = readConf()
	return Configuration {
		Bucketname: conf.Bucketname, 
		Mailhost: conf.Mailhost, 
		Mailport: conf.Mailport, 
		User: conf.User, 
		Password: conf.Password, 
	}
}