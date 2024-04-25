#!/bin/bash

aws dynamodb put-item --table-name church_dev --item "$(cat runbook/local_church.json)" --endpoint http://127.0.0.1:4566
aws dynamodb put-item --table-name user_dev --item "$(cat runbook/local_user.json)"  --endpoint http://127.0.0.1:4566