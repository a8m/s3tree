s3tree
---
> s3tree is a tree command for [Amazon S3](https://aws.amazon.com/s3/).

### Installation:
```sh
$ go install github.com/kk17/s3tree@latest
```
for golang version less than 1.8:
```sh
$ go get github.com/kk17/s3tree
```

### How to use ?
```sh
$ s3tree -b bucket-name -p prefix(optional) [options...]
```
Remember, your credentials should located at `~/.aws/credentials` or as an environment variables: 
`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`

### Develop
```sh
go get
go build
./s3tree xxx
```

### License
MIT
