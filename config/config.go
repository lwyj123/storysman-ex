package config

import (
	"io/ioutil"

	"github.com/json-iterator/go"

	"github.com/bitly/go-simplejson"
)

var (
	json = jsoniter.ConfigDefault
)

// Config 配置项
type Config struct {
	DBConfigs              map[string]DBConfig `json:"dbs"`
	PSM                    string              `json:"psm"`
	Env                    string              `json:"env"`
	RedisClusterName       string              `json:"redis_cluster_name"`
	RedisHosts             []string            `json:"redis_hosts"`
	KafkaService           string              `json:"kafka_service"`       //kafka服务名, 通过consul发现
	KafkaRiscService       string              `json:"kafka_risc_service"`  //kafka服务名, 通过consul发现
	KafkaRisc2Service      string              `json:"kafka_risc2_service"` //kafka服务名, 通过consul发现
	NsqConnStr             string              `json:"nsq_conn_str"`
	NsqLookupdService      string              `json:"nsq_lookupd_service"`
	NsqLookupCommonService string              `json:"nsq_lookupd_common_service"`
	NsqdService            string              `json:"nsqd_service"`
}

// DBConfig 数据库配置
type DBConfig struct {
	Database string        `json:"database"`
	Settings string        `json:"settings"`
	WriteDB  DBConnectInfo `json:"write"`
	ReadDB   DBConnectInfo `json:"read"`
}

// DBConnectInfo 数据库连接信息
type DBConnectInfo struct {
	AuthKey         string `json:"auth_key"`
	Consul          string `json:"consul"`
	UserName        string `json:"username"`
	Password        string `json:"password"`
	DefaultHostPort string `json:"default_host_port"`
}

// ENV环境变量
const (
	Prod    = "prod"    // 线上环境
	Dev     = "dev"     // 开发环境
	Staging = "staging" //测试环境
)

// Product 检查当前环境是否是线上环境
func (c *Config) Product() bool {
	return c.Env == Prod
}

// Instance 默认配置实例
var Instance *Config

// DBSettings 将conf文件载入成json形式
var DBSettings *simplejson.Json

// NewConfig 从文件中加载一个配置实例
func NewConfig(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}

	DBSettings, err = simplejson.NewJson(content)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Init 初始化加载配置信息
func Init(file string) error {
	if Instance != nil {
		return nil
	}
	conf, err := NewConfig(file)
	if err != nil {
		return err
	}
	// 如果没有配置env，默认是开发环境
	if len(conf.Env) == 0 {
		conf.Env = Dev
	}
	Instance = conf
	return nil
}
