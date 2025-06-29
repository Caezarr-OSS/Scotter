module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'body-max-line-length': [1, 'always', 100],
    // Allow capitalized subject (our current commits use capitalization)
    'subject-case': [0, 'always', []],
    // Make these warnings instead of errors
    'type-enum': [1, 'always', ['build', 'chore', 'ci', 'docs', 'feat', 'fix', 'perf', 'refactor', 'revert', 'style', 'test']],
    'subject-empty': [1, 'never'],
    'type-empty': [1, 'never']
  },
};
