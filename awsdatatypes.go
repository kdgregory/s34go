// A library for accessing Amazon S3 from Go
package s34go


// holder for raw AWS error information
type AWSError struct {
	Code                string
	Message             string
	HostId              string
	RequestId           string
    Resource            string
    AWSAccessKeyId      string
	ArgumentName        string
	ArgumentValue       string
    CanonicalRequest    string
    StringToSign        string
}
