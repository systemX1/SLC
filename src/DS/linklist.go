package DS

type ElemType interface {
}

type LinkNode struct {
	Val  ElemType
	Next *LinkNode
}

type LinkList interface {
	Add(val ElemType)
	Delete(index int) ElemType
	Insert(index int, val ElemType)
	GetLength() int
	Search(val ElemType) int
	GetAll(index int) ElemType
	Reverse() *LinkNode
}


