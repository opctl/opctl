module.exports = {
  preset: 'ts-jest',
  verbose: true,
  coverageThreshold: {
    global: {
      functions: 60,
      lines: 84,
      statements: 83,
      branch: 91
    }
  },
  roots: ['src']
}
