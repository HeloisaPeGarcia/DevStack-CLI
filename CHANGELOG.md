# 📜 Changelog

Todas as alterações notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/) e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

---

## [1.0.0] - 2026-08-05

### 🎉 Adicionado
- **Arquitetura CLI:** Refatoração de arquitetura seguindo Go Project Layout (`cmd/devstack`, `internal/bootstrap`, `pkg/`).
- **3 Pilhas Pré-Configuradas:**
  - `go-react`: Go Backend (API HTTP) + React Frontend (Vite/TypeScript) + PostgreSQL + Docker Compose + Makefile.
  - `python-fastapi`: Python FastAPI REST API + Streamlit Dashboard + Docker + requirements.txt.
  - `node-vue`: Node.js Express em TypeScript + Vue 3 Frontend (Vite) + Docker + Makefile.
- **Motor de Auditoria Concorrente:** `pkg/audit` inspeciona ferramentas no `%PATH%` usando goroutines e `context.WithTimeout` de 3s por binário.
- **Automação do Winget:** `pkg/winget` com suporte a instalações silenciosas, detecção de privilégios UAC do Windows e modo simulação (`--dry-run`).
- **Resolução Determinística:** `pkg/resolver` com regras de prioridade explicita para resolver stacks por ID, aliases e palavras-chave.
- **Orquestrador de IDE:** `pkg/ide` para auto-instalação de extensões do VS Code e geração de `.vscode/settings.json`, `launch.json` e `extensions.json`.
- **Suporte a `NO_COLOR`:** Suporte à especificação `NO_COLOR` e abstração `io.Writer` em `pkg/ui`.
- **Saída Estruturada JSON:** Flag `--output json` no comando `audit` para integração com pipelines de CI/CD.
- **Automação CI/CD:** GitHub Actions em `.github/workflows/ci.yml` para testes e compilação automatizada em Windows.
- **Documentação & Comunidade:** Adicionados `LICENSE` (MIT), `SECURITY.md`, `CONTRIBUTING.md`, `docs/WIKI.md` e `README.md` com banner dinâmico Typing SVG.

### 🛡️ Segurança & Correções
- Prevenção total contra Command Injection usando isolamento de argumentos `os/exec.CommandContext`.
- Tratamento explícito de erros de E/S ao gerar configurações de IDE e workspace.
- Exit codes padronizados (`ExitOK`, `ExitError`, `ExitConfigError`, `ExitInstallError`).
