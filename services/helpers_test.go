package services_test

func makeTestString(len int) string {
	b := make([]byte, len)
	for i := range b {
		b[i] = 't'
	}
	return string(b)
}
