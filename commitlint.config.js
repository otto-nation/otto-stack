// Commitlint configuration for local validation
// This mirrors .github/.commitlintrc.json but uses inline rules instead of extends
// to avoid module resolution issues with npx
//
// IMPORTANT: Keep this in sync with .github/.commitlintrc.json
// The rules below are from @commitlint/config-conventional plus our custom overrides

module.exports = {
  rules: {
    // From @commitlint/config-conventional
    'body-leading-blank': [1, 'always'],
    'body-max-line-length': [2, 'always', 100],
    'footer-leading-blank': [1, 'always'],
    'footer-max-line-length': [2, 'always', 100],
    'header-max-length': [2, 'always', 72],
    'subject-case': [
      2,
      'never',
      ['sentence-case', 'start-case', 'pascal-case', 'upper-case'],
    ],
    'subject-empty': [2, 'never'],
    'subject-full-stop': [2, 'never', '.'],
    'type-case': [2, 'always', 'lower-case'],
    'type-empty': [2, 'never'],
    'type-enum': [
      2,
      'always',
      [
        'feat',
        'fix',
        'perf',
        'deps',
        'revert',
        'docs',
        'style',
        'refactor',
        'test',
        'build',
        'ci',
        'chore',
      ],
    ],
  },
};
