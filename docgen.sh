#!/usr/bin/env bash

go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

gomarkdoc --output ".\docs\doc.md" |
                            ".\cmd" |
                            ".\internal\utils" |
                            ".\internal\responses" |
                            ".\internal\models" |
                            ".\internal\middlewares" |
                            ".\internal\keypair" |
                            ".\internal\db" |
                            ".\internal\config" |
                            ".\internal\collections" |
                            ".\internal\api\content" |
                            ".\internal\api\info" |
                            ".\internal\api\license" 
                            
