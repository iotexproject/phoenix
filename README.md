# phoenix-gem

## install

```
make build
make run #build + run
```

## configure
### use minio instead of s3
```
docker run -p 9000:9000 \
  -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE" \
  -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  minio/minio:edge server /data
```

### s3 storage config
```
storage:
  provider: s3
s3:
  endpoint: http://localhost:9001
  accessKey: AKIAIOSFODNN7EXAMPLE
  secretKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  region: 
```

## test
`test11` means bucket
`foo.txt` means object in bucket
### create pod
```
curl -H "Content-type: application/json" -d '{ "name": "test11"}' http://localhost:8080/pods
```
### delete pod
```
curl -X DELETE http://localhost:8080/pods/test11
```
### upload pea
```
curl -d 'foobar' http://localhost:8080/pea/test11/foo.txt
```
### get pea
```
curl http://localhost:8080/pea/test11/foo.txt
```
### delete pea
```
curl -X DELETE http://localhost:8080/pea/test11/foo.txt
```

### list pea
```
curl http://localhost:8080/pea/test11
```