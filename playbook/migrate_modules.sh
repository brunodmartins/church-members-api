#!/bin/bash
cd ../terraform
terraform state mv aws_dynamodb_table.member-table module.dynamodb.aws_dynamodb_table.member_table
terraform state mv aws_dynamodb_table.member-history-table module.dynamodb.aws_dynamodb_table.member_history_table
terraform state mv aws_iam_policy.church-members-api-policy module.iam.aws_iam_policy.church_members_api_policy
terraform state mv aws_iam_role.church-members-api-role module.iam.aws_iam_role.church_members_api_role
terraform state mv aws_iam_role_policy_attachment.attach-policy module.iam.aws_iam_role_policy_attachment.attach_policy
terraform state mv aws_ecr_repository.repo module.ecr.aws_ecr_repository.repo
terraform state mv aws_lambda_function.lambda module.lambda.aws_lambda_function.lambda
terraform state mv aws_api_gateway_rest_api.api-gateway module.gateway.aws_api_gateway_rest_api.api_gateway
terraform state mv aws_api_gateway_deployment.api-deployment module.gateway.aws_api_gateway_deployment.api_deployment
terraform state mv aws_api_gateway_stage.api-stage module.gateway.aws_api_gateway_stage.api_stage
terraform state mv aws_lambda_permission.policy-post-members-create module.gateway.aws_lambda_permission.policy_post_members_create
terraform state mv aws_lambda_permission.policy-get-members module.gateway.aws_lambda_permission.policy_get_members
terraform state mv aws_lambda_permission.policy-put-members module.gateway.aws_lambda_permission.policy_put_members
terraform state mv aws_lambda_permission.policy-post-members module.gateway.aws_lambda_permission.policy_post_members
terraform state mv aws_lambda_permission.policy-get-reports-members module.gateway.aws_lambda_permission.policy_get_reports_members
terraform state mv aws_lambda_permission.policy-get-reports-members-birthday module.gateway.aws_lambda_permission.policy_get_reports_members_birthday
terraform state mv aws_lambda_permission.policy-get-reports-members-marriage module.gateway.aws_lambda_permission.policy_get_reports_members_marriage
terraform state mv aws_lambda_permission.policy-get-reports-members-legal module.gateway.aws_lambda_permission.policy_get_reports_members_legal
terraform state mv aws_lambda_permission.policy-get-reports-members-classification module.gateway.aws_lambda_permission.policy_get_reports_members_classification
terraform state mv aws_api_gateway_authorizer.authorizer module.gateway.aws_api_gateway_authorizer.authorizer
terraform state mv aws_cognito_user_pool.user-pool module.cognito.aws_cognito_user_pool.user_pool
