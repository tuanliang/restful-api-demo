package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 全局config实例对象
// 也就是我们程序，在内存中的配置对象
// 程序内部获取配置，都通过读取该对象
// 该Config对象，什么时候被初始化？
//
//	       配置加载时:
//					    LoadConfigFromToml
//					    LoadConfigFromEnv
//
// 为了不被程序在运行时被恶意修改，设置为私有变量
var config *Config

// 全局Mysql客户端实例
var db *sql.DB

// 想获取config配置，单独提供函数
func C() *Config {
	return config
}

// 初始化一个有默认值的config对象
func NewDefaultConfig() *Config {
	return &Config{
		App:   NewDefaultApp(),
		Log:   NewDefaultLog(),
		MySQL: NewDefaultMySQL(),
	}
}

// Config 应用配置
// 通过封装为一个对象，来与外部配置进行对接
type Config struct {
	App   *App   `toml:"app"`
	Log   *Log   `toml:"log"`
	MySQL *MySQL `toml:"mysql"`
}

func NewDefaultApp() *App {
	return &App{
		Name: "dome",
		Host: "127.0.0.1",
		Port: "8050",
	}
}

type App struct {
	Name string `toml:"name" env:"APP_NAME"`
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func (a *App) HttpAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func NewDefaultMySQL() *MySQL {
	return &MySQL{
		Host:        "127.0.0.1",
		Port:        "3306",
		UserName:    "root",
		Password:    "root",
		Database:    "gotest1",
		MaxOpenConn: 10,
		MaxIdleConn: 5,
	}
}

// MySQL todo
type MySQL struct {
	Host     string `toml:"host" env:"MYSQL_HOST"`
	Port     string `toml:"port" env:"MYSQL_PORT"`
	UserName string `toml:"username" env:"MYSQL_USERNAME"`
	Password string `toml:"password" env:"MYSQL_PASSWORD"`
	Database string `toml:"database" env:"MYSQL_DATABASE"`
	// 使用的是Mysql连接池，需要池做一些规划配置：
	// 控制当前程序的Mysql打开的连接数
	MaxOpenConn int `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	// 控制Mysql复用，比如5，最多运行5个来复用
	MaxIdleConn int `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	// 一个连接的生命周期，和mysql server配置有关，必须小于server配置
	MaxLifeTime int `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	// Idle连接，最多允许存活多久
	MaxIdleTime int `toml:"max_idle_time" env:"MYSQL_MAX_idle_TIME"`
	// 作为私有变量，用于控制GetDB
	lock sync.Mutex
}

// 1.第一种方式，使用LoadGlobal 在加载时 初始化全局db实例
// 2.第二种方式，惰性加载,获取DB时，动态判断在初始化
func (m *MySQL) GetDB() *sql.DB {
	// 直接加锁，锁住临界区
	m.lock.Lock()
	defer m.lock.Unlock()

	if db == nil {
		// 如果实例不存在，就初始化一个新的实例
		conn, err := m.getDBConn()
		if err != nil {
			fmt.Println("getdb 错误:", err)
		}
		db = conn
	}
	return db
}

// 连接池，driverConn具体的连接对象，他维护着一个Socket
// pool []*driverConn,维护pool里面的连接都是可用的，定期检查conn的健康状况
// 某一个driverConn已经失效，driverConn.Reset()，清空该结构体的数据，
// 避免driverConn结构体的内存申请和释放的成本
func (m *MySQL) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

func NewDefaultLog() *Log {
	return &Log{
		Level:  "info",
		Format: TextFormat,
		To:     ToStdout,
	}
}

// Log todo
type Log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}
