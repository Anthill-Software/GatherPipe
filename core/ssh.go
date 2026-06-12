package core

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

func (m *PluginManager) StartAdminShell(config *Config) {
	server := &ssh.Server{
		Addr:             fmt.Sprintf("%s:%d", config.Server.Ssh.Host, config.Server.Ssh.Port),
		PasswordHandler:  nil,
		PublicKeyHandler: m.publicKeyHandler,
		Handler: func(s ssh.Session) {

			// Affichage du MOTD
			motd := m.getColoredMOTD(s.User(), config.Server.Ssh.Port, config.Server.Version)
			fmt.Fprint(s, motd)

			// 1. On crée un terminal interactif sur la session SSH
			term := term.NewTerminal(s, fmt.Sprintf("%s%s@%s%s$ ",
				Colors.Green,
				s.User(),
				"gatherpipe",
				Colors.Reset,
			))

			for {
				// 2. On lit la ligne (ReadLine gère le buffer et la touche Entrée)
				line, err := term.ReadLine()
				if err != nil {
					break
				}

				args := strings.Fields(line) // Découpe proprement (gère les espaces multiples)
				if len(args) == 0 {
					continue
				}

				cmd := args[0]

				switch cmd {
				case "list":
					m.renderList(term)
				case "stop":
					m.renderStop(term, args)
				case "start":
					m.renderStart(term, args)
				case "restart":
					m.renderRestart(term, args)
				case "reload":
					m.renderReload(term, config)
				case "install":
					m.renderInstall(term, config, args)
				case "uninstall":
					m.renderUninstall(term, args)
				case "config":
					m.renderConfig(term, args)
				case "stats":
					m.renderStats(term)
				case "catalog":
					m.renderCatalog(term)
				case "clear":
					m.renderClear(term)
				case "help":
					m.renderHelp(term)
				case "exit", "quit":
					term.Write([]byte("Fermeture de la session...\n"))
					return
				case "":
					continue
				default:
					term.Write([]byte(fmt.Sprintf("Commande inconnue : %s\n", cmd)))
					m.renderHelp(term)
				}
			}
		},
	}

	slog.Info(fmt.Sprintf("🔐 Console admin dispo sur le port %s", server.Addr))
	// On lance le serveur
	err := server.ListenAndServe()
	if err != nil {
		slog.Error(fmt.Sprintf("Erreur serveur SSH : %v\n", err))
	}
}
