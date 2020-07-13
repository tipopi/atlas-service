package api

import (
	"atlus-service/pkg/api/dao"
	"atlus-service/pkg/api/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"atlus-service/pkg/api/router"
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
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "./config/in-local.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Tcp port server listening on")
	StartCmd.PersistentFlags().Uint8VarP(&loglevel, "loglevel", "l", 0, "Log level")
	StartCmd.PersistentFlags().BoolVarP(&cors, "cors", "x", false, "Enable cors headers")
	StartCmd.PersistentFlags().BoolVarP(&cluster, "cluster", "s", false, "cluster-alone mode or distributed mod")
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
func zkCenterSetup(){
	viper.SetConfigFile("./zkconfig.yaml")
}
func setup()  {
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))
	//2.Set up configuration
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
	//3.Set up run mode
	mode := viper.GetString("mode")
	gin.SetMode(mode)
	//4.Set up database connection
	dao.Setup()
}