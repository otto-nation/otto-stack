---
title: "Troubleshooting"
description: "Common issues and solutions for otto-stack problems"
lead: "Quick solutions to the most common otto-stack issues"
date: "2025-10-01"
lastmod: "2025-11-10"
draft: false
weight: 70
toc: true
---

# Troubleshooting

## 🚨 Quick Fixes

### Docker Issues

```bash
# Check if Docker is running
docker info

# Start Docker (macOS with Colima)
colima start

# Restart Docker services
otto-stack down && otto-stack up
```

### Port Conflicts

```bash
# Find what's using a port
lsof -i :5432

# Kill process using port
kill -9 $(lsof -t -i:5432)
```

### Service Won't Start

```bash
# Check service logs
otto-stack logs postgres

# Restart specific service
otto-stack restart postgres

# Full reset
otto-stack down && otto-stack up
```

### Configuration Issues

```bash
# Validate configuration
otto-stack doctor

# Debug mode
otto-stack --verbose up
```

## 🔍 Common Issues

### "Docker not found"

- Install Docker Desktop or Colima
- Ensure Docker daemon is running
- Check PATH includes Docker binaries

### "Port already in use"

- Another service is using the same port
- Use `lsof -i :PORT` to identify the process
- Either stop the conflicting service or change ports

### "Service unhealthy"

- Service failed health checks
- Check logs with `otto-stack logs SERVICE_NAME`
- Verify service configuration in YAML files

### "Out of memory"

- Increase Docker memory limits
- Check `docker stats` for resource usage
- Consider reducing number of running services

## 📞 Getting Help

1. **Check logs**: `otto-stack logs SERVICE_NAME`
2. **Run diagnostics**: `otto-stack doctor`
3. **Search issues**: [GitHub Issues](https://github.com/otto-nation/otto-stack/issues)
4. **Ask questions**: [GitHub Issues](https://github.com/otto-nation/otto-stack/issues)

## 🧹 Reset Everything

If all else fails:

```bash
# Stop all services
otto-stack down

# Clean up containers and volumes
docker system prune -a --volumes

# Restart fresh
otto-stack up
```
