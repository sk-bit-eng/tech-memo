package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	memorepo "tech-memo/internal/adapter/gateway/memo"
	todorepo "tech-memo/internal/adapter/gateway/todo"
	"tech-memo/internal/domain"
	sqlserverinfra "tech-memo/internal/infrastructure/persistence/sqlserver"

	"github.com/google/uuid"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	testUserID = "test-user-001"
)

func TestAllEndpoints(t *testing.T) {
	baseDSN := getTestDSN()
	dbName := fmt.Sprintf("tech_memo_api_test_%d", time.Now().UnixNano())
	dsn := replaceDatabaseInDSN(baseDSN, dbName)

	db, err := sqlserverinfra.Open(dsn)
	if err != nil {
		t.Skipf("sqlserver unavailable for integration test: %v", err)
	}
	defer cleanupTestDatabase(t, db, dbName, baseDSN)

	if err := seedTestData(db); err != nil {
		t.Fatalf("seed test data: %v", err)
	}

	router := newRouter(db)

	body := performJSONRequest(t, router, http.MethodGet, "/", nil)
	assertContains(t, body, `"message":"teck-memo server is running"`)

	memoID := testMemoEndpoints(t, router)
	testMemoUserEndpoints(t, router, memoID)

	todoID := testTodoEndpoints(t, router)
	testTodoUserEndpoints(t, router, todoID)
}

func testMemoEndpoints(t *testing.T, router http.Handler) string {
	createBody := map[string]any{
		"userID":     testUserID,
		"title":      "memo created title",
		"content":    "memo created content",
		"categoryID": "memo-cat-1",
		"parameters": []map[string]any{},
		"isPinned":   false,
	}

	body := performJSONRequestWithStatus(t, router, http.MethodPost, "/memos", createBody, http.StatusCreated)
	memoID := decodeID(t, body)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/memos/"+memoID, nil, http.StatusOK)
	assertContains(t, body, `"ID":"`+memoID+`"`)

	updateBody := map[string]any{
		"title":      "memo updated title",
		"content":    "memo updated content",
		"categoryID": "memo-cat-2",
		"parameters": []map[string]any{},
		"isPinned":   false,
	}
	body = performJSONRequestWithStatus(t, router, http.MethodPut, "/memos/"+memoID, updateBody, http.StatusOK)
	assertContains(t, body, `"Title":"memo updated title"`)

	body = performJSONRequestWithStatus(t, router, http.MethodPatch, "/memos/"+memoID+"/pin", nil, http.StatusOK)
	assertContains(t, body, `"IsPinned":true`)

	performJSONRequestWithStatus(t, router, http.MethodDelete, "/memos/"+memoID, nil, http.StatusNoContent)

	return memoID
}

func testMemoUserEndpoints(t *testing.T, router http.Handler, deletedMemoID string) {
	body := performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/memos", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-memo-1"`)
	assertNotContains(t, body, deletedMemoID)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/memos/search?q=seed%20memo", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-memo-1"`)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/memos/category/memo-cat-seed", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-memo-1"`)
}

func testTodoEndpoints(t *testing.T, router http.Handler) string {
	createBody := map[string]any{
		"userID":     testUserID,
		"title":      "todo created title",
		"content":    "todo created content",
		"categoryID": "todo-cat-1",
		"parameters": []map[string]any{},
		"isPinned":   false,
	}

	body := performJSONRequestWithStatus(t, router, http.MethodPost, "/todos", createBody, http.StatusCreated)
	todoID := decodeID(t, body)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/todos/"+todoID, nil, http.StatusOK)
	assertContains(t, body, `"ID":"`+todoID+`"`)

	updateBody := map[string]any{
		"title":      "todo updated title",
		"content":    "todo updated content",
		"categoryID": "todo-cat-2",
		"parameters": []map[string]any{},
		"isPinned":   false,
	}
	body = performJSONRequestWithStatus(t, router, http.MethodPut, "/todos/"+todoID, updateBody, http.StatusOK)
	assertContains(t, body, `"Title":"todo updated title"`)

	body = performJSONRequestWithStatus(t, router, http.MethodPatch, "/todos/"+todoID+"/pin", nil, http.StatusOK)
	assertContains(t, body, `"IsPinned":true`)

	performJSONRequestWithStatus(t, router, http.MethodPatch, "/todos/"+todoID+"/complete", nil, http.StatusNoContent)
	performJSONRequestWithStatus(t, router, http.MethodPatch, "/todos/"+todoID+"/incomplete", nil, http.StatusNoContent)
	performJSONRequestWithStatus(t, router, http.MethodDelete, "/todos/"+todoID, nil, http.StatusNoContent)

	return todoID
}

func testTodoUserEndpoints(t *testing.T, router http.Handler, deletedTodoID string) {
	body := performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/todos", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-todo-pending"`)
	assertContains(t, body, `"ID":"seed-todo-completed"`)
	assertNotContains(t, body, deletedTodoID)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/todos/search?q=pending", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-todo-pending"`)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/todos/pending", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-todo-pending"`)
	assertNotContains(t, body, `"ID":"seed-todo-completed"`)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/todos/completed", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-todo-completed"`)

	body = performJSONRequestWithStatus(t, router, http.MethodGet, "/users/"+testUserID+"/todos/category/todo-cat-seed", nil, http.StatusOK)
	assertContains(t, body, `"ID":"seed-todo-pending"`)
}

func seedTestData(db *gorm.DB) error {
	memoRepo := memorepo.NewRepository(db)
	todoRepo := todorepo.NewRepository(db)
	now := time.Now().UTC()
	completedAt := now.Add(-time.Hour)

	if err := memoRepo.Save(&domain.Memo{
		ID:         "seed-memo-1",
		UserID:     testUserID,
		Title:      "seed memo title",
		Content:    "seed memo content",
		CategoryID: "memo-cat-seed",
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return err
	}

	if err := todoRepo.Save(&domain.Todo{
		ID:         "seed-todo-pending",
		UserID:     testUserID,
		Title:      "seed pending todo",
		Content:    "seed todo content",
		CategoryID: "todo-cat-seed",
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return err
	}

	if err := todoRepo.Save(&domain.Todo{
		ID:          "seed-todo-completed",
		UserID:      testUserID,
		Title:       "seed completed todo",
		Content:     "seed todo completed content",
		CategoryID:  "todo-cat-done",
		CompletedAt: &completedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		return err
	}

	return nil
}

func cleanupTestDatabase(t *testing.T, db *gorm.DB, dbName, baseDSN string) {
	t.Helper()

	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	master, err := gorm.Open(sqlserver.Open(replaceDatabaseInDSN(baseDSN, "master")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open master db for cleanup: %v", err)
	}

	sql := fmt.Sprintf(`
IF DB_ID(N'%s') IS NOT NULL
BEGIN
	ALTER DATABASE [%s] SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
	DROP DATABASE [%s];
END
`, dbName, dbName, dbName)

	if err := master.Exec(sql).Error; err != nil {
		t.Fatalf("drop test database: %v", err)
	}
}

func performJSONRequest(t *testing.T, router http.Handler, method, path string, body any) string {
	t.Helper()
	return performJSONRequestWithStatus(t, router, method, path, body, http.StatusOK)
}

func performJSONRequestWithStatus(t *testing.T, router http.Handler, method, path string, body any, wantStatus int) string {
	t.Helper()

	var reqBody *bytes.Reader
	if body == nil {
		reqBody = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != wantStatus {
		t.Fatalf("%s %s: status=%d body=%s", method, path, rec.Code, rec.Body.String())
	}

	return rec.Body.String()
}

func decodeID(t *testing.T, body string) string {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	id, ok := payload["ID"].(string)
	if !ok || id == "" {
		t.Fatalf("response id missing: %s", body)
	}

	if _, err := uuid.Parse(id); err != nil {
		t.Fatalf("response id is not uuid: %s", id)
	}

	return id
}

func assertContains(t *testing.T, body, substr string) {
	t.Helper()
	if !strings.Contains(body, substr) {
		t.Fatalf("response does not contain %q: %s", substr, body)
	}
}

func assertNotContains(t *testing.T, body, substr string) {
	t.Helper()
	if strings.Contains(body, substr) {
		t.Fatalf("response unexpectedly contains %q: %s", substr, body)
	}
}

func getTestDSN() string {
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		return dsn
	}
	return "sqlserver://sa:Test@1234@localhost:1433?database=tech_memo"
}

func replaceDatabaseInDSN(dsn, dbName string) string {
	re := regexp.MustCompile(`(?i)([?&]database=)([^&]+)`)
	if re.MatchString(dsn) {
		return re.ReplaceAllString(dsn, "${1}"+dbName)
	}

	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	return dsn + sep + "database=" + dbName
}
