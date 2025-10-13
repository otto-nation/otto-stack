---
title: "Troubleshooting"
description: "Common issues and solutions for otto-stack problems"
lead: "Quick solutions to the most common otto-stack issues and problems"
date: "2025-10-01"
lastmod: "2025-10-11"
draft: false
weight: 70
toc: true
---

# Troubleshooting Guide

## 🚨 Top 5 Issues (Quick Reference)

1. **Docker Not Running**
   - Run `docker info` to check status.
   - Start Docker/Colima if not running.

2. **Port Conflicts**
   - Run `lsof -i :PORT_NUMBER` to find conflicts.
   - Use `--cleanup-existing` to resolve.

3. **Service Won't Start**
   - Check logs: `otto-stack logs SERVICE_NAME`
   - Restart service: `otto-stack down && otto-stack up`

4. **Memory Issues**
   - Check usage: `docker stats`
   - Increase Docker memory limit.

5. **Invalid Configuration**
   - Validate YAML: `otto-stack doctor`
   - Run setup with debug: `otto-stack --verbose up`

---

This guide covers common issues, debugging techniques, and solutions for the Local Development Framework.

## 📋 Overview

Most issues with the framework fall into these categories:

- Docker and container issues
- Service connectivity problems
- Configuration errors
- Resource constraints
- Port conflicts

## 🚨 Quick Diagnosis

### Health Check Commands

```bash
# Quick system check
docker info                              # Docker daemon status
otto-stack status                       # Service status and connection information
otto-stack doctor                       # Comprehensive system health check

# Resource check
docker system df                        # Docker disk usage
free -h                                 # Available memory (Linux)
vm_stat                                 # Memory stats (macOS)

# Network check
lsof -i :5432                          # PostgreSQL port
lsof -i :6379                          # Redis port
lsof -i :9092                          # Kafka port
```

### Log Analysis

```bash
# View all service logs
otto-stack logs

# Recent errors only
otto-stack logs --since=1h | grep -i error

# Service-specific logs
otto-stack logs postgres -f
otto-stack logs redis --tail=50
```

## 🐳 Docker Issues

### Docker Not Running

**Symptoms:**

- "Cannot connect to the Docker daemon" error
- `docker info` fails

**Solutions:**

**macOS with Colima:**

```bash
# Check Colima status
colima status

# Start Colima
colima start --cpu 4 --memory 8

# If stuck, reset Colima
colima stop
colima delete
colima start --cpu 4 --memory 8 --vm-type=vz --mount-type=virtiofs
```

**macOS with Docker Desktop:**

```bash
# Restart Docker Desktop through the application
# Or via command line:
killall Docker && open /Applications/Docker.app
```

**Linux:**

```bash
# Check Docker service
sudo systemctl status docker

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# If permission denied
sudo usermod -aG docker $USER
# Logout and login again
```

### Docker Out of Space

**Symptoms:**

- "No space left on device" errors
- Container creation fails

**Solutions:**

```bash
# Check Docker disk usage
docker system df

# Clean up unused resources
docker system prune -a

# Remove unused volumes
docker volume prune

# Remove unused networks
docker network prune

# Clean up framework resources specifically
otto-stack cleanup
```

### Docker Memory Issues

**Symptoms:**

- Services crash randomly
- Slow performance
- "Cannot allocate memory" errors

**Solutions:**

```bash
# Check current memory usage
docker stats

# Increase Docker memory limit
# Docker Desktop: Settings > Resources > Memory (8GB+)
# Colima: colima start --memory 8

# Reduce service memory usage in config
vim otto-stack-config.yaml
# overrides:
#   postgres:
#     memory_limit: "256m"
#   redis:
#     memory_limit: "128m"
```

## 🔌 Service Connectivity Issues

### Cannot Connect to Database

**Symptoms:**

- Connection refused errors
- Timeout connecting to PostgreSQL/MySQL
- Application startup fails

**Diagnosis:**

```bash
# Check if service is running
otto-stack status

# Check if port is open
telnet localhost 5432               # PostgreSQL
telnet localhost 3306               # MySQL

# Check service logs
otto-stack logs postgres
otto-stack logs mysql

# Test connection directly
psql -h localhost -U postgres
mysql -h localhost -u root -p
```

**Solutions:**

```bash
# Restart database service
otto-stack down
otto-stack up

# Check for port conflicts
lsof -i :5432
# Kill conflicting process if found
kill -9 PID

# Verify configuration
otto-stack status

# Recreate service
otto-stack down
otto-stack up
```

### Redis Connection Issues

**Symptoms:**

- "Connection refused" to Redis
- Authentication failures
- Timeout errors

**Diagnosis:**

```bash
# Test Redis connection
redis-cli -h localhost -p 6379 ping

# With password
redis-cli -h localhost -p 6379 -a your-password ping

# Check Redis logs
otto-stack logs redis

# Check Redis info
otto-stack exec redis redis-cli INFO
```

**Solutions:**

```bash
# Restart Redis
otto-stack restart redis

# Clear Redis data if corrupted
otto-stack exec redis redis-cli FLUSHALL

# Check Redis configuration
otto-stack exec redis redis-cli CONFIG GET "*"

# Verify password in configuration
otto-stack services
```

### Kafka Connection Issues

**Symptoms:**

- Cannot connect to Kafka broker
- Topic creation fails
- Consumer/producer errors

**Diagnosis:**

```bash
# Check Kafka status
otto-stack logs kafka
otto-stack logs zookeeper

# Test Kafka connection
otto-stack exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Check Kafka UI
curl http://localhost:8080
```

**Solutions:**

```bash
# Restart Kafka stack
otto-stack restart kafka

# Clear Kafka data if needed
otto-stack down kafka
docker volume rm $(docker volume ls -q | grep kafka)
otto-stack up kafka

# Check Zookeeper connectivity
otto-stack exec zookeeper zkCli.sh -server localhost:2181
```

## 🌐 Network and Port Issues

### Port Already in Use

**Symptoms:**

- "Port is already allocated" errors
- "Address already in use" errors
- Services fail to start

**Diagnosis:**

```bash
# Find what's using the port
lsof -i :5432                          # PostgreSQL
lsof -i :6379                          # Redis
lsof -i :9092                          # Kafka
netstat -tulpn | grep :PORT

# Check for other framework instances
otto-stack services
```

**Solutions:**

```bash
# Kill process using the port
kill -9 PID

# Use different ports in configuration
vim otto-stack-config.yaml
# overrides:
#   postgres:
#     port: 5433
#   redis:
#     port: 6380

# Let framework handle conflicts automatically
otto-stack up --cleanup-existing
otto-stack up --force
```

### DNS Resolution Issues

**Symptoms:**

- Cannot resolve service hostnames
- "Name or service not known" errors

**Solutions:**

```bash
# Use localhost instead of service names
# In application configuration:
# spring.datasource.url=jdbc:postgresql://localhost:5432/db

# Check Docker network
docker network ls
docker network inspect otto-stack-framework_default

# Recreate network
otto-stack down --volumes
otto-stack up
```

## ⚙️ Configuration Issues

### Invalid Configuration File

**Symptoms:**

- YAML parsing errors
- "Configuration file not found"
- Setup script fails with validation errors

**Diagnosis:**

```bash
# Validate YAML syntax
python -c "import yaml; yaml.safe_load(open('otto-stack-config.yaml'))"

# Check configuration with framework
otto-stack up --validate-only

# Show resolved configuration
otto-stack up --debug --dry-run
```

**Solutions:**

```bash
# Create new configuration from sample
otto-stack up --init --force

# Fix YAML syntax errors
vim otto-stack-config.yaml
# Common issues:
# - Incorrect indentation (use spaces, not tabs)
# - Missing quotes around strings with special characters
# - Incorrect list syntax

# Use online YAML validator
# Copy content to https://yaml-online-parser.appspot.com/
```

### Service Configuration Errors

**Symptoms:**

- Services start but behave incorrectly
- Authentication failures
- Wrong database/cache settings

**Solutions:**

```bash
# Check generated configuration
cat docker-compose.generated.yml
cat .env.generated
cat application-local.yml.generated

# Compare with sample configuration
otto-stack init --dry-run > sample-config.yaml
diff otto-stack-config.yaml sample-config.yaml

# Reset to defaults
otto-stack up --init
# Manually merge your changes
```

## 🔄 Service-Specific Issues

### PostgreSQL Issues

**Connection refused:**

```bash
# Check if PostgreSQL is ready
otto-stack exec postgres pg_isready -h localhost

# Check PostgreSQL logs
otto-stack logs postgres

# Reset PostgreSQL data
otto-stack down postgres
docker volume rm $(docker volume ls -q | grep postgres)
otto-stack up postgres
```

**Database doesn't exist:**

```bash
# Create database manually
otto-stack exec postgres createdb -U postgres my_app_dev

# Or recreate with correct configuration
vim otto-stack-config.yaml
# overrides:
#   postgres:
#     database: "my_app_dev"
otto-stack up postgres
```

**Permission denied:**

```bash
# Check user and permissions
otto-stack exec postgres psql -U postgres -c "\du"

# Create user if missing
otto-stack exec postgres psql -U postgres -c "CREATE USER app_user WITH PASSWORD 'password';"
otto-stack exec postgres psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE my_app_dev TO app_user;"
```

### Redis Issues

**Memory issues:**

```bash
# Check Redis memory usage
otto-stack exec redis redis-cli INFO memory

# Clear Redis data
otto-stack exec redis redis-cli FLUSHALL

# Increase memory limit
vim otto-stack-config.yaml
# overrides:
#   redis:
#     memory_limit: "512m"
```

**Persistence issues:**

```bash
# Check Redis persistence
otto-stack exec redis redis-cli LASTSAVE

# Disable persistence for development
vim otto-stack-config.yaml
# overrides:
#   redis:
#     config: |
#       save ""
```

### LocalStack Issues

**Services not available:**

```bash
# Check LocalStack logs
otto-stack logs localstack

# Check LocalStack health
curl http://localhost:4566/health

# Restart LocalStack
otto-stack restart localstack

# Check enabled services
curl http://localhost:4566/_localstack/health | jq
```

**SQS/SNS issues:**

```bash
# List SQS queues
aws --endpoint-url=http://localhost:4566 sqs list-queues

# List SNS topics
aws --endpoint-url=http://localhost:4566 sns list-topics

# Recreate queues/topics
otto-stack down localstack
otto-stack up localstack
```

**DynamoDB issues:**

```bash
# List DynamoDB tables
aws --endpoint-url=http://localhost:4566 dynamodb list-tables

# Check table status
aws --endpoint-url=http://localhost:4566 dynamodb describe-table --table-name my-table

# Recreate tables
otto-stack exec localstack awslocal dynamodb delete-table --table-name my-table
otto-stack up localstack
```

## 🚀 Performance Issues

### Slow Service Startup

**Symptoms:**

- Services take a long time to start
- Timeouts during startup
- Application fails to connect initially

**Solutions:**

```bash
# Check resource usage
otto-stack status

# Reduce memory limits for faster startup
vim otto-stack-config.yaml
# overrides:
#   global:
#     memory_limit: "256m"

# Disable unnecessary services
# services:
#   enabled:
#     - redis
#     - postgres
#     # - kafka      # Comment out if not needed
#     # - localstack # Comment out if not needed

# Pre-pull images
docker pull postgres:15-alpine
docker pull redis:7-alpine
```

### High Memory Usage

**Symptoms:**

- System becomes slow
- Out of memory errors
- Services crash randomly

**Solutions:**

```bash
# Monitor memory usage
docker stats
otto-stack status

# Reduce service memory limits
vim otto-stack-config.yaml
# overrides:
#   postgres:
#     memory_limit: "256m"
#   redis:
#     memory_limit: "128m"
#   kafka:
#     memory_limit: "512m"

# Increase system memory allocation
# Docker Desktop: Settings > Resources > Memory
# Colima: colima start --memory 8
```

### Slow Database Performance

**Symptoms:**

- Long query execution times
- Application timeouts
- High database CPU usage

**Solutions:**

```bash
# Check PostgreSQL performance
otto-stack exec postgres psql -U postgres -c "SELECT * FROM pg_stat_activity;"

# Optimize PostgreSQL for development
vim otto-stack-config.yaml
# overrides:
#   postgres:
#     config: |
#       shared_buffers = 256MB
#       effective_cache_size = 1GB
#       work_mem = 4MB
#       maintenance_work_mem = 64MB
#       wal_buffers = 16MB
#       checkpoint_completion_target = 0.9
#       random_page_cost = 1.1

# For development only (data safety disabled):
#       fsync = off
#       synchronous_commit = off
#       full_page_writes = off
```

## 🔍 Advanced Debugging

### Container Inspection

```bash
# Inspect running containers
docker ps
docker inspect CONTAINER_ID

# Check container resource usage
docker stats CONTAINER_NAME

# Execute shell in container
otto-stack exec postgres bash
otto-stack exec redis sh

# Check container logs with timestamps
docker logs --timestamps CONTAINER_NAME
```

### Network Debugging

```bash
# Check Docker networks
docker network ls
docker network inspect otto-stack-framework_default

# Test network connectivity between containers
otto-stack exec app-container ping postgres
otto-stack exec app-container telnet redis 6379

# Check DNS resolution
otto-stack exec app-container nslookup postgres
```

### File System Debugging

```bash
# Check Docker volumes
docker volume ls
docker volume inspect VOLUME_NAME

# Check file permissions
otto-stack exec postgres ls -la /var/lib/postgresql/data

# Copy files for debugging
docker cp $(docker ps -q -f name=postgres):/var/log/postgresql/ ./postgres-logs/
```

## 🔧 Advanced Solutions

### Complete Reset

When all else fails, perform a complete reset:

```bash
# Stop all services
otto-stack down

# Remove all framework resources
otto-stack down --volumes

# Clean Docker system
docker system prune -a
docker volume prune

# Remove configuration and start fresh
rm otto-stack-config.yaml
rm docker-compose.generated.yml
rm .env.generated
rm application-local.yml.generated

# Initialize new configuration
otto-stack up --init
vim otto-stack-config.yaml
otto-stack up
```

### Framework Recovery

If the framework itself is corrupted:

```bash
# Update framework (if using git submodule)
git submodule update --remote otto-stack-framework

# Or re-copy framework files
rm -rf otto-stack-framework
cp -r /path/to/fresh/otto-stack-framework ./

# Make scripts executable
# Scripts are no longer needed - use otto-stack CLI directly

# Regenerate configuration
otto-stack init
```

### System Resource Recovery

```bash
# Free up system resources
docker system prune -a --volumes

# Clear system caches (Linux)
sudo sync && echo 3 | sudo tee /proc/sys/vm/drop_caches

# Restart Docker daemon (Linux)
sudo systemctl restart docker

# Reset Colima completely (macOS)
colima stop
colima delete
rm -rf ~/.colima
colima start --cpu 4 --memory 8
```

## 📊 Monitoring and Prevention

### Health Monitoring

Create a health check script:

```bash
#!/bin/bash
# health-check.sh

echo "=== Framework Health Check ==="
echo "Docker Status:"
docker info > /dev/null 2>&1 && echo "✓ Docker running" || echo "✗ Docker not running"

echo "Services Status:"
otto-stack status

echo "Resource Usage:"
otto-stack status | head -10

echo "Disk Usage:"
docker system df

echo "=== End Health Check ==="
```

### Preventive Maintenance

Weekly maintenance routine:

```bash
#!/bin/bash
# weekly-maintenance.sh

# Backup databases
# Backup postgres (use docker exec or service-specific tools)
# Backup mysql (use docker exec or service-specific tools)

# Clean up Docker resources
docker system prune

# Update service images
otto-stack up --build

# Validate configuration
otto-stack up --validate-only

echo "Maintenance complete"
```

## 📞 Getting Help

### Self-Help Checklist

Before seeking help, try these steps:

1. **Check the basics:**
   - [ ] Docker is running: `docker info`
   - [ ] Services are running: `otto-stack status`
   - [ ] No port conflicts: `lsof -i :5432 :6379 :9092`
   - [ ] Sufficient resources: `docker stats`

2. **Review logs:**
   - [ ] Framework logs: `otto-stack logs`
   - [ ] Service-specific logs: `otto-stack logs SERVICE_NAME`
   - [ ] System logs: `dmesg | tail` (Linux)

3. **Validate configuration:**
   - [ ] YAML syntax: `python -c "import yaml; yaml.safe_load(open('otto-stack-config.yaml'))"`
   - [ ] Framework validation: `otto-stack up --validate-only`

4. **Try simple fixes:**
   - [ ] Restart services: `otto-stack restart`
   - [ ] Recreate services: `otto-stack up --force`
   - [ ] Clear caches: `otto-stack down --volumes-docker`

### Information to Collect

When reporting issues, include:

```bash
# System information
uname -a
docker --version
docker compose version

# Framework status
otto-stack services
otto-stack status

# Configuration
cat otto-stack-config.yaml

# Recent logs
otto-stack logs --since=1h

# Resource usage
docker stats --no-stream
docker system df
```

### Debug Mode

Enable debug mode for detailed information:

```bash
# Debug setup
otto-stack up --debug

# Debug with dry run
otto-stack up --debug --dry-run

# Verbose logging
export DEBUG=1
otto-stack up
```

## 📚 Related Documentation

- **[Setup Guide](setup.md)** - Initial installation and configuration
- **[Configuration Guide](configuration.md)** - Detailed configuration options
- **[Usage Guide](usage.md)** - Daily commands and workflows
- **[Services Guide](services.md)** - Service-specific information
- **[Quick Reference](reference.md)** - Commands cheatsheet

## 🏥 Emergency Procedures

---

## 📚 See Also

- [README](../README.md)
- [Setup Guide](setup.md)
- [Configuration Guide](configuration.md)
- [Services Guide](services.md)
- [Usage Guide](usage.md)
- [Integration Guide](integration.md)
- [Contributing Guide](contributing.md)

### Complete System Recovery

```bash
# 1. Stop everything
otto-stack down --volumes-all

# 2. Clean Docker completely
docker system prune -a --volumes

# 3. Restart Docker
# macOS: colima stop && colima start
# Linux: sudo systemctl restart docker

# 4. Start fresh
otto-stack up --init
otto-stack up

# 5. Verify
otto-stack status
```

### Data Recovery

```bash
# If you have backups
# Restore postgres (use docker exec or pg_restore)

# If no backups, check Docker volumes
docker volume ls | grep postgres
# Mount volume to recover data
docker run --rm -v VOLUME_NAME:/data -v $(pwd):/backup alpine cp -r /data /backup/
```

Remember: Most issues can be resolved by restarting services or recreating them with `otto-stack up --force`. When in doubt, start with the simplest solutions first.
