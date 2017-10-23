module.exports = {
  verbose: true,
  coverageThreshold: {
    global: {
      functions: 74,
      lines: 85,
      statements: 80,
    },
  },
  roots: ["src"],
  coveragePathIgnorePatterns: ["vnd/"]
};
