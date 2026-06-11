package api

import (
	"context"
	"fmt"
	"net/url"
)

// CreateMergeRequestOptions 是创建合并请求的请求体。
type CreateMergeRequestOptions struct {
	Title              string   `json:"title"`
	Description        string   `json:"description,omitempty"`
	SourceBranch       string   `json:"sourceBranch"`
	SourceProjectID    int64    `json:"sourceProjectId"`
	TargetBranch       string   `json:"targetBranch"`
	TargetProjectID    int64    `json:"targetProjectId"`
	ReviewerUserIDs    []string `json:"reviewerUserIds,omitempty"`
	TriggerAIReviewRun bool     `json:"triggerAIReviewRun,omitempty"`
	WorkItemIDs        string   `json:"workItemIds,omitempty"`
}

// UserInfo 是接口返回中的用户信息。
type UserInfo struct {
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	State    string `json:"state"`
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

// Reviewer 是评审人信息。
type Reviewer struct {
	UserInfo
	HasCommented        bool   `json:"hasCommented"`
	HasReviewed         bool   `json:"hasReviewed"`
	ReviewOpinionStatus string `json:"reviewOpinionStatus"`
	ReviewTime          string `json:"reviewTime"`
}

// MergeRequest 是创建合并请求接口的返回结果。
type MergeRequest struct {
	Ahead                          int        `json:"ahead"`
	AllRequirementsPass            bool       `json:"allRequirementsPass"`
	Author                         UserInfo   `json:"author"`
	Behind                         int        `json:"behind"`
	CanRevertOrCherryPick          bool       `json:"canRevertOrCherryPick"`
	ConflictCheckStatus            string     `json:"conflictCheckStatus"`
	CreateFrom                     string     `json:"createFrom"`
	CreateTime                     string     `json:"createTime"`
	Description                    string     `json:"description"`
	DetailURL                      string     `json:"detailUrl"`
	HasReverted                    bool       `json:"hasReverted"`
	LocalID                        int64      `json:"localId"`
	MergedRevision                 string     `json:"mergedRevision"`
	MrType                         string     `json:"mrType"`
	ProjectID                      int64      `json:"projectId"`
	Reviewers                      []Reviewer `json:"reviewers"`
	SourceBranch                   string     `json:"sourceBranch"`
	SourceProjectID                int64      `json:"sourceProjectId"`
	Status                         string     `json:"status"`
	Subscribers                    []UserInfo `json:"subscribers"`
	SupportMergeFastForwardOnly    bool       `json:"supportMergeFastForwardOnly"`
	TargetBranch                   string     `json:"targetBranch"`
	TargetProjectID                int64      `json:"targetProjectId"`
	TargetProjectNameWithNamespace string     `json:"targetProjectNameWithNamespace"`
	TargetProjectPathWithNamespace string     `json:"targetProjectPathWithNamespace"`
	Title                          string     `json:"title"`
	TotalCommentCount              int        `json:"totalCommentCount"`
	UnResolvedCommentCount         int        `json:"unResolvedCommentCount"`
	UpdateTime                     string     `json:"updateTime"`
	WebURL                         string     `json:"webUrl"`
}

// CreateMergeRequest 创建合并请求。
// repositoryID 可以是代码库数字 ID，也可以是全路径（如 org/repo，内部会做 URL 编码）。
func (c *Client) CreateMergeRequest(ctx context.Context, repositoryID string, opts *CreateMergeRequestOptions) (*MergeRequest, error) {
	path := fmt.Sprintf("%s/repositories/%s/changeRequests", c.codeupBasePath(), url.PathEscape(repositoryID))
	var mr MergeRequest
	if err := c.post(ctx, path, opts, &mr); err != nil {
		return nil, err
	}
	return &mr, nil
}
