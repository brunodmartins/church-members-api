locals {
  ecr_image_tag       = "latest"
  root                = "../../"
  ecr_repository_name = "church-members-api-container"
}

resource "aws_ecr_repository" "repo" {
  name = local.ecr_repository_name
}

data "archive_file" "internal" {
  type        = "zip"
  source_dir  = "${local.root}internal/"
  output_path = "internal.zip"
}

data "archive_file" "cmd" {
  type        = "zip"
  source_dir  = "${local.root}cmd/"
  output_path = "cmd.zip"
}


data "aws_ecr_image" "lambda_image" {
  repository_name = local.ecr_repository_name
  image_tag       = local.ecr_image_tag
}

output "image_id" {
  value = "${aws_ecr_repository.repo.repository_url}:latest"
}