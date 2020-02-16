#!/bin/bash

set -euo pipefail

source .env

sed -i  's/<MONGO_USER>/'"$MONGO_USER"'/g' app.yaml
sed -i  's/<MONGO_PASSWORD>/'"$MONGO_PASSWORD"'/g' app.yaml

gcloud app deploy -q


sed -i  's/'"$MONGO_USER"'/<MONGO_USER>/g' app.yaml
sed -i  's/'"$MONGO_PASSWORD"'/<MONGO_PASSWORD>/g' app.yaml
