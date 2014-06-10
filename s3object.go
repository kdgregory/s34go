// A library for accessing Amazon S3 from Go
package s34go

import "io"

// Represents an object within a specific bucket
type S3Object struct {
    bucket      *S3Bucket
    path         string
}


// returns the path of this object
func (object *S3Object) Path() string {
    return object.path;
}


// returns the URL that may be used to retrieve this object via HTTP
func (object *S3Object) Url() string {
    return "not implemented yet";
}


// standard function for formatting an object -- intended for debugging
func (object *S3Object) String() string {
    return "S3Object(" + object.bucket.String() + ", " + object.path + ")"
}


// retrieves this object's content into a byte array
func (object *S3Object) Read() ([]byte,error) {
    return nil,S3Error{"not implemented yet"}
}


// retrieves this object's content and writes it to the passed writer
func (object *S3Object) ReadToWriter(w io.Writer) error {
    return S3Error{"not implemented yet"}
}


// writes the specified bytes to this object, replacing any previous content
// and creating the object if it doesn't alread exist
func (object *S3Object) Write([]byte) error {
    return S3Error{"not implemented yet"}
}


// writes the content of the passed reader, replacing any previous content
// and creating the object if it doesn't already exist
func (object *S3Object) WriteFromReader(r io.Reader) error {
    return S3Error{"not implemented yet"}
}


// deletes this object from the server
func (object *S3Object) Delete() error {
    return S3Error{"not implemented yet"}
}
