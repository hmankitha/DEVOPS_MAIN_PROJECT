variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "ecommerce"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "eks_version" {
  description = "EKS cluster version"
  type        = string
  default     = "1.28"
}

variable "eks_node_groups" {
  description = "EKS node groups configuration"
  type = map(object({
    desired_size   = number
    min_size       = number
    max_size       = number
    instance_types = list(string)
    capacity_type  = string
  }))
  default = {
    general = {
      desired_size   = 3
      min_size       = 2
      max_size       = 5
      instance_types = ["t3.large"]
      capacity_type  = "ON_DEMAND"
    }
  }
}

variable "databases" {
  description = "RDS databases configuration"
  type = map(object({
    engine         = string
    engine_version = string
    instance_class = string
    allocated_storage = number
    db_name        = string
  }))
  default = {
    user_management = {
      engine            = "postgres"
      engine_version    = "15.4"
      instance_class    = "db.t3.medium"
      allocated_storage = 100
      db_name           = "usermanagement"
    }
    product_catalog = {
      engine            = "postgres"
      engine_version    = "15.4"
      instance_class    = "db.t3.medium"
      allocated_storage = 100
      db_name           = "productcatalog"
    }
    order_management = {
      engine            = "postgres"
      engine_version    = "15.4"
      instance_class    = "db.t3.medium"
      allocated_storage = 100
      db_name           = "ordermanagement"
    }
  }
}

variable "s3_buckets" {
  description = "S3 buckets to create"
  type        = list(string)
  default     = ["ecommerce-users", "ecommerce-products", "ecommerce-orders"]
}

variable "secrets" {
  description = "Secrets to store in AWS Secrets Manager"
  type        = map(string)
  default     = {}
  sensitive   = true
}

variable "redis_node_type" {
  description = "Redis node type"
  type        = string
  default     = "cache.t3.medium"
}
