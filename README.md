# Domain Monitoring Tool

## Program description:
The Domain Monitoring Tool is a Go program designed for monitoring the health of specified domains. It checks the availability of endpoints, measures response times, and calculates the availability percentage. The tool reads endpoint information from a YAML file and performs health checks at regular intervals.

# Usage
## Prerequisites
- Ensure that Go is installed on your system  [Ubuntu installation example](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-20-04),[Mac OS installation example](https://www.digitalocean.com/community/tutorials/how-to-install-go-and-set-up-a-local-programming-environment-on-macos)

- Program requires YAML file to read http requests

# Project folder strcuture
```
domain-monitoring-tool
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── sample_inputs
    │   └── sample-input.yml
    └── utils
        └── utils.go

4 directories, 6 files
```

# How to run the program
1. Clone the repo:
```
git clone git@github.com:vraj9845/Domain-Monitoring-Tool.git
cd Domain-Monitoring-Tool
```
2. Paste the below command to run the tool:
```
go run main.go -file ./sample_inputs/sample-input.yml
```


Note: To run a custom file replace `./sample_inputs/sample-input.yml` with `path/to/your/custom-file.yaml` containing endpoint information.
