type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.
func (r MyReader) Read(b []byte) (int, error) {
	b[0] = 'A'
	fmt.Printf("%v", b[0])
	return 1, nil
}

func main() {
	reader.Validate(MyReader{})
}
