package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Anthill-Software/GatherPipe/core"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var usersFile string
var appConfig *core.Config

var rootCmd = &cobra.Command{
	Use:   "gatherpipe",
	Short: "GatherPipe - Station météo modulaire",
	Long:  `GatherPipe est un serveur de station météo modulaire.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "fichier de config (default: "+core.DefaultConfigPath+"/"+core.DefaultConfigFilename+")")
	rootCmd.PersistentFlags().StringVarP(&usersFile, "users", "u", core.DefaultUsersFile, "fichier d'utilisateurs (default: "+core.DefaultUsersFile+")")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(core.DefaultConfigPath)
		viper.SetConfigName(core.DefaultConfigFilename)
	}

	viper.SetDefault("server.interval", "5s")
	viper.SetDefault("server.log_level", "WARN")

	viper.SetDefault("server.plugin.dir", "/etc/gatherpipe/plugins")

	viper.SetDefault("server.ssh.host", "localhost")
	viper.SetDefault("server.ssh.port", 2233)

	viper.SetDefault("server.doc.enabled", true)
	viper.SetDefault("server.doc.host", "localhost")
	viper.SetDefault("server.doc.port", 8080)
	viper.SetDefault("server.doc.dir", "/usr/share/doc/gatherpipe")

	viper.ReadInConfig()

	if err := viper.Unmarshal(&appConfig); err != nil {
		slog.Error(fmt.Sprintf("Erreur décodage config: %v", err))
	}

	initLogger(appConfig.Server.LogLevel)

	appConfig.Server.Version = core.Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initLogger(levelStr string) {
	var level slog.Level

	// 1. Parsing du niveau textuel (souvent lu depuis ta config)
	switch levelStr {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Configuration du handler "Pretty"
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: "2006-01-02 15:04:05",
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)

				switch level {
				case slog.LevelDebug:
					return slog.String(slog.LevelKey, core.Prefix.Debug)
				case slog.LevelInfo:
					return slog.String(slog.LevelKey, core.Prefix.Information)
				case slog.LevelWarn:
					return slog.String(slog.LevelKey, core.Prefix.Warning)
				case slog.LevelError:
					return slog.String(slog.LevelKey, core.Prefix.Error)
				}
			}

			return a
		},
	})

	slog.SetDefault(slog.New(handler))
}
