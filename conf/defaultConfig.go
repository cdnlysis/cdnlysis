package conf

const defaultYaml = `
---
engine:
    verbose: true
    threads: 10
syncprogress:
    path: /tmp/cdn_sync_progress
s3:
    region: us-east-1
logs:
    prefix: cdn
    location: /tmp
`
