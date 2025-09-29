package config

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "log"
    )

type DBConfig struct {
    Port int `yaml:"port"`
    Host string `yaml:"host"`
    User string `yaml:"user"`
    DBname string `yaml:"dbname"`
    Password string `yaml:"password"`
}

type ServerConfig struct {
    Port int `yaml:"port"`
    Ipaddress string `yaml:"ipaddress"`
}

type Config struct {
    DB DBConfig `yaml:"database"`
    Server ServerConfig `yaml:"server"`
}

func CreateConfig(configPath string) Config {

    //data, err := ioutil.ReadFile("config.yaml")
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    return config
}

func CreateConnectString(config Config) string {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBname)
    return psqlInfo
}
