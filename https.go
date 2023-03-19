package https

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type Service struct {
	manager *autocert.Manager
}

// New creates autocert service
// cacheDir should be a path to the folder when you want to store your certificates
// (to save them when the app stops)
// email is used only for Let's Encrypt registration purposes
// hosts is the list of domain names of your site (could be just one)
func New(cacheDir string, email string, hosts ...string) *Service {
	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache(cacheDir),
		Email:      email,
	}

	return &Service{
		manager: manager,
	}
}

// GetHTTPSServer should be used when you want some manual control over created https server
// after getting the server and applying your settings just call .ListenAndServeTLS("", "")
func (s *Service) GetHTTPSServer(hostPort string, mux http.Handler) *http.Server {
	httpsSrv := makeServerFromMux(mux)
	httpsSrv.Addr = hostPort
	httpsSrv.TLSConfig = &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
		MinVersion: tls.VersionTLS12,
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := s.manager.GetCertificate(hello)
			if err != nil {
				log.Println("GetCertificate error:", err)
			}
			return cert, err
		},
	}
	return httpsSrv
}

// ListenHTTPS is the easiest way to start your https server
// hostPort may look like :443 (listens all interfaces) or 1.2.3.4:443 (listens only provided ip)
// mux should be a router of your choice (for example github.com/julienschmidt/httprouter)
func (s *Service) ListenHTTPS(hostPort string, mux http.Handler) error {
	httpsSrv := s.GetHTTPSServer(hostPort, mux)
	return httpsSrv.ListenAndServeTLS("", "")
}

// ListenHTTPRedirect listens plain http and redirects all requests to https version of the site
// works only when your https server listens at 443 port
func (s *Service) ListenHTTPRedirect(listenHostPort string) error {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	httpSrv := makeServerFromMux(mux)

	// allow autocert handle Let's Encrypt callbacks over http
	httpSrv.Handler = s.manager.HTTPHandler(httpSrv.Handler)

	httpSrv.Addr = listenHostPort
	return httpSrv.ListenAndServe()
}

func makeServerFromMux(mux http.Handler) *http.Server {
	return &http.Server{
		IdleTimeout: 10 * time.Second,
		Handler:     mux,
	}
}
