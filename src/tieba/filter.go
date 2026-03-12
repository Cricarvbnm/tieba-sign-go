package tieba

func ForumsFilterUnsigned(forums []Forum) []Forum {
	var result []Forum
	for _, forum := range forums {
		if forum.IsSign == 0 {
			result = append(result, forum)
		}
	}
	return result
}
