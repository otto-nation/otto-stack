module.exports = {
  generators: [
    { name: 'cli-reference', enabled: true, output: 'reference.md' },
    { name: 'services-guide', enabled: true, output: 'services.md' },
    { name: 'configuration-guide', enabled: true, output: 'configuration.md' },
    { name: 'homepage', enabled: true, output: '_index.md' }
  ],

  validation: {
    enabled: true,
    strict: false // Set to true to fail on warnings
  },

  templates: './templates',
  servicesDir: '../internal/config/services',
  outputDir: './content'
};
