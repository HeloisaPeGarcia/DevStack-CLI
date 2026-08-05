package audit

import (
	"context"
	"devstack/pkg/config"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type AuditStatus string

const (
	StatusInstalled AuditStatus = "INSTALLED"
	StatusMissing   AuditStatus = "MISSING"
	StatusError     AuditStatus = "ERROR"
)

type AuditResult struct {
	Package      config.SystemPackage `json:"package"`
	Status       AuditStatus          `json:"status"`
	VersionFound string               `json:"version_found"`
	PathFound    string               `json:"path_found"`
	ErrMessage   string               `json:"err_message,omitempty"`
}

type Inspector struct{}

func NewInspector() *Inspector {
	return &Inspector{}
}

func (i *Inspector) CheckPackage(pkg config.SystemPackage) AuditResult {
	result := AuditResult{
		Package: pkg,
		Status:  StatusMissing,
	}

	path, err := exec.LookPath(pkg.CheckCmd)
	if err != nil {
		result.Status = StatusMissing
		result.ErrMessage = fmt.Sprintf("Executável '%s' não encontrado no PATH", pkg.CheckCmd)
		return result
	}
	result.PathFound = path

	if len(pkg.CheckArgs) == 0 {
		result.Status = StatusInstalled
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, pkg.CheckCmd, pkg.CheckArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Status = StatusInstalled
		result.VersionFound = "Versão desconhecida"
		return result
	}

	outStr := strings.TrimSpace(string(output))
	lines := strings.Split(outStr, "\n")
	if len(lines) > 0 {
		result.VersionFound = strings.TrimSpace(lines[0])
	} else {
		result.VersionFound = outStr
	}

	result.Status = StatusInstalled
	return result
}

func (i *Inspector) AuditRecipe(recipe config.StackRecipe) []AuditResult {
	results := make([]AuditResult, len(recipe.SystemPackages))
	var wg sync.WaitGroup

	for index, pkg := range recipe.SystemPackages {
		wg.Add(1)
		go func(idx int, p config.SystemPackage) {
			defer wg.Done()
			results[idx] = i.CheckPackage(p)
		}(index, pkg)
	}

	wg.Wait()
	return results
}
