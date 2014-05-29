// A library for accessing Amazon S3 from Go
package s34go

// A common error return
type S3Error struct {
    desc    string
}

func (error S3Error) Error() string {
    return error.desc
}
