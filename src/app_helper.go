package src

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (a *App) LoadOBJFile(fileName string) ([]*Vec3, []*Face, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	vertices := []*Vec3{}
	faces := []*Face{}
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		if strings.HasPrefix(line, "v") {
			vals := strings.Split(line, " ")
			x, _ := strconv.ParseFloat(vals[1], 64)
			y, _ := strconv.ParseFloat(vals[2], 64)
			z, _ := strconv.ParseFloat(vals[3], 64)
			vertices = append(vertices, &Vec3{X: x, Y: y, Z: z})
		} else if strings.HasPrefix(line, "f") {
			vals := strings.Split(line, " ")

			tmpA := strings.Split(vals[1], "/")[0]
			a, err := strconv.ParseInt(tmpA, 10, 32)
			if err != nil {
				return nil, nil, err
			}
			tmpB := strings.Split(vals[2], "/")[0]
			b, err := strconv.ParseInt(tmpB, 10, 32)
			if err != nil {
				return nil, nil, err
			}
			tmpC := strings.Split(vals[3], "/")[0]
			c, err := strconv.ParseInt(tmpC, 10, 32)
			if err != nil {
				return nil, nil, err
			}

			faces = append(faces, &Face{A: int(a), B: int(b), C: int(c)})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return vertices, faces, nil
}
