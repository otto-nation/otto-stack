# Documentation Site

Hugo-based documentation site for otto-stack.

## Quick Start

### Prerequisites

- Hugo Extended v0.146.0+
- Node.js 18+

### Setup

```bash
cd docs-site
npm install
git submodule update --init --recursive
```

### Development

```bash
# Generate docs and start dev server
npm run docs:dev

# Or from project root
task docs-serve
```

Visit <http://localhost:1313>

## Documentation Generation

CLI reference and services guide are auto-generated from source:

```bash
# Generate all documentation
npm run docs:generate

# Generate with options
node scripts/generate-docs.js --skip-build --generator=homepage
```

### Generated Files

- `content/cli-reference.md` - From `internal/config/commands.yaml`
- `content/services.md` - From service configurations
- `content/configuration.md` - Configuration examples
- `content/_index.md` - From root `README.md`
- `content/contributing.md` - From root `CONTRIBUTING.md`

Do not edit these files manually - they are regenerated on each build.

## Available Commands

### npm Scripts

- `npm run docs:generate` - Generate documentation from source
- `npm run docs:dev` - Start Hugo dev server
- `npm run docs:build` - Build production site
- `npm run docs:serve` - Serve built site
- `npm run docs:clean` - Clean build artifacts
- `npm run format` - Format with Prettier
- `npm run lint:md` - Lint markdown files
- `npm run lint:links` - Check for broken links

### Task Commands (from project root)

- `task docs` - Generate documentation
- `task docs-serve` - Generate and serve docs
- `task docs-fmt` - Format documentation
- `task lint:docs` - Lint documentation

## Structure

```
docs-site/
├── config/_default/   # Hugo configuration
│   ├── hugo.toml      # Main config
│   ├── params.toml    # Site parameters
│   └── menus/         # Navigation
├── content/           # Markdown content
├── generators/        # Doc generators
├── scripts/           # Build scripts
├── templates/         # Handlebars templates
├── themes/PaperMod/   # Hugo theme (submodule)
└── utils/             # Generator utilities
```

## Theme

Uses [PaperMod](https://github.com/adityatelange/hugo-PaperMod) theme via git submodule.

Update theme:
```bash
cd themes/PaperMod
git pull origin master
```

## Deployment

Site deploys automatically to GitHub Pages via `.github/workflows/pages.yml` on push to `main`.

See [docs/GITHUB_PAGES.md](../docs/GITHUB_PAGES.md) for setup details.

## Troubleshooting

**Hugo build fails:**
```bash
hugo config  # Check configuration
```

**Theme missing:**
```bash
git submodule update --init --recursive
```

**Stale content:**
```bash
npm run docs:clean
npm run docs:build
```

**Link checker fails:**
Internal Hugo links (starting with `/`) are ignored by the link checker - this is expected.
