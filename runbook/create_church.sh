#!/bin/bash

echo "Church Name:"
read -r name
echo "Church Abbreviation:"
read -r abbreviation
echo "Language:"
read -r language

id=$(uuidgen | tr '[:upper:]' '[:lower:]')
item="{\"id\":{\"S\":\"$id\"},\"name\":{\"S\":\"$name\"},\"abbreviation\":{\"S\":\"$abbreviation\"},\"language\":{\"S\":\"$language\"} }"
aws dynamodb put-item --table-name church --item "$item"
