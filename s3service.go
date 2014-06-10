// A library for accessing Amazon S3 from Go
package s34go

import "net/http"

// Represents a connection to S3 as a specific user
type S3Service struct {
    endpoint    string
    accessKey   string
    secretKey   string
    client      *http.Client;
}


// Creates a new service instance. Does not attempt to verify that the instance can
// be used to make a request.
func NewS3Service(accessKey string, secretKey string) (*S3Service,error) {
    return &S3Service{DEFAULT_S3_ENDPOINT, accessKey, secretKey, &http.Client{}}, nil
}


// returns the endpoint of this service
func (service *S3Service) Endpoint() string {
    return service.endpoint;
}


// standard function for formatting an object -- intended for debugging
func (service *S3Service) String() string {
    return "S3Service(" + service.endpoint + ")"
}


// Returns a list of the buckets accessible from this service; each bucket
// can then be used to access the objects it contains
func (service *S3Service) ListBuckets() ([]*S3Bucket,error) {
    request := service.newS3Request().prepare()
    response,err := request.execute(service.client)
    if err != nil {
        return nil, err
    }

    defer response.Close()

    // FIXME - process list of buckets
    return make([]*S3Bucket,0),nil
}


// Returns a reference to the bucket with the given name, if it exists, nil otherwise
func (service *S3Service) GetBucket(bucketName string) (*S3Bucket,error) {
    return nil,S3Error{"not implemented yet"}
}


// Returns a reference to the bucket with the given name, creating it if necessary
func (service *S3Service) GetOrCreateBucket(bucketName string) (*S3Bucket,error) {
    return nil,S3Error{"not implemented yet"}
}


//----------------------------------------------------------------------------------------------
//  Package methods
//----------------------------------------------------------------------------------------------

// Creates a new request from this service; this method exists primarily so that buckets and
// objects won't be required to dig into the service.
func (service *S3Service) newS3Request() *S3Request {
    return NewS3Request(service.accessKey, service.secretKey)
}


//----------------------------------------------------------------------------------------------
//  Internals
//----------------------------------------------------------------------------------------------
