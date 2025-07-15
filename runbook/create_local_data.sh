#!/bin/bash

aws dynamodb put-item --table-name church_dev --item "$(cat local_church.json)" --endpoint http://127.0.0.1:4566
aws dynamodb put-item --table-name user_dev --item "$(cat local_user.json)"  --endpoint http://127.0.0.1:4566

echo "
Church Abbreviation LOCAL
Login test_user
Password 1234567890
"