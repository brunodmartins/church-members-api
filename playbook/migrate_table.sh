#!/bin/bash
items=($(aws dynamodb scan --table-name member | jq '.Items[] | .id["S"]'))
for id in "${items[@]}"
do
  item=$(aws dynamodb get-item --table-name member --key "{\"id\":{\"S\":${id}}}" | jq -r '.Item')
  aws dynamodb put-item --table-name member_v3 --item "$item"
done
