package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Repository 是代码库信息。
type Repository struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PathWithNamespace string `json:"pathWithNamespace"`
	Visibility        string `json:"visibility"`
	Archived          bool   `json:"archived"`
	HttpUrlToRepo     string `json:"httpUrlToRepo"`
	SshUrlToRepo      string `json:"sshUrlToRepo"`
	WebURL            string `json:"webUrl"`
	LastActivityAt    string `json:"lastActivityAt"`
}

// ListRepositoriesOptions 是列出代码库的查询参数。
type ListRepositoriesOptions struct {
	Search  string
	Page    int
	PerPage int
	// Archived 为 nil 时不传该参数（返回全部）
	Archived *bool
}

// ListRepositories 列出组织下的代码库（仅中心版，需配置 OrganizationID）。
func (c *Client) ListRepositories(ctx context.Context, opts ListRepositoriesOptions) ([]Repository, error) {
	if c.OrganizationID == "" {
		return nil, fmt.Errorf("list repositories 需要配置 org-id（仅中心版支持）")
	}
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	perPage := opts.PerPage
	if perPage <= 0 {
		perPage = 20
	}

	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("perPage", strconv.Itoa(perPage))
	if opts.Search != "" {
		q.Set("search", opts.Search)
	}
	if opts.Archived != nil {
		q.Set("archived", strconv.FormatBool(*opts.Archived))
	}

	path := fmt.Sprintf("%s/repositories?%s", c.codeupBasePath(), q.Encode())

	var items []Repository
	if err := c.get(ctx, path, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetRepository 按 ID 或全路径获取单个代码库。
func (c *Client) GetRepository(ctx context.Context, repositoryID string) (*Repository, error) {
	path := fmt.Sprintf("%s/repositories/%s", c.codeupBasePath(), url.PathEscape(repositoryID))
	var repo Repository
	if err := c.get(ctx, path, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}
