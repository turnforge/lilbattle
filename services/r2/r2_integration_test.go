package r2

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	v1 "github.com/turnforge/lilbattle/gen/go/lilbattle/v1/models"
)

// testR2Client creates a real R2 client from environment variables.
// Returns nil if credentials are not set (skips test).
func testR2Client(t *testing.T) *R2Client {
	t.Helper()
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" || accessKey == "" || secretKey == "" {
		t.Skip("R2 credentials not set, skipping integration test")
	}

	client, err := NewR2Client(R2Config{
		AccountID:      accountID,
		AccessKeyID:    accessKey,
		SecretAccessKey: secretKey,
		Bucket:         "lilbattle-assets",
	})
	if err != nil {
		t.Fatalf("Failed to create R2 client: %v", err)
	}
	return client
}

func TestR2UploadGetDelete(t *testing.T) {
	client := testR2Client(t)
	ctx := context.Background()
	svc := NewR2FileStoreService(client)

	testKey := fmt.Sprintf("test/ci-integration-%d.txt", time.Now().UnixNano())
	testContent := []byte("integration test content")

	// Upload
	putResp, err := svc.PutFile(ctx, &v1.PutFileRequest{
		File:    &v1.File{Path: testKey, ContentType: "text/plain"},
		Content: testContent,
	})
	if err != nil {
		t.Fatalf("PutFile: %v", err)
	}
	if putResp.File.FileSize != uint64(len(testContent)) {
		t.Errorf("PutFile size = %d, want %d", putResp.File.FileSize, len(testContent))
	}

	// Get metadata
	getResp, err := svc.GetFile(ctx, &v1.GetFileRequest{Path: testKey})
	if err != nil {
		t.Fatalf("GetFile: %v", err)
	}
	if getResp.File.FileSize != uint64(len(testContent)) {
		t.Errorf("GetFile size = %d, want %d", getResp.File.FileSize, len(testContent))
	}
	if getResp.File.ContentType != "text/plain" {
		t.Errorf("GetFile content type = %q, want %q", getResp.File.ContentType, "text/plain")
	}

	// List
	listResp, err := svc.ListFiles(ctx, &v1.ListFilesRequest{Path: "test/"})
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	found := false
	for _, f := range listResp.Items {
		if f.Path == testKey {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ListFiles did not return uploaded file %s", testKey)
	}

	// Delete
	_, err = svc.DeleteFile(ctx, &v1.DeleteFileRequest{Path: testKey})
	if err != nil {
		t.Fatalf("DeleteFile: %v", err)
	}

	// Verify deleted
	_, err = svc.GetFile(ctx, &v1.GetFileRequest{Path: testKey})
	if err == nil {
		t.Error("GetFile after delete should have failed")
	}
}

func TestR2PresignedURL(t *testing.T) {
	client := testR2Client(t)
	ctx := context.Background()
	svc := NewR2FileStoreService(client)

	testKey := fmt.Sprintf("test/ci-presign-%d.txt", time.Now().UnixNano())
	testContent := []byte("presign test")

	// Upload
	_, err := svc.PutFile(ctx, &v1.PutFileRequest{
		File:    &v1.File{Path: testKey, ContentType: "text/plain"},
		Content: testContent,
	})
	if err != nil {
		t.Fatalf("PutFile: %v", err)
	}
	defer svc.DeleteFile(ctx, &v1.DeleteFileRequest{Path: testKey})

	// Get with signed URLs
	getResp, err := svc.GetFile(ctx, &v1.GetFileRequest{
		Path:              testKey,
		IncludeSignedUrls: true,
	})
	if err != nil {
		t.Fatalf("GetFile: %v", err)
	}

	if len(getResp.File.SignedUrls) == 0 {
		t.Error("Expected signed URLs, got none")
	}
	for duration, url := range getResp.File.SignedUrls {
		if url == "" {
			t.Errorf("Signed URL for %s is empty", duration)
		}
		t.Logf("Signed URL (%s): %s...%s", duration, url[:40], url[len(url)-20:])
	}
}
