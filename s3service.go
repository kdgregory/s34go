// A library for accessing Amazon S3 from Go
package s34go

// Represents a connection to S3 as a specific user
type S3Service struct {
    endpoint    string
    publicKey   string
    secretKey   string
}


// Creates a new instance of an S3Service, and verifies that it can be used (this
// implies a call to S3)
func NewS3Service(publicKey string, secretKey string) (*S3Service,error) {
    return &S3Service{DEFAULT_ENDPOINT, publicKey, secretKey}, nil
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
    return nil,S3Error{"not implemented yet"}
}


// Returns a reference to the bucket with the given name, if it exists, nil otherwise
func (service *S3Service) GetBucket(bucketName string) (*S3Bucket,error) {
    return nil,S3Error{"not implemented yet"}
}


// Returns a reference to the bucket with the given name, creating it if necessary
func (service *S3Service) GetOrCreateBucket(bucketName string) (*S3Bucket,error) {
    return nil,S3Error{"not implemented yet"}
}
