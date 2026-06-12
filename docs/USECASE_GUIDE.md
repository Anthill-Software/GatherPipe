# 📖 Guide des cas d'utilisation

Voici des exemples de cas d'utilisation. Les drivers et exporters cités sont potentiellement inexistant. 

## Environnement & Stations météorologiques

- __Description__ : Collecte de données environnementales complexes depuis du matériel professionnel ou des API.

- __Drivers__ : modbus-tcp, http-api-weather

- __Exporters__ : influxdb, weather-dashboard 

- __Valeur ajoutée__ : Permet d'agréger des capteurs physiques locaux et des prévisions cloud dans un seul flux normalisé.

## Domotique & Gestion de l'énergie

- __Description__ : Centralisation des flux énergétiques d'une maison ou d'un atelier (téléinformation Linky, production de panneaux solaires, sondes de température).

- __Drivers__ : mqtt-subscriber (flux Zigbee2MQTT), serial-tic (liaison série pour compteur électrique)

- __Exporters__ : home-assistant-ws (intégration directe), timescale-db

## Supervision d'infrastructure & SysAdmin

- __Description__ : Monitoring léger de parcs de serveurs, conteneurs ou clusters (Raspberry Pi, VMs) là où des agents lourds consommeraient trop de ressources.

- __Drivers__ : system-metrics (CPU, RAM, Disque), ntp-check (dérive d'horloge)

- __Exporters__ : prometheus-exporter, mqtt-forwarder

- __Valeur ajoutée__ : Idéal pour le Edge Computing et les infrastructures "bare-metal" minimalistes, offrant une collecte centralisée sans surcharge CPU.

## Télémétrie industrielle & Robotique

- __Description__ : Suivi en temps réel de l'état de santé de systèmes automatisés ou robotisés (cycles de fonctionnement, tensions, température des cartes électroniques).

- __Drivers__ : modbus-tcp, gpiod-monitor

- __Exporters__ : websocket-stream (dashboard web temps réel), stdout-exporter

- __Valeur ajoutée__ : Permet de découpler complètement la logique de scrutation matérielle basse température de la couche d'affichage ou d'analyse.

---

[Précédent](SSH_GUIDE.md)

© 2026 F. Colinet - **GatherPipe** est distribué sous licence **CC BY-NC-SA 4.0**.
