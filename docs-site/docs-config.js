module.exports = {
  generators: [
    { name: 'cli-reference', enabled: true, output: 'cli-reference.md' },
    { name: 'services-guide', enabled: true, output: 'services.md' },
    { name: 'configuration-guide', enabled: true, output: 'configuration.md' },
    { name: 'homepage', enabled: true, output: '_index.md' },
    { name: 'contributing-guide', enabled: true, output: 'contributing.md' }
  ],

  validation: {
    enabled: true,
    strict: false // TODO: Re-enable after fixing remaining schemas
  },

  templates: './templates',
  servicesDir: '../internal/config/services',
  outputDir: './content'
};
