package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestClient 返回指向 httptest 服务器的客户端。
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)

	c := NewClient(strings.TrimPrefix(srv.URL, "https://"), "pt-test-token", "")
	c.httpClient = srv.Client()
	return c
}

func TestCreateMergeRequestCentral(t *testing.T) {
	var gotPath, gotToken string
	var gotBody map[string]any

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotToken = r.Header.Get("x-yunxiao-token")
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"localId": 7,
			"title":   "mr title",
			"status":  "UNDER_REVIEW",
			"webUrl":  "https://example.com/change/7",
		})
	})
	client.OrganizationID = "60d54f3daccf2bbd6659f3ad"

	mr, err := client.CreateMergeRequest(context.Background(), "myorg/DemoRepo", &CreateMergeRequestOptions{
		Title:           "mr title",
		SourceBranch:    "demo-branch",
		SourceProjectID: 2813489,
		TargetBranch:    "master",
		TargetProjectID: 2813489,
		ReviewerUserIDs: []string{"62c795xxxb468af8"},
	})
	if err != nil {
		t.Fatal(err)
	}

	wantPath := "/oapi/v1/codeup/organizations/60d54f3daccf2bbd6659f3ad/repositories/myorg%2FDemoRepo/changeRequests"
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
	if gotToken != "pt-test-token" {
		t.Errorf("token header = %q", gotToken)
	}
	if gotBody["sourceProjectId"] != float64(2813489) || gotBody["targetBranch"] != "master" {
		t.Errorf("unexpected body: %v", gotBody)
	}
	if _, ok := gotBody["description"]; ok {
		t.Error("空 description 不应出现在请求体中")
	}
	if mr.LocalID != 7 || mr.Status != "UNDER_REVIEW" {
		t.Errorf("unexpected result: %+v", mr)
	}
}

func TestCreateMergeRequestRegion(t *testing.T) {
	var gotPath string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		_ = json.NewEncoder(w).Encode(map[string]any{"localId": 1})
	})

	if _, err := client.CreateMergeRequest(context.Background(), "2813489", &CreateMergeRequestOptions{
		Title: "t", SourceBranch: "a", TargetBranch: "b", SourceProjectID: 1, TargetProjectID: 1,
	}); err != nil {
		t.Fatal(err)
	}

	wantPath := "/oapi/v1/codeup/repositories/2813489/changeRequests"
	if gotPath != wantPath {
		t.Errorf("path = %q, want %q", gotPath, wantPath)
	}
}

func TestCreateMergeRequestAPIError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"errorCode":"AccessDenied","errorMessage":"无权限","requestId":"req-123"}`))
	})

	_, err := client.CreateMergeRequest(context.Background(), "1", &CreateMergeRequestOptions{
		Title: "t", SourceBranch: "a", TargetBranch: "b", SourceProjectID: 1, TargetProjectID: 1,
	})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 403 || apiErr.Code != "AccessDenied" || apiErr.RequestID != "req-123" {
		t.Errorf("unexpected APIError: %+v", apiErr)
	}
}
