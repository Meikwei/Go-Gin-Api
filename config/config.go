package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gopkg.in/ini.v1"
)

//app struct
type App struct {
	Template      string
	PageSize      int
	JwtSecret     string
	JwtExpiresAt  time.Duration
	SigningMethod string
}

//server struct
type Server struct {
	HttpAddress  string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

//Mysql struct
type Mysql struct {
	LogMode       bool
	MysqlUser     string
	MysqlPassword string
	MysqlHost     string
	MysqlName     string
	MysqlPrefix   string
	MaxLifetime   time.Duration
}

type Log struct {
	Level        string
	Formatter    string
	ReportCaller bool
}

var (
	Cfg             *ini.File
	RunMode         string
	configPathError = errors.New("Can not get current file info")
	currentPath     = currentFile()
	AppSetting      = &App{}
	ServerSetting   = &Server{}
	MysqlSetting    = &Mysql{}
	LogSetting      = &Log{}
)

//get config path Single
func getConfigPath(path string) (file string) {
	return fmt.Sprintf("%s/%s", path, "app.ini")
}

func init() {
	InitConfig()
}

func currentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(configPathError)
	}
	return fmt.Sprintf("%s/..", file)
}

func InitConfig() {
	iniPath := getConfigPath(currentPath)
	var err error
	Cfg, err = ini.Load(iniPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
	LoadBasice()
	LoadApp()
	LoadServer()
	LoadDatabase()
	loadLog()
}

//加载基础配置
func LoadBasice() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

//加载app配置
func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	err = sec.MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}

	AppSetting.JwtExpiresAt = time.Duration(sec.Key("JWT_EXPIRE_TIME").MustInt(10)) * time.Minute
}

//加载http服务配置
func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	err = sec.MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeout = time.Duration(sec.Key("ReadTimeout").MustInt(60)) * time.Second
	ServerSetting.WriteTimeout = time.Duration(sec.Key("WriteTimeout").MustInt(60)) * time.Second
}

//加载数据库配置
func LoadDatabase() {
	sec, err := Cfg.GetSection("database-mysql")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	err = sec.MapTo(MysqlSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo MysqlSetting err: %v", err)
	}

	MysqlSetting.MaxLifetime = time.Duration(sec.Key("MaxLifetime").MustInt(60)) * time.Second
}

func loadLog() {
	sec, err := Cfg.GetSection("log")
	if err != nil {
		log.Fatalf("Fail to get section 'log': %v", err)
	}

	err = sec.MapTo(LogSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo LogSetting err: %v", err)
	}

}
