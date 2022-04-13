module.exports = {
    branches: ['main'],
    plugins: [
        "@semantic-release/commit-analyzer",
        ['@semantic-release/exec', {
            publishCmd: "export NEXT_VERSION=${nextRelease.version}"
          }
        ]
    ]
}
