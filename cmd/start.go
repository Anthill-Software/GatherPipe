package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Anthill-Software/GatherPipe/core"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Lancer l'orchestrateur météo",
	Run: func(cmd *cobra.Command, args []string) {
		manager := core.NewPluginManager(appConfig, usersFile)
		manager.StartTime = time.Now()

		interval := appConfig.Server.Interval

		go func() {
			manager.StartAdminShell(appConfig)
		}()

		// Lancement du serveur de doc dans une goroutine dédiée
		if appConfig.Server.Doc.Enabled {
			go func() {
				addr := fmt.Sprintf("%s:%d", appConfig.Server.Doc.Host, appConfig.Server.Doc.Port)
				slog.Info(fmt.Sprintf("Démarrage du serveur de documentation sur http://%s", addr))

				mux := http.NewServeMux()
				mux.HandleFunc("/", core.MakeDocHandler(appConfig.Server.Doc.Dir))

				// ListenAndServe est bloquant, mais comme il est dans sa goroutine,
				// il n'impacte pas le reste de l'application.
				if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
					slog.Error(fmt.Sprintf("Échec du serveur de documentation: %v", err))
				}
			}()
		}

		// 1. Détection et Chargement dynamique
		slog.Info("Chargement des plugins...")

		// On scanne le dossier ./plugins
		manager.AutoLoad(appConfig)

		// 2. Nettoyage à l'arrêt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c
			slog.Info("Fermeture propre de GatherPipe...")
			manager.StopAll()
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}()

		// 3. Boucle de monitoring
		ticker := time.NewTicker(interval)
		slog.Info(fmt.Sprintf("GatherPipe est à l'écoute (Intervalle: %v) (Ctrl+C pour arrêter)", interval))

		for range ticker.C {
			for _, pluginDriver := range manager.Plugins() {
				if pluginDriver.Driver() == nil || pluginDriver.Status() != core.Status.Running {
					continue
				}

				data, err := pluginDriver.Driver().Fetch()
				if err != nil {
					slog.Error(fmt.Sprintf("Erreur Fetch: %v", err))
					continue
				}

				for _, pluginExporter := range manager.Plugins() {
					if pluginExporter.Exporter() == nil || pluginExporter.Status() != core.Status.Running {
						continue
					}

					err := pluginExporter.Exporter().Export(data)
					if err != nil {
						slog.Error(fmt.Sprintf("Erreur Export vers un plugin : %v", err))
					}

				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
