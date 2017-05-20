package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/ebitgo/ebitgo.com/models"
	"github.com/ebitgo/ebitgo.com/models/logger"
	_ "github.com/ebitgo/ebitgo.com/routers"
)

func main() {
	logmodels.InitLogger()
	logmodels.Logger.Info(logmodels.SPrintInfo("main", "** Application service start **"))

	if setCurrentArgs(os.Args) == 0 {
		models.MainInit(models.DatabaseConfig)
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
		beego.Run()
	} else {
		printHelp()
		logmodels.Logger.Error(logmodels.SPrintError("main", "Missing required parameters, can not start!"))
		os.Exit(0)
	}

}

func setCurrentArgs(args []string) int {
	lenArgs := len(args)
	if lenArgs < 2 {
		return -1
	}
	for i := 1; i < lenArgs; i++ {
		switch strings.ToLower(args[i]) {
		case "--dbconfig":
			i += setDatabaseConfig(args[i:])
		case "--loglevel":
			i += setLogLevel(args[i:])
		}
	}
	return 0
}

func setDatabaseConfig(args []string) int {
	// fmt.Println(args)
	// fmt.Println(len(args))
	if args == nil || len(args) < 2 {
		return 0
	}
	tmps := strings.Split(args[1], ";")
	if len(tmps) != 4 {
		return 0
	}
	models.DatabaseConfig.Host = tmps[0]
	models.DatabaseConfig.Port = tmps[1]
	models.DatabaseConfig.UserName = tmps[2]
	models.DatabaseConfig.Password = tmps[3]
	logmodels.Logger.Trace(logmodels.SPrintTrace("main -> setDatabaseConfig", "%v", *models.DatabaseConfig))
	return 1
}

func setLogLevel(args []string) int {
	// fmt.Println(args)
	// fmt.Println(len(args))
	if args == nil || len(args) < 2 {
		return 0
	}
	arg := args[1]
	strings.ToUpper(arg)
	length := len(arg)
	logmodels.LogLevelConfig.Reset()
	for i := 0; i < length; i++ {
		switch arg[i] {
		case 'D':
			logmodels.LogLevelConfig.Debug = true
		case 'I':
			logmodels.LogLevelConfig.Info = true
		case 'T':
			logmodels.LogLevelConfig.Trace = true
		case 'W':
			logmodels.LogLevelConfig.Warning = true
		case 'E':
			logmodels.LogLevelConfig.Error = true
		case 'C':
			logmodels.LogLevelConfig.Critical = true
		}
	}
	logmodels.Logger.Trace(logmodels.SPrintTrace("main -> setLogLevel", "%v", *logmodels.LogLevelConfig))
	return 1
}

func printHelp() {
	fmt.Println(" ebitgo params\r\n")
	fmt.Println(" params include :")
	fmt.Println(" --dbconfig\t configuration for databse")
	fmt.Println("           \t format \"host;port;UserName;Password\"")
	fmt.Println("           \t Exsample --dbconfig \"localhost;321;;\"")
	fmt.Println("\r\n --loglevel\t configuration for log output")
	fmt.Println("           \t format \"DITWEC\" and default is all")
	fmt.Println("           \t D is Debug")
	fmt.Println("           \t I is Information")
	fmt.Println("           \t T is Trace")
	fmt.Println("           \t W is Warning")
	fmt.Println("           \t E is Error")
	fmt.Println("           \t C is Critical")
	fmt.Println("           \t Exsample --loglevel \"IEW\"")
}
