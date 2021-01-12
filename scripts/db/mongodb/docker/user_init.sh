#!/bin/bash
set -e

# This script creates a new non-root mongodb user with readWrite access to the initialized database

if [ -n "${MONGO_INITDB_ROOT_USERNAME:-}" ] && [ -n "${MONGO_INITDB_ROOT_PASSWORD:-}" ] && [ -n "${DATABASE_USER:-}" ] && [ -n "${DATABASE_PASSWORD:-}" ] && [ -n "${MONGO_INITDB_DATABASE:-}" ]; then
mongo -u $MONGO_INITDB_ROOT_USERNAME -p $MONGO_INITDB_ROOT_PASSWORD<<EOF
db=db.getSiblingDB('$MONGO_INITDB_DATABASE');
db.createUser({
  user:  '$DATABASE_USER',
  pwd: '$DATABASE_PASSWORD',
  roles: [{
    role: 'readWrite',
    db: '$MONGO_INITDB_DATABASE'
  }]
});
EOF
else
    echo "Failed to initialize database. Ensure that environment variables MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD, MONGO_INITDB_DATABASE, DATABASE_USER, and DATABASE_PASSWORD are supplied."
    exit 1
fi