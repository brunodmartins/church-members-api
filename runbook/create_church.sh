#!/bin/bash

echo "Church Name:"
read -r name
echo "Church Abbreviation:"
read -r abbreviation
echo "Language:"
read -r language
echo "UUID"
read -r id

item="{\"id\":{\"S\":\"$id\"},\"name\":{\"S\":\"$name\"},\"abbreviation\":{\"S\":\"$abbreviation\"},\"language\":{\"S\":\"$language\"} }"
aws dynamodb put-item --table-name church --item "$item"
