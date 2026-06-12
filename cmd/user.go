package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Anthill-Software/GatherPipe/core"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Gérer les utilisateurs de l'orchestrateur météo",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "Liste les utilisateurs de la console SSH",
	Run: func(cmd *cobra.Command, args []string) {
		usersFilePath := usersFile
		if usersFilePath == "" {
			usersFilePath = core.DefaultUsersFile
		}

		users, err := core.ListUsers(usersFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, core.Prefix.Error+" Erreur lors du chargement des utilisateurs : %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Utilisateurs dans %s :\n", usersFilePath)
		for _, u := range users {
			fmt.Printf(" - %s\n", u.Username)
		}
	},
}

var userAddCmd = &cobra.Command{
	Use:   "add [username] [ssh-public-key]",
	Short: "Ajoute ou met à jour un utilisateur de la console SSH",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		pubKey := args[1]

		usersFilePath := usersFile
		if usersFilePath == "" {
			usersFilePath = core.DefaultUsersFile
		}

		err := core.AddUser(usersFilePath, username, pubKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, core.Prefix.Error+" Erreur lors de l'ajout de l'utilisateur : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf(core.Prefix.Success+" Utilisateur '%s' ajouté ou mis à jour avec succès dans %s.\n", username, usersFilePath)
	},
}

var userRmCmd = &cobra.Command{
	Use:   "rm [username]",
	Short: "Supprime un utilisateur de la console SSH",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		usersFilePath := usersFile
		if usersFilePath == "" {
			usersFilePath = core.DefaultUsersFile
		}

		err := core.RemoveUser(usersFilePath, username)
		if err != nil {
			fmt.Fprintf(os.Stderr, core.Prefix.Error+" Erreur lors de la suppression de l'utilisateur : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf(core.Prefix.Success+" Utilisateur '%s' supprimé avec succès de %s.\n", username, usersFilePath)
	},
}

var userKeyRmCmd = &cobra.Command{
	Use:   "rm-key [username]",
	Short: "Supprime une clé publique d'un utilisateur de la console SSH",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		usersFilePath := usersFile
		if usersFilePath == "" {
			usersFilePath = core.DefaultUsersFile
		}

		// 1. On récupère l'utilisateur pour afficher ses clés
		user, err := core.GetUser(usersFilePath, username)
		if err != nil {
			fmt.Printf(core.Prefix.Error+" Erreur : %v\n", err)
			return
		}

		if len(user.PubKeys) == 0 {
			fmt.Printf(core.Prefix.Information+" L'utilisateur '%s' n'a actuellement aucune clé SSH enregistrée.\n", username)
			return
		}

		// 2. Affichage de la liste interactive
		fmt.Printf("Clés SSH enregistrées pour l'utilisateur '%s' :\n", username)
		for i, key := range user.PubKeys {
			// On tronque visuellement la clé pour que l'affichage reste propre au terminal
			parts := strings.Fields(key)
			displayKey := key
			if len(parts) >= 3 {
				// parts[0] = type (ssh-ed25519)
				// parts[1] = la clé b64
				// parts[2] = le commentaire (home@dev)
				displayKey = fmt.Sprintf("%s ... %s", key[:30], parts[2])
			} else if len(key) > 35 {
				// Repli de sécurité si la clé n'a pas de commentaire
				displayKey = key[:35] + "..."
			}
			fmt.Printf("  [%d] %s\n", i, displayKey)
		}
		fmt.Println()

		// 3. Prompt de sélection
		fmt.Print("Entrez l'index de la clé à supprimer (ou Entrée pour annuler) : ")
		var input string
		fmt.Scanln(&input)

		// Si l'utilisateur fait juste "Entrée", on quitte proprement sans erreur
		if input == "" {
			fmt.Println(core.Prefix.Error + " Opération annulée.")
			return
		}

		// Convertir la saisie en entier
		var selectedIndex int
		_, err = fmt.Sscanf(input, "%d", &selectedIndex)
		if err != nil || selectedIndex < 0 || selectedIndex >= len(user.PubKeys) {
			fmt.Println(core.Prefix.Error + " Choix invalide. L'index doit être un nombre de la liste.")
			return
		}

		// 4. Confirmation et suppression effective
		err = core.RemoveKeyByIndex(usersFilePath, username, selectedIndex)
		if err != nil {
			fmt.Printf(core.Prefix.Error+" Erreur lors de la suppression : %v\n", err)
			return
		}

		fmt.Printf(core.Prefix.Success+" La clé [%d] a été supprimée avec succès pour l'utilisateur '%s'.\n", selectedIndex, username)
	},
}

func init() {
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userKeyRmCmd)
	userCmd.AddCommand(userRmCmd)
	rootCmd.AddCommand(userCmd)
}
