package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/authremote"
)

var (
	BuiltinCert    []byte          // set on server startup
	BuiltinKey     []byte          // set on server startup
	BuiltinTLSCert tls.Certificate // set on server startup
	certmap        = &certmapper{}
	dynamic_certs  = flag.Bool("ge_retrieve_missing_https_certificates", false, "if true, retrieve missing https certificates")
)

type certmapper struct {
	sync.Mutex
	certmap map[string]*tls.Certificate // hostname->cert
}

func (cm *certmapper) ByHostname(hostname string) *tls.Certificate {
	hm := strings.ToLower(hostname)
	cm.Lock()
	if cm.certmap == nil {
		cm.certmap = make(map[string]*tls.Certificate)
	}
	tcert := cm.certmap[hm]
	if tcert != nil {
		cm.Unlock()
		return tcert
	}
	cm.Unlock()
	ctx := authremote.Context()
	cert, err := certmanager.GetCertManagerClient().GetLocalCertificate(ctx, &certmanager.LocalCertificateRequest{Subject: hm})
	if err != nil {
		fmt.Printf("[go-easyops] certs failed to get cert for \"%s\": %s\n", hm, err)
		return nil
	}

	cm.Lock()
	tcert = cm.certmap[hm]
	if tcert != nil {
		cm.Unlock()
		return tcert
	}

	// convert into tls certificate
	tc, err := tls.X509KeyPair([]byte(cert.PemCertificate), []byte(cert.PemPrivateKey))
	if err != nil {
		fmt.Printf("[go-easyops] certs Failed to parse cert %s: %s\n", hm, err)
		cm.Unlock()
		return nil
	}
	// add the ca:
	block, _ := pem.Decode([]byte(cert.PemCA))
	if block == nil {
		fmt.Printf("[go-easyops] certs certificate %s has no CA certificate\n", cert.Host)
	} else {
		xcert, xerr := x509.ParseCertificate(block.Bytes)
		if xerr != nil {
			fmt.Printf("[go-easyops] certs Cannot parse certificate %s: %s\n", cert.Host, err)
			cm.Unlock()
			return nil
		}
		now := time.Now()
		if now.After(xcert.NotAfter) {
			fmt.Printf("[go-easyops] certs certificate for \"%s\" expired on %v\n", hm, xcert.NotAfter)
			cm.Unlock()
			return nil
		}

		b := &bytes.Buffer{}
		err = pem.Encode(b, block)
		if err != nil {
			cm.Unlock()
			fmt.Printf("[go-easyops] certs cert for \"%s\" failed to encode: %s\n", hm, err)
			return nil
		}
		tc.Certificate = append(tc.Certificate, block.Bytes)
	}
	fmt.Printf("[go-easyops] certs retrieved certificate for host \"%s\"\n", hm)
	cm.certmap[hm] = &tc
	cm.Unlock()

	return &tc
}

func getcert(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	hostname := chi.ServerName
	if hostname == "rfc-client" { // call from go-easyops
		return &BuiltinTLSCert, nil
	}
	if !*dynamic_certs {
		fmt.Printf("[go-easyops] certs not retrieving for host \"%s\", because flag -ge_retrieve_missing_https_certificates is false\n", hostname)
		return &BuiltinTLSCert, nil
	}
	cert := certmap.ByHostname(hostname)
	if cert != nil {
		return cert, nil
	}
	fmt.Printf("[go-easyops] certs No cert for host \"%s\"\n", hostname)
	return &BuiltinTLSCert, nil
}
