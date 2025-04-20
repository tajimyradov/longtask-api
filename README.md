## 🚧 Архитектура

Ключевые принципы:
- Четкое разделение **веб-интерфейса** и **обработки задач**
- Использование **очереди задач**
- Возможность **расширения бизнес-логики** через простое добавление новых обработчиков
- **Хранение состояния** задач (status + result)
- **Простой и читаемый код** для быстрого подключения новых разработчиков

---

## 🔧 Технологии

- Язык: **Go**
- Web: `net/http` + `gorilla/mux` или `chi`
- Очередь задач: in-memory (на старте, можно заменить на Redis/NSQ/RabbitMQ позже)
- Хранение задач: in-memory map или BoltDB/sqlite для MVP
- Обработка: фоновые воркеры (go routines + канал)
- Расширение задач — через интерфейсы

---

## ✅ Поведение API

1. `POST /tasks` — создать задачу  
   - Request: `{ "type": "long_task", "payload": {...} }`
   - Response: `{ "id": "123", "status": "queued" }`

2. `GET /tasks/{id}` — получить статус или результат  
   - Response:  
     - `{ "id": "123", "status": "in_progress" }`  
     - `{ "id": "123", "status": "done", "result": {...} }`  
     - `{ "id": "123", "status": "failed", "error": "reason" }`

---

## 📦 Структура проекта

```bash
.
├── main.go
├── task/
│   ├── manager.go       # Task manager (create, track, store)
│   ├── types.go         # Interfaces and task types
│   ├── worker.go        # Background processing
└── api/
    └── handler.go       # HTTP endpoints
```

---

## 🧩 Интерфейсы и Типы Задач

```go
type Task interface {
	ID() string
	Type() string
	Payload() any
	Run(ctx context.Context) (any, error)
}
```

Каждый новый тип задачи реализует `Task`.

---

## 💡 Пример: Простая долгосрочная задача

### task/types.go

```go
type LongTask struct {
	id      string
	payload map[string]interface{}
}

func (t *LongTask) ID() string          { return t.id }
func (t *LongTask) Type() string        { return "long_task" }
func (t *LongTask) Payload() any        { return t.payload }
func (t *LongTask) Run(ctx context.Context) (any, error) {
	// Пример долгой операции
	select {
	case <-time.After(3 * time.Minute):
		return map[string]string{"result": "success"}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
```

---

## 🚀 main.go (MVP)

```go
func main() {
	manager := task.NewManager()
	go manager.StartWorkerPool(5) // запустить 5 воркеров

	r := chi.NewRouter()
	r.Post("/tasks", api.CreateTaskHandler(manager))
	r.Get("/tasks/{id}", api.GetTaskHandler(manager))

	http.ListenAndServe(":8080", r)
}
```

---

## 🧠 Возможности для масштабирования

- Использовать Redis как очередь и хранилище задач
- Вынести воркеры в отдельный сервис
- Добавить авторизацию
- Поддержка WebSocket/Callback для уведомлений
- Обработка зависимостей между задачами (планирование)
- Панель мониторинга (например через Prometheus)

---

## ❗️Важно

- Воркеры обрабатывают задачи через канал, и задачи хранятся в мапе `id -> TaskStatus`
- Всё проектируется так, чтобы **добавление нового типа задачи** занимало <5 минут: просто новая структура + реализация интерфейса `Task`

---
