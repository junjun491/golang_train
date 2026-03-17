module "ecr" {
  source = "../../infrastructure_modules/ecr"

  repositories = [
    "otayori-frontend",
    "otayori-backend-go",
  ]

  tags = {
    Project     = "otayori"
    Environment = "dev"
  }
}
