# 📖 Guide de démarrage - Administration de GatherPipe

👉 [Accueil](../README.md)

__GatherPipe__ intègre une console d'administration sécurisée accessible via SSH. Ce guide explique comment démarrer le serveur, gérer les utilisateurs et configurer leurs accès par clé publique.

## 🚀 Démarrage du serveur

Par défaut, le serveur d'administration cherche sa configuration des utilisateurs dans `/etc/gatherpipe/users.yaml`. En mode développement, vous pouvez spécifier un chemin personnalisé.

```Bash
# Démarrage standard (utilise /etc/gatherpipe/users.yaml)
gatherpipe start

# Démarrage en mode développement avec un fichier local
gatherpipe start -c config.yaml -u ./data/users.yaml
```

## ⚙️ Configuration globale du serveur

Le serveur __GatherPipe__ s'appuie sur un fichier de configuration principal, généralement nommé `config.yaml`. Par défaut, il est recherché dans `/etc/gatherpipe/config.yaml`, mais vous pouvez spécifier un autre fichier au démarrage via le flag `--config` ou `-c`.

### Structure du fichier `config.yaml`

Voici l'anatomie complète du fichier de configuration du serveur :

```YAML
# ==========================================
# Configuration de base du serveur GatherPipe
# ==========================================

server:
  interval: 10s               # Intervale de rafraichissement des données (défaut: 5s)
  log_level: DEBUG            # Niveau de log (défaut: WARN) DEBUG|INFO|WARN|ERROR
  plugin:
    dir: "./plugins"          # Dossier d'installation des plugins (défaut: /etc/gatherpipe/plugins)
  ssh:
    host: "0.0.0.0"           # Adresse d'écoute (0.0.0.0 pour écouter sur tout le réseau) (défaut: localhost)
    port: 2222                # Port de la console d'administration SSH (défaut: 2233)
  doc:
    enabled: false            # Active.Déactive le serveur de documentation (défaut: true)
    host: "0.0.0.0"           # Adresse d'écoute (0.0.0.0 pour écouter sur tout le réseau) (défaut: localhost)
    port: 80                  # Port de la documentation (défaut: 8080)
    dir: "."                  # Dossier de la documentation (Position du README.md) (défaut: /usr/share/doc/gatherpipe)
```

### Options et Flags de la ligne de commande (CLI)

Le binaire `gatherpipe` accepte plusieurs arguments globaux pour configurer le comportement du serveur au runtime sans modifier le fichier YAML.

| Flag court | Flag long | Description | Valeur par défaut |
| :--- | :--- | :--- | :--- |
| `-c` | `--config` | Chemin vers le fichier `config.yaml` principal | `/etc/gatherpipe/config.yaml` |
| `-u` | `--users` | Chemin vers le fichier des clés SSH (`users.yaml`) | `/etc/gatherpipe/users.yaml` |

#### Exemples de commandes :

- Démarrer avec une configuration spécifique :

```Bash
gatherpipe start --config ./development/config.yaml
```

## 👥 Gestion des utilisateurs et des clés SSH

Toutes les commandes de gestion des accès sont regroupées sous la commande principale `user`.

### Ajouter un utilisateur ou une clé publique

Pour créer un utilisateur ou ajouter une nouvelle clé SSH à un utilisateur existant :

```Bash
gatherpipe user add [username] "[argument_cle_publique]"
```

_Exemple_ :

```Bash
gatherpipe user add admin "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIVLjkgCFsumFUYQkes6CnvmMsnObNbYTrfmWR7XJNqf vscode@c0783fcca6a0"
```

_Note_ : Si l'utilisateur existe déjà, la clé est simplement ajoutée à sa liste d'accès (multi-clés).

### Supprimer une clé spécifique (Mode Interactif)

Pour nettoyer les clés obsolètes (comme une clé de Dev Container ou d'un ancien PC) sans supprimer le compte :

```Bash
gatherpipe user rm-key [username]
```

L'outil affiche la liste des clés avec leur origine et demande l'index à supprimer :

```Plaintext
Clés SSH enregistrées pour l'utilisateur 'admin' :
  [0] ssh-ed25519 AAAAC3NzaC1lZDI1NT ... home@dev
  [1] ssh-ed25519 AAAAC3NzaC1lZDI1NT ... florent@Laptop-de-Florent
  [2] ssh-ed25519 AAAAC3NzaC1lZDI1NT ... vscode@c0783fcca6a0

Entrez l'index de la clé à supprimer (ou Entrée pour annuler) : 2
✅ La clé [2] a été supprimée avec succès pour l'utilisateur 'admin'.
```

### Supprimer complètement un utilisateur

Pour révoquer totalement tous les accès d'un utilisateur :

```Bash
gatherpipe user rm [username]
```

## 🛠️ Structure du fichier de configuration (`users.yaml`)

Le fichier généré et lu par GatherPipe adopte la structure suivante. Il est modifiable manuellement, mais il est recommandé de passer par la CLI pour éviter les erreurs de syntaxe YAML.

```YAML
- username: admin
  pub_keys:
    - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIVLjkgCFsumFUYQkes6CnvmMsnObNbYTrfmWR7XJNqf vscode@c0783fcca6a0
```

---

[Précédent](../README.md) - [Suivant](SSH_GUIDE.md)

© 2026 F. Colinet - **GatherPipe** est distribué sous licence **CC BY-NC-SA 4.0**.
