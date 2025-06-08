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
        'deps'
      ]
    ],
    'scope-empty': [2, 'never'],
    'subject-case': [0],
    'body-max-line-length': [0]
  }
};
