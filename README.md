# 🌦️ GatherPipe

**GatherPipe** est un orchestrateur modulaire écrit en Go. Il permet de collecter, traiter et exporter des données métriques via un système de plugins dynamiques pilotables en temps réel depuis une console SSH sécurisée.

[![Go CI](https://github.com/Anthill-Software/GatherPipe/actions/workflows/go.yml/badge.svg)](https://github.com/Anthill-Software/GatherPipe/actions/workflows/go.yml)

## 🚀 Caractéristiques

* **Architecture Modulaire** : Drivers (entrées) et Exporters (sorties) interchangeables et isolés.
* **Console SSH Intégrée** : Gestion des plugins et surveillance des métriques à distance.
* **Gestionnaire de Plugins** : Installation, mise à jour et suppression automatisées via un catalogue distant.
* **Documentation embarquée** : Consultation locale complète et hors-ligne de la doc technique directement depuis un navigateur web, sans dépendance externe.
* **Multi-plateforme** : Binaires natifs légers pour Linux (AMD64/ARM), Windows et macOS.
* **Installeur** : Paquets natifs pour Debian et Alpine (AMD64/ARM).
* **Service** : Intégration native Debian (systemd).
* **Sécurité Native** : Authentification stricte par paires de clés SSH (RSA/Ed25519) avec gestion multi-clés par utilisateur.

### Flux d'Orchestration

```plaintext
Capteurs (Drivers)  ──▶  GatherPipe Core  ──▶  Sorties (Exporters)
   (ARM/x64)            (SSH Admin)           (JSON/Terminal)
```

## 🛠️ Installation & Configuration

### Télécharger le package

Récupérez la version correspondant à votre architecture directement depuis les [Releases](https://github.com/Anthill-Software/GatherPipe/releases).

```Bash
# Exemple pour Debian AMD64
wget gatherpipe_v1.0.0-RC1_linux_amd64.deb
dpkg -i gatherpipe_v1.0.0-RC1_linux_amd64.deb

# Exemple pour Alpine AMD64
wget gatherpipe_v1.0.0-RC1_linux_amd64.apk
apk add --allow-untrusted gatherpipe_v1.0.0-RC1_linux_amd64.apk
```

### Télécharger le binaire

Récupérez la version correspondant à votre architecture directement depuis les [Releases](https://github.com/Anthill-Software/GatherPipe/releases).

```Bash
# Exemple pour Debian AMD64
wget https://github.com/Anthill-Software/GatherPipe/releases/download/v1.0.0-RC1/gatherpipe_1.0.0-RC1_linux_amd64.tar.gz
tar -xf gatherpipe_1.0.0-RC1_linux_amd64.tar.gz
```

## 📖 Documentation utilisateur

Pour la configuration, le déploiement en production ou la gestion avancée des utilisateurs (comme la suppression interactive de clés obsolètes), consultez le guide dédié :

👉 [Guide de démarrage](docs/START_GUIDE.md)

👉 [Guide d'utilisation de la console SSH](docs/SSH_GUIDE.md)

## 🏗️ Compilation locale (Développement)

Si vous souhaitez compiler et tester le projet vous-même :

```Bash
git clone [https://github.com/Anthill-Software/GatherPipe.git](https://github.com/Anthill-Software/GatherPipe.git)
cd gatherpipe
go build -o gatherpipe ./cmd/server/main.go
go test ./core/...
```

---

[Suivant](docs/START_GUIDE.md)

© 2026 F. Colinet - **GatherPipe** est distribué sous licence **CC BY-NC-SA 4.0**.
