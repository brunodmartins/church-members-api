#!/bin/bash
items=($(aws dynamodb scan --table-name member | jq '.Items[] | .id["S"]'))
for id in "${items[@]}"
do
  item=$(aws dynamodb get-item --table-name member --key "{\"id\":{\"S\":${id}}}" | jq -r '.Item')
  maritalStatus="MARRIED"
  single=($(echo "$item" | jq '. | {marriageDateShort}' | grep '"NULL": true'))
  if [ -n "$single" ]; then
    maritalStatus="SINGLE"
  fi
  item=`echo "$item" | jq --arg status "$maritalStatus" '. + {"maritalStatus":{"S":$status}}'`
  aws dynamodb put-item --table-name member_v2 --item "$item"
done