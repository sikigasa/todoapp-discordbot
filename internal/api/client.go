package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client はgithub-task-controller APIクライアント
type Client struct {
	baseURL    string
	authCookie string
	httpClient *http.Client
}

// NewClient は新しいAPIクライアントを作成する
func NewClient(baseURL, authCookie string) *Client {
	return &Client{
		baseURL:    baseURL,
		authCookie: authCookie,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ─── Todo models ────────────────────────────────────────

// Todo はTODOアイテムを表す
type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTodoRequest はTODO作成リクエスト
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTodoRequest はTODO更新リクエスト
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// ─── Task models ────────────────────────────────────────

// Task はタスクを表す
type Task struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      int     `json:"status"`
	Priority    int     `json:"priority"`
	EndDate     *string `json:"end_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateTaskRequest はタスク作成リクエスト
type CreateTaskRequest struct {
	ProjectID   string  `json:"project_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      int     `json:"status"`
	Priority    int     `json:"priority"`
	EndDate     *string `json:"end_date,omitempty"`
}

// UpdateTaskRequest はタスク更新リクエスト
type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      int     `json:"status"`
	Priority    int     `json:"priority"`
	EndDate     *string `json:"end_date,omitempty"`
}

// ─── Project models ─────────────────────────────────────

// Project はプロジェクトを表す
type Project struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateProjectRequest はプロジェクト作成リクエスト
type CreateProjectRequest struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateProjectRequest はプロジェクト更新リクエスト
type UpdateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ProblemDetail はRFC 9457エラーレスポンス
type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// ─── HTTP helpers ───────────────────────────────────────

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.AddCookie(&http.Cookie{
		Name:  "auth-session",
		Value: c.authCookie,
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// ─── Todo API ───────────────────────────────────────────

// CreateTodo はTODOを作成する
func (c *Client) CreateTodo(req *CreateTodoRequest) (*Todo, error) {
	body, status, err := c.doRequest("POST", "/api/v1/todos", req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, parseError(body, status)
	}

	var todo Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &todo, nil
}

// ListTodos は全TODOを取得する
func (c *Client) ListTodos() ([]*Todo, error) {
	body, status, err := c.doRequest("GET", "/api/v1/todos", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var todos []*Todo
	if err := json.Unmarshal(body, &todos); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return todos, nil
}

// GetTodo はIDでTODOを取得する
func (c *Client) GetTodo(id string) (*Todo, error) {
	body, status, err := c.doRequest("GET", "/api/v1/todos/"+id, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var todo Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &todo, nil
}

// UpdateTodo はTODOを更新する
func (c *Client) UpdateTodo(id string, req *UpdateTodoRequest) (*Todo, error) {
	body, status, err := c.doRequest("PUT", "/api/v1/todos/"+id, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var todo Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &todo, nil
}

// DeleteTodo はTODOを削除する
func (c *Client) DeleteTodo(id string) error {
	_, status, err := c.doRequest("DELETE", "/api/v1/todos/"+id, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", status)
	}
	return nil
}

// ─── Task API ───────────────────────────────────────────

// CreateTask はタスクを作成する
func (c *Client) CreateTask(req *CreateTaskRequest) (*Task, error) {
	body, status, err := c.doRequest("POST", "/api/v1/tasks", req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, parseError(body, status)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &task, nil
}

// ListTasks はプロジェクトIDでタスク一覧を取得する
func (c *Client) ListTasks(projectID string) ([]*Task, error) {
	body, status, err := c.doRequest("GET", "/api/v1/tasks?project_id="+projectID, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var tasks []*Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return tasks, nil
}

// GetTask はIDでタスクを取得する
func (c *Client) GetTask(id string) (*Task, error) {
	body, status, err := c.doRequest("GET", "/api/v1/tasks/"+id, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &task, nil
}

// UpdateTask はタスクを更新する
func (c *Client) UpdateTask(id string, req *UpdateTaskRequest) (*Task, error) {
	body, status, err := c.doRequest("PUT", "/api/v1/tasks/"+id, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &task, nil
}

// DeleteTask はタスクを削除する
func (c *Client) DeleteTask(id string) error {
	_, status, err := c.doRequest("DELETE", "/api/v1/tasks/"+id, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", status)
	}
	return nil
}

// ─── Project API ────────────────────────────────────────

// CreateProject はプロジェクトを作成する
func (c *Client) CreateProject(req *CreateProjectRequest) (*Project, error) {
	body, status, err := c.doRequest("POST", "/api/v1/projects", req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, parseError(body, status)
	}

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &project, nil
}

// ListProjects はユーザーIDでプロジェクト一覧を取得する
func (c *Client) ListProjects(userID string) ([]*Project, error) {
	body, status, err := c.doRequest("GET", "/api/v1/projects?user_id="+userID, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var projects []*Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return projects, nil
}

// GetProject はIDでプロジェクトを取得する
func (c *Client) GetProject(id string) (*Project, error) {
	body, status, err := c.doRequest("GET", "/api/v1/projects/"+id, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &project, nil
}

// UpdateProject はプロジェクトを更新する
func (c *Client) UpdateProject(id string, req *UpdateProjectRequest) (*Project, error) {
	body, status, err := c.doRequest("PUT", "/api/v1/projects/"+id, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, parseError(body, status)
	}

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &project, nil
}

// DeleteProject はプロジェクトを削除する
func (c *Client) DeleteProject(id string) error {
	_, status, err := c.doRequest("DELETE", "/api/v1/projects/"+id, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", status)
	}
	return nil
}

// ─── Helpers ────────────────────────────────────────────

// StatusText はタスクステータスの数値を文字列に変換する
func StatusText(status int) string {
	switch status {
	case 0:
		return "To Do"
	case 1:
		return "In Progress"
	case 2:
		return "Done"
	default:
		return "Unknown"
	}
}

// PriorityText はタスク優先度の数値を文字列に変換する
func PriorityText(priority int) string {
	switch priority {
	case 0:
		return "Low"
	case 1:
		return "Medium"
	case 2:
		return "High"
	default:
		return "Unknown"
	}
}

// StatusEmoji はタスクステータスのEmojiを返す
func StatusEmoji(status int) string {
	switch status {
	case 0:
		return "📋"
	case 1:
		return "🔄"
	case 2:
		return "✅"
	default:
		return "❓"
	}
}

// PriorityEmoji はタスク優先度のEmojiを返す
func PriorityEmoji(priority int) string {
	switch priority {
	case 0:
		return "🟢"
	case 1:
		return "🟡"
	case 2:
		return "🔴"
	default:
		return "⚪"
	}
}

func parseError(body []byte, status int) error {
	var problem ProblemDetail
	if err := json.Unmarshal(body, &problem); err == nil && problem.Detail != "" {
		return fmt.Errorf("[%d] %s: %s", status, problem.Title, problem.Detail)
	}
	return fmt.Errorf("API error (status %d): %s", status, string(body))
}
