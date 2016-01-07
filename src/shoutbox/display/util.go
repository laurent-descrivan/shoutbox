package display

func sameRows(rowsa, rowsb [][]bool) bool {
	for i := range rowsa {
		rowa := rowsa[i]
		rowb := rowsb[i]
		for j := range rowa {
			if rowa[j] != rowb[j] {
				return false
			}
		}
	}
	return true
}
