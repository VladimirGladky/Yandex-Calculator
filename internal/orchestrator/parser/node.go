package parser

type ExpressionNode struct {
	IsLeaf        bool
	Value         float64
	Operator      string
	Left, Right   *ExpressionNode
	TaskScheduled bool
}
