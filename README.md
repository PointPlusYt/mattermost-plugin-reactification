# Reactification - Plugin Mattermost

Envoie une notification par message direct (DM) lorsqu'une reaction emoji est ajoutee a l'un de vos messages.

## Fonctionnement

- **Active par defaut** pour tous les utilisateurs des l'installation du plugin.
- La notification inclut :
  - L'emoji utilise
  - Le nom de l'utilisateur ayant reagi
  - Un lien vers le message original
  - Un apercu du message (complet si court, tronque a 80 caracteres sinon)
- Aucune notification n'est envoyee si vous reagissez a votre propre message.

## Commandes

| Commande | Action |
|---|---|
| `/reactification on` | Active les notifications (defaut) |
| `/reactification off` | Desactive les notifications |

## Installation

### Prerequis
- Go 1.21+
- Acces administrateur Mattermost
- Mattermost Server 6.0+

### Compilation et installation

```bash
# Cloner / extraire le projet
cd reactification

# Compiler et creer le bundle
make

# Le fichier com.example.reactification-1.0.0.zip est genere
# L'uploader dans : Administration Mattermost > Plugins > Installer un plugin
```

### Installation manuelle (sans compilation)
1. Compiler pour votre plateforme cible avec `go build`
2. Placer le binaire dans `server/dist/`
3. Zipper avec `plugin.json`
4. Uploader dans l'interface d'administration Mattermost

## Structure du projet

```
reactification/
├── plugin.json          # Manifeste du plugin
├── Makefile             # Build & packaging
├── README.md
└── server/
    ├── go.mod
    ├── main.go          # Point d'entree
    ├── plugin.go        # Logique principale + hooks + commandes
    └── store.go         # Persistance KV Store
```

## Exemple de notification recue

```
@marcel a réagi :+1: à "oui je sais bien mais c'est possible" dans **général** 
```
