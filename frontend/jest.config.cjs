module.exports = {
  testEnvironment: "jsdom",
  transform: {
    "^.+\\.tsx?$": "ts-jest",
  },
  moduleNameMapper: {
    "^.+\\.svg$": "jest-svg-transformer",
    "\\.(css|less|sass|scss)$": "identity-obj-proxy",
    "^.+\\.(jpg|jpeg|png|gif|webp|avif|svg|ttf|woff|woff2)$":
      "jest-svg-transformer",
  },
  setupFilesAfterEnv: ["<rootDir>/jest.setup.ts"],
};
