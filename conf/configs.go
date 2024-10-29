package conf

type Configs struct {
	Server Server `mapstructure:"server"`
}

// Server 服务器配置
type Server struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}
