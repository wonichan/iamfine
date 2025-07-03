package utils

func IntPtr(i int32) *int32 {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func StringPtr(s string) *string {
	return &s
}
