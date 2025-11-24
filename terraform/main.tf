provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project     = "NexusFlow"
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

module "vpc" {
  source = "./modules/vpc"

  environment  = var.environment
  vpc_cidr     = var.vpc_cidr
  cluster_name = "nexusflow-${var.environment}"
}

module "eks" {
  source = "./modules/eks"

  environment  = var.environment
  cluster_name = "nexusflow-${var.environment}"
  vpc_id       = module.vpc.vpc_id
  subnet_ids   = module.vpc.private_subnets
}

module "rds" {
  source = "./modules/rds"

  environment        = var.environment
  vpc_id             = module.vpc.vpc_id
  subnet_ids         = module.vpc.private_subnets
  security_group_ids = [module.eks.cluster_security_group_id]
}

module "s3" {
  source = "./modules/s3"

  environment = var.environment
}
