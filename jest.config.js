module.exports = {
  verbose: true,
  coverageThreshold: {
    global: {
      functions: 100,
      lines: 100,
      statements: 100,
    },
  },
  roots: ["src"],
  coveragePathIgnorePatterns: ["vendored/"]
};
