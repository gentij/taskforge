package api

type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"pageSize"`
	Total      int  `json:"total"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

type Paginated[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Workflow struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	IsActive        bool    `json:"isActive"`
	LatestVersionID *string `json:"latestVersionId"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

type WorkflowVersion struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflowId"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

type Trigger struct {
	ID         string `json:"id"`
	WorkflowID string `json:"workflowId"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	Config     any    `json:"config"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type WorkflowRun struct {
	ID                string  `json:"id"`
	WorkflowID        string  `json:"workflowId"`
	WorkflowVersionID string  `json:"workflowVersionId"`
	Status            string  `json:"status"`
	CreatedAt         string  `json:"createdAt"`
	StartedAt         *string `json:"startedAt"`
	FinishedAt        *string `json:"finishedAt"`
}

type StepRun struct {
	ID            string  `json:"id"`
	WorkflowRunID string  `json:"workflowRunId"`
	StepKey       string  `json:"stepKey"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"createdAt"`
	StartedAt     *string `json:"startedAt"`
	FinishedAt    *string `json:"finishedAt"`
}

type Secret struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}
