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
    // La règle scope-empty est configurée pour exiger un scope pour tous les types de commits SAUF docs
    'scope-empty': [
      2,
      'never',
      {
        'exceptions': ['docs']
      }
    ],
    'subject-case': [0],
    'body-max-line-length': [0]
  }
};
