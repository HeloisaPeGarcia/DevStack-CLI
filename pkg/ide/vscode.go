package ide

import (
	"context"
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

func (c *Configurator) SetupWorkspaceConfigs(extensions []string) error {
	vscodeDir := filepath.Join(c.ProjectDir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta .vscode: %w", err)
	}

	settings := map[string]interface{}{
		"editor.formatOnSave":     true,
		"editor.defaultFormatter": "esbenp.prettier-vscode",
		"go.useLanguageServer":    true,
		"go.formatTool":           "gofmt",
		"typescript.tsdk":         "node_modules/typescript/lib",
		"files.exclude": map[string]bool{
			"**/.git":         true,
			"**/.DS_Store":    true,
			"**/node_modules": true,
		},
	}
	settingsBytes, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar settings.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(vscodeDir, "settings.json"), settingsBytes, 0644); err != nil {
		return fmt.Errorf("erro ao escrever settings.json: %w", err)
	}

	recMap := map[string][]string{
		"recommendations": extensions,
	}
	recBytes, err := json.MarshalIndent(recMap, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar extensions.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(vscodeDir, "extensions.json"), recBytes, 0644); err != nil {
		return fmt.Errorf("erro ao escrever extensions.json: %w", err)
	}

	launch := map[string]interface{}{
		"version": "0.2.0",
		"configurations": []map[string]interface{}{
			{
				"name":    "Debug Go Backend",
				"type":    "go",
				"request": "launch",
				"mode":    "auto",
				"program": "${workspaceFolder}/backend/cmd/api",
				"env":     map[string]string{"PORT": "8080"},
			},
		},
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
