package src

type Mesh struct {
	Vertices []*Vec3
	Faces    []*Face
}

type Face struct {
	A, B, C int
}
