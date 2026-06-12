package core

import (
	"fmt"
	"sort"

	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func (m *PluginManager) renderConfig(term *term.Terminal, args []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ssCmd := args[1]

	switch ssCmd {

	case "help":
		format := "%-30s %s"

		helpText := fmt.Sprintf("%s", Colors.Yellow+"COMMANDES DISPONIBLES :"+Colors.Reset+"\n") +
			fmt.Sprintf(format, Colors.Cyan+"config show"+Colors.Reset, "Affiche la configuration actuelle du fichier (en attente de redémarrage).\n") +
			fmt.Sprintf(format, Colors.Cyan+"set [clé.sousclé] [valeur]"+Colors.Reset, "Modifie une valeur de configuration dans le fichier (en attente de redémarrage).\n") +
			fmt.Sprintf(format, Colors.Cyan+"config compare"+Colors.Reset, "Compare la configuration en mémoire avec celle du fichier et affiche les différences.\n") +
			fmt.Sprintf(format, Colors.Cyan+"config get [clé.sousclé]"+Colors.Reset, "Affiche la valeur actuelle d'une clé de configuration spécifique.\n") +
			fmt.Sprintf(format, Colors.Cyan+"config plugin [commande] [arguments]"+Colors.Reset, "Exécute une commande de configuration spécifique à un plugin actif.\n")

		term.Write([]byte(helpText))

	case "show":
		// Affiche la configuration du fichier (Pending)
		settings := viper.AllSettings()
		yamlData, err := yaml.Marshal(settings)
		if err != nil {
			term.Write([]byte(fmt.Sprintf("%s%s Erreur de mise en forme : %v%s\n", Colors.Red, Prefix.Error, err, Colors.Reset)))
			return
		}

		term.Write([]byte(fmt.Sprintf("\n%s--- # Configuration actuelle d'GatherPipe (Fichier) ---%s\n", Colors.White, Colors.Reset)))
		term.Write(yamlData)
		term.Write([]byte("---\n"))

	case "set":
		if len(args) < 4 {
			term.Write([]byte(fmt.Sprintf("%sUsage: config set [clé.sousclé] [valeur]%s\n", Colors.Yellow, Colors.Reset)))
			return
		}
		key := args[2]
		value := args[3]

		// Mise à jour de Viper et écriture immédiate sur le disque
		viper.Set(key, value)
		if err := viper.WriteConfig(); err != nil {
			term.Write([]byte(fmt.Sprintf("%s%s Erreur lors de la sauvegarde du fichier : %v%s\n", Colors.Red, Prefix.Error, err, Colors.Reset)))
			return
		}

		term.Write([]byte(fmt.Sprintf("%sModifié avec succès : %s = %s (En attente de redémarrage)%s\n", Colors.Green, key, value, Colors.Reset)))

	case "compare":
		if m.liveConfig == nil {
			term.Write([]byte(fmt.Sprintf("%s%s Impossible de comparer : la configuration en mémoire n'est pas initialisée.%s\n", Colors.Red, Prefix.Error, Colors.Reset)))
			return
		}

		// 1. Extraction et aplatissement de la config PENDING (Viper)
		pendingFlat := make(map[string]interface{})
		flattenMap("", viper.AllSettings(), pendingFlat)

		// 2. Extraction et aplatissement de la config LIVE (AppConfig en mémoire)
		liveBytes, err := yaml.Marshal(m.liveConfig)
		if err != nil {
			term.Write([]byte(fmt.Sprintf("%s%s Erreur de lecture de la config mémoire : %v%s\n", Colors.Red, Prefix.Error, err, Colors.Reset)))
			return
		}
		var liveSettings map[string]interface{}
		_ = yaml.Unmarshal(liveBytes, &liveSettings)
		liveFlat := make(map[string]interface{})
		flattenMap("", liveSettings, liveFlat)

		// 3. Comparaison
		term.Write([]byte(fmt.Sprintf("\n%s%s Changements en attente de redémarrage :%s\n", Colors.White, Prefix.Information, Colors.Reset)))
		term.Write([]byte("----------------------------------------------------------------------------\n"))

		hasChanges := false
		var keys []string
		for k := range pendingFlat {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			liveVal := liveFlat[key]
			pendingVal := pendingFlat[key]

			if fmt.Sprintf("%v", liveVal) != fmt.Sprintf("%v", pendingVal) {
				if liveVal == nil {
					liveVal = "(non défini)"
				}
				term.Write([]byte(fmt.Sprintf("%s%s %s%s\n", Colors.Yellow, Prefix.Warning, key, Colors.Reset)))
				term.Write([]byte(fmt.Sprintf("    • [Valeur actuelle]         : %v\n", liveVal)))
				term.Write([]byte(fmt.Sprintf("    • [Valeur après redémarrage]: %v\n\n", pendingVal)))
				hasChanges = true
			}
		}

		if !hasChanges {
			term.Write([]byte(" ✨ La configuration en mémoire est parfaitement synchronisée avec le fichier.\n"))
		} else {
			term.Write([]byte(fmt.Sprintf("%s%s Ces modifications seront appliquées au prochain redémarrage de l'orchestrateur.%s\n", Colors.White, Prefix.Information, Colors.Reset)))
		}
		term.Write([]byte("----------------------------------------------------------------------------\n"))

	case "get":
		if len(args) < 3 {
			term.Write([]byte(fmt.Sprintf("%sUsage: config get [clé.sousclé]%s\n", Colors.Yellow, Colors.Reset)))
			return
		}
		key := args[2]
		value := viper.Get(key)
		term.Write([]byte(fmt.Sprintf("%sParamètre '%s' : %s%s\n", Colors.Green, key, value, Colors.Reset)))

	case "plugin":
		if len(args) < 3 {
			term.Write([]byte(fmt.Sprintf("%sUsage: config plugin [commande] [arguments]%s\n", Colors.Yellow, Colors.Reset)))
			return
		}

		pluginCmd := args[2]
		pluginArgs := args[3:]

		term.Write([]byte(fmt.Sprintf("%sConfiguration des plugins%s\n\n", Colors.White, Colors.Reset)))

		output, intercepted := m.renderPluginCommand(pluginCmd, pluginArgs)
		if !intercepted {
			term.Write([]byte(fmt.Sprintf("%s[Erreur] Aucune commande de plugin active ne correspond à '%s'%s\n", Colors.Red, pluginCmd, Colors.Reset)))
		} else {
			term.Write([]byte(output + "\n"))
		}

	case "exporter":
		term.Write([]byte("Configuration serveur non implémentée.\n"))

	default:
		term.Write([]byte(fmt.Sprintf("%sSous-commande de configuration inconnue: %s%s\n", Colors.Red, args[1], Colors.Reset)))
	}
}

// renderPluginCommand vérifie si la commande saisie appartient à un plugin actif
func (m *PluginManager) renderPluginCommand(cmd string, args []string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.plugins {
		// On ne peut requêter que les plugins démarrés qui ont des commandes
		if p.status == Status.Running && p.commander != nil {
			cmds, err := p.commander.SupportedCommands()
			if err != nil {
				continue
			}

			for _, c := range cmds {
				// Si le nom de la commande correspond à ce que l'utilisateur a écrit
				if c.Name == cmd {
					output, err := p.commander.ExecuteCommand(cmd, args)
					if err != nil {
						return fmt.Sprintf("Erreur lors de l'exécution dans le plugin: %v", err), true
					}
					return output, true
				}
			}
		}
	}

	return "", false // Aucune commande de plugin trouvée
}

// Fonction utilitaire pour aplatir les structures de maps imbriquées
func flattenMap(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	for k, v := range src {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		if nested, ok := v.(map[string]interface{}); ok {
			flattenMap(fullKey, nested, dest)
		} else if nestedMapInterface, ok := v.(map[interface{}]interface{}); ok {
			converted := make(map[string]interface{})
			for nk, nv := range nestedMapInterface {
				converted[fmt.Sprintf("%v", nk)] = nv
			}
			flattenMap(fullKey, converted, dest)
		} else {
			dest[fullKey] = v
		}
	}
}
