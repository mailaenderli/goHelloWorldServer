module.exports = {
    branches: ['main'],
    plugins: [
        "@semantic-release/commit-analyzer",
        ['@codedependant/semantic-release-docker', {
                dockerTags: ['latest', '{{version}}', '{{major}}-latest', '{{major}}.{{minor}}'],
                dockerImage: 'gohelloworldserver',
                dockerFile: 'Dockerfile',
                dockerProject: 'rmailaender',
                dockerArgs: {
                    RELEASE_DATE: new Date().toISOString()
                    , RELEASE_VERSION: '{{next.version}}'
                }
            }
        ],
    ]
}
