package main

import (

	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/couchbase/gocb/v2"
	"gopkg.in/yaml.v3"
)

type CACertificate struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

type ConnectionDetails struct {
	User     string       `yaml:"user"`
	Password string       `yaml:"password"`
	URL      string       `yaml:"url"`
	CA       CACertificate `yaml:"ca_certificate"`
}

type QueryDetails struct {
	QueryName string `yaml:"query_name"`
	N1QL      string `yaml:"n1ql"`
	File      string `yaml:"file"`
}

type Configuration struct {
	ConnectionDetails ConnectionDetails `yaml:"connection_details"`
	QueryDetails      []QueryDetails    `yaml:"query_details"`
}

func readConfig(filePath string) (Configuration, error) {
	var config Configuration
	file, err := os.Open(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("failed to parse YAML: %v", err)
	}

	// Handle CA certificate if enabled
	if config.ConnectionDetails.CA.Enabled {
		err = os.WriteFile(config.ConnectionDetails.CA.Name, []byte(config.ConnectionDetails.CA.Content), 0644)
		if err != nil {
			return config, fmt.Errorf("failed to write CA certificate: %v", err)
		}
		log.Printf("CA certificate written to: %s", config.ConnectionDetails.CA.Name)
	}

	return config, nil
}

func executeAndDump(config Configuration, cluster *gocb.Cluster) {
	for _, query := range config.QueryDetails {
		log.Printf("Executing query: %s", query.QueryName)
		outFile, err := os.Create(query.File)
		if err != nil {
			log.Printf("Failed to create output file '%s': %v", query.File, err)
			continue
		}
		defer outFile.Close()

		result, err := cluster.Query(query.N1QL, &gocb.QueryOptions{})
		if err != nil {
			log.Printf("Query execution failed for '%s': %v", query.QueryName, err)
			continue
		}

		for result.Next() {
			var row map[string]interface{}
			err = result.Row(&row)
			if err != nil {
				log.Printf("Failed to parse query result row: %v", err)
				continue
			}
			jsonData, _ := json.Marshal(row)
			outFile.WriteString(string(jsonData) + "\n")
		}

		if err := result.Err(); err != nil {
			log.Printf("Query result iteration error for '%s': %v", query.QueryName, err)
		}

		log.Printf("Results written to: %s", query.File)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <config.yml>\n", os.Args[0])
		os.Exit(1)
	}

	configPath := os.Args[1]
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("Configuration loaded:")
	fmt.Printf("%+v\n", config)

	authenticator := gocb.PasswordAuthenticator{
		Username: config.ConnectionDetails.User,
		Password: config.ConnectionDetails.Password,
	}

	clusterOpts := gocb.ClusterOptions{
		Authenticator: authenticator,
	}

	// Handle CA certificate if enabled
	if config.ConnectionDetails.CA.Enabled {
		tlsConfig := &tls.Config{}
		caCert := []byte(config.ConnectionDetails.CA.Content)
		if !strings.Contains(string(caCert), "-----BEGIN CERTIFICATE-----") {
			log.Fatalf("Invalid CA certificate content")
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		clusterOpts.SecurityConfig.TLSRootCAs = caCertPool
	}

	cluster, err := gocb.Connect(config.ConnectionDetails.URL, clusterOpts)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}
	defer cluster.Close(nil)

	executeAndDump(config, cluster)
}
