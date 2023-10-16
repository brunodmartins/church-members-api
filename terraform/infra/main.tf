module "email" {
  source = "./email"
  email  = var.email_sender
}

module "repository" {
  source = "./image_repository"
}