# Go Mongo

Go + Gorilla/Mux + MongoDB-driver Crud Api

# Usage

Run locally using `go run go-mongo.go`

# Secret

Secret are stored in an ignored `.env` file.

An example schema for the file can be found at `.env.example`

# Deployment

Deploy to App Engine through:

```bash
chmod +x scripts/deploy.sh
scripts/deploy.sh
```

**Note: This script requires .env file to exist and contain valid credentials**
