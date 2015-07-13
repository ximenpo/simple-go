package simple

type HandleID uint32

type HandleMap struct {
	data     map[HandleID]interface{} // 信息列表
	seq_seed HandleID                 // 序号种子
}

func NewHandleMap() *HandleMap {
	return &HandleMap{make(map[HandleID]interface{}), 0}
}

func (self *HandleMap) Add(info interface{}) (seq HandleID) {
	seq = self.newSeq()
	self.data[seq] = info
	return
}

func (self *HandleMap) Set(seq HandleID, info interface{}) {
	self.data[seq] = info
}

func (self *HandleMap) Get(seq HandleID) interface{} {
	if info, ok := self.data[seq]; ok {
		return info
	}
	return nil
}

func (self *HandleMap) Del(seq HandleID) interface{} {
	if info, ok := self.data[seq]; ok {
		delete(self.data, seq)
		return info
	}
	return nil
}

func (self *HandleMap) Clear(seq HandleID) {
	self.data = make(map[HandleID]interface{})
	self.seq_seed = 0
}

// 永远不会返回 0 或 0xFFFFFFFF(-1)
func (self *HandleMap) newSeq() HandleID {
	self.seq_seed++
	seq := self.seq_seed
	for {
		if _, ok := self.data[seq]; !ok {
			break
		}
		seq++
	}
	self.seq_seed = seq

	// 0/-1 为非法序号
	for (0 == seq) || (0xFFFFFFFF == seq) {
		seq = self.newSeq()
	}
	return seq
}
