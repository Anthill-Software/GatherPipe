# 🔐 Guide d'utilisation - Console d'Administration SSH

👉 [Accueil](../README.md)

La console d'administration de __GatherPipe__ est une interface en ligne de commande (CLI) sécurisée accessible par SSH. Elle permet de piloter l'orchestrateur en temps réel, de superviser les flux météo, de gérer le cycle de vie des plugins et de modifier la configuration du serveur sans interrompre le processus principal.

## 🔒 Connexion à la console d'administration

Une fois les clés configurées, la connexion s'effectue avec n'importe quel client SSH standard sur le port dédié (2233 par défaut).

```Bash
ssh -i ~/.ssh/votre_cle_privee -p 2233 [username]@localhost
```

### 💡 Astuce de configuration SSH

Pour éviter de spécifier le port et la clé à chaque connexion, ajoutez ceci à votre fichier `~/.ssh/config` local :

```Plaintext
Host gatherpipe-admin
    HostName localhost
    Port 2233
    User admin
    IdentityFile ~/.ssh/gatherpipe
    IdentitiesOnly yes
```

Vous pourrez alors vous connecter instantanément avec la commande :

```Bash
ssh gatherpipe-admin
```

## 🛠️ Commandes générales & supervision

Ces commandes permettent de vérifier l'état de santé général du système et de naviguer dans la console.

- `help` : Affiche le menu d'aide principal contenant la liste des commandes.

- `stats` : Affiche les métriques de collecte de l'orchestrateur en temps réel (ex: nombre de records météo traités successfully, volume d'erreurs rencontrées par les drivers).

- `clear` : Nettoie l'écran du terminal pour une meilleure lisibilité.

- `exit` ou `quit` : Ferme proprement la session SSH en cours.

## 🔌 Gestion du cycle de vie des plugins

__GatherPipe__ repose sur une architecture modulaire composée de _Drivers_ (collecte) et d'_Exporters_ (diffusion).

- `list` : La commande centrale. Elle liste l'ensemble des drivers et exporters actuellement chargés dans l'orchestrateur, ainsi que leur statut d'exécution (`RUNNING`, `STOPPED`, etc.).

- `start [name]` : Démarre un plugin spécifique.

- `stop [name]` : Arrête un plugin spécifique et libère proprement ses ressources système (processus enfants).

- `restart [name]` : Effectue un arrêt suivi d'un redémarrage immédiat d'un plugin.

💡 __Astuce__ : Pour les commandes `start`, `stop` et `restart`, veillez à utiliser le __nom exact__ du plugin tel qu'il apparaît lors de l'exécution de la commande `list` (ex: `DummyDriver`).

- `reload` : Redémarre l'intégralité du système de plugins.

### 📦 Catalogue et installation

- `catalog` : Explore le catalogue distant ou local des plugins disponibles à l'installation.

- `install [id]` : Télécharge et installe un nouveau plugin à partir de son identifiant du catalogue.

- `uninstall [name]` : Supprime définitivement un plugin spécifique du système.

## ⚙️ Configuration dynamique (config)

La commande `config` permet d'inspecter et de modifier le comportement de GatherPipe. Les modifications de fichiers sont gérées de manière asynchrone (en attente de application).

### Gestion de la configuration globale du serveur

- `config show` : Affiche l'état actuel du fichier de configuration (`config.yaml`). Les modifications non appliquées y sont visibles.

- `config get [clé.sousclé]` : Permet de lire précisément la valeur actuelle d'un paramètre.

  - Exemple : `config.get server.interval`

- `config set [clé.sousclé] [valeur]` : Écrit une nouvelle valeur dans le fichier de configuration.

  - _Note_ : La modification est enregistrée dans le fichier mais reste __en attente de redémarrage__ pour être chargée en mémoire active.

- `config compare` : Outil de diagnostic très utile. Compare la configuration actuellement exécutée en mémoire vive (RAM) avec l'état réel du fichier sur le disque, puis liste les différences en attente d'application.

### Configuration dédiée aux plugins

- `config plugin [commande] [arguments]` : Transmet une commande de configuration spécifique directement à un plugin actif qui implémente cette capacité. Permet d'ajuster le comportement interne d'un driver ou d'un exporter à chaud.

---

[Précédent](START_GUIDE.md) - [Suivant](USECASE_GUIDE.md)

© 2026 F. Colinet - **GatherPipe** est distribué sous licence **CC BY-NC-SA 4.0**.
