package zstorager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/denghuo98/zfast/zconsts"
	webSchema "github.com/denghuo98/zfast/zschema/zweb"
)

func TestCosDriveUploadIntegration(t *testing.T) {
	if os.Getenv("ZZ_STORAGE_INTEGRATION") != "1" {
		t.Skip("set ZZ_STORAGE_INTEGRATION=1 to run COS integration test")
	}

	secretId := os.Getenv("ZZ_COS_SECRET_ID")
	secretKey := os.Getenv("ZZ_COS_SECRET_KEY")
	bucketURL := normalizeCosBucketURL(os.Getenv("ZZ_COS_BUCKET_URL"))
	cosPath := os.Getenv("ZZ_COS_PATH")
	if cosPath == "" {
		cosPath = "zzframe-test/attachment/"
	}
	if secretId == "" || secretKey == "" || bucketURL == "" {
		t.Fatal("ZZ_COS_SECRET_ID, ZZ_COS_SECRET_KEY and ZZ_COS_BUCKET_URL are required")
	}

	SetConfig(&webSchema.UploadConfig{
		Drive:        zconsts.UploadDriveCos,
		CosSecretId:  secretId,
		CosSecretKey: secretKey,
		CosBucketURL: bucketURL,
		CosPath:      cosPath,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file := newUploadFileForTest(t, "zzframe-cos-integration.txt", []byte("zzframe cos integration test\n"))
	fullPath, err := New(zconsts.UploadDriveCos).Upload(ctx, file)
	if err != nil {
		t.Fatalf("upload to COS failed: %v", err)
	}
	if fullPath == "" {
		t.Fatal("upload returned empty object path")
	}
	if !strings.HasPrefix(fullPath, cosPath) {
		t.Fatalf("object path %q does not use configured prefix %q", fullPath, cosPath)
	}

	client := newCosClientForTest(t, bucketURL, secretId, secretKey)
	defer func() {
		if _, err = client.Object.Delete(ctx, fullPath); err != nil {
			t.Logf("cleanup COS object %q failed: %v", fullPath, err)
		}
	}()

	if _, err = client.Object.Head(ctx, fullPath, nil); err != nil {
		t.Fatalf("uploaded object HEAD failed: %v", err)
	}

	if os.Getenv("ZZ_COS_PUBLIC_READ") == "1" {
		publicURL := strings.TrimRight(bucketURL, "/") + "/" + fullPath
		resp, err := http.Get(publicURL)
		if err != nil {
			t.Fatalf("public object GET failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("public object GET returned status %d", resp.StatusCode)
		}
	}

}

func newUploadFileForTest(t *testing.T, filename string, data []byte) *ghttp.UploadFile {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create multipart file failed: %v", err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatalf("write multipart file failed: %v", err)
	}
	if err = writer.Close(); err != nil {
		t.Fatalf("close multipart writer failed: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://example.test/upload", &body)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err = req.ParseMultipartForm(int64(body.Len())); err != nil {
		t.Fatalf("parse multipart form failed: %v", err)
	}
	files := req.MultipartForm.File["file"]
	if len(files) == 0 {
		t.Fatal("multipart file was not parsed")
	}
	return &ghttp.UploadFile{FileHeader: files[0]}
}

func newCosClientForTest(t *testing.T, bucketURL, secretId, secretKey string) *cos.Client {
	t.Helper()

	parsedURL, err := url.Parse(bucketURL)
	if err != nil {
		t.Fatalf("parse COS bucket URL failed: %v", err)
	}
	return cos.NewClient(&cos.BaseURL{BucketURL: parsedURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
}

func normalizeCosBucketURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return strings.TrimRight(raw, "/")
	}
	return fmt.Sprintf("https://%s", strings.TrimRight(raw, "/"))
}
