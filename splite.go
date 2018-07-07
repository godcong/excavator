package excavator

//
//func StringSplite(base, sta, end string, keep bool) []string {
//	sa := make([]string, 0)
//	size := len(base)
//	for i := 0; i < size; {
//		s := strings.Index(base, sta)
//		if s < 0 {
//			break
//		}
//		if !keep {
//			base = base[s+len(sta):]
//		} else {
//			base = base[s:]
//		}
//
//		e := strings.Index(base, end)
//		if e < 0 {
//			break
//		}
//		if keep {
//			e = e + len(end)
//		}
//
//		sa = append(sa, base[0:e])
//
//		i += e
//		base = base[e:]
//	}
//	return sa
//}
