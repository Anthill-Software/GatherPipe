package core

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

func (m *PluginManager) renderHelp(term *term.Terminal) {
	format := "%-30s %s"
	pluginFormat := "%-45s %s"

	// 1. Liste des commandes natives d'GatherPipe
	helpText := fmt.Sprintf("%s", Colors.Yellow+"COMMANDES DISPONIBLES :"+Colors.Reset+"\n") +
		fmt.Sprintf(format, Colors.Cyan+"list"+Colors.Reset, "Affiche les drivers et exporters chargés et leur statut.\n") +
		fmt.Sprintf(format, Colors.Cyan+"stop [name]"+Colors.Reset, "Arrête un plugin spécifique et libère les ressources.\n") +
		fmt.Sprintf(format, Colors.Cyan+"start [name]"+Colors.Reset, "Démarre un plugin spécifique.\n") +
		fmt.Sprintf(format, Colors.Cyan+"restart [name]"+Colors.Reset, "Arrête et redémarre un plugin spécifique.\n") +
		fmt.Sprintf(format, Colors.Cyan+"reload"+Colors.Reset, "Redémarre l'intégralité du système de plugins.\n") +
		fmt.Sprintf(format, Colors.Cyan+"install [id]"+Colors.Reset, "Installe un plugin depuis le catalogue.\n") +
		fmt.Sprintf(format, Colors.Cyan+"uninstall [name]"+Colors.Reset, "Désinstalle un plugin spécifique.\n") +
		fmt.Sprintf(format, Colors.Cyan+"config [target]"+Colors.Reset, "Configure le 'server' ou un 'plugin' actif. 'help' pour plus d'informations.\n") +
		fmt.Sprintf(format, Colors.Cyan+"stats"+Colors.Reset, "Affiche les métriques de collecte (records, erreurs).\n") +
		fmt.Sprintf(format, Colors.Cyan+"catalog"+Colors.Reset, "Affiche le catalogue de plugins.\n") +
		fmt.Sprintf(format, Colors.Cyan+"clear"+Colors.Reset, "Efface l'écran de la console.\n") +
		fmt.Sprintf(format, Colors.Cyan+"help"+Colors.Reset, "Affiche ce menu d'aide.\n") +
		fmt.Sprintf(format, Colors.Cyan+"exit|quit"+Colors.Reset, "Ferme la session SSH.\n")

	// 2. Génération dynamique de l'aide des plugins connectés via RPC
	m.mu.RLock()
	var pluginHelp strings.Builder
	hasPluginCmds := false

	for _, p := range m.plugins {
		if p.status == Status.Running && p.commander != nil {
			cmds, err := p.commander.SupportedCommands()
			if err == nil && len(cmds) > 0 {
				if !hasPluginCmds {
					pluginHelp.WriteString("\n" + Colors.Yellow + "CONFIGURATIONS EXTENDED (PLUGINS ACTIFS) :" + Colors.Reset + "\n")
					hasPluginCmds = true
				}
				for _, c := range cmds {
					// On construit la syntaxe complète : config plugin [usage_ou_nom]
					syntax := fmt.Sprintf("config plugin %s", c.Name)
					if c.Usage != "" {
						syntax = fmt.Sprintf("config plugin %s", c.Usage)
					}
					pluginHelp.WriteString(fmt.Sprintf(pluginFormat, Colors.Cyan+syntax+Colors.Reset, c.Description+"\n"))
				}
			}
		}
	}
	m.mu.RUnlock()

	// Assemblage final
	helpText += pluginHelp.String() + "\n" + fmt.Sprintf("%s", Colors.Yellow+"Astuce: Utilisez le nom exact affiché dans 'list' pour les commandes stop/start."+Colors.Reset+"\n")

	// Nettoyage des sauts de ligne pour le protocole SSH/Terminal
	cleanedText := strings.ReplaceAll(helpText, "\r\n", "\n")
	term.Write([]byte(strings.ReplaceAll(cleanedText, "\n", "\r\n")))
}
