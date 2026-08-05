# 🛡️ Política de Segurança — DevStack CLI

A equipe do **DevStack** leva a segurança do usuário e a integridade do sistema muito a sério. Esta política detalha nossas práticas de segurança, como reportar vulnerabilidades e as garantias do nosso software.

---

## 📋 Versões Suportadas

Atualmente, fornecemos atualizações de segurança para as seguintes versões do DevStack:

| Versão | Suportada |
| :--- | :--- |
| `v1.0.x` | ✅ Sim |
| `< 1.0.0` | ❌ Não |

---

## 🔒 Controles de Segurança do DevStack

O DevStack CLI executa automações de sistema e instalação de softwares. Para garantir que seu ambiente continue seguro, implementamos os seguintes controles:

### 1. Prevenção Total contra Injeção de Comandos (Subprocess Isolation)
- **Execução Segura:** Nenhuma chamada ao sistema operacional utiliza interpretação de strings em shell (`cmd.exe /c` ou `powershell -c`).
- Todas as execuções de processos utilizam a API nativa do Go `os/exec.CommandContext(name, arg1, arg2...)`, onde cada argumento é estritamente isolado no nível da syscall do sistema operacional.

### 2. Validação Estrita do Winget
Ao instalar pacotes via Winget, o DevStack obriga a execução de verificações do fabricante:
- `--accept-package-agreements`
- `--accept-source-agreements`
- Verificação de hash SHA-256 e assinatura digital (Code Signing Certificate) pelo próprio subsistema do Windows Winget.

### 3. Princípio de Menor Privilégio & UAC
- O DevStack **não exige rodar como Administrador por padrão**. Ele roda com permissões normais de usuário.
- Apenas o subprocesso do Winget solicitará elevação UAC quando um instalador específico (ex: `.msi` de sistema) exigir permissões de escrita em `C:\Program Files`.

---

## 🚨 Relatando uma Vulnerabilidade (Responsible Disclosure)

Se você encontrou uma falha de segurança ou vulnerabilidade no **DevStack CLI**, **NÃO abra uma Issue pública no GitHub**.

Por favor, siga o processo de divulgação responsável:

1. **Email Privado:** Envie os detalhes da vulnerabilidade para `security@devstack.local` (ou abra um *Draft Security Advisory* privado no repositório do GitHub).
2. **Conteúdo do Relatório:**
   - Descrição detalhada da vulnerabilidade.
   - Passo a passo para reproduzir o problema (PoC - Proof of Concept).
   - Versão do DevStack e do Windows testadas.
   - Impacto potencial estimado.
3. **Prazo de Resposta:**
   - Responderemos à sua mensagem em até **48 horas**.
   - Forneceremos uma estimativa de correção em até **7 dias úteis**.
   - Notificaremos você assim que o patch de segurança for lançado para publicação de créditos.

---

## 🔐 Boas Práticas Recomendadas para Usuários

- **Baixe apenas de fontes oficiais:** Utilize as releases oficiais no GitHub ou instale via `go install` / `winget`.
- **Verifique os Hashes:** Sempre verifique o hash SHA-256 dos binários `.exe` baixados contra a lista disponibilizada na página de Releases do GitHub.

---

<div align="center">
  <sub>Obrigado por ajudar a manter o DevStack seguro para toda a comunidade dev!</sub>
</div>
