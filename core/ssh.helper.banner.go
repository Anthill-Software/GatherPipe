package core

import (
	"fmt"
	"strconv"
	"strings"
)

var MOTD = BoldColors.Cyan + `┌─────────────────────────────────────────────────────┐` + "\r\n" +
	BoldColors.Cyan + `│     ___   _  _____ _  _ ___ ___  ___ ___ ___ ___    │` + "\r\n" +
	BoldColors.Cyan + `│   / __| /_\ |_   _| || | __| _ \| _ \_ _| _ \ __|   │` + "\r\n" +
	BoldColors.Cyan + `│  | (_ |/ _ \  | | | __ | _||   /|  _/| ||  _/ _|    │` + "\r\n" +
	BoldColors.Cyan + `│   \___/_/ \_\ |_| |_||_|___|_|_\|_| |___|_| |___|   │` + "\r\n" +
	BoldColors.Cyan + `│                                                     │` + "\r\n" +
	BoldColors.Cyan + `│           ` + Colors.Yellow + `>> GatherPipe SYSTEM CONSOLE <<` + BoldColors.Cyan + `           │` + "\r\n" +
	BoldColors.Cyan + `│                                                     │` + "\r\n" +
	BoldColors.Cyan + `│ ` + Colors.Reset + ` Version: %-17s` + Colors.Green + `Status: Operational` + BoldColors.Cyan + `      │` + "\r\n" +
	BoldColors.Cyan + `│ ` + Colors.Reset + ` Admin Port: %-14s` + Colors.Reset + `User: %-19s` + BoldColors.Cyan + `│` + "\r\n" +
	BoldColors.Cyan + `└─────────────────────────────────────────────────────┘` + Colors.Reset + "\r\n"

func (m *PluginManager) getColoredMOTD(user string, port int, version string) string {
	sPort := strconv.Itoa(port)
	return strings.ReplaceAll(fmt.Sprintf(MOTD, version, sPort, user), "\n", "\r\n")
}
