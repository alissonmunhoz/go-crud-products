# 🛠 go-crud-products

API REST em Go para gerenciamento de produtos — operações CRUD (Create, Read, Update, Delete).

Este projeto serve como backend para uma aplicação de gerenciamento de inventário, podendo ser usado isoladamente ou em conjunto com um frontend (ex: Next.js).

---

## 🔧 Tecnologias usadas

- Linguagem: **Go (Golang)**
- Módulos Go (`go.mod`)
- Docker / Docker Compose
- Organização de código por pacotes internos
- Middlewares, roteamento e tratamento de erros
- MySQL (ou SQLite, dependendo da configuração)

---

## 📁 Estrutura do projeto

```
go-crud-products/
├── cmd/                   # Ponto de entrada da aplicação (main, servidores, inicialização)
├── internal/              # Pacotes internos (produtos, repositório, serviço)
├── docs/                  # Documentação (ex: arquivo OpenAPI, exemplos)
├── Dockerfile             # Configuração para container Docker
├── docker-compose.yml     # Orquestração com serviços auxiliares (banco, etc.)
├── go.mod                 # Módulo Go
├── go.sum                 # Verificações de dependências
└── README.md              # Documentação principal
```

---

## 🚀 Como rodar localmente

### Pré-requisitos

- Go (1.25)
- Docker & Docker Compose (opcional, mas recomendado para facilitar execução)
- Banco configurado (MySQL ou SQLite conforme ajuste no código)

### Rodando sem Docker

1. Clone o repositório

   ```bash
   git clone https://github.com/alissonmunhoz/go-crud-products.git
   cd go-crud-products
   ```

2. A API estará disponível em `http://localhost:8080/v1`

### Rodando com Docker / Docker Compose

1. Certifique-se de ter o Docker ativo
2. Na raiz do projeto, execute:
   ```bash
   docker-compose up
   ```
3. Isso deverá levantar containers — API + banco (dependendo da configuração).

---

## 📦 Endpoints da API

| Método   | Rota               | Descrição                     | Corpo (JSON) / Parâmetros                                                  |
| -------- | ------------------ | ----------------------------- | -------------------------------------------------------------------------- |
| `GET`    | `/v1/products`     | Lista todos os produtos       | —                                                                          |
| `GET`    | `/v1/product?id=1` | Retorna um produto pelo ID    | Query param `id`                                                           |
| `POST`   | `/v1/product`      | Cria um novo produto          | `{ "name": "...", "price": 123.45, "quantity": 10, "description": "..." }` |
| `PUT`    | `/v1/product?id=1` | Atualiza um produto existente | Query param `id` + corpo JSON com campos a mudar                           |
| `DELETE` | `/v1/product?id=1` | Remove um produto pelo ID     | Query param `id`                                                           |

### Exemplo de JSON para criação/atualização

```json
{
  "name": "Teclado Mecânico",
  "price": 299.99,
  "quantity": 20,
  "description": "Teclado com switches mecânicos AZUL"
}
```

---

## 📄 Documentação adicional / Swagger

Se configurado, a documentação Swagger estará disponível em:

```
http://localhost:8080/swagger/index.html
```

---

## 🧪 Testes

Execute os testes com:

```bash
go test ./...
```
