//local the configs for server
package conf

import (
	"os"

	"github.com/spf13/viper"
	"log"
	"os/exec"
	"path/filepath"
)

//default params for config, perfer to load value from the file named "local.conf", the value is json style
type config struct {
	TLSCert    string //server certificate
	TLSKey     string //server private key
	TLS        bool
	Domain     string
	Port       int //port
	LogDir     string
	Data       string //dir to store the data
	Database   database
	Redis      redis
	SmtpServer smtpServer
}

// Postgres is here for embedded struct feature
type database struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

type redis struct {
	Host     string
	Port     int
	Num      int
	Password string
	MaxConn  int
}

type smtpServer struct {
	Account  string
	Password string
	Host     string
	Port     int
}

//default
var Config *config

func InitConfig() {
	viper.SetConfigName(ConfigName)             // name of config file (without extension)
	viper.AddConfigPath("/etc/" + ServerName)   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + ServerName) // call multiple times to add many search paths
	viper.AddConfigPath(getExecDir()) // call multiple times to add many search paths
	viper.AddConfigPath(".")                    // optionally look for config in the working directory
	viper.SetConfigType("toml")

	setDefault()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatal(err)
	}
	Config = &config{}
	err = viper.Unmarshal(Config)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(Config.LogDir, 0666)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(Config.Data, 0666)
	if err != nil {
		panic(err)
	}
}


func setDefault(){
	viper.SetDefault("tls",false)
	viper.SetDefault("domain","localhost")
	viper.SetDefault("port",8080)
	viper.SetDefault("logdir","/tmp/log")
	viper.SetDefault("data","/tmp/data")

	viper.SetDefault("database.host","127.0.0.1")
	viper.SetDefault("database.port",5432)
	viper.SetDefault("database.user","postgres")
	viper.SetDefault("database.sslmode","disable")

	viper.SetDefault("redis.host","127.0.0.1")
	viper.SetDefault("redis.port",6379)
	viper.SetDefault("redis.num",0)
	viper.SetDefault("redis.maxconn",1000)

	viper.SetDefault("smtpserver.port",25)
}

func getExecDir() string {
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	//    Is Symlink
	fi, err := os.Lstat(execPath)
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		execPath, err = os.Readlink(execPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	execDir := filepath.Dir(execPath)
	if execDir == "." {
		execDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	return execDir
}
