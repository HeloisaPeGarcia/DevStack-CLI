package winget

import (
	"context"
	"devstack/pkg/config"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type InstallCallback func(pkg config.SystemPackage, status string, err error)

type Manager struct {
	DryRun bool
}

func NewManager(dryRun bool) *Manager {
	return &Manager{
		DryRun: dryRun,
	}
}

func (m *Manager) IsWingetAvailable() bool {
	_, err := exec.LookPath("winget")
	return err == nil
}

func (m *Manager) IsAdminMode() bool {
	if runtime.GOOS != "windows" {
		return true
	}
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

func (m *Manager) InstallPackage(ctx context.Context, pkg config.SystemPackage, callback InstallCallback) error {
	if m.DryRun {
		if callback != nil {
			callback(pkg, fmt.Sprintf("[DRY-RUN] Simulação: winget install --id %s --silent --accept-package-agreements --accept-source-agreements", pkg.WingetID), nil)
		}
		return nil
	}

	if !m.IsWingetAvailable() {
		err := fmt.Errorf("winget não está disponível neste sistema Windows")
		if callback != nil {
			callback(pkg, "FALHOU", err)
		}
		return err
	}

	if callback != nil {
		callback(pkg, "Iniciando download e instalação silenciosa...", nil)
	}

	args := []string{
		"install",
		"--id", pkg.WingetID,
		"--silent",
		"--accept-package-agreements",
		"--accept-source-agreements",
	}

	cmd := exec.CommandContext(ctx, "winget", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outStr := strings.TrimSpace(string(output))
		execErr := fmt.Errorf("erro ao instalar %s (%s): %v. Saída: %s", pkg.Name, pkg.WingetID, err, outStr)
		if callback != nil {
			callback(pkg, "FALHOU", execErr)
		}
		return execErr
	}

	if callback != nil {
		callback(pkg, "Instalado com sucesso!", nil)
	}
	return nil
}

func (m *Manager) InstallPackages(packages []config.SystemPackage, callback InstallCallback) map[string]error {
	errorsMap := make(map[string]error)

	for _, pkg := range packages {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		err := m.InstallPackage(ctx, pkg, callback)
		if err != nil {
			errorsMap[pkg.WingetID] = err
		}
		cancel()
	}

	return errorsMap
}
