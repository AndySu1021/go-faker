package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-faker/codegen"
	"go-faker/config"
	"go-faker/db"
	"go-faker/logger"
	"go-faker/model"
	"io"
	"os"
)

var Num int

func init() {
	appConfig, err := config.InitConfig()
	if err != nil {
		fmt.Println("init config error: ", err)
		os.Exit(1)
	}

	// init logger
	if err = logger.InitZapLogger(appConfig.Logger); err != nil {
		fmt.Println("init logger error: ", err)
		os.Exit(1)
	}

	// init database
	if err = db.InitDatabase(appConfig.Database); err != nil {
		logger.Logger.Errorf("init database error: %s", err)
		os.Exit(1)
	}
}

func Do(stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:   "faker",
		Short: "Generate fake data",
		Long:  "This is a CLI tool for generating fake data to database",
	}

	createCmd.Flags().IntVarP(&Num, "num", "n", 1, "record num")
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(refreshCmd)

	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate fake data",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		m := model.ModelMap[modelName]

		if cmd.Flags().ArgsLenAtDash() > 0 {
			args = cmd.Flags().Args()[cmd.Flags().ArgsLenAtDash():]
		} else {
			args = args[1:]
		}

		total, err := db.DB.Create(m, Num, args)
		if err != nil {
			logger.Logger.Errorf("create fake data error: %s", err)
			return
		}

		if total > 1 {
			logger.Logger.Infof("%s -> create %d records", m.TableName(), total)
		} else {
			logger.Logger.Infof("%s -> create %d record", m.TableName(), total)
		}
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear fake data",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		m := model.ModelMap[modelName]

		if err := db.DB.Truncate(m); err != nil {
			logger.Logger.Errorf("clear fake data error: %s", err)
			return
		}

		logger.Logger.Infof("%s -> clear success", m.TableName())
	},
}

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Make model template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		if err := codegen.GenTableFile(modelName); err != nil {
			logger.Logger.Errorf("make template error: %s", err)
			return
		}
	},
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh model map",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := codegen.GenModelFile(); err != nil {
			logger.Logger.Errorf("make template error: %s", err)
			return
		}
	},
}
