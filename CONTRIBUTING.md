# 🤝 Guia de Contribuição — DevStack

Obrigado pelo seu interesse em contribuir com o **DevStack**! Este projeto é mantido pela comunidade e ficamos felizes em receber contribuições de todos os níveis.

---

## 🛠️ Como Começar

### Pré-requisitos
- [Go 1.22 ou superior](https://go.dev/dl/)
- [Git](https://git-scm.com/)
- Windows 10/11 (para testes do Winget)

### Configuração do Ambiente de Desenvolvimento

```bash
# 1. Faça o fork do repositório no GitHub
# 2. Clone o seu fork
git clone https://github.com/HeloisaPeGarcia/DevStack-CLI.git
cd DevStack-CLI

# 3. Rode os testes unitários
make test

# 4. Compile o binário localmente
make build
```

---

## ➕ Como Adicionar uma Nova Stack

Para adicionar uma nova receita de stack pré-configurada:

1. Abra o arquivo [`pkg/config/recipes.go`](pkg/config/recipes.go).
2. Adicione uma nova struct `StackRecipe` no slice retornado por `GetPredefinedRecipes()`.
3. Adicione o mapeamento de palavras-chave correspondente em [`pkg/resolver/resolver.go`](pkg/resolver/resolver.go).
4. Adicione um novo template sob a pasta [`pkg/scaffold/templates/`](pkg/scaffold/templates/).
5. Adicione um teste em [`pkg/resolver/resolver_test.go`](pkg/resolver/resolver_test.go) para garantir que a stack é resolvida corretamente.

---

## 🧪 Rodando os Testes e Linter

Antes de submeter um Pull Request, certifique-se de que todos os testes passam sem erros e que a taxa de concorrência é válida:

```bash
# Rodar testes unitários com race detector
go test ./... -v -race -count=1

# Rodar linter (se tiver golangci-lint instalado)
golangci-lint run ./...
```

---

## 📋 Padrão de Commits

Seguimos o padrão de **Conventional Commits**:

- `feat: adiciona receita para Elixir Phoenix`
- `fix: corrige timeout no inspector de dependências`
- `docs: atualiza guia de instalação no README`
- `test: adiciona teste unitário para resolver`

---

## 📄 Licença

Ao contribuir para o DevStack, você concorda que suas contribuições serão licenciadas sob a licença MIT do projeto.
