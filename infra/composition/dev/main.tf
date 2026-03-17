data "aws_caller_identity" "current" {}

data "aws_ecr_repository" "backend" {
  name = "otayori-backend-go"
}

data "aws_ecr_repository" "frontend" {
  name = "otayori-frontend"
}

module "network" {
  source = "../../infrastructure_modules/network"

  name       = "otayori-dev"
  cidr_block = "10.0.0.0/16"
  az_count   = 3

  public_subnet_cidrs  = ["10.0.0.0/20", "10.0.16.0/20", "10.0.32.0/20"]
  private_subnet_cidrs = ["10.0.128.0/20", "10.0.144.0/20", "10.0.160.0/20"]

  enable_nat        = true
  nat_gateway_count = 1

  tags = {
    Environment = "dev"
    Project     = "otayori"
  }
}

# module "ecr" {
#   source = "../../infrastructure_modules/ecr"
#
#   repositories = [
#     "otayori-backend-go",
#     "otayori-frontend",
#   ]
#
#   tags = {
#     Environment = "dev"
#     Project     = "otayori"
#   }
# }

module "database" {
  source = "../../infrastructure_modules/database"

  name = "otayori-dev-db"

  vpc_id     = module.network.vpc_id
  subnet_ids = module.network.private_subnet_ids

  tags = {
    Environment = "dev"
    Project     = "otayori"
  }
}

module "alb" {
  source = "../../infrastructure_modules/alb"

  name = "otayori-dev-alb"

  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids

  frontend_port = 3000
  backend_port  = 3001

  tags = {
    Environment = "dev"
    Project     = "otayori"
  }
}

module "ecs" {
  source = "../../infrastructure_modules/ecs"

  name_prefix        = "otayori-dev"
  vpc_id             = module.network.vpc_id
  private_subnet_ids = module.network.private_subnet_ids

  alb_security_group_id = module.alb.security_group_id
  frontend_tg_arn       = module.alb.target_group_arns["frontend"]
  backend_tg_arn        = module.alb.target_group_arns["backend"]
  jwt_secret_arn        = "arn:aws:secretsmanager:ap-northeast-1:102464981360:secret:otayori/dev/go/jwt_secret-mF9KyT"

  backend_image  = "${data.aws_ecr_repository.backend.repository_url}:dev"
  frontend_image = "${data.aws_ecr_repository.frontend.repository_url}:dev"

  database_url = format(
    "postgres://%s:%s@%s:%d/%s",
    module.database.username,
    module.database.password,
    module.database.endpoint,
    module.database.port,
    module.database.database_name,
  )

  rds_security_group_id = module.database.security_group_id

  tags = {
    Environment = "dev"
    Project     = "otayori"
  }
}

# module "iam_github_oidc" {
#   source = "../../infrastructure_modules/iam_github_oidc"
#
#   name_prefix       = "otayori-dev"
#   github_repository = "junjun491/otayori_app:ref:refs/heads/main"
#
#   tags = {
#     Environment = "dev"
#     Project     = "otayori"
#   }
# }