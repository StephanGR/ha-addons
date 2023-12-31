package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mdlayher/wol"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Domains []DomainConfig `json:"domains"`
}

type DomainConfig struct {
	Url              string   `json:"url"`
	MacAddress       string   `json:"macAddress"`
	BroadcastAddress string   `json:"broadcastAddress"`
	WakeUpPort       int      `json:"wakeUpPort"`
	WakeUpIp         string   `json:"wakeUpIp"`
	ForwardIp        string   `json:"forwardIp"`
	ForwardPort      int      `json:"forwardPort"`
	WakeUpEndpoints  []string `json:"wakeUpEndpoints"`
}

type ServerState struct {
	wakingUp bool
	lock     sync.Mutex
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}

func logRequest(logger *logrus.Logger, r *http.Request) {

	logger.WithFields(logrus.Fields{
		"client":     r.Header.Get("X-Forwarded-For"), // Replace the value of client by X-Forwarded-For
		"method":     r.Method,
		"user-agent": r.UserAgent(),
		"path":       r.URL.Path,
	}).Info()
}

func (s *ServerState) IsWakingUp() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.wakingUp
}

func (s *ServerState) StartWakingUp() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.wakingUp {
		return false // Already waking up
	}
	s.wakingUp = true
	return true
}

func (s *ServerState) DoneWakingUp() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.wakingUp = false
}

func loadConfig(configPath string) (*Config, error) {
	var config Config
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func isServerUp(logger *logrus.Logger, address string) bool {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logger.Warn("Failed to connect:", err) // Use your logger to log this error appropriately.
		return false
	}
	conn.Close()
	return true
}

func wakeServer(logger *logrus.Logger, macAddress string, broadcastAddress string, serverState *ServerState) {
	if !serverState.StartWakingUp() {
		logger.Info("There is already a wake up in progress")
		return
	}
	defer serverState.DoneWakingUp()

	client, err := wol.NewClient()
	if err != nil {
		logger.Warn("Error when creating WOL client : %v", err)
		return
	}
	defer func(client *wol.Client) {
		err := client.Close()
		if err != nil {
			logger.Warn("Unable to close the WOL client")
		}
	}(client)

	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		logger.Warn("Invalid mac address : %v", err)
		return
	}
	if err := client.Wake(broadcastAddress, mac); err != nil {
		logger.Warn("Error when sending magic packet : %v", err)
	} else {
		logger.Info("Magic packet sent")
	}
}

func handleDomainProxy(w http.ResponseWriter, r *http.Request, domain DomainConfig) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", domain.ForwardIp, domain.ForwardPort),
	})

	r.URL.Host = fmt.Sprintf("%s:%d", domain.ForwardIp, domain.ForwardPort)
	r.URL.Scheme = "http"
	r.Host = fmt.Sprintf("%s:%d", domain.ForwardIp, domain.ForwardPort)
	proxy.ServeHTTP(w, r)
}

func handler(logger *logrus.Logger, w http.ResponseWriter, r *http.Request, config *Config, serverState *ServerState) {
	domainConfig, found := findDomainConfig(config.Domains, r.Host)
	if !found {
		http.Error(w, "Domain not configured", http.StatusNotFound)
		return
	}

	if shouldWakeServer(r.URL.Path, domainConfig.WakeUpEndpoints) {
		serverAddress := fmt.Sprintf("%s:%d", domainConfig.WakeUpIp, domainConfig.WakeUpPort)
		if !isServerUp(logger, serverAddress) {
			logger.Info("Server is offline, trying to wake up using Wake On Lan")
			wakeServer(logger, domainConfig.MacAddress, domainConfig.BroadcastAddress, serverState)

			waitServerOnline(logger, serverAddress, w, r)
			return
		} else {
			handleDomainProxy(w, r, *domainConfig)
		}
	} else {
		handleDomainProxy(w, r, *domainConfig)
	}
}

func waitServerOnline(logger *logrus.Logger, serverAddress string, w http.ResponseWriter, r *http.Request) {
	timeout := time.After(1 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if isServerUp(logger, serverAddress) {
				logger.Info("Server is up !")
				return
				//externalUrl := fmt.Sprintf("http://%s", r.Host)
				//http.Redirect(w, r, externalUrl, http.StatusSeeOther)
				//return
			}
			//else {
			//	logger.Info("Waiting for server to wake up...")
			//}
		case <-timeout:
			logger.Info("Timeout reached, server did not wake up.")
			fmt.Fprintf(w, "Server is still offline. Please try again later.")
			return
		}
	}
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func requestLoggerMiddleware(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(logger, r)
		next.ServeHTTP(w, r)
	})
}

func findDomainConfig(domains []DomainConfig, host string) (*DomainConfig, bool) {
	for _, domain := range domains {
		parsedUrl, err := url.Parse(domain.Url)
		if err != nil {
			continue
		}

		domainName := parsedUrl.Host
		if domainName == "" {
			domainName = domain.Url
		}

		if domainName == host {
			return &domain, true
		}
	}
	return nil, false
}

func shouldWakeServer(endpoint string, wakeUpEndpoints []string) bool {
	for _, pattern := range wakeUpEndpoints {
		if pattern == endpoint || matchesPattern(endpoint, pattern) {
			return true
		}
	}
	return false
}

func matchesPattern(endpoint, pattern string) bool {
	splitPattern := strings.Split(pattern, "*")
	if len(splitPattern) != 2 {
		return false // ou retourner une erreur si le motif n'est pas valide
	}

	prefix, suffix := splitPattern[0], splitPattern[1]
	return strings.HasPrefix(endpoint, prefix) && strings.HasSuffix(endpoint, suffix)
}

func main() {
	logger := initLogger()
	serverState := &ServerState{}
	config, err := loadConfig("/config.json")
	if err != nil {
		logger.Fatal("Error loading config file: ", err)
	}

	logger.Info("Configuration successfully loaded")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(logger, w, r, config, serverState)
	})

	mux.HandleFunc("/ping", PingHandler)

	loggedMux := requestLoggerMiddleware(logger, mux)

	logger.Info("Starting app..")
	logger.Fatal(http.ListenAndServe(":3881", loggedMux))
}
