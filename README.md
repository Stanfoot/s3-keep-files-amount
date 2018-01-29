# S3 Keep Files Amount

S3のオブジェクトを常に指定数に保ちます。

```bash
AWS_ACCESS_KEY_ID=xxxxx AWS_SECRET_ACCESS_KEY=xxxx s3-keep-files-amount keep [amount] [region] [bucket]
```

## Image Test
### Pull image

```bash
docker pull gcr.io/gcp-runtimes/container-structure-test:v0.1.3
```

###

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v `pwd`/container-structure-test.yaml:/test.yaml gcr.io/gcp-runtimes/container-structure-test:v0.1.3 -pull -test.v --image stanfoot/s3-keep-files-amount:latest test.yaml
```