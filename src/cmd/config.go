package cmd

import (
	log "github.com/sirupsen/logrus"
)

type Config struct {

	LogCon LogConf `mapstructure:"log"`
}



type LogConf struct {
	FilePath			string
	FileName 			string
}

var Conf = new(Config)

func InitConfig() {
	//viper.AddConfigPath("./conf")
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Errorf("Fatal error config file: %s \n", err)
	//}
	//// 将读取的配置信息保存至全局变量Conf
	//if err := viper.Unmarshal(Conf); err != nil {
	//	log.Errorf("unmarshal conf failed, err:%s \n", err)
	//}
	//// 监控配置文件变化
	//viper.WatchConfig()
	////配置文件发生变化后同步到全局变量Conf
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("Config file changed, reloading...")
	//	if err := viper.Unmarshal(Conf); err != nil {
	//		log.Errorf("unmarshal conf failed, err:%s \n", err)
	//	}
	//})

	log.SetFormatter(&log.TextFormatter{DisableTimestamp : true})
	log.Info("Config init successfully")
}