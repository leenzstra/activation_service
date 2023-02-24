#!/usr/bin/env bash

go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

packagelist=(
    ".\cmd"
    ".\internal\utils"
    ".\internal\responses"
    ".\internal\models"                  
    ".\internal\middlewares"
    ".\internal\keypair" 
    ".\internal\db" 
    ".\internal\config" 
    ".\internal\collections" 
    ".\internal\api\content" 
    ".\internal\api\info"
    ".\internal\api\license" 
)


gomarkdoc --output ".\docs\doc.md" ${packagelist[@]}
                            
                            
