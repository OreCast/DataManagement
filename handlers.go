package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET handlers
type SiteParams struct {
	Site string `uri:"site" binding:"required"`
}
type BucketParams struct {
	SiteParams
	Bucket string `uri:"bucket" binding:"required"`
}
type ObjectParams struct {
	BucketParams
	Object string `uri:"object" binding:"required"`
}

// StorageHandler provides access to GET /storage end-point
func StorageHandler(c *gin.Context) {
	// TODO: implement return site meta-data info
	c.JSON(200, gin.H{"status": "ok"})
}

// SiteHandler provides access to GET /storage/:site end-point
func SiteHandler(c *gin.Context) {
	var params SiteParams
	if err := c.ShouldBindUri(&params); err == nil {
		data := S3Content(params.Site, "")
		c.JSON(200, gin.H{"status": "ok", "data": data})
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// BucketHandler provides access to GET /storate/:site/:bucket end-point
func BucketHandler(c *gin.Context) {
	var params BucketParams
	if err := c.ShouldBindUri(&params); err == nil {
		data := S3Content(params.Site, params.Bucket)
		c.JSON(200, gin.H{"status": "ok", "data": data})
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// FileHandler provides access to GET /storage/:site/:bucket/:object end-point
func FileHandler(c *gin.Context) {
	var params ObjectParams
	if err := c.ShouldBindUri(&params); err == nil {
		if data, err := getObject(params.Site, params.Bucket, params.Object); err == nil {
			header := fmt.Sprintf("attachment; filename=%s", params.Object)
			c.Header("Content-Disposition", header)
			c.Data(http.StatusOK, "application/octet-stream", data)
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// POST handlers

// BucketPostHandler provides access to POST /storate/:site/:bucket end-point
/*
 * Example of curl call
 * `curl -X POST http://localhost:8340/storage/cornell/s3-bucket`
 */
func BucketPostHandler(c *gin.Context) {
	var params BucketParams
	if err := c.ShouldBindUri(&params); err == nil {
		if err := createBucket(params.Site, params.Bucket); err == nil {
			msg := fmt.Sprintf("Bucket %s/%s created successfully", params.Site, params.Bucket)
			c.JSON(200, gin.H{"status": "ok", "msg": msg})
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// FilePostHandler provides access to POST /storate/:site/:bucket/:object end-point
/*
 * Example of curl call
```
 curl -X POST http://localhost:8340/storage/cornell/s3-bucket/archive.zip \
  -F "file=@/path/test.zip" \
  -H "Content-Type: multipart/form-data"
```
*/
func FilePostHandler(c *gin.Context) {
	var params ObjectParams
	if err := c.ShouldBindUri(&params); err == nil {
		// TODO: read data part from HTTP request body

		// single file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
			return
		}
		log.Printf("INFO: uploading file %s", file.Filename)

		// Upload the file to specific dst.
		reader, err := file.Open()
		if err != nil {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
			return
		}
		defer reader.Close()
		size := file.Size
		ctype := "" // TODO: decide on how to read content-type

		if err := uploadObject(
			params.Site,
			params.Bucket,
			params.Object,
			ctype,
			reader,
			size); err == nil {
			msg := fmt.Sprintf("File %s/%s/%s uploaded successfully", params.Site, params.Bucket, params.Object)
			c.JSON(200, gin.H{"status": "ok", "msg": msg})
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// PUT handlers

// BucketPutHandler provides access to PUT /storate/:site/:bucket end-point
func BucketPutHandler(c *gin.Context) {
}

// FilePutHandler provides access to PUT /storate/:site/:bucket/:object end-point
func FilePutHandler(c *gin.Context) {
	var params ObjectParams
	if err := c.ShouldBindUri(&params); err == nil {
		// TODO: read data from HTTP request body
		var data []byte
		if err := updateFile(params.Site, params.Bucket, params.Object, data); err == nil {
			msg := fmt.Sprintf("Bucket %s/%s/%s updated successfully", params.Site, params.Bucket, params.Object)
			c.JSON(200, gin.H{"status": "ok", "msg": msg})
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// DELETE handlers

// BucketDeleteHandler provides access to DELETE /storate/:site/:bucket end-point
/*
 * Example of curl call
 * `curl -X DELETE http://localhost:8340/storage/cornell/s3-bucket`
 */
func BucketDeleteHandler(c *gin.Context) {
	var params BucketParams
	if err := c.ShouldBindUri(&params); err == nil {
		if err := deleteBucket(params.Site, params.Bucket); err == nil {
			msg := fmt.Sprintf("Bucket %s/%s deleted successfully", params.Site, params.Bucket)
			c.JSON(200, gin.H{"status": "ok", "msg": msg})
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}

// BucketDeleteHandler provides access to DELETE /storate/:site/:bucket end-point
func FileDeleteHandler(c *gin.Context) {
	var params ObjectParams
	if err := c.ShouldBindUri(&params); err == nil {
		if err := deleteFile(params.Site, params.Bucket, params.Object); err == nil {
			msg := fmt.Sprintf("File %s/%s/%s deleted successfully", params.Site, params.Bucket, params.Object)
			c.JSON(200, gin.H{"status": "ok", "msg": msg})
		} else {
			c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
		}
	} else {
		c.JSON(400, gin.H{"status": "fail", "error": err.Error()})
	}
}
