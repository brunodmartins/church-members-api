locals {
  ecr_image_tag       = "latest"
  ecr_repository_name = "church-members-api-container"
}

resource "aws_ecr_repository" "repo" {
  name = local.ecr_repository_name
}

data "aws_ecr_image" "lambda_image" {
  repository_name = local.ecr_repository_name
  image_tag       = local.ecr_image_tag
}

output "image_id" {
  value = "${aws_ecr_repository.repo.repository_url}:latest"
}