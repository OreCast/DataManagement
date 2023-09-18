# DataManagement
Data Management for OreCast provides RESTfull access to site's S3 storage.
It supports the following APIs:
```
# create bucket (s3-bucket) at site cornell
curl -X POST http://localhost:8340/storage/cornell/s3-bucket

# delete bucket
 curl -X DELETE http://localhost:8340/storage/cornell/s3-bucket

# upload file:
# take local file at /path/test.zip and upload it to
# S3 object: cornell/s3-bucket/archive.zip
curl -x POST http://localhost:8340/storage/cornell/s3-bucket/archive.zip \
  -f "file=@/path/test.zip" \
  -h "content-type: multipart/form-data"

# get file
curl http://localhost:8340/storage/cornell/s3-bucket/archive.zip > archive.zip

# delete file
curl -X DELETE http://localhost:8340/storage/cornell/s3-bucket/archive.zip

# update file
curl -x PUT http://localhost:8340/storage/cornell/s3-bucket/archive.zip \
  -f "file=@/path/test.zip" \
  -h "content-type: multipart/form-data"
```
