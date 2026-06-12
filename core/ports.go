package core

import (
	"time"

	"github.com/hashicorp/go-plugin"
)

// La donnée universelle pour GatherPipe
type Metric struct {
	Timestamp time.Time
	ID        string
	Value     string
	Format    string
	Unit      string
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GatherPipe_PLUGIN",
	MagicCookieValue: "GatherPipe",
}

// ========== Interface driver ==========
// L'interface qu'un plugin "Driver" doit implémenter
// Driver est le port d'entrée (Infrastructure -> Core)
type Driver interface {
	Name() (string, error)
	Init(config map[string]string) error
	Fetch() ([]Metric, error)
}

// ========== Interface exporter ==========
// L'interface qu'un plugin "Exporter" doit implémenter
// Exporter est le port de sortie (Core -> Infrastructure)
type Exporter interface {
	Name() (string, error)
	Init(config map[string]string) error
	Export(metrics []Metric) error
}

// ========== Interface commander ==========
// CommandArg définit la structure d'une commande exposée par un plugin
type CommandArg struct {
	Name        string
	Usage       string
	Description string
}

// Commander définit l'interface que les plugins doivent implémenter pour étendre la console
type Commander interface {
	SupportedCommands() ([]CommandArg, error)
	ExecuteCommand(cmd string, args []string) (string, error)
}
