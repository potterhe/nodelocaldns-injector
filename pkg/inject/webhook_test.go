package inject

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInject(t *testing.T) {

	whParams := WebhookParameters{}
	wh, _ := NewWebhook(whParams)

	postJson := `
	{
		"apiVersion": "admission.k8s.io/v1beta1",
		"kind": "AdmissionReview",
		"request": {
		  "uid": "705ab4f5-6393-11e8-b7cc-42010a800002"
		}
	  }
	`
	reader := strings.NewReader(postJson)
	r := httptest.NewRequest(http.MethodPost, "/inject", reader)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wh.serveInject(w, r)
	t.Log(w)
}
