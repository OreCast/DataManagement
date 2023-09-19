package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	minio "github.com/minio/minio-go/v7"
	cryptoutils "github.com/vkuznet/cryptoutils"
)

// Site represents Site object returned from discovery service
type Site struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// helper function to fetch sites info from discovery service
func sites() []Site {
	var out []Site
	rurl := fmt.Sprintf("%s/sites", Config.DiscoveryURL)
	if Config.Verbose > 0 {
		log.Println("query Discovery service", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		log.Println("ERROR:", err)
		return out
	}
	defer resp.Body.Close()
	var results []Site
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&results); err != nil {
		log.Println("ERROR:", err)
		return out
	}
	return results
}

// type BucketObject struct {
//     Name string `json:"name"`
// }

// SiteObject represents site object
type SiteObject struct {
	Site    string             `json:"site"`
	Buckets []minio.BucketInfo `json:"buckets"`
}

// BucketObject represents site object
type BucketObject struct {
	Site    string             `json:"site"`
	Bucket  string             `json:"bucket"`
	Objects []minio.ObjectInfo `json:"objects"`
}

// DiscoveryRecord represents structure of data discovery record
type DiscoveryRecord struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	Endpoint     string `json:"endpoint""`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	UseSSL       bool   `json:"use_ssl"`
}

// helper function to return S3 object for given site
func site2s3(site string) (S3, error) {
	var s3 S3
	if Config.Verbose > 0 {
		log.Printf("looking for site:%s", site)
	}
	surl := fmt.Sprintf("%s/sites", Config.DiscoveryURL)
	if Config.Verbose > 0 {
		log.Println("query Discovery service", surl)
	}
	resp, err := http.Get(surl)
	if err != nil {
		log.Printf("ERROR: unable to contact DataDiscovery service %s, error %v", surl, err)
		return s3, err
	}
	// read data discovery content
	var records []DiscoveryRecord
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: unable to read DataDiscovery response, error %v", err)
		return s3, err
	}
	err = json.Unmarshal(body, &records)
	if err != nil {
		log.Printf("ERROR: unable to unmarshal DataDiscovery response, error %v", err)
		return s3, err
	}
	if Config.Verbose > 0 {
		log.Printf("site records %+v", records)
	}

	for _, rec := range records {
		if rec.Name == site {
			log.Printf("INFO: found %s in DataDiscovery records, will access its s3 via %s", rec.Name, rec.URL)
			// bingo: we got desired site, now we can query its s3 storage for datasets
			if Config.Verbose > 0 {
				log.Println("###", rec.AccessKey, Config.DiscoveryPassword, Config.DiscoveryCipher)
			}
			akey, err := cryptoutils.HexDecrypt(rec.AccessKey, Config.DiscoveryPassword, Config.DiscoveryCipher)
			if err != nil {
				log.Printf("ERROR: unable to decrypt data discovery access key, error %v", err)
				return s3, nil

			}
			apwd, err := cryptoutils.HexDecrypt(rec.AccessSecret, Config.DiscoveryPassword, Config.DiscoveryCipher)
			if err != nil {
				log.Printf("ERROR: unable to decrypt data discovery acess secret, error %v", err)
				return s3, nil

			}
			s3 = S3{
				Endpoint:     rec.Endpoint,
				AccessKey:    string(akey),
				AccessSecret: string(apwd),
				UseSSL:       rec.UseSSL,
			}
			return s3, nil
		}
	}
	return s3, errors.New("No matching site found")
}

// siteContent provides content on given site
func siteContent(site string) (SiteObject, error) {
	if Config.Verbose > 0 {
		log.Printf("looking for site:%s", site)
	}
	surl := fmt.Sprintf("%s/sites", Config.DiscoveryURL)
	if Config.Verbose > 0 {
		log.Println("query Discovery service", surl)
	}
	resp, err := http.Get(surl)
	var siteObj SiteObject
	if err != nil {
		log.Printf("ERROR: unable to contact DataDiscovery service %s, error %v", surl, err)
		return siteObj, err
	}
	// read data discovery content
	var records []DiscoveryRecord
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: unable to read DataDiscovery response, error %v", err)
		return siteObj, err
	}
	err = json.Unmarshal(body, &records)
	if err != nil {
		log.Printf("ERROR: unable to unmarshal DataDiscovery response, error %v", err)
		return siteObj, err
	}
	if Config.Verbose > 0 {
		log.Printf("site records %+v", records)
	}

	for _, rec := range records {
		if rec.Name == site {
			log.Printf("INFO: found %s in DataDiscovery records, will access its s3 via %s", rec.Name, rec.URL)
			// bingo: we got desired site, now we can query its s3 storage for datasets
			if Config.Verbose > 0 {
				log.Println("###", rec.AccessKey, Config.DiscoveryPassword, Config.DiscoveryCipher)
			}
			akey, err := cryptoutils.HexDecrypt(rec.AccessKey, Config.DiscoveryPassword, Config.DiscoveryCipher)
			if err != nil {
				log.Printf("ERROR: unable to decrypt data discovery access key, error %v", err)
				return siteObj, err

			}
			apwd, err := cryptoutils.HexDecrypt(rec.AccessSecret, Config.DiscoveryPassword, Config.DiscoveryCipher)
			if err != nil {
				log.Printf("ERROR: unable to decrypt data discovery acess secret, error %v", err)
				return siteObj, err

			}
			s3 := S3{
				Endpoint:     rec.Endpoint,
				AccessKey:    string(akey),
				AccessSecret: string(apwd),
				UseSSL:       rec.UseSSL,
			}
			if Config.Verbose > 0 {
				log.Printf("INFO: accessing %+v", s3)
			}
			buckets, err := listBuckets(s3)
			if err != nil {
				log.Printf("ERROR: unabel to list buckets at %s, error %v", site, err)
			}
			obj := SiteObject{
				Site:    site,
				Buckets: buckets,
			}
			return obj, nil
		}
	}
	return siteObj, errors.New("Not reachable")
}

// bucketContent provides content on given site and bucket
func bucketContent(site, bucket string) (BucketObject, error) {
	if Config.Verbose > 0 {
		log.Printf("looking for site:%s bucket:%s", site, bucket)
	}
	surl := fmt.Sprintf("%s/sites", Config.DiscoveryURL)
	if Config.Verbose > 0 {
		log.Println("query Discovery service", surl)
	}
	resp, err := http.Get(surl)
	var bucketObj BucketObject
	if err != nil {
		log.Printf("ERROR: unable to contact DataDiscovery service %s, error %v", surl, err)
		return bucketObj, err
	}
	// read data discovery content
	var records []DiscoveryRecord
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: unable to read DataDiscovery response, error %v", err)
		return bucketObj, err
	}
	err = json.Unmarshal(body, &records)
	if err != nil {
		log.Printf("ERROR: unable to unmarshal DataDiscovery response, error %v", err)
		return bucketObj, err
	}
	if Config.Verbose > 0 {
		log.Printf("site records %+v", records)
	}

	for _, rec := range records {
		if rec.Name == site {
			log.Printf("INFO: found %s in DataDiscovery records, will access its s3 via %s", rec.Name, rec.URL)
			// bingo: we got desired site, now we can query its s3 storage for datasets
			if Config.Verbose > 0 {
				log.Println("###", rec.AccessKey, Config.DiscoveryPassword, Config.DiscoveryCipher)
			}
			akey, err := cryptoutils.HexDecrypt(rec.AccessKey, Config.DiscoveryPassword, Config.DiscoveryCipher)
			if err != nil {
				log.Printf("ERROR: unable to decrypt data discovery access key, error %v", err)
				return bucketObj, err

			}
			apwd, err := cryptoutils.HexDecrypt(rec.AccessSecret, Config.DiscoveryPassword, Config.DiscoveryCipher)
			if err != nil {
				log.Printf("ERROR: unable to decrypt data discovery acess secret, error %v", err)
				return bucketObj, err

			}
			s3 := S3{
				Endpoint:     rec.Endpoint,
				AccessKey:    string(akey),
				AccessSecret: string(apwd),
				UseSSL:       rec.UseSSL,
			}
			if Config.Verbose > 0 {
				log.Printf("INFO: accessing %+v", s3)
			}
			objects, err := listObjects(s3, bucket)
			if err != nil {
				log.Printf("ERROR: unabel to list objects at %s/%s, error %v", site, bucket, err)
			}
			obj := BucketObject{
				Site:    site,
				Bucket:  bucket,
				Objects: objects,
			}
			return obj, nil
		}
	}
	return bucketObj, errors.New("Not reacheable")
}
