package auto

type AutoRecord struct {
	IsListen chan bool
}

func NewAutoRecord() *AutoRecord {
	ans := new(AutoRecord)
	ans.IsListen = make(chan bool)
	return ans
}

func (ar *AutoRecord) OnListen(isListen bool) {
	ar.IsListen <- isListen
}
