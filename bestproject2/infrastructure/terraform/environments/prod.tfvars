environment = "prod"
aws_region  = "us-east-1"

vpc_cidr = "10.0.0.0/16"

availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]

eks_version = "1.28"

eks_node_groups = {
  general = {
    desired_size   = 3
    min_size       = 2
    max_size       = 10
    instance_types = ["t3.large"]
    capacity_type  = "ON_DEMAND"
  }
  spot = {
    desired_size   = 2
    min_size       = 1
    max_size       = 5
    instance_types = ["t3.large", "t3a.large"]
    capacity_type  = "SPOT"
  }
}

databases = {
  user_management = {
    engine            = "postgres"
    engine_version    = "15.4"
    instance_class    = "db.r6g.xlarge"
    allocated_storage = 500
    db_name           = "usermanagement"
  }
  product_catalog = {
    engine            = "postgres"
    engine_version    = "15.4"
    instance_class    = "db.r6g.xlarge"
    allocated_storage = 500
    db_name           = "productcatalog"
  }
  order_management = {
    engine            = "postgres"
    engine_version    = "15.4"
    instance_class    = "db.r6g.xlarge"
    allocated_storage = 500
    db_name           = "ordermanagement"
  }
}

redis_node_type = "cache.r6g.large"
