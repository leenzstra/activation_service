#!/usr/bin/env bash

go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

sudo apt install markdown

ls

packagelist=(
    "./cmd"
    "./internal/utils"
    "./internal/responses"
    "./internal/models"                  
    "./internal/middlewares/access"
    "./internal/keypair" 
    "./internal/db" 
    "./internal/config" 
    "./internal/collections" 
    "./internal/api/content" 
    "./internal/api/info"
    "./internal/api/license" 
)

gomarkdoc --output "./docs/doc.md" ${packagelist[@]}

markdown ./docs/doc.md > ./docs/README.md
                            
                            
