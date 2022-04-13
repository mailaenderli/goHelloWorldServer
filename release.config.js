module.exports = {
    branches: ['main'],
    plugins: [
        "@semantic-release/commit-analyzer",
        ['@semantic-release/exec', {
            publishCmd: ". ./setenvVar.sh ${nextRelease.version}"
          }
        ]
    ]
}
