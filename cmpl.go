package otto

import (
	"github.com/extrame/otto/ast"
	"github.com/extrame/otto/file"
)

type compiler struct {
	file    *file.File
	program *ast.Program
}
