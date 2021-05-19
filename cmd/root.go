package cmd

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rifqiakrm/chat-grpc-go/app"
	trace "github.com/rifqiakrm/go-microservice-lib/tracer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

type server interface {
	Run(int) error
}

var (
	cfgFile   string
)

var rootCMD = &cobra.Command{
	Use:   "chat-grpc-go",
	Short: "GRPC chat with client and server side streaming",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize(splash, initconfig, GRPCService)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCMD.PersistentFlags().StringVar(&cfgFile, "configs", "configs/config.toml", "configs file (example is $HOME/configs.toml)")
}

// splash print plain text message to console
func splash() {
	fmt.Print(`
     _         _                       
  __| |_  __ _| |_   __ _ _ _ _ __  __ 
 / _| ' \/ _' |  _| / _' | '_| '_ \/ _|
 \__|_||_\__,_|\__| \__, |_| | .__/\__|
                    |___/    |_|       
`)
}

func initconfig() {
	viper.SetConfigType("toml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// search configs in home directory with name "configs" (without extension)
		viper.AddConfigPath("./configs")
		viper.SetConfigName(os.Getenv("CONFIG_FILE"))
	}

	//read env
	viper.AutomaticEnv()

	// if a configs file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("config application:", err)
	}

	log.Println("starting microservice using configs file:", viper.ConfigFileUsed())
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GRPCService() {
	var port = viper.GetInt("app.port")
	var jaegeraddr = fmt.Sprintf("%v:%v",
		viper.GetString("jaeger.host"),
		viper.GetString("jaeger.port"),
	)
	tracer, err := trace.New(viper.GetString("app.name"), jaegeraddr)
	if err != nil {
		log.Fatalf("trace new error: %v", err)
	}
	log.Println("jaeger initiated!")

	var srv server

	srv = app.NewChat(tracer)

	if err := srv.Run(port); err != nil {
		log.Fatalf("failed to start rpc server : %v", err)
	}
}
