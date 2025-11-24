variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "security_group_ids" {
  type = list(string)
}

resource "aws_db_subnet_group" "main" {
  name       = "nexusflow-${var.environment}-db-subnet-group"
  subnet_ids = var.subnet_ids

  tags = {
    Name = "nexusflow-${var.environment}-db-subnet-group"
  }
}

resource "aws_security_group" "rds" {
  name        = "nexusflow-${var.environment}-rds-sg"
  description = "Allow inbound traffic from EKS"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = var.security_group_ids
  }

  tags = {
    Name = "nexusflow-${var.environment}-rds-sg"
  }
}

resource "aws_db_instance" "main" {
  identifier           = "nexusflow-${var.environment}-db"
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "postgres"
  engine_version       = "15.4"
  instance_class       = "db.t3.micro"
  db_name              = "nexusflow"
  username             = "nexusflow"
  password             = "change-me-in-production" # Use Secrets Manager in prod
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  tags = {
    Name = "nexusflow-${var.environment}-db"
  }
}

output "db_endpoint" {
  value = aws_db_instance.main.endpoint
}
