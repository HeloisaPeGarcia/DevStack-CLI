package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

type Generator struct {
	ProjectName string
	OutputDir   string
}

func NewGenerator(projectName, outputDir string) *Generator {
	if outputDir == "" {
		outputDir = "."
	}
	return &Generator{
		ProjectName: projectName,
		OutputDir:   filepath.Join(outputDir, projectName),
	}
}

func (g *Generator) GenerateGoReactMonorepo() error {
	dirs := []string{
		filepath.Join(g.OutputDir, "backend", "cmd", "api"),
		filepath.Join(g.OutputDir, "backend", "internal", "config"),
		filepath.Join(g.OutputDir, "backend", "internal", "handler"),
		filepath.Join(g.OutputDir, "frontend", "src"),
		filepath.Join(g.OutputDir, "frontend", "public"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("falha ao criar diretório %s: %w", d, err)
		}
	}

	files := map[string]string{
		filepath.Join(g.OutputDir, "backend", "go.mod"): fmt.Sprintf(`module %s-backend

go 1.22
`, g.ProjectName),

		filepath.Join(g.OutputDir, "backend", "Dockerfile"): `FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
`,

		filepath.Join(g.OutputDir, "frontend", "package.json"): fmt.Sprintf(`{
  "name": "%s-frontend",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.1",
    "typescript": "^5.4.5",
    "vite": "^5.2.11"
  }
}
`, g.ProjectName),

		filepath.Join(g.OutputDir, "frontend", "vite.config.ts"): `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
`,

		filepath.Join(g.OutputDir, "frontend", "src", "main.tsx"): `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
`,

		filepath.Join(g.OutputDir, "frontend", "src", "App.tsx"): fmt.Sprintf(`import { useState, useEffect } from 'react'

export default function App() {
  const [health, setHealth] = useState<string>('Carregando...')

  useEffect(() => {
    fetch('/api/health')
      .then(res => res.json())
      .then(data => setHealth(data.status + ' (' + data.service + ')'))
      .catch(() => setHealth('Off-line'))
  }, [])

  return (
    <div style={{ fontFamily: 'system-ui, sans-serif', padding: '2rem', maxWidth: '800px', margin: '0 auto' }}>
      <h1>🚀 App %s</h1>
      <p>Gerado pelo <strong>devstack</strong> - AI-Powered Dev-Stack Installer</p>
      
      <div style={{ background: '#f4f4f5', padding: '1rem', borderRadius: '8px', marginTop: '1rem' }}>
        <h3>Status do Backend Go:</h3>
        <p><strong>%s</strong></p>
      </div>
    </div>
  )
}
`, g.ProjectName, "{health}"),

		filepath.Join(g.OutputDir, "frontend", "public", "index.html"): fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-BR">
  <head>
    <meta charset="UTF-8" />
    <title>%s</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`, g.ProjectName),

		filepath.Join(g.OutputDir, "docker-compose.yml"): fmt.Sprintf(`version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: %s_postgres
    environment:
      POSTGRES_USER: devuser
      POSTGRES_PASSWORD: devpassword
      POSTGRES_DB: %s_db
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  backend:
    build: ./backend
    container_name: %s_backend
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DATABASE_URL=postgres://devuser:devpassword@postgres:5432/%s_db?sslmode=disable
    depends_on:
      - postgres

volumes:
  pgdata:
`, g.ProjectName, g.ProjectName, g.ProjectName, g.ProjectName),

		filepath.Join(g.OutputDir, "Makefile"): `dev-backend:
	cd backend && go run ./cmd/api

dev-frontend:
	cd frontend && pnpm install && pnpm dev

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
`,

		filepath.Join(g.OutputDir, "README.md"): fmt.Sprintf(`# %s

Projeto gerado automaticamente pelo **devstack CLI**.
`, g.ProjectName),
	}

	mainTmpl, err := templateFS.ReadFile("templates/go-react/backend/cmd/api/main.go.tmpl")
	if err == nil {
		content := strings.ReplaceAll(string(mainTmpl), "{{.ProjectName}}", g.ProjectName)
		files[filepath.Join(g.OutputDir, "backend", "cmd", "api", "main.go")] = content
	}

	tsTmpl, err := templateFS.ReadFile("templates/go-react/frontend/tsconfig.json.tmpl")
	if err == nil {
		files[filepath.Join(g.OutputDir, "frontend", "tsconfig.json")] = string(tsTmpl)
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("falha ao criar arquivo %s: %w", path, err)
		}
	}

	return nil
}

func (g *Generator) GeneratePythonFastAPIApp() error {
	dirs := []string{
		filepath.Join(g.OutputDir, "app"),
		filepath.Join(g.OutputDir, "dashboard"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("falha ao criar diretório %s: %w", d, err)
		}
	}

	files := map[string]string{
		filepath.Join(g.OutputDir, "Dockerfile"): `FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8000 8501
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
`,
		filepath.Join(g.OutputDir, "README.md"): fmt.Sprintf(`# %s (Python FastAPI + Streamlit)

Projeto gerado pelo **DevStack CLI**.

## 🚀 Executar Localmente
`+"```bash"+`
pip install -r requirements.txt
uvicorn app.main:app --reload
streamlit run dashboard/app.py
`+"```"+`
`, g.ProjectName),
	}

	mainPy, err := templateFS.ReadFile("templates/python-fastapi/app/main.py.tmpl")
	if err == nil {
		content := strings.ReplaceAll(string(mainPy), "{{.ProjectName}}", g.ProjectName)
		files[filepath.Join(g.OutputDir, "app", "main.py")] = content
	}

	dashPy, err := templateFS.ReadFile("templates/python-fastapi/dashboard/app.py.tmpl")
	if err == nil {
		content := strings.ReplaceAll(string(dashPy), "{{.ProjectName}}", g.ProjectName)
		files[filepath.Join(g.OutputDir, "dashboard", "app.py")] = content
	}

	reqTxt, err := templateFS.ReadFile("templates/python-fastapi/requirements.txt.tmpl")
	if err == nil {
		files[filepath.Join(g.OutputDir, "requirements.txt")] = string(reqTxt)
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("falha ao criar arquivo %s: %w", path, err)
		}
	}

	return nil
}

func (g *Generator) GenerateNodeVueApp() error {
	dirs := []string{
		filepath.Join(g.OutputDir, "backend", "src"),
		filepath.Join(g.OutputDir, "frontend", "src"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("falha ao criar diretório %s: %w", d, err)
		}
	}

	files := map[string]string{
		filepath.Join(g.OutputDir, "backend", "package.json"): fmt.Sprintf(`{
  "name": "%s-backend",
  "version": "0.1.0",
  "main": "dist/index.js",
  "scripts": {
    "dev": "ts-node-dev --respawn src/index.ts",
    "build": "tsc"
  },
  "dependencies": {
    "express": "^4.19.2",
    "cors": "^2.8.5"
  },
  "devDependencies": {
    "@types/express": "^4.17.21",
    "@types/cors": "^2.8.17",
    "typescript": "^5.4.5",
    "ts-node-dev": "^2.0.0"
  }
}
`, g.ProjectName),
		filepath.Join(g.OutputDir, "frontend", "package.json"): fmt.Sprintf(`{
  "name": "%s-frontend",
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build"
  },
  "dependencies": {
    "vue": "^3.4.27"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.4",
    "typescript": "^5.4.5",
    "vite": "^5.2.11",
    "vue-tsc": "^2.0.19"
  }
}
`, g.ProjectName),
	}

	idxTs, err := templateFS.ReadFile("templates/node-vue/backend/src/index.ts.tmpl")
	if err == nil {
		content := strings.ReplaceAll(string(idxTs), "{{.ProjectName}}", g.ProjectName)
		files[filepath.Join(g.OutputDir, "backend", "src", "index.ts")] = content
	}

	appVue, err := templateFS.ReadFile("templates/node-vue/frontend/src/App.vue.tmpl")
	if err == nil {
		content := strings.ReplaceAll(string(appVue), "{{.ProjectName}}", g.ProjectName)
		files[filepath.Join(g.OutputDir, "frontend", "src", "App.vue")] = content
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("falha ao criar arquivo %s: %w", path, err)
		}
	}

	return nil
}
