#!/bin/bash
echo Building domainfinder...
go build -o domainfinder.bin
echo Building synonyms...
cd ../synonyms
go build -o ../domainfinder/lib/synonyms.bin
echo Building available...
cd ../available
go build -o ../domainfinder/lib/available.bin
echo Building sprinkle...
cd ../sprinkle
go build -o ../domainfinder/lib/sprinkle.bin
echo Building coolify...
cd ../coolify
go build -o ../domainfinder/lib/coolify.bin
echo Building domainify...
cd ../domainify
go build -o ../domainfinder/lib/domainify.bin
