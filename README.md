# 🎯 Monitoring Lab

Laboratório completo de monitoramento com aplicações Go e Spring Boot conectadas ao MongoDB, expondo métricas para Prometheus e visualização no Grafana.

## 📁 Estrutura do Projeto

```
/monitoring-lab
    /go-app              # Aplicação Go
        main.go
        go.mod
        Dockerfile
    /spring-app          # Aplicação Spring Boot
        src/main/java/...
        pom.xml
        Dockerfile
    /prometheus          # Configuração do Prometheus
        prometheus.yml
    /grafana            # Provisioning do Grafana
        provisioning/datasources/
    docker-compose.yml
    README.md
```

## 🚀 Como Executar

### Pré-requisitos

- Docker
- Docker Compose

### Iniciar todos os serviços

```bash
docker-compose up --build
```

Este comando irá:
1. Construir as imagens das aplicações Go e Spring Boot
2. Iniciar todos os serviços (MongoDB, aplicações, Prometheus, Grafana, exporters)
3. Configurar a rede entre os containers

### Parar os serviços

```bash
docker-compose down
```

Para remover também os volumes (dados):

```bash
docker-compose down -v
```

## 🌐 Acessos

### Aplicações

- **Go App**: http://localhost:8080/ping
- **Spring App**: http://localhost:8081/ping

### Ferramentas de Monitoramento

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
  - Usuário: `admin`
  - Senha: `admin`

## 📊 Métricas Expostas

### Go App (`/metrics`)

- `mongodb_pings_total` - Total de operações ping realizadas
- `mongodb_inserts_total` - Total de operações de inserção
- `mongodb_operation_duration_seconds` - Latência das operações MongoDB

### Spring App (`/actuator/prometheus`)

- `mongodb_total_operations` - Total de operações MongoDB
- `mongodb_count_velocity` - Velocidade de operações de contagem
- `mongodb_operation_latency` - Latência das operações MongoDB

### MongoDB Exporter

- Métricas do MongoDB (conexões, operações, performance, etc.)

### Node Exporter

- Métricas do sistema (CPU, memória, disco, rede, etc.)

## 📈 Visualizando Métricas

### No Prometheus

1. Acesse http://localhost:9090
2. Use a interface de query para testar métricas:
   - `mongodb_pings_total`
   - `mongodb_inserts_total`
   - `mongodb_operation_duration_seconds`
   - `mongodb_total_operations`
   - `mongodb_operation_latency_seconds`

### No Grafana

1. Acesse http://localhost:3000
2. Faça login com `admin/admin`
3. O datasource Prometheus já está configurado automaticamente
4. **Dashboard pré-configurado**: O dashboard "Monitoring Lab - Go vs Java" já está disponível automaticamente!
   - Acesse: Dashboards > Monitoring Lab - Go vs Java
   - O dashboard mostra:
     - **Parte Superior**: Métricas da aplicação Go (CPU, Memória, Latência MongoDB)
     - **Parte Inferior**: Métricas da aplicação Java Spring Boot (CPU, Memória, Latência MongoDB)
     - **Comparações**: Gráficos comparativos de latência e memória entre Go e Java

#### Exemplo de Queries para Dashboard

**Taxa de pings por segundo (Go App):**
```
rate(mongodb_pings_total[1m])
```

**Taxa de inserts por segundo (Go App):**
```
rate(mongodb_inserts_total[1m])
```

**Latência média (Go App):**
```
rate(mongodb_operation_duration_seconds_sum[1m]) / rate(mongodb_operation_duration_seconds_count[1m])
```

**Total de operações (Spring App):**
```
mongodb_total_operations
```

**Latência (Spring App):**
```
mongodb_operation_latency_seconds
```

## 🔧 Configurações

### Variáveis de Ambiente

**Go App:**
- `MONGO_URI`: URI de conexão MongoDB (padrão: `mongodb://mongo:27017`)

**Spring App:**
- `SPRING_DATA_MONGODB_URI`: URI de conexão MongoDB (padrão: `mongodb://mongo:27017/monitoring`)

### Prometheus

O arquivo `prometheus/prometheus.yml` configura:
- Intervalo de scrape: 5 segundos
- Targets: go-app, spring-app, mongodb-exporter, node-exporter

## 📝 Funcionalidades

### Go App

- ✅ Conecta ao MongoDB usando mongo-driver
- ✅ Loop infinito em goroutine executando a cada 5 segundos:
  - Ping no MongoDB (`db.RunCommand({"ping":1})`)
  - Inserção de documento na collection `events`
- ✅ Expõe métricas Prometheus:
  - Contador de pings
  - Contador de inserts
  - Histograma de latência
- ✅ Endpoint `/ping` para health check
- ✅ Endpoint `/metrics` para Prometheus
- ✅ Dockerfile multi-stage

### Spring App

- ✅ Java 17 + Spring Boot 3.x
- ✅ Spring Data MongoDB
- ✅ Scheduler automático (`@Scheduled`) executando a cada 5 segundos:
  - Count na collection `events`
  - Inserção de documento
- ✅ Métricas Micrometer:
  - Total de operações
  - Velocidade de contagem
  - Latência
- ✅ Endpoint `/ping` para health check
- ✅ Endpoint `/actuator/prometheus` para métricas
- ✅ Dockerfile funcional

## 🐳 Serviços Docker Compose

1. **MongoDB** - Banco de dados na porta 27017
2. **Go App** - Aplicação Go na porta 8080
3. **Spring App** - Aplicação Spring Boot na porta 8081
4. **MongoDB Exporter** - Exportador de métricas do MongoDB na porta 9216
5. **Node Exporter** - Exportador de métricas do sistema na porta 9100
6. **Prometheus** - Coletor e armazenador de métricas na porta 9090
7. **Grafana** - Visualização de métricas na porta 3000

## 🔍 Verificando o Funcionamento

1. **Verificar aplicações:**
   ```bash
   curl http://localhost:8080/ping  # Go App
   curl http://localhost:8081/ping  # Spring App
   ```

2. **Verificar métricas:**
   ```bash
   curl http://localhost:8080/metrics  # Go App
   curl http://localhost:8081/actuator/prometheus  # Spring App
   ```

3. **Verificar logs:**
   ```bash
   docker-compose logs -f go-app
   docker-compose logs -f spring-app
   ```

## 📚 Próximos Passos

- Criar dashboards personalizados no Grafana
- Adicionar alertas no Prometheus
- Configurar alertas no Grafana
- Adicionar mais métricas customizadas
- Implementar health checks mais robustos

## 🐛 Troubleshooting

### Aplicações não conectam ao MongoDB

Verifique se o MongoDB está saudável:
```bash
docker-compose ps mongo
docker-compose logs mongo
```

### Prometheus não coleta métricas

Verifique os targets no Prometheus:
- Acesse http://localhost:9090/targets
- Verifique se todos os targets estão "UP"

### Grafana não acessa Prometheus

Verifique se o datasource está configurado:
- Acesse http://localhost:3000/connections/datasources
- Verifique se o Prometheus está configurado e testado

