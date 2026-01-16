# GitHub Pages Setup

Guide for setting up and maintaining the otto-stack documentation site.

## Quick Setup

### 1. Enable GitHub Pages

1. Go to **Settings** → **Pages**
2. Set **Source** to "GitHub Actions"
3. Save

### 2. Initialize Theme

```bash
git submodule update --init --recursive
```

### 3. Deploy

Push to `main` branch. Site deploys automatically via `.github/workflows/pages.yml`.

Site URL: `https://[username].github.io/otto-stack`

## How It Works

**Deployment Process:**
1. Push to `main` triggers workflow
2. Hugo builds static site from `docs-site/`
3. Documentation generators run (CLI reference, homepage, etc.)
4. Site deploys to GitHub Pages

**Key Files:**
- `.github/workflows/pages.yml` - Deployment workflow
- `docs-site/hugo.toml` - Hugo configuration
- `docs-site/content/` - Documentation content
- `docs-site/themes/PaperMod/` - Hugo theme (submodule)
- `docs-site/scripts/generate-docs.js` - Documentation generator

## Local Development

### Setup

```bash
# Install Hugo Extended
brew install hugo  # macOS

# Start development server
cd docs-site
hugo server --buildDrafts
# Site at http://localhost:1313
```

### Test Build

```bash
cd docs-site
hugo --gc --minify
```

### Generate Documentation

```bash
# Generate all docs
cd docs-site
node scripts/generate-docs.js

# Generate specific docs only
node scripts/generate-docs.js --generator=homepage
node scripts/generate-docs.js --skip-build --skip-format
```

See improved generate-docs.js with `--skip-build`, `--skip-format`, and `--generator` flags.

## Troubleshooting

### Site Not Deploying

**Check:**
- Settings → Pages → Source = "GitHub Actions"
- Actions tab for workflow errors
- Repository is public

**Fix:**
- Manually trigger workflow from Actions tab
- Check workflow logs for build errors

### Theme Missing

```bash
# Re-initialize submodule
git submodule update --init --recursive

# Verify theme exists
ls docs-site/themes/PaperMod/
```

### Content Not Updating

**Common causes:**
- Browser cache (hard refresh: Cmd+Shift+R)
- GitHub Pages cache (wait 5-10 minutes)
- Changes not pushed to `main` branch

### Build Failures

**Check Hugo version:**
```bash
hugo version
# Requires Hugo Extended v0.146.0+
```

**Validate configuration:**
```bash
cd docs-site
hugo config
```

## Manual Deployment

If automatic deployment fails:

1. Go to **Actions** tab
2. Select "Deploy Documentation to GitHub Pages"
3. Click "Run workflow"
4. Select `main` branch

## Maintenance

### Update Hugo Version

Edit `.github/workflows/pages.yml`:

```yaml
- name: Setup Hugo
  uses: peaceiris/actions-hugo@v3
  with:
    hugo-version: "0.151.0"  # Update version here
```

### Update Theme

```bash
cd docs-site/themes/PaperMod
git pull origin master
cd ../../..
git add docs-site/themes/PaperMod
git commit -m "chore: update PaperMod theme"
```

## Resources

- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [PaperMod Theme](https://github.com/adityatelange/hugo-PaperMod)
