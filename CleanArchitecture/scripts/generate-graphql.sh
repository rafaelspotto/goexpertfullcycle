#!/bin/bash

# Install gqlgen
go install github.com/99designs/gqlgen@latest

# Generate GraphQL code
gqlgen generate 