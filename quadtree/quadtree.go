package quadtree

import "github.com/alldroll/spatial-index/geometry"

type QuadTree struct {
	root   *node
	length int
}

type node struct {
	*shape.BoundaryBox
	points   []*shape.Point
	children [totalChild]*node
	length   int
	level    int
	capacity int
}

const (
	totalChild = 4
)

type QuadTreeError struct {
	msg string
}

func (e *QuadTreeError) Error() string {
	return e.msg
}

func NewQuadTree(x1, y1, x2, y2 float64, capacity int) (*QuadTree, error) {
	if x1 > x2 || y1 > y2 {
		return nil, &QuadTreeError{"Invalid Points for BoundaryBox construct"}
	}

	global := shape.NewBoundaryBox(
		shape.NewPoint(x1, y1),
		shape.NewPoint(x2, y2),
	)
	root := newNode(global, 0, capacity)
	return &QuadTree{root, 0}, nil
}

func (qt *QuadTree) Insert(x, y float64) bool {
	p := shape.NewPoint(x, y)
	return qt.InsertPoint(p)
}

func (qt *QuadTree) InsertPoint(point *shape.Point) bool {
	res := qt.root.insertPoint(point)
	if res {
		qt.length += 1
	}

	return res
}

func (qt *QuadTree) GetPoints(x1, y1, x2, y2 float64) ([]*shape.Point, error) {
	if x1 > x2 || y1 > y2 {
		return nil, &QuadTreeError{"Invalid Points for BoundaryBox construct"}
	}

	area := shape.NewBoundaryBox(shape.NewPoint(x1, y1), shape.NewPoint(x2, y2))
	return qt.root.getPointsFromArea(area), nil
}

func (qt *QuadTree) GetLength() int {
	return qt.length
}

func newNode(box *shape.BoundaryBox, level int, capacity int) *node {
	return &node{
		box,
		[]*shape.Point{},
		[totalChild]*node{nil, nil, nil, nil},
		0,
		level,
		capacity,
	}
}

func (self *node) insertPoint(point *shape.Point) bool {
	if !self.ContainsPoint(point) {
		return false
	}

	if self.length < self.capacity {
		self.points = append(self.points, point)
		self.length += 1
		return true
	}

	if self.isLeaf() {
		self.splitNode()
	}

	success := false
	for i := 0; i < totalChild && !success; i++ {
		success = self.children[i].insertPoint(point)
	}

	return success
}

func (self *node) splitNode() {
	boxes := self.Quarter()
	nlevel := self.level + 1
	capacity := self.capacity
	for i := 0; i < totalChild; i++ {
		self.children[i] = newNode(boxes[i], nlevel, capacity)
	}
}

func (self *node) getPointsFromArea(area *shape.BoundaryBox) []*shape.Point {
	if !self.Intersect(area) {
		return []*shape.Point{}
	}

	result := []*shape.Point{}
	for _, point := range self.getPoints() {
		if area.ContainsPoint(point) {
			result = append(result, point)
		}
	}

	if !self.isLeaf() {
		children := self.children
		for i := 0; i < totalChild; i++ {
			result = append(result, children[i].getPointsFromArea(area)...)
		}
	}

	return result
}

func (self *node) getPoints() []*shape.Point {
	return self.points
}

func (self *node) isLeaf() bool {
	return self.children[0] == nil
}
