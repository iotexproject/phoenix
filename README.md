# phoenix
Proxy for managing and securely delegating your data for trusted party's use 

## How it works
Let's take a simple example: you have certain data stored in Amazon S3 storage that you wanted to share with your trusted user so they can use it. The workflow consists of 3 steps:

First, register your data access endpoint into phoenix's database, in this example Amazon S3 storage (explained [here](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys)). This allows phoenix to access the data on behalf of delegated user later on.

To register, send an HTTP POST request to <https://phoenix.iotex.io:8080/register>. You can find details in API section [here](#register)

Second, to allow your trusted user to access the data, you will need to issue them a piece of authentication token called JWT. It clearly specifies *what, by when, and how* the access is granted. For example, say you registered a data endpoint named **weather** and you want to grant user to be able to *read in next 12 hours*. The JWT will contain claims like:

```
{
  "exp": "1607772249", // 12 hours from issue time
  "iat": "1605730464", // token issue time
  "iss": "0x04c1bd03c5974777e512a38ff2037ea89772ac81dee6d204afc6dee60ac73238a127fdf383bd7d63387ec5e221026b2b314d3049cf304528d548ea926ea8a834c2",
  "scope": "Read",  // can only read
  "sub": "weather", // can only access data in 'weather'
}  
```

For how to sign and issue JWT, see [here](https://docs.iotex.io/developer/ioctl/jwt.html)

Finally, upon receiving the JWT, your trusted user can embed it into their HTTP request to access or operate on data. For details, see API section [here](#get)

## Install

```
make build
make run #build + run
```

## Deploy
1. Pull the docker image:
```
docker pull iotex/phoenix:latest
```

2. Set the environment with the following commands:
```
mkdir -p ~/iotex-phoenix
cd ~/iotex-phoenix

export IOTEX_HOME=$PWD
```

3. Download the default config file:
```
curl https://github.com/iotexproject/phoenix/blob/main/config.yaml > $IOTEX_HOME/config.yaml
```

4. Run the following command to start a node:
```
docker run -d --name phoenix \
        -p 8080:8080 \
        -v=$IOTEX_HOME:/var/data:rw \
        iotex/phoenix:latest \
        iotex-phoenix
```

## API Documentation

### <a name="register"/>Register storage 

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
### <a name="get"/>Get object

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
