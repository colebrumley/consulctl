package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var certDirectories = []string{
	"/system/etc/security/cacerts",     // Android
	"/usr/local/share/ca-certificates", // Debian derivatives
	"/etc/pki/ca-trust/source/anchors", // RedHat derivatives
	"/etc/ca-certificates",             // Misc alternatives
	"/usr/share/ca-certificates",       // Misc alternatives
}

type AppConfig struct {
	client    *api.Client
	queryOpts *api.QueryOptions
	writeOpts *api.WriteOptions
}

func addCACert(path string, roots *x509.CertPool) *x509.CertPool {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Could not open CA cert: %v", err)
		return roots
	}

	fBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Failed to read CA cert: %v", err)
		return roots
	}

	if !roots.AppendCertsFromPEM(fBytes) {
		log.Printf("Could not add client CA to pool: %v", err)
	}
	return roots
}

func loadSystemRootCAs() (systemRoots *x509.CertPool, err error) {
	systemRoots = x509.NewCertPool()

	for _, directory := range certDirectories {
		fis, err := ioutil.ReadDir(directory)
		if err != nil {
			continue
		}
		for _, fi := range fis {
			data, err := ioutil.ReadFile(directory + "/" + fi.Name())
			if err == nil && systemRoots.AppendCertsFromPEM(data) {
				log.Printf("Loaded Root CA %s", fi.Name())
			}
		}
	}

	return
}

func NewAppConfig(c *cli.Context) (cfg *AppConfig, err error) {
	// Start with the default Consul API config
	config := api.DefaultConfig()

	// Create a TLS config to be populated with flag-defined certs if applicable
	tlsConf := &tls.Config{}

	consulUrl, err := url.Parse(c.GlobalString("addr"))
	if err != nil {
		log.Errorf("Invalid Consul URL: %v", err)
		return
	}

	config.Scheme = consulUrl.Scheme
	config.Address = consulUrl.Host
	config.Datacenter = c.GlobalString("datacenter")

	// Check for insecure flag
	if c.GlobalBool("insecure") {
		tlsConf.InsecureSkipVerify = true
	}

	// Load default system root CAs
	// ignore errors since the TLS config
	// will only be applied if SSL is used
	tlsConf.ClientCAs, _ = loadSystemRootCAs()

	// If --cert and --key are defined, load them and apply the TLS config
	if len(c.GlobalString("cert")) > 0 && len(c.GlobalString("key")) > 0 {
		// Make sure scheme is HTTPS when certs are used, regardless of the flag
		config.Scheme = "https"

		// Load cert and key files
		cert, err := tls.LoadX509KeyPair(c.GlobalString("cert"), c.GlobalString("key"))
		if err != nil {
			log.Errorf("Could not parse SSL cert: %v", err)
		}
		tlsConf.Certificates = append(tlsConf.Certificates, cert)

		// If cacert is defined, add it to the cert pool
		// else just use system roots
		if len(c.GlobalString("cacert")) > 0 {
			tlsConf.ClientCAs = addCACert(c.GlobalString("cacert"), tlsConf.ClientCAs)
			tlsConf.RootCAs = tlsConf.ClientCAs
		}
	}

	if config.Scheme == "https" {
		// Set Consul's transport to the TLS config
		config.HttpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConf,
		}
	}

	// Check for HTTP auth flags
	if len(c.GlobalString("user")) > 0 && len(c.GlobalString("pass")) > 0 {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: c.GlobalString("user"),
			Password: c.GlobalString("pass"),
		}
	}

	// Generate and return the API client
	cl, err := api.NewClient(config)
	cfg = &AppConfig{
		client: cl,
		queryOpts: &api.QueryOptions{
			Datacenter: c.GlobalString("datacenter"),
			Token:      c.GlobalString("token"),
		},
		writeOpts: &api.WriteOptions{
			Datacenter: c.GlobalString("datacenter"),
			Token:      c.GlobalString("token"),
		},
	}
	return
}
