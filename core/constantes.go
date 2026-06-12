package core

import "fmt"

var Version = "1.0.0-RC1"

var DefaultConfigPath = "/etc/gatherpipe"
var DefaultConfigFilename = "config"
var DefaultUsersFile = DefaultConfigPath + "/users.yaml"

var Status = struct {
	Running string
	Stopped string
}{
	Running: "RUNNING",
	Stopped: "STOPPED",
}

var prefixFormat = "%-15s %s"

var Prefix = struct {
	Success     string
	Error       string
	Warning     string
	Information string
	Debug       string
	Restart     string
	Install     string
	Delete      string
}{
	Success:     Colors.Green + fmt.Sprintf(prefixFormat, "[Succès]") + Colors.Reset,
	Error:       Colors.Red + fmt.Sprintf(prefixFormat, "[Erreur]") + Colors.Reset,
	Warning:     Colors.Yellow + fmt.Sprintf(prefixFormat, "[Attention]") + Colors.Reset,
	Information: Colors.Blue + fmt.Sprintf(prefixFormat, "[Info]") + Colors.Reset,
	Debug:       Colors.Magenta + fmt.Sprintf(prefixFormat, "[Debug]") + Colors.Reset,
	Restart:     Colors.Cyan + fmt.Sprintf(prefixFormat, "[Redémarrage]") + Colors.Reset,
	Install:     Colors.Green + fmt.Sprintf(prefixFormat, "[Installation]") + Colors.Reset,
	Delete:      Colors.Red + fmt.Sprintf(prefixFormat, "[Suppression]") + Colors.Reset,
}

var Colors = struct {
	Blue    string
	Cyan    string
	Green   string
	Red     string
	Reset   string
	White   string
	Yellow  string
	Magenta string
}{
	Blue:    "\033[34m",
	Cyan:    "\033[36m",
	Green:   "\033[32m",
	Red:     "\033[31m",
	Reset:   "\033[0m",
	White:   "\033[1m",
	Yellow:  "\033[33m",
	Magenta: "\033[35m",
}

var BoldColors = struct {
	Blue   string
	Cyan   string
	Green  string
	Red    string
	Reset  string
	White  string
	Yellow string
}{
	Blue:   "\033[1;34m",
	Cyan:   "\033[1;36m",
	Green:  "\033[1;32m",
	Red:    "\033[1;31m",
	Reset:  "\033[1;0m",
	White:  "\033[1;1m",
	Yellow: "\033[1;33m",
}
