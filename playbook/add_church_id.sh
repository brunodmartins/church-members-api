#!/bin/bash
church_id=""
jq_response=$(aws dynamodb scan --table-name member | jq '.Items[] | .id["S"]')
ids=($jq_response)

for id in "${ids[@]}"
do
  aws dynamodb update-item \
    --table-name member \
    --key "{\"id\":{\"S\":${id}}}" \
    --update-expression "SET church_id = :church_id" \
    --expression-attribute-values "{\":church_id\": {\"S\": \"${church_id}\"}}"
done
