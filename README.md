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

## API Documentation

### Register storage 

**URL**

`POST` http://localhost:8000/register

**Description**

register storage. create a storage which is to be used for accessing and modifying user content.

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| name | storage driver name, current support `s3`,`minio` | body |
| region | Amazon S3 region | body |
| endpoint | Amazon S3 endpoint | body |
| key | Amazon S3 access key | body |
| token | Amazon S3 access token | body |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model : json containing message successful

**Example**
```
curl --request POST \
  --url http://localhost:8000/register \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <jwt token>' \
  --data '{ 
    "name": "s3", 
    "region":"www", 
    "endpoint":"xxx", 
    "key":"yyy", 
    "token":"zzz"
}'
```  

### UnRegister storage 

**URL**

`DELETE` http://localhost:8000/register/<name>

**Description**

remove registered storage from storage db.

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| name | storage driver name, current support `s3`,`minio` | url |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model : json containing message successful

**Example**
```
curl --request DELETE \
  --url http://localhost:8000/register/s3 \
  --header 'Authorization: Bearer <jwt token>' 
```  

### Create bucket

**URL**

`POST` http://localhost:8000/pods

**Description**

create bucket in storage .

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| name | bucket name | body |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model : json containing message successful

**Example**
```
curl --request POST \
  --url http://localhost:8000/pods \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer <jwt token>' \
  --data '{ 
    "name": "test"
}'
```  
### Delete bucket

**URL**

`DELETE` http://localhost:8000/pods/<bucket_name>

**Description**

delete bucket in storage .

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| bucket_name | bucket name | url |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model : json containing message successful

**Example**
```
curl --request DELETE \
  --url http://localhost:8000/pods/test \
  --header 'Authorization: Bearer <jwt token>' 
```  

### Create object

**URL**

`POST` http://localhost:8000/pea/<bucket_name>/<object_name>

**Description**

upload a object to bucket in storage .

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| bucket_name | bucket name, in example: `test` | url |
| object_name | object name, in example: `foobar.txt` | url |
| object_content | object content | body |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model : json containing message successful

**Example**
```
curl --request POST \
  --url http://localhost:8000/pods/test/foobar.txt \
  --header 'Authorization: Bearer <jwt token>' 
  --data 'hello, foobar'
```  
### Get object

**URL**

`GET` http://localhost:8000/pea/<bucket_name>/<object_name>

**Description**

fetch content of a object .

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| bucket_name | bucket name, in example: `test` | url |
| object_name | object name, in example: `foobar.txt` | url |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model :  object content

**Example**
```
curl --request GET \
  --url http://localhost:8000/pods/test/foobar.txt \
  --header 'Authorization: Bearer <jwt token>' 
```  

### Get objects

**URL**

`GET` http://localhost:8000/pea/<bucket_name>

**Description**

fetch object list of a bucket .

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| bucket_name | bucket name, in example: `test` | url |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model :  object list content

**Example**
```
curl --request GET \
  --url http://localhost:8000/pea/test \
  --header 'Authorization: Bearer <jwt token>' 
```  


### Delete object

**URL**

`DELETE` http://localhost:8000/pea/<bucket_name>/<object_name>

**Description**

delete a object of bucket .

**Parameters**

| Parameter | Description | Parameter type |
| --- | --- | --- |
| jwt token | authentication jwt token | header |
| bucket_name | bucket name, in example: `test` | url |
| object_name | object name, in example: `foobar.txt` | url |

**Response Messages**

- Response Code : `400`
  - Response model : json containing error message

- Response Code : `403` 
  - Response model : json containing error message
  - Reason: User don't have permission for this

- Response Code : `200`
  - Response model :  object content

**Example**
```
curl --request DELETE \
  --url http://localhost:8000/pods/test/foobar.txt \
  --header 'Authorization: Bearer <jwt token>' 
```  
