module.exports = {
  verbose: true,
  collectCoverage: true,
  coverageThreshold: {
    global: {
      functions: 18,
      lines: 46,
      statements: 45,
      branch: 17
    }
  },
  moduleNameMapper: {
    '\\.(css|less|scss|sss|styl)$': '<rootDir>/../../node_modules/jest-css-modules'
  },
  roots: ['src'],
  preset: 'ts-jest',
  testEnvironment: 'jsdom'
};