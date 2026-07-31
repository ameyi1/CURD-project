package model

//这是解决硬编码问题的
type AppConfig struct {
	Server ServerConfig `yaml:"server"` // 小箱子1：装服务配置
	Mysql  MysqlConfig  `yaml:"mysql"`  // 小箱子2：装数据库配置
	Jwt    JwtConfig    `yaml:"jwt"`    // 小箱子3：装JWT配置
}
type ServerConfig struct {
	Port int `yaml:"port"`
}

// 小箱子2
type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

// 小箱子3
type JwtConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}
