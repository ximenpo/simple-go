package simple

type HandleID uint32

type HandleMap struct {
	data     map[HandleID]interface{} // 信息列表
	seq_seed HandleID                 // 序号种子
}

func NewHandleMap() *HandleMap {
	return &HandleMap{make(map[HandleID]interface{}), 0}
}

func (m *HandleMap) Add(info interface{}) (seq HandleID) {
	seq = m.newSeq()
	m.data[seq] = info
	return
}

func (m *HandleMap) Set(seq HandleID, info interface{}) {
	m.data[seq] = info
}

func (m *HandleMap) Get(seq HandleID) interface{} {
	if info, ok := m.data[seq]; ok {
		return info
	}
	return nil
}

func (m *HandleMap) Del(seq HandleID) interface{} {
	if info, ok := m.data[seq]; ok {
		delete(m.data, seq)
		return info
	}
	return nil
}

func (m *HandleMap) Clear(seq HandleID) {
	m.data = make(map[HandleID]interface{})
	m.seq_seed = 0
}

// 永远不会返回 0 或 0xFFFFFFFF(-1)
func (m *HandleMap) newSeq() HandleID {
	m.seq_seed++
	seq := m.seq_seed
	for {
		if _, ok := m.data[seq]; !ok {
			break
		}
		seq++
	}
	m.seq_seed = seq

	// 0/-1 为非法序号
	for (0 == seq) || (0xFFFFFFFF == seq) {
		seq = m.newSeq()
	}
	return seq
}
