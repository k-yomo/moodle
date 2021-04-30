package maputil

func MergeStrMap(m ...map[string]string) map[string]string {
	ans := make(map[string]string)

	for _, c := range m {
		for k, v := range c {
			ans[k] = v
		}
	}
	return ans
}
