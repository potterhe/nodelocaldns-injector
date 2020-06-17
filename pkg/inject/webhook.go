package inject

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

// WebhookParameters configures parameters for the dnsConfig injection webhook
type WebhookParameters struct {

	// CertFile is the path to the x509 certificate for https.
	CertFile string

	// KeyFile is the path to the x509 private key matching `CertFile`.
	KeyFile string

	// Port is the webhook port, e.g. typically 443 for https.
	Port int
}

// Webhook implements a mutating webhook for automatic dnsConfig injection.
type Webhook struct {
	server *http.Server
}

// NewWebhook creates a new instance of a mutating webhook for automatic dnsConfig injection.
func NewWebhook(p WebhookParameters) (*Webhook, error) {

	sCert, err := tls.LoadX509KeyPair(p.CertFile, p.KeyFile)
	if err != nil {
		return nil, err
	}

	wh := &Webhook{
		server: &http.Server{
			Addr: fmt.Sprintf(":%v", p.Port),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{sCert},
				// TODO: uses mutual tls after we agree on what cert the apiserver should use.
				// ClientAuth:   tls.RequireAndVerifyClientCert,
			},
		},
	}

	http.HandleFunc("/inject", wh.serveInject)
	return wh, nil
}

func (wh *Webhook) serveInject(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	glog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1beta1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		glog.Error(err)
		responseAdmissionReview.Response = toAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = toAdmissionResponse(errors.New("test"))
		//responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	glog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		glog.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		glog.Error(err)
	}
}

// Serve start the webhook server
func (wh *Webhook) Serve() error {
	err := wh.server.ListenAndServeTLS("", "")
	return err
}

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
