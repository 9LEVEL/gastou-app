<div align="center">
  <img src="frontend/public/logo.svg" alt="Gastou.app" width="240" />
  <p><strong>Controle inteligente de gastos no supermercado</strong></p>
  <p>
    <a href="#funcionalidades">Funcionalidades</a> •
    <a href="#instalacao">Instalação</a> •
    <a href="#api">API</a> •
    <a href="#contribuindo">Contribuindo</a>
  </p>
  <p>
    <img src="https://img.shields.io/badge/license-MIT-green" alt="MIT License" />
    <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go" alt="Go" />
    <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs" alt="Vue 3" />
    <img src="https://img.shields.io/badge/Docker-ready-2496ED?logo=docker" alt="Docker" />
  </p>
</div>

---

**Gastou.app** é open source e self-hosted: ajuda famílias brasileiras a planejar compras de supermercado, registrar gastos e manter o orçamento sob controle. Funciona no celular como PWA, com tema escuro e modo offline.

## Por que usar?

- **Acabou o mês e o dinheiro também?** Planeje o mês inteiro com orçamento por categoria antes de ir ao mercado
- **Preços que mudam toda semana?** Registre o preço real de cada compra e acompanhe a evolução ao longo do tempo
- **Família desorganizada na hora das compras?** Compartilhe a lista, marque itens em tempo real e finalize direto do carrinho

## Funcionalidades

### Lista de Compras
- Lista mensal pré-preenchida com produtos e preços de referência
- Itens agrupados por categoria (Proteínas, Hortifruti, Limpeza, etc.)
- Categorias colapsáveis com estado salvo no navegador
- Busca de produto por texto (sem precisar rolar)
- Preço de referência auto-preenchido ao selecionar produto
- Filtros rápidos: **Todos / Faltam / Comprados**
- Marcar, adicionar, editar e remover itens
- Suporte a itens que rendem vários meses (ex: frango congelado por 2 meses)

### Compras
- Registrar cada ida ao mercado com local, data e itens comprados
- Vincular compra a uma lista mensal
- Total automático calculado dos itens
- Finalizar compra direto da lista (carrinho)
- Comparar total da NF-e com o total calculado

### Dashboard
- **Resumo**: planejado vs gasto vs saldo restante
- **Progresso**: barra visual de itens comprados
- **Estratégia semanal**: quanto gastar por semana até o fim do mês
- **Planejado vs Real**: gráfico de barras por categoria
- **Evolução**: histórico de gastos mensais

### Configurações
- Editar renda mensal
- CRUD de categorias com cores
- CRUD de produtos com preço de referência
- Gerenciar listas (criar, copiar de anterior, ativar, excluir)
- Seções colapsáveis para navegação rápida

### PWA e Offline
- Tema claro/escuro com detecção automática do sistema
- PWA instalável no celular
- Funciona offline com cache de assets e dados de API
- Banner de status de conexão

## Stack

| Camada | Tecnologia |
|--------|------------|
| Frontend | Vue 3 + TypeScript + Vite |
| Backend | Go 1.26 + Chi router |
| Banco | SQLite (modernc.org/sqlite, pure Go) |
| Deploy | Docker multi-stage (container único) |
| PWA | vite-plugin-pwa + Workbox |

## Instalação

### Requisitos

- Docker + Docker Compose

### Subir em 1 comando

```bash
git clone https://github.com/9LEVEL/gastou-app.git
cd gastou-app
docker compose up --build -d
```

Acesse: http://localhost:8081

O comando acima executa automaticamente:
- Build do frontend Vue com Vite
- Build do backend Go
- Criação do banco SQLite em `./data/`
- Execução das migrations
- Seed com produtos e preços de referência

### Instalar no celular (PWA)

1. Abra `http://SEU-IP:8081` no Chrome do celular (mesma rede Wi-Fi)
2. Toque nos 3 pontinhos > **Adicionar a tela inicial**
3. O app fica com ícone próprio e abre em tela cheia

Para descobrir o IP da máquina:
```bash
# Linux / macOS
ip addr show | grep "inet "

# Windows
ipconfig
```

### Desenvolvimento local

**Backend:**
```bash
cd backend
go mod tidy
go run ./cmd/server
# Roda em http://localhost:3001
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# Roda em http://localhost:5173 com proxy para :3001
```

**Comandos úteis:**
```bash
# Parar
docker compose down

# Ver logs
docker compose logs -f app

# Rebuild após alterar código
docker compose up --build -d

# Limpar tudo (apaga o banco)
docker compose down -v && rm -f data/mercado.db
```

## API

Base URL: `/api`

<details>
<summary>Ver todos os endpoints (31 rotas)</summary>

### Categorias

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/categorias` | Listar categorias |
| POST | `/api/categorias` | Criar categoria |
| PUT | `/api/categorias/:id` | Editar categoria |
| DELETE | `/api/categorias/:id` | Remover categoria |

### Produtos

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/produtos` | Listar produtos (`?categoria_id=X&ativo=true`) |
| POST | `/api/produtos` | Criar produto |
| PUT | `/api/produtos/:id` | Editar produto |
| GET | `/api/produtos/:id/precos` | Histórico de preços do produto |

### Listas

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/listas` | Listar listas |
| POST | `/api/listas` | Criar lista (`?copiar_de=ID` para copiar itens) |
| PUT | `/api/listas/:id` | Editar lista |
| DELETE | `/api/listas/:id` | Remover lista |

### Itens da Lista

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/listas/:id/itens` | Itens da lista |
| POST | `/api/listas/:id/itens` | Adicionar item |
| PUT | `/api/listas/:id/itens/:itemId` | Editar item |
| PATCH | `/api/listas/:id/itens/:itemId/check` | Toggle comprado |
| DELETE | `/api/listas/:id/itens/:itemId` | Remover item |

### Compras

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/compras` | Listar compras (`?lista_id=X`) |
| POST | `/api/compras` | Registrar compra |
| PUT | `/api/compras/:id` | Editar compra |
| DELETE | `/api/compras/:id` | Remover compra |
| POST | `/api/compras/:id/itens` | Adicionar item a compra |
| PUT | `/api/compras/:id/itens/:itemId` | Editar item da compra |
| DELETE | `/api/compras/:id/itens/:itemId` | Remover item da compra |

### Dashboard

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/listas/:id/dashboard/resumo` | Resumo financeiro (planejado vs gasto) |
| GET | `/api/listas/:id/dashboard/comparativo` | Planejado vs Real por categoria |
| GET | `/api/dashboard/evolucao` | Evolução mensal de gastos |

</details>

## Estrutura

```
gastou/
├── docker-compose.yml              # Self-hosted (SQLite)
├── Makefile                        # Atalhos: make up, make down, make logs
├── backend/                        # Go + Chi
│   ├── Dockerfile                  # Multi-stage: Vue build + Go build + Alpine runtime
│   ├── cmd/server/                 # Entry point
│   ├── internal/
│   │   ├── handlers/               # HTTP handlers (categorias, produtos, listas, compras, dashboard)
│   │   ├── repository/             # Camada de dados (SQLite)
│   │   ├── models/                 # Structs e tipos
│   │   └── middleware/             # Logger e Recovery
│   ├── migrations/                 # DDL do banco
│   └── seed/                       # Dados iniciais
├── frontend/                       # Vue 3 + TypeScript
│   ├── src/
│   │   ├── views/                  # 4 telas: Lista, Compras, Dashboard, Config
│   │   ├── components/             # Componentes reutilizáveis
│   │   ├── composables/            # Estado reativo: useLista, useCompras, useTheme, useOffline
│   │   ├── api/client.ts           # HTTP client com cache offline
│   │   └── styles/main.css         # Design system (light + dark via CSS variables)
│   └── public/
│       ├── logo.svg
│       └── manifest.json           # PWA manifest
├── data/                           # Volume Docker (banco SQLite)
└── LICENSE                         # MIT
```

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `3001` | Porta do servidor HTTP |
| `HOST` | `0.0.0.0` | Interface de rede |
| `DB_PATH` | `./data/mercado.db` | Caminho do banco SQLite |
| `STATIC_DIR` | `./static` | Diretório dos assets do frontend |
| `MIGRATIONS` | `./migrations` | Diretório das migrations |
| `SEED_FILE` | `./seed/seed.sql` | Arquivo de seed inicial |

## Contribuindo

1. **Fork** o repositório
2. **Crie uma branch** descritiva: `git checkout -b feat/minha-feature`
3. **Commit** com mensagem clara: `git commit -m "feat: adiciona exportação para CSV"`
4. **Abra um PR** descrevendo o problema que resolve e como testar
5. **Seja legal** — todo feedback e contribuição são bem-vindos

Issues e sugestões também são contribuições. Se encontrar um bug ou tiver uma ideia, abra uma issue.

## Licença

MIT — 9LEVEL 2026

---

<div align="center">
  <p><sub>Feito com café e Go por <a href="https://github.com/9LEVEL">9LEVEL</a></sub></p>
</div>
