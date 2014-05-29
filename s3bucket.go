// A library for accessing Amazon S3 from Go
package s34go


// Represents a specific bucket, accessed via a specific service
type S3Bucket struct {
    name        string
    service     *S3Service
}


// returns the name of this bucket, in case String() isn't sufficient
func (bucket *S3Bucket) Name() string {
    return bucket.name;
}


// standard function for formatting an object -- intended for debugging
func (bucket *S3Bucket) String() string {
    return "S3Bucket(" + bucket.service.String() + ", " + bucket.name + ")"
}


// returns a listing of all objects within the bucket
func (bucket *S3Bucket) ListObjects() (*[]S3Object,error) {
    return nil,S3Error{"not implemented yet"}
}


// returns a listing of all objects within the bucket that have the specified leading path components
func (bucket *S3Bucket) ListObjectsWithin(path string) (*[]S3Object,error) {
    return nil,S3Error{"not implemented yet"}
}


// creates an object reference for the specified name; does not make a server call
func (bucket *S3Bucket) NewObject(name string) (*S3Object,error) {
    return &S3Object{name, bucket},nil
}


// deletes this bucket; will fail if the bucket is non-empty
func (bucket *S3Bucket) Delete() error {
    return S3Error{"not implemented yet"}
}


// deletes this bucket, after first deleting all of the objects it contains
func (bucket *S3Bucket) DeleteCascade() error {
    return S3Error{"not implemented yet"}
}
