package core

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gliderlabs/ssh"
	crypto_ssh "golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type AuthorizedUser struct {
	Username string   `yaml:"username"`
	PubKeys  []string `yaml:"pub_keys"`
}

// loadAuthorizedKeys charge les utilisateurs depuis le fichier YAML
func loadAuthorizedKeys(path string) ([]AuthorizedUser, error) {
	slog.Debug(fmt.Sprintf("[SSH] Chargement du fichier d'utilisateurs depuis : %s", path))

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []AuthorizedUser{}, nil
		}
		return nil, err
	}

	var users []AuthorizedUser
	if err := yaml.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// publicKeyHandler est la fonction de callback pour ton serveur SSH
func (m *PluginManager) publicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	incomingUser := ctx.User()
	slog.Debug(fmt.Sprintf("[SSH] Analyse de la clé pour l'utilisateur : %s", incomingUser))

	users, err := loadAuthorizedKeys(m.usersFile)
	if err != nil {
		slog.Error(fmt.Sprintf("Échec critique du chargement des clés : %v", err))
		return false
	}

	for _, u := range users {
		if u.Username != incomingUser {
			continue
		}

		fingerprint := crypto_ssh.FingerprintSHA256(key)
		slog.Debug("[SSH] Analyse de la clé",
			"user", incomingUser,
			"fingerprint", fingerprint[:15]+"...",
			"type", key.Type(),
		)

		slog.Debug(fmt.Sprintf("[SSH] Utilisateur '%s' trouvé. Analyse de ses %d clé(s) enregistrée(s)...", u.Username, len(u.PubKeys)))

		for i, keyStr := range u.PubKeys {
			authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyStr))
			if err != nil {
				slog.Debug(fmt.Sprintf("[SSH] Clé index %d malformée pour %s : %v", i, u.Username, err))
				continue
			}

			if ssh.KeysEqual(key, authorizedKey) {
				slog.Debug(fmt.Sprintf("[SSH] Clé valide (index %d) ! Authentification réussie pour : %s", i, u.Username))
				ctx.SetValue("username", u.Username)
				return true
			}
		}

		slog.Debug(fmt.Sprintf("[SSH] Aucune des clés stockées pour %s ne correspond.", u.Username))
		return false
	}

	slog.Error(fmt.Sprintf("Authentification refusée pour '%s' (utilisateur inconnu).", incomingUser))
	return false
}

// ListUsers retourne la liste des utilisateurs autorisés à se connecter à la console SSH
func ListUsers(filePath string) ([]AuthorizedUser, error) {
	return loadAuthorizedKeys(filePath)
}

// AddUser ajoute ou met à jour un utilisateur dans le fichier YAML cible
func AddUser(filePath string, username string, pubKey string) error {
	users, err := loadAuthorizedKeys(filePath)
	if err != nil {
		return err
	}

	found := false
	for i, u := range users {
		if u.Username == username {
			// L'utilisateur existe, on vérifie si la clé y est déjà
			keyExists := false
			for _, existingKey := range u.PubKeys {
				if existingKey == pubKey {
					keyExists = true
					break
				}
			}
			// Si la clé n'existe pas, on l'ajoute à son tableau
			if !keyExists {
				users[i].PubKeys = append(users[i].PubKeys, pubKey)
			}
			found = true
			break
		}
	}

	// Si l'utilisateur n'existe pas du tout, on le crée avec sa première clé
	if !found {
		users = append(users, AuthorizedUser{
			Username: username,
			PubKeys:  []string{pubKey},
		})
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(&users)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}

// RemoveUser supprime un utilisateur du fichier YAML cible
func RemoveUser(filePath string, username string) error {
	users, err := loadAuthorizedKeys(filePath)
	if err != nil {
		return err
	}

	newUsers := []AuthorizedUser{}
	found := false

	for _, u := range users {
		if u.Username == username {
			found = true
			continue // On ignore l'utilisateur à supprimer
		}
		newUsers = append(newUsers, u)
	}

	if !found {
		return errors.New("utilisateur introuvable")
	}

	data, err := yaml.Marshal(&newUsers)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}

// GetUser retourne un utilisateur spécifique s'il existe
func GetUser(filePath string, username string) (*AuthorizedUser, error) {
	users, err := loadAuthorizedKeys(filePath)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("utilisateur introuvable")
}

// RemoveKeyByIndex supprime la clé à l'index donné et sauvegarde le fichier
func RemoveKeyByIndex(filePath string, username string, index int) error {
	users, err := loadAuthorizedKeys(filePath)
	if err != nil {
		return err
	}

	userFound := false
	for i := range users {
		if users[i].Username == username {
			userFound = true

			// Validation de l'index
			if index < 0 || index >= len(users[i].PubKeys) {
				return errors.New("index de clé invalide")
			}

			// Suppression de l'élément dans le slice (gère le début, le milieu et la fin)
			users[i].PubKeys = append(users[i].PubKeys[:index], users[i].PubKeys[index+1:]...)
			break
		}
	}

	if !userFound {
		return errors.New("utilisateur introuvable")
	}

	// Sauvegarde du fichier mis à jour
	data, err := yaml.Marshal(&users)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}
