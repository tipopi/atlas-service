package api

import (
	error2 "atlas-service/pkg/api/error"
	"atlas-service/pkg/api/log"
	"atlas-service/pkg/api/router"
	"atlas-service/pkg/api/zk"
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var (
	config   string
	port     string
	loglevel uint8
	cors     bool
	cluster  bool
	useZk    bool
	zkHost   string
	zkPath   string
	//StartCmd : set up restful api server
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Start atlas API server",
		Example: "atlas server -c config/in-local.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			//组件加载
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "./config/local.json", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Tcp port server listening on")
	StartCmd.PersistentFlags().Uint8VarP(&loglevel, "loglevel", "l", 0, "Log level")
	StartCmd.PersistentFlags().BoolVarP(&cors, "cors", "x", false, "Enable cors headers")
	StartCmd.PersistentFlags().BoolVarP(&cluster, "cluster", "s", false, "cluster-alone mode or distributed mod")
	StartCmd.PersistentFlags().BoolVarP(&useZk, "useZk", "u", false, "config mode : zk mode:true ; local mode:false")
	StartCmd.PersistentFlags().StringVarP(&zkHost, "zkHost", "", "", "Zk host of configuration center")
	StartCmd.PersistentFlags().StringVarP(&zkPath, "zkPath", "", "/atlas/config", "Zk path of configuration center")
}

func run() error {
	engine := gin.Default()
	router.SetUp(engine)
	return engine.Run(":" + port)
}
func usage() {
	usageStr := `
    _       _____    _         _      ____     
U  /"\  u  |_ " _|  |"|    U  /"\  u / __"| u  
 \/ _ \/     | |  U | | u   \/ _ \/ <\___ \/   
 / ___ \    /| |\  \| |/__  / ___ \  u___) |   
/_/   \_\  u |_|U   |_____|/_/   \_\ |____/>>  
 \\    >>  _// \\_  //  \\  \\    >>  )(  (__) 
(__)  (__)(__) (__)(_")("_)(__)  (__)(__)                         
`
	fmt.Printf("%s\n", usageStr)
}
func ConfigSetup() {
	defer loadError()
	if useZk {
		zkLoad()
	} else {
		localLoad()
	}
}
func loadError() {
	if r := recover(); r != nil {
		zk.GetConfig().Close()
		switch r.(type) {
		case error2.ZkConfigError:
			log.Error(r.(error).Error())
			localLoad()
		default:
			panic(r)
		}
	}
}

//基本config
func localLoad() {
	viper.SetConfigFile(config)
	content, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}
	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//viper配置发生变化了 执行响应的操作
		log.Info(viper.GetString("test"))
	})
}
func zkLoad() {
	e := func(err error) {
		panic(error2.ZkConfigError{Msg: err.Error()})
	}
	if zkHost != "" {
		viper.Set("config-center.host", zkHost)
		viper.Set("config-center.path", zkPath)
	} else {
		viper.SetConfigFile("./config/zkConfig.yaml")
		err := viper.ReadInConfig()
		error2.CheckError(err, false, e)
	}
	host := viper.GetString("config-center.host")
	path := viper.GetString("config-center.path")
	//连接zk
	config, err := zk.SetConfig(host, path)
	error2.CheckError(err, false, e)
	viper.SetConfigType("json")
	err = viper.ReadConfig(bytes.NewBuffer(config))
	error2.CheckError(err, false, e)

}
func setup() {
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))
	ConfigSetup()
	//3.Set up run mode
	mode := viper.GetString("mode")
	gin.SetMode(mode)
	//4.Set up database connection
	//dao.Setup()
}
