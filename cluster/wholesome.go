package cluster

//
//import "github.com/olesho/classify"
//
//func wholesomeVolume(a *classify.Arena, matrix *RateMatrix, members []int) float32 {
//	if len(members) > 0 {
//
//		m := make([][]int, len(members))
//		for i, memberIdx := range members {
//			m[i] = a.Wholesome(memberIdx)
//		}
//
//		identicalCnt := 0
//		wholesomeCnt := 0
//		for _, idx := range m[0] {
//			if allMatchingWholesome(matrix, m[1:], idx) {
//				wholesomeCnt++
//				if allIdenticalWholesome(a, m[1:], idx) {
//					identicalCnt++
//				}
//
//			}
//		}
//
//		return float32((wholesomeCnt-identicalCnt) * len(members))
//	}
//	return 0
//}
//
//func allMatchingWholesome(matrix *RateMatrix, columns [][]int, idx int) bool {
//	for _, col := range columns {
//		if !hasMatchingWholesomeIn(matrix, col, idx) {
//			return false
//		}
//	}
//	return true
//}
//
//func hasMatchingWholesomeIn(matrix *RateMatrix, indexes []int, idx int) bool {
//	for _, idx2 := range indexes {
//		if matrix.Cmp(idx, idx2) > 0 {
//			return true
//		}
//	}
//	return false
//}
//
//func allIdenticalWholesome(a *classify.Arena, columns [][]int, idx int) bool {
//	for _, col := range columns {
//		if !hasIdenticalWholesomeIn(a, col, idx) {
//			return false
//		}
//	}
//	return true
//}
//
//func hasIdenticalWholesomeIn(a *classify.Arena, indexes []int, idx int) bool {
//	for _, idx2 := range indexes {
//		info1, _ := a.Get(idx).WholesomeInfo()
//		info2, _ := a.Get(idx2).WholesomeInfo()
//		if info1 == info2 {
//			return true
//		}
//	}
//	return false
//}
