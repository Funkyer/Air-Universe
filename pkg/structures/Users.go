package structures

type UserInfo struct {
	Id      int
	Uuid    string
	AlertId uint32
	Level   uint32
	InTag   string
}

type UserTraffic struct {
	Id   int   `json:"user_id"`
	Up   int64 `json:"u"`
	Down int64 `json:"d"`
}

func FindUserDiffer(before, now *[]UserInfo) (remove, add *[]UserInfo, err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		remove = nil
	//		add = nil
	//		err = errors.New(fmt.Sprintf("model FindUserDiffer cause error - %s", r))
	//	}
	//}()

	remove = new([]UserInfo)
	add = new([]UserInfo)
	// 对于空的对象要处理下，因为会死循环
	if len(*before) == 0 {
		return nil, now, err
	} else if len(*now) == 0 {
		return before, nil, err
	}

	n := 0
	b := 0
	nLastAppear := false
	bLastAppear := false
	for true {
		if n == len(*now) {
			nLastAppear = true
			n--
		} else if b == len(*before) {
			bLastAppear = true
			b--
		} else if (*before)[b] == (*now)[n] {
			n++
			b++
		} else if (*before)[b].Id < (*now)[n].Id {
			// (*before)[b] has been removed
			*remove = append(*remove, (*before)[b])
			b++
		} else if (*before)[b].Id > (*now)[n].Id {
			// (*now)[n] has been inserted
			*add = append(*add, (*now)[n])
			n++
		} else if (*before)[b].Id == (*now)[n].Id && (*before)[b].Uuid != (*now)[n].Uuid {
			//user (*before)[b] changed uuid
			*remove = append(*remove, (*before)[b])
			*add = append(*add, (*now)[n])
			n++
			b++
			// Last one will tagged
			continue
		}
		// The last element has not been processed in the loop. we will process after loop.
		if n == len(*now)-1 && b == len(*before)-1 {
			break
		}
	}

	// Process last one
	if (*before)[len(*before)-1] != (*now)[len(*now)-1] {
		if nLastAppear == false {
			*add = append(*add, (*now)[len(*now)-1])
		}
		if bLastAppear == false {
			*remove = append(*remove, (*before)[len(*before)-1])
		}
	}

	return
}