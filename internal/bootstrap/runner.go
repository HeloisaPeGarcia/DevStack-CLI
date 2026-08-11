package bootstrap

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"devstack/pkg/audit"
	"devstack/pkg/config"
	"devstack/pkg/ide"
	"devstack/pkg/resolver"
	"devstack/pkg/scaffold"
	"devstack/pkg/ui"
	"devstack/pkg/winget"
)

const (
	ExitOK           = 0
	ExitError        = 1
	ExitConfigError  = 2
	ExitInstallError = 3
)

type Runner struct {
	Version   string
	BuildDate string
}

func NewRunner(version, buildDate string) *Runner {
	return &Runner{
		Version:   version,
		BuildDate: buildDate,
	}
}

func (r *Runner) Run(args []string) int {
	ui.PrintBanner()

	if len(args) < 2 {
		r.PrintUsage()
		return ExitOK
	}

	command := args[1]

	switch command {
	case "bootstrap":
		return r.RunBootstrap(args[2:])
	case "audit":
		return r.RunAudit(args[2:])
	case "list-stacks":
		return r.RunListStacks()
	case "version", "--version", "-v":
		fmt.Printf("devstack %s (compilado em %s)\n", r.Version, r.BuildDate)
		return ExitOK
	case "help", "-h", "--help":
		r.PrintUsage()
		return ExitOK
	default:
		ui.PrintError(fmt.Sprintf("Comando desconhecido: '%s'", command))
		r.PrintUsage()
		return ExitError
	}
}

func (r *Runner) PrintUsage() {
	fmt.Println("Uso:")
	fmt.Println("  devstack bootstrap [flags]    Provisiona uma stack completa de desenvolvimento")
	fmt.Println("  devstack audit [flags]        Audita dependências instaladas no sistema")
	fmt.Println("  devstack list-stacks          Lista todas as pilhas pré-configuradas")
	fmt.Println("  devstack version              Exibe a versão do CLI")
	fmt.Println("\nFlags para 'bootstrap':")
	fmt.Println("  --stack string                Nome da stack ou descrição (ex: 'Go Backend + React Frontend')")
	fmt.Println("  --project-name string         Nome do projeto a ser criado (padrão: 'my-dev-app')")
	fmt.Println("  --output-dir string           Diretório de saída (padrão: diretório atual)")
	fmt.Println("  --dry-run                     Modo simulação (não altera o sistema nem instala pacotes)")
	fmt.Println("  --skip-install                Pula a instalação dos pacotes via Winget")
	fmt.Println("  --skip-scaffold               Pula a geração de templates de código do projeto")
	fmt.Println("  --output string               Formato de saída: 'text' ou 'json' (padrão: 'text')")
}

func (r *Runner) RunListStacks() int {
	ui.PrintHeader("Pilhas de Desenvolvimento Pré-Configuradas")
	recipes := config.GetPredefinedRecipes()
	for _, rec := range recipes {
		lines := []string{
			fmt.Sprintf("ID: %s", rec.ID),
			fmt.Sprintf("Aliases: %v", rec.Aliases),
			fmt.Sprintf("Descrição: %s", rec.Description),
			fmt.Sprintf("Pacotes Winget (%d):", len(rec.SystemPackages)),
		}
		for _, p := range rec.SystemPackages {
			lines = append(lines, fmt.Sprintf("  • %s (%s)", p.Name, p.WingetID))
		}
		ui.PrintBox(rec.Name, lines)
		fmt.Println()
	}
	return ExitOK
}

type AuditJSONOutput struct {
	Stack          string              `json:"stack"`
	InstalledCount int                 `json:"installed_count"`
	TotalCount     int                 `json:"total_count"`
	Results        []audit.AuditResult `json:"results"`
}

func (r *Runner) RunAudit(args []string) int {
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	stackFlag := fs.String("stack", "go-react", "Nome da stack a ser auditada")
	outputFlag := fs.String("output", "text", "Formato de saída: text ou json")
	_ = fs.Parse(args)

	res := resolver.NewResolver()
	recipe, exactMatch := res.ResolveStack(*stackFlag)
	if recipe == nil {
		ui.PrintError(fmt.Sprintf("Stack '%s' não encontrada.", *stackFlag))
		return ExitConfigError
	}
	if !exactMatch {
		ui.PrintWarning(fmt.Sprintf("Stack '%s' não encontrada exatamente. Auditando receita padrão '%s'.", *stackFlag, recipe.Name))
	}

	inspector := audit.NewInspector()
	results := inspector.AuditRecipe(*recipe)

	installedCount := 0
	for _, res := range results {
		if res.Status == audit.StatusInstalled {
			installedCount++
		}
	}

	if *outputFlag == "json" {
		jsonOut := AuditJSONOutput{
			Stack:          recipe.Name,
			InstalledCount: installedCount,
			TotalCount:     len(results),
			Results:        results,
		}
		b, _ := json.MarshalIndent(jsonOut, "", "  ")
		fmt.Println(string(b))
		return ExitOK
	}

	ui.PrintHeader(fmt.Sprintf("Auditoria do Ambiente: %s", recipe.Name))
	for _, res := range results {
		if res.Status == audit.StatusInstalled {
			ui.PrintStatusBadge(res.Package.Name, res.VersionFound, true)
		} else {
			ui.PrintStatusBadge(res.Package.Name, res.ErrMessage, false)
		}
	}

	fmt.Printf("\nResultado: %d de %d dependências encontradas no sistema.\n", installedCount, len(results))
	return ExitOK
}

func (r *Runner) RunBootstrap(args []string) int {
	fs := flag.NewFlagSet("bootstrap", flag.ExitOnError)
	stackFlag := fs.String("stack", "Go Backend + React Frontend", "Nome da stack a ser provisionada")
	projectNameFlag := fs.String("project-name", "my-dev-app", "Nome do projeto a ser criado")
	outputDirFlag := fs.String("output-dir", ".", "Diretório raiz para criação do projeto")
	dryRunFlag := fs.Bool("dry-run", false, "Simular sem realizar instalações ou criar arquivos")
	skipInstallFlag := fs.Bool("skip-install", false, "Pular etapa de instalação pelo Winget")
	skipScaffoldFlag := fs.Bool("skip-scaffold", false, "Pular criação de arquivos do projeto")
	_ = fs.Parse(args)

	res := resolver.NewResolver()
	recipe, exactMatch := res.ResolveStack(*stackFlag)
	if recipe == nil {
		ui.PrintError(fmt.Sprintf("Não foi possível resolver a stack '%s'", *stackFlag))
		return ExitConfigError
	}
	if !exactMatch {
		ui.PrintWarning(fmt.Sprintf("Stack '%s' não encontrada exatamente. Provisionando receita padrão '%s'.", *stackFlag, recipe.Name))
	}

	ui.PrintHeader(fmt.Sprintf("Iniciando Bootstrap: %s", recipe.Name))

	ui.PrintInfo("Fase 1/4: Auditando ferramentas no sistema...")
	inspector := audit.NewInspector()
	auditResults := inspector.AuditRecipe(*recipe)

	var missingPackages []config.SystemPackage
	for _, res := range auditResults {
		if res.Status != audit.StatusInstalled {
			missingPackages = append(missingPackages, res.Package)
			ui.PrintWarning(fmt.Sprintf("Faltando: %s (%s)", res.Package.Name, res.Package.WingetID))
		} else {
			ui.PrintSuccess(fmt.Sprintf("Já instalado: %s (%s)", res.Package.Name, res.VersionFound))
		}
	}

	if !*skipInstallFlag && len(missingPackages) > 0 {
		ui.PrintHeader("Fase 2/4: Instalação de Pacotes via Winget")
		mgr := winget.NewManager(*dryRunFlag)

		if !mgr.IsAdminMode() {
			ui.PrintWarning("Atenção: O processo não está rodando como Administrador. Algumas instalações do Winget podem solicitar permissão UAC.")
		}

		for _, pkg := range missingPackages {
			ui.PrintInfo(fmt.Sprintf("Instalando %s...", pkg.Name))
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			err := mgr.InstallPackage(ctx, pkg, func(p config.SystemPackage, status string, err error) {
				if err != nil {
					ui.PrintError(fmt.Sprintf("Status %s: %v", status, err))
				} else {
					ui.PrintSuccess(fmt.Sprintf("Status: %s", status))
				}
			})
			cancel()
			if err != nil {
				ui.PrintError(fmt.Sprintf("Erro ao instalar %s: %v", pkg.Name, err))
			}
		}
	} else if len(missingPackages) == 0 {
		ui.PrintSuccess("Fase 2/4: Todas as dependências já estão instaladas no SO!")
	}

	projectPath := filepath.Join(*outputDirFlag, *projectNameFlag)
	if !*skipScaffoldFlag {
		ui.PrintHeader(fmt.Sprintf("Fase 3/4: Criando Estrutura de Código em '%s'", projectPath))

		if *dryRunFlag {
			ui.PrintInfo(fmt.Sprintf("[DRY-RUN] Simulação: Criando arquivos do projeto em %s", projectPath))
		} else {
			gen := scaffold.NewGenerator(*projectNameFlag, *outputDirFlag)
			var err error
			switch recipe.TemplateType {
			case "go-react-monorepo":
				err = gen.GenerateGoReactMonorepo()
			case "python-fastapi-app":
				err = gen.GeneratePythonFastAPIApp()
			case "node-vue-app":
				err = gen.GenerateNodeVueApp()
			default:
				err = gen.GenerateGoReactMonorepo()
			}

			if err != nil {
				ui.PrintError(fmt.Sprintf("Erro ao gerar projeto: %v", err))
				return ExitInstallError
			} else {
				ui.PrintSuccess(fmt.Sprintf("Estrutura do projeto '%s' criada com sucesso!", recipe.Name))
			}
		}
	}

	ui.PrintHeader("Fase 4/4: Configurando Ambiente de IDE (VS Code)")
	ideConfig := ide.NewConfigurator(projectPath, *dryRunFlag)
	if !*dryRunFlag {
		if err := ideConfig.SetupWorkspaceConfigs(*recipe); err != nil {
			ui.PrintError(fmt.Sprintf("Erro ao configurar VS Code: %v", err))
		} else {
			ui.PrintSuccess("Configurações `.vscode/settings.json`, `launch.json` e `extensions.json` criadas!")
		}
	} else {
		ui.PrintInfo("[DRY-RUN] Configurações de IDE seriam criadas em .vscode/")
	}

	ui.PrintHeader("🎉 BOOTSTRAP CONCLUÍDO COM SUCESSO!")
	boxLines := []string{
		fmt.Sprintf("Projeto: %s", *projectNameFlag),
		fmt.Sprintf("Caminho: %s", projectPath),
		fmt.Sprintf("Stack:   %s", recipe.Name),
		"",
		"Para iniciar o desenvolvimento:",
		fmt.Sprintf("  cd %s", *projectNameFlag),
		"  code .             # Abrir no VS Code",
	}
	ui.PrintBox("Próximos Passos", boxLines)
	return ExitOK
}
