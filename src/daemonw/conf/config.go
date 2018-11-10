//local the configs for server
package conf

import (
	"os"

	"path/filepath"

	"github.com/koding/multiconfig"
	"os/exec"
	"log"
)

//default params for config, perfer to load value from the file named "local.conf", the value is json style
type config struct {
	TLSCert     string     `toml:"tls_cert"` //server certificate
	TLSKey      string     `toml:"tls_key"`  //server private key
	TLS         bool       `toml:"tls" default:"true"`
	UseAutoCert bool       `toml:"use_auto_cert" default:"false"`
	Domain      string     `toml:"domain" default:"localhost"`
	Port        int        `toml:"port" default:"8080"` //port
	LogDir      string     `toml:"log_dir" default:"/tmp/log"`
	Data        string     `toml:"data" default:"./data"` //dir to store the data
	Database    database   `toml:"database"`
	Redis       redis      `toml:"redis"`
	SmtpServer  smtpServer `toml:"smtp_server"`
}

// Postgres is here for embedded struct feature
type database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port" default:"5432"`
	Name     string `toml:"name" required:"true"`
	User     string `toml:"user" default:"postgres"`
	Password string `toml:"password"`
	SSLMode  string `toml:"ssl_mode" default:"disable"`
}

type redis struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port" default:"6379"`
	Num      int    `toml:"num" default:"0"`
	Password string `toml:"password"`
	MaxConn  int    `toml:"max_conn"`
}

type smtpServer struct {
	Account  string `toml:"account"`
	Password string `toml:"password"`
	Host     string `toml:"host"` //smtp server to send email
	Port     int    `toml:"port" default:"25"`
}

//default
var Config *config
var BinDir = getExecDir()

func ParseConfig(path string) error {
	if path == "" {
		path = BinDir + "/config.toml"
	}
	loader := multiconfig.NewWithPath(path) // supports TOML and JSON
	Config = &config{}
	err := loader.Load(Config)
	if err != nil {
		return err
	}
	//init value
	initConfig(Config)
	return nil
}

func initConfig(c *config) {
	if !filepath.IsAbs(c.TLSCert) {
		c.TLSCert = filepath.Join(BinDir, c.TLSCert)
	}

	if !filepath.IsAbs(c.TLSKey) {
		c.TLSKey = filepath.Join(BinDir, c.TLSKey)
	}

	if !filepath.IsAbs(c.LogDir) {
		c.LogDir = filepath.Join(BinDir, c.LogDir)
	}

	if !filepath.IsAbs(c.Data) {
		c.Data = filepath.Join(BinDir, c.Data)
	}

	err := os.MkdirAll(c.LogDir, 0777)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(c.Data, 0777)
	if err != nil {
		log.Fatal(err)
	}
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
