package ide

import (
	"context"
	"devstack/pkg/config"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Configurator struct {
	ProjectDir string
	DryRun     bool
}

func NewConfigurator(projectDir string, dryRun bool) *Configurator {
	return &Configurator{
		ProjectDir: projectDir,
		DryRun:     dryRun,
	}
}

func (c *Configurator) IsVSCodeInstalled() bool {
	_, err := exec.LookPath("code")
	return err == nil
}

func (c *Configurator) InstallExtension(extID string) error {
	if c.DryRun {
		fmt.Printf("[DRY-RUN] Simulação: code --install-extension %s\n", extID)
		return nil
	}

	if !c.IsVSCodeInstalled() {
		return fmt.Errorf("VS Code ('code') não encontrado no PATH")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "code", "--install-extension", extID, "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao instalar extensão %s: %v (%s)", extID, err, string(output))
	}
	return nil
}

func (c *Configurator) SetupWorkspaceConfigs(recipe config.StackRecipe) error {
	vscodeDir := filepath.Join(c.ProjectDir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta .vscode: %w", err)
	}

	defaultFormatter := "esbenp.prettier-vscode"
	if recipe.ID == "python-fastapi" {
		defaultFormatter = "ms-python.python"
	}

	settings := map[string]interface{}{
		"editor.formatOnSave":     true,
		"editor.defaultFormatter": defaultFormatter,
		"files.exclude": map[string]bool{
			"**/.git":         true,
			"**/.DS_Store":    true,
			"**/node_modules": true,
			"**/__pycache__":  true,
		},
	}

	if recipe.ID == "go-react" {
		settings["go.useLanguageServer"] = true
		settings["go.formatTool"] = "gofmt"
		settings["typescript.tsdk"] = "frontend/node_modules/typescript/lib"
	} else if recipe.ID == "node-vue" {
		settings["typescript.tsdk"] = "frontend/node_modules/typescript/lib"
	}

	settingsBytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar settings.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(vscodeDir, "settings.json"), settingsBytes, 0644); err != nil {
		return fmt.Errorf("erro ao escrever settings.json: %w", err)
	}

	recMap := map[string][]string{
		"recommendations": recipe.VSCodeExtensions,
	}
	recBytes, err := json.MarshalIndent(recMap, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar extensions.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(vscodeDir, "extensions.json"), recBytes, 0644); err != nil {
		return fmt.Errorf("erro ao escrever extensions.json: %w", err)
	}

	var debugConfig map[string]interface{}
	switch recipe.ID {
	case "python-fastapi":
		debugConfig = map[string]interface{}{
			"name":    "Debug FastAPI",
			"type":    "python",
			"request": "launch",
			"module":  "uvicorn",
			"args":    []string{"app.main:app", "--reload"},
		}
	case "node-vue":
		debugConfig = map[string]interface{}{
			"name":    "Debug Node Express",
			"type":    "node",
			"request": "launch",
			"program": "${workspaceFolder}/backend/src/index.ts",
			"runtimeArgs": []string{"-r", "ts-node/register"},
		}
	default:
		debugConfig = map[string]interface{}{
			"name":    "Debug Go Backend",
			"type":    "go",
			"request": "launch",
			"mode":    "auto",
			"program": "${workspaceFolder}/backend/cmd/api",
			"env":     map[string]string{"PORT": "8080"},
		}
	}

	launch := map[string]interface{}{
		"version":        "0.2.0",
		"configurations": []map[string]interface{}{debugConfig},
	}
	launchBytes, err := json.MarshalIndent(launch, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar launch.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(vscodeDir, "launch.json"), launchBytes, 0644); err != nil {
		return fmt.Errorf("erro ao escrever launch.json: %w", err)
	}

	return nil
}
