package s3

import "sort"

func SortFileHeadsByTime(fileHeads []FileHead) {
	sort.SliceStable(fileHeads, func(i, j int) bool {
		if fileHeads[i].LastModified.Equal(fileHeads[j].LastModified) {
			return fileHeads[i].Key < fileHeads[j].Key
		} else {
			return fileHeads[i].LastModified.Before(fileHeads[j].LastModified)
		}
	})
}

type FilterByFileHeadFunc func(fh *FileHead) bool
