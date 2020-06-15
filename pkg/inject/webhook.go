package inject

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type WebhookParameters struct {

	// CertFile is the path to the x509 certificate for https.
	CertFile string

	// KeyFile is the path to the x509 private key matching `CertFile`.
	KeyFile string

	// Port is the webhook port, e.g. typically 443 for https.
	Port int
}

type Webhook struct {
	server *http.Server
}

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
	w.Write([]byte("hello"))
}

func (wh *Webhook) Serve() error {
	err := wh.server.ListenAndServeTLS("", "")
	return err
}
