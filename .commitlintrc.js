module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat',
        'fix',
        'docs',
        'style',
        'refactor',
        'perf',
        'test',
        'build',
        'ci',
        'chore',
        'revert'
      ]
    ],
    'scope-enum': [
      2,
      'always',
      [
        'core',
        'model',
        'prompt',
        'generator',
        'config',
        'init',
        'cli',
        'ci',
        'docs',
        'deps',
        'build'
      ]
    ],
    'scope-empty': [
      2,
      'never',
      {
        'exceptions': ['docs']
      }
    ],
    // Règle scope-required qui confirme l'exception pour docs
    'scope-required': [
      2,
      'always',
      {
        'exceptions': ['docs']
      }
    ],
    'subject-case': [0],
    'body-max-line-length': [0]
  }
};
