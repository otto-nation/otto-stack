# Release Process

Automated releases using Release-Please and conventional commits.

## Quick Reference

**Process**: Push to `main` → Release-Please creates PR → Merge PR → Release created

**Manual trigger**: GitHub Actions → Release Please → Run workflow

## Conventional Commits

### Format
```
type(scope): description
```

### Types & Version Bumps

| Type | Changelog | Version Bump |
|------|-----------|--------------|
| `feat` | ✅ Features | Minor |
| `fix` | ✅ Bug Fixes | Patch |
| `perf` | ✅ Performance | Patch |
| `feat!` or `BREAKING CHANGE:` | ✅ Breaking | Major |
| `docs`, `refactor`, `test`, `ci`, `chore` | ❌ Hidden | None |

### Examples

```bash
# Feature
git commit -m "feat(cli): add backup command"

# Bug fix
git commit -m "fix(docker): resolve port conflict"

# Breaking change
git commit -m "feat!: change config location"
```

## Release Workflow

1. **Commits pushed to main** → Release-Please analyzes
2. **Release PR created** → Contains version bump + changelog
3. **Review and merge PR** → Triggers release:
   - Builds binaries (all platforms)
   - Creates GitHub release
   - Updates package managers (if enabled)

## Configuration

**Location**: `.github/config/release-config.yaml`

**After editing**, regenerate:
```bash
task generate-release-configs
```

## Troubleshooting

**No release PR created?**
- Ensure commits use conventional format
- Check if there are releasable commits since last release

**Wrong version bump?**
- Verify commit types are correct
- Use `feat!:` or `BREAKING CHANGE:` footer for major bumps

**Build failures?**
- Fix issues and push again
- Release-Please will update the PR automatically

## Best Practices

- Use descriptive commit messages
- Keep subject lines under 50 characters
- Test before merging to main
- Review release PRs carefully
- Document breaking changes with migration steps

## Resources

- [Release-Please docs](https://github.com/googleapis/release-please)
- [Conventional Commits spec](https://www.conventionalcommits.org/)
- View recent releases for examples
