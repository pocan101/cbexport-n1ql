# Couchbase Query Executor

This project is a Go-based tool for executing N1QL queries against a Couchbase cluster and exporting the results to JSON files.

------

## Features

- Connects to a Couchbase cluster using TLS if a CA certificate is provided.
- Executes multiple N1QL queries defined in a configuration file.
- Exports query results to JSON files for easy processing or sharing.

------

## Prerequisites

1. Couchbase Server with appropriate permissions for the provided user credentials.
2. Go installed (version 1.20+ recommended).
3. A `config.yml` file (sample provided below).

------

## Usage

### Building the Application

1. Clone the repository or download the source code.

2. Run the following command to build the binary:

   ```
   
   go build -o cbQueryExecutor main.go
   ```

### Running the Application

1. Execute the binary with the configuration file as an argument:

   ```
   
   ./cbQueryExecutor config.yml
   ```

------

## Configuration File (`config.yml`)

### Example

```yaml

connection_details:
  user: Administrator
  password: password
  url: couchbases://ec2-3-88-72-54.compute-1.amazonaws.com
  ca_certificate:
    enabled: true
    name: "ca_certificate.pem" # Name of the generated file
    content: |
      -----BEGIN CERTIFICATE-----
      MIIDDDCCAfSgAwIBAgIIGA4YNVtbAwYwDQYJKoZIhvcNAQELBQAwJDEiMCAGA1UE
      AxMZQ291Y2hiYXNlIFNlcnZlciAzNzJlZWY3MTAeFw0xMzAxMDEwMDAwMDBaFw00
      OTEyMzEyMzU5NTlaMCQxIjAgBgNVBAMTGUNvdWNoYmFzZSBTZXJ2ZXIgMzcyZWVm
      NzEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDeMWK3CWJcbxXqD818
      8JMIEz689ByxQpClCTqQn3IRLpQufV5EWOKoeQrIrLoxfBHy93j0LOaR57eztdXr
      8R7SN4TNwK/hocjpbtjTic2nzAecMAQqhs8mc9aPNDtbJvxS/0IyxPnyEK0dN1ib
      gUKX38NBW5uiAU5bugU+4mZX/KMFuI+DMvnKqItANB+Q7Opcwtna1Ke121zFntWV
      z0TYReFgW4lgirvoC0gxmi61E0Jtr0ZXSMLv0L2SP5kLelKqxjYU4tmgdJvrrHD1
      1pnVwhUrA6XpXy+PYehEUrRHPIa+j/3cv7ng/2D/XY7DS3kMxe4GzWQ89cKatcYZ
      q+3TAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBhjAPBgNVHRMBAf8EBTADAQH/MB0G
      A1UdDgQWBBR8JCP7Kll4RS8Eh6awFFrdS1sGFjANBgkqhkiG9w0BAQsFAAOCAQEA
      yK0TKRTxZXMHFAJ3JBCHOipVPlbCtaRBjCh6UKchcNXZESdA2BbDjlJWwP19x8xG
      sBvpp1GWTkMzuRNilsWBterCnUW/HjPRagAXUZ6dx9gCZoRacMzEfKid6rz+1Drj
      kX2syG3b0rfHcQqc/CMzrLruk+mLobHTxB5CRjtIJoIOMn7+laLwwnC+luZm5Z66
      wGJMaY7BagzZk4GrWQXoeq2IdJyPEX2gnMOljpS7QqyodBseXw0u6RrohX0dhZX2
      nzFSHWXJHYCYY9iYAFjJ8dUpffPBxYMwDHsz7XJwPooa+ylTXXcAGQb/irpEAb3S
      uiXs9TY6QIDez+pxkDH0fw==
      -----END CERTIFICATE-----

query_details:

  - query_name: ratings
    n1ql: SELECT r.airline, r.sourceairport, r.destinationairport, a.name, a.callsign FROM `travel-sample` r JOIN `travel-sample` a ON r.airlineid = META(a).id WHERE r.type = "route" AND a.type = "airline" AND r.sourceairport = "SFO"
    file: query1.json

  - query_name: airports
    n1ql: SELECT airportname, city, country FROM `travel-sample` WHERE type = "airport" AND country = "United States";
    file: query2.json

  - query_name: flights
    n1ql: SELECT flight, sourceairport, destinationairport, distance FROM `travel-sample` WHERE type = "route" AND sourceairport = "LAX" AND destinationairport = "SFO";
    file: query3.json

  - query_name: hotels
    n1ql: SELECT name, address, city, state FROM `travel-sample` WHERE type = "hotel" AND city = "San Francisco";
    file: query4.json
```

------

## Configuration Template Fields

### `connection_details`

- `user`: Username for Couchbase authentication.

- `password`: Password for Couchbase authentication.

- `url`: Connection string for the Couchbase cluster.

- ```
  ca_certificate
  ```

  - `enabled`: Boolean to enable/disable CA certificate validation.
  - `name`: Name of the CA certificate file.
  - `content`: PEM-encoded CA certificate content.

### `query_details`

- `query_name`: Name/description of the query.
- `n1ql`: The N1QL query to execute.
- `file`: The output JSON file for query results.

------

## Output

- JSON files specified in the `file` fields of the `query_details` section will be generated, containing the results of the executed queries.

------

## Example Execution

```

./cbexport-n1ql config.yml
```