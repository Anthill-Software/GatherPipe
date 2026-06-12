package core

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//go:embed github-markdown.min.css
var markdownStyle string

// Le handler qui convertit le markdown à la volée
func MakeDocHandler(docsDir string) http.HandlerFunc {

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)

	return func(w http.ResponseWriter, r *http.Request) {

		filePath := filepath.Join(docsDir, r.URL.Path)

		// Gestion de la page d'accueil par défaut
		if r.URL.Path == "/" || r.URL.Path == "" {
			filePath = filepath.Join(docsDir, "README.md")
		}

		mdContent, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Document non trouvé", http.StatusNotFound)
			return
		}

		var buf bytes.Buffer
		if err := md.Convert(mdContent, &buf); err != nil {
			http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html lang="fr">
			<head>
				<meta charset="UTF-8">
				<title>GatherPipe - Documentation</title>
				<style>
					/* Le CSS de GitHub embarqué nativement */
					%s
					
					/* Tes surcharges de style locales */
					body { box-sizing: border-box; min-width: 200px; max-width: 980px; margin: 0 auto; padding: 45px; }
				</style>
			</head>
			<body class="markdown-body">%s</body>
			</html>
		`, markdownStyle, buf.String())
	}
}
