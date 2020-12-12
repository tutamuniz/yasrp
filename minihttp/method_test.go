package minihttp

import "testing"

func Test(t *testing.T) {
	testCases := []struct {
		desc     string
		method   string
		expected bool
	}{
		{
			desc:     "GET",
			method:   "get",
			expected: true,
		},
		{
			desc:     "POST",
			method:   "post",
			expected: true,
		},
		{
			desc:     "PUT",
			method:   "put",
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := isValidMethod(tC.method); got != tC.expected {
				t.Errorf("Invalid method:%v %s, %s", got, tC.method, ErrNotImplementedMethod.Error())
			}
		})
	}
}
