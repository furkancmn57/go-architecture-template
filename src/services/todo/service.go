package todo

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/furkancmn57/go-architecture-template/src/common/apperr"
	"github.com/furkancmn57/go-architecture-template/src/constants"
	"github.com/furkancmn57/go-architecture-template/src/models/requests"
	"github.com/furkancmn57/go-architecture-template/src/models/responses"
	"github.com/furkancmn57/go-architecture-template/src/services/todo/validations"
)

// Todo is the in-memory domain model.
type Todo struct {
	ID          uuid.UUID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Service implements todo use-cases with a process-local map (no DB on this branch).
type Service struct {
	mu   sync.RWMutex
	byID map[uuid.UUID]*Todo
}

// NewService wires a Service with an empty in-memory collection.
func NewService() *Service {
	return &Service{byID: make(map[uuid.UUID]*Todo)}
}

// Create validates and stores a new todo.
func (s *Service) Create(_ context.Context, req requests.CreateTodoRequest) (*responses.TodoResponse, *apperr.Error) {
	if err := validations.CreateTodoRequest(req); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	row := &Todo{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.mu.Lock()
	s.byID[row.ID] = row
	s.mu.Unlock()

	resp := toResponse(row)
	return &resp, nil
}

// Todos lists every todo, most recent first.
func (s *Service) Todos(_ context.Context) ([]responses.TodoResponse, *apperr.Error) {
	s.mu.RLock()
	rows := make([]*Todo, 0, len(s.byID))
	for _, row := range s.byID {
		rows = append(rows, row)
	}
	s.mu.RUnlock()

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CreatedAt.After(rows[j].CreatedAt)
	})
	out := make([]responses.TodoResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, toResponse(row))
	}
	return out, nil
}

// TodoById fetches a single todo by its ID.
func (s *Service) TodoById(_ context.Context, id string) (*responses.TodoResponse, *apperr.Error) {
	row, appErr := s.findByID(id)
	if appErr != nil {
		return nil, appErr
	}
	resp := toResponse(row)
	return &resp, nil
}

// Update validates and applies changes to an existing todo.
func (s *Service) Update(_ context.Context, id string, req requests.UpdateTodoRequest) (*responses.TodoResponse, *apperr.Error) {
	if err := validations.UpdateTodoRequest(req); err != nil {
		return nil, err
	}
	row, appErr := s.findByID(id)
	if appErr != nil {
		return nil, appErr
	}

	s.mu.Lock()
	row.Title = req.Title
	row.Description = req.Description
	row.Completed = req.Completed
	row.UpdatedAt = time.Now().UTC()
	s.mu.Unlock()

	resp := toResponse(row)
	return &resp, nil
}

// Complete marks a todo as completed.
func (s *Service) Complete(_ context.Context, id string) (*responses.TodoResponse, *apperr.Error) {
	row, appErr := s.findByID(id)
	if appErr != nil {
		return nil, appErr
	}
	if row.Completed {
		return nil, apperr.Conflict(constants.TodoAlreadyCompleted, "todo is already completed")
	}

	s.mu.Lock()
	row.Completed = true
	row.UpdatedAt = time.Now().UTC()
	s.mu.Unlock()

	resp := toResponse(row)
	return &resp, nil
}

// Delete removes a todo.
func (s *Service) Delete(_ context.Context, id string) *apperr.Error {
	todoID, err := uuid.Parse(id)
	if err != nil {
		return apperr.BadRequest(constants.TodoInvalidID, "invalid todo id")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[todoID]; !ok {
		return apperr.NotFound(constants.TodoNotFound, "todo not found")
	}
	delete(s.byID, todoID)
	return nil
}

func (s *Service) findByID(id string) (*Todo, *apperr.Error) {
	todoID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.BadRequest(constants.TodoInvalidID, "invalid todo id")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	row, ok := s.byID[todoID]
	if !ok {
		return nil, apperr.NotFound(constants.TodoNotFound, "todo not found")
	}
	return row, nil
}

func toResponse(row *Todo) responses.TodoResponse {
	return responses.TodoResponse{
		ID:          row.ID.String(),
		Title:       row.Title,
		Description: row.Description,
		Completed:   row.Completed,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
