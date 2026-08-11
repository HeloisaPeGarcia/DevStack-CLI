package scaffold_test

import (
	"devstack/pkg/scaffold"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerator_GenerateGoReactMonorepo(t *testing.T) {
	tmpDir := t.TempDir()
	gen := scaffold.NewGenerator("test-go-app", tmpDir)

	err := gen.GenerateGoReactMonorepo()
	if err != nil {
		t.Fatalf("esperava sucesso ao gerar go-react, obteve erro: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tmpDir, "test-go-app", "backend", "cmd", "api", "main.go"),
		filepath.Join(tmpDir, "test-go-app", "backend", "go.mod"),
		filepath.Join(tmpDir, "test-go-app", "frontend", "package.json"),
		filepath.Join(tmpDir, "test-go-app", "frontend", "src", "App.tsx"),
		filepath.Join(tmpDir, "test-go-app", "docker-compose.yml"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("esperava que o arquivo '%s' existisse, mas não foi encontrado", f)
		}
	}
}

func TestGenerator_GeneratePythonFastAPIApp(t *testing.T) {
	tmpDir := t.TempDir()
	gen := scaffold.NewGenerator("test-python-app", tmpDir)

	err := gen.GeneratePythonFastAPIApp()
	if err != nil {
		t.Fatalf("esperava sucesso ao gerar python-fastapi, obteve erro: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tmpDir, "test-python-app", "app", "main.py"),
		filepath.Join(tmpDir, "test-python-app", "dashboard", "app.py"),
		filepath.Join(tmpDir, "test-python-app", "requirements.txt"),
		filepath.Join(tmpDir, "test-python-app", "docker-compose.yml"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("esperava que o arquivo '%s' existisse, mas não foi encontrado", f)
		}
	}
}

func TestGenerator_GenerateNodeVueApp(t *testing.T) {
	tmpDir := t.TempDir()
	gen := scaffold.NewGenerator("test-node-app", tmpDir)

	err := gen.GenerateNodeVueApp()
	if err != nil {
		t.Fatalf("esperava sucesso ao gerar node-vue, obteve erro: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(tmpDir, "test-node-app", "backend", "src", "index.ts"),
		filepath.Join(tmpDir, "test-node-app", "frontend", "src", "App.vue"),
		filepath.Join(tmpDir, "test-node-app", "frontend", "index.html"),
		filepath.Join(tmpDir, "test-node-app", "frontend", "vite.config.ts"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("esperava que o arquivo '%s' existisse, mas não foi encontrado", f)
		}
	}
}
