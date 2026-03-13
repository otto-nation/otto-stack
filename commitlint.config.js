// Commitlint configuration for local git hooks (commit-msg).
// Uses inline rules instead of `extends` to avoid module resolution issues
// when run via npx in the commit-msg hook.
//
// IMPORTANT: Keep commit types and header-max-length in sync with
// .github/.commitlintrc.mjs (used by CI).

module.exports = {
  ignores: [
    (commit) => commit.includes('dependabot[bot]'),
  ],
  rules: {
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
