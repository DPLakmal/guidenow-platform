output "alb_dns_name" {
  description = "Application Load Balancer DNS endpoint"
  value       = aws_lb.app.dns_name
}

output "ecs_cluster_name" {
  description = "ECS cluster identifier"
  value       = aws_ecs_cluster.main.name
}

output "redis_primary_endpoint" {
  description = "ElastiCache Redis primary endpoint"
  value       = aws_elasticache_replication_group.redis.primary_endpoint_address
}

output "rds_endpoint" {
  description = "Amazon RDS PostgreSQL connection endpoint"
  value       = aws_db_instance.postgres.endpoint
}
