package config

type SystemPackage struct {
	Name        string   `json:"name"`
	WingetID    string   `json:"winget_id"`
	CheckCmd    string   `json:"check_cmd"`
	CheckArgs   []string `json:"check_args"`
	MinVersion  string   `json:"min_version"`
	Description string   `json:"description"`
}

type StackRecipe struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Aliases          []string        `json:"aliases"`
	Description      string          `json:"description"`
	SystemPackages   []SystemPackage `json:"system_packages"`
	VSCodeExtensions []string        `json:"vscode_extensions"`
	TemplateType     string          `json:"template_type"`
}

func GetPredefinedRecipes() []StackRecipe {
	return []StackRecipe{
		{
			ID:          "go-react",
			Name:        "Go Backend + React Frontend",
			Aliases:     []string{"go-react", "react-go", "fullstack-go"},
			Description: "Backend em Go (REST API) + Frontend React (TypeScript/Vite) + PostgreSQL + Docker",
			SystemPackages: []SystemPackage{
				{
					Name:        "Go SDK",
					WingetID:    "GoLang.Go",
					CheckCmd:    "go",
					CheckArgs:   []string{"version"},
					MinVersion:  "1.22.0",
					Description: "Linguagem de programação Go e ferramentas padrão",
				},
				{
					Name:        "Node.js LTS",
					WingetID:    "NodeJS.NodeJS.LTS",
					CheckCmd:    "node",
					CheckArgs:   []string{"-v"},
					MinVersion:  "20.0.0",
					Description: "Runtime JavaScript para frontend e ferramentas de build",
				},
				{
					Name:        "PNPM Package Manager",
					WingetID:    "PNPM.PNPM",
					CheckCmd:    "pnpm",
					CheckArgs:   []string{"-v"},
					MinVersion:  "9.0.0",
					Description: "Gerenciador de pacotes JS ultra-rápido e eficiente",
				},
				{
					Name:        "Docker Desktop",
					WingetID:    "Docker.DockerDesktop",
					CheckCmd:    "docker",
					CheckArgs:   []string{"--version"},
					MinVersion:  "24.0.0",
					Description: "Containerização para ambiente local e banco de dados",
				},
				{
					Name:        "VS Code",
					WingetID:    "Microsoft.VisualStudioCode",
					CheckCmd:    "code",
					CheckArgs:   []string{"--version"},
					Description: "Editor de código recomendado com suporte a extensões",
				},
				{
					Name:        "Git",
					WingetID:    "Git.Git",
					CheckCmd:    "git",
					CheckArgs:   []string{"--version"},
					Description: "Controle de versão de código",
				},
				{
					Name:        "PostgreSQL 16",
					WingetID:    "PostgreSQL.PostgreSQL.16",
					CheckCmd:    "psql",
					CheckArgs:   []string{"--version"},
					Description: "Banco de dados relacional",
				},
			},
			VSCodeExtensions: []string{
				"golang.go",
				"dbaeumer.vscode-eslint",
				"esbenp.prettier-vscode",
				"humao.rest-client",
				"ms-azuretools.vscode-docker",
				"eamodio.gitlens",
			},
			TemplateType: "go-react-monorepo",
		},
		{
			ID:          "python-fastapi",
			Name:        "Python FastAPI + Streamlit",
			Aliases:     []string{"python-fastapi", "fastapi", "python-ai"},
			Description: "API REST assíncrona em Python com FastAPI + Dashboard Streamlit",
			SystemPackages: []SystemPackage{
				{
					Name:        "Python 3.12",
					WingetID:    "Python.Python.3.12",
					CheckCmd:    "python",
					CheckArgs:   []string{"--version"},
					MinVersion:  "3.11.0",
					Description: "Linguagem Python e pip",
				},
				{
					Name:        "Git",
					WingetID:    "Git.Git",
					CheckCmd:    "git",
					CheckArgs:   []string{"--version"},
					Description: "Controle de versão",
				},
				{
					Name:        "VS Code",
					WingetID:    "Microsoft.VisualStudioCode",
					CheckCmd:    "code",
					CheckArgs:   []string{"--version"},
					Description: "Editor de código",
				},
			},
			VSCodeExtensions: []string{
				"ms-python.python",
				"ms-python.vscode-pylance",
				"njpwerner.autopep8",
			},
			TemplateType: "python-fastapi-app",
		},
		{
			ID:          "node-vue",
			Name:        "Node.js Express + Vue 3",
			Aliases:     []string{"node-vue", "vue-express", "mevn"},
			Description: "Backend Node.js Express em TypeScript + Frontend Vue 3 (Vite)",
			SystemPackages: []SystemPackage{
				{
					Name:        "Node.js LTS",
					WingetID:    "NodeJS.NodeJS.LTS",
					CheckCmd:    "node",
					CheckArgs:   []string{"-v"},
					Description: "Runtime JS",
				},
				{
					Name:        "Git",
					WingetID:    "Git.Git",
					CheckCmd:    "git",
					CheckArgs:   []string{"--version"},
					Description: "Controle de versão",
				},
				{
					Name:        "VS Code",
					WingetID:    "Microsoft.VisualStudioCode",
					CheckCmd:    "code",
					CheckArgs:   []string{"--version"},
					Description: "Editor de código",
				},
			},
			VSCodeExtensions: []string{
				"Vue.volar",
				"dbaeumer.vscode-eslint",
				"esbenp.prettier-vscode",
			},
			TemplateType: "node-vue-app",
		},
	}
}
