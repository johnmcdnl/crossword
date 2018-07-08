package crossword

import (
	"sort"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Crossword struct {
	Words     []string    `json:"words"`
	NotPlaced []string    `json:"not_placed"`
	Grid      *[][]string `json:"grid"`
}

func New(size int) *Crossword {
	logrus.Debugln(`(c *Crossword) New`)

	var grid [][]string

	for i := 0; i <= size; i++ {
		grid = append(grid, make([]string, size))
	}

	crossword := &Crossword{
		Words: []string{"hello", "goodbye", "hate", "hat"},
		Grid:  &grid,
	}

	sort.Slice(crossword.Words, func(i, j int) bool {
		return len(crossword.Words[i]) > len(crossword.Words[j])
	})

	crossword.Generate()

	return crossword
}

func (c *Crossword) String() string {
	logrus.Debugln(`(c *Crossword) String`)

	for ri, row := range *c.Grid {
		for ci, col := range row {
			if col == "" {

				col = fmt.Sprintf(" %02d,%02d ", ri, ci)
			} else {
				col = "   " + col + "   "
			}
			fmt.Print(col)
		}
		fmt.Println()
	}
	j, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (c *Crossword) Generate() {
	for _, w := range c.Words {
		c.Insert(w)
	}
}

func (c *Crossword) Insert(word string) {
	logrus.Debugln(`(c *Crossword) Insert`)

	if c.IsEmpty() {
		middleRow := len(*c.Grid) / 2
		grid := *c.Grid
		middleCol := len(grid[middleRow]) / 2
		colIndex := middleCol - len(word)/2
		fmt.Println(middleRow, middleCol, colIndex)
		c.InsertHorizontal(middleRow, colIndex, word)
		return
	}

	if !c.FindSpace(word) {
		logrus.Errorln("no space for word ", word)
		c.NotPlaced = append(c.NotPlaced, word)
	}

}

func (c *Crossword) IsEmpty() bool {
	logrus.Debugln(`(c *Crossword) IsEmpty`)

	for _, row := range *c.Grid {
		for _, col := range row {
			if col != "" {
				return false
			}
		}
	}
	return true
}

func (c *Crossword) InsertHorizontal(row int, col int, word string) {
	logrus.Debugln(`(c *Crossword) InsertHorizontal`)

	for i, char := range word {
		grid := *c.Grid
		grid[row][col+i] = string(char)
		c.Grid = &grid
	}
}

func (c *Crossword) InsertVertical(row int, col int, word string) {
	logrus.Debugln(`(c *Crossword) InsertVertical`, row, col, word)

	for i, char := range word {
		grid := *c.Grid
		grid[row+i][col] = string(char)
		c.Grid = &grid
	}
}

func (c *Crossword) FindSpace(word string) (valid bool) {
	logrus.Debugln(`(c *Crossword) FindSpace`)

	letters := strings.Split(word, "")

	if randBool() {
		if valid := c.FindSpaceVertical(letters); valid {
			return true
		}

		if valid := c.FindSpaceHorizontal(letters); valid {
			return true
		}
	} else {
		if valid := c.FindSpaceHorizontal(letters); valid {
			return true
		}

		if valid := c.FindSpaceVertical(letters); valid {
			return true
		}
	}

	return false
}

func (c *Crossword) FindSpaceHorizontal(letters []string) (valid bool) {
	logrus.Debugln(`(c *Crossword) FindSpaceVertical`)
	points := c.FindIntersectionPoints(letters)
	logrus.Infoln(points)

	var validPoints []IntersectPoint

	for _, point := range points {
		var insertPoints []int

		for i := -point.LetterIndex; i < (len(letters) - point.LetterIndex); i++ {
			insertPoints = append(insertPoints, point.ColIndex+i)
		}

		var isValid = true
		for i, col := range insertPoints {
			grid := *c.Grid
			if col != point.ColIndex {
				if grid[point.RowIndex+1][col] != "" {
					isValid = false
					break
				}

				if grid[point.RowIndex-1][col] != "" {
					isValid = false
					break
				}
			}

			if grid[point.RowIndex][col] == "" {
				continue
			}

			if grid[point.RowIndex][col] == letters[i] {
				continue
			}

			isValid = false
			break
		}

		if isValid {
			point.insertPoints = insertPoints
			validPoints = append(validPoints, point)
		}
	}

	if len(validPoints) == 0 {
		return false
	}

	c.InsertHorizontal(
		randomPoint(validPoints).RowIndex,
		randomPoint(validPoints).insertPoints[0],
		randomPoint(validPoints).Word,
	)

	return true
}

func (c *Crossword) FindSpaceVertical(letters []string) (valid bool) {
	logrus.Debugln(`(c *Crossword) FindSpaceVertical`)
	points := c.FindIntersectionPoints(letters)
	logrus.Infoln(points)

	var validPoints []IntersectPoint

	for _, point := range points {

		var insertPoints []int

		for i := -point.LetterIndex; i < (len(letters) - point.LetterIndex); i++ {
			insertPoints = append(insertPoints, point.RowIndex+i)
		}

		var isValid = true
		for i, r := range insertPoints {
			grid := *c.Grid

			if r != point.RowIndex {
				if grid[r][point.ColIndex+1] != "" {
					isValid = false
					break
				}

				if grid[r][point.ColIndex-1] != "" {
					isValid = false
					break
				}
			}

			if grid[r][point.ColIndex] == "" {
				continue
			}

			if grid[r][point.ColIndex] == letters[i] {
				continue
			}

			isValid = false
			break
		}

		if isValid {
			point.insertPoints = insertPoints
			validPoints = append(validPoints, point)
		}
	}

	if len(validPoints) == 0 {
		return false
	}

	c.InsertVertical(
		randomPoint(validPoints).insertPoints[0],
		randomPoint(validPoints).ColIndex,
		randomPoint(validPoints).Word,
	)

	return true
}

func randomPoint(points []IntersectPoint) IntersectPoint {

	return points[rand.Intn(len(points))]
}

type IntersectPoint struct {
	LetterIndex  int
	RowIndex     int
	ColIndex     int
	Word         string
	insertPoints []int
}

func (i *IntersectPoint) String() string {
	j, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (c *Crossword) FindIntersectionPoints(letters []string) []IntersectPoint {
	logrus.Debugln(`(c *Crossword) FindIntersectionPoints`)

	var points []IntersectPoint
	for letterIndex, letter := range letters {
		grid := *c.Grid
		for rowIndex, row := range grid {
			for colIndex, col := range row {
				if letter == col {
					logrus.Infoln("intersection", col, letter, letterIndex, rowIndex, colIndex)
					intersectPoint := IntersectPoint{
						LetterIndex: letterIndex,
						RowIndex:    rowIndex,
						ColIndex:    colIndex,
						Word:        strings.Join(letters, ""),
					}
					points = append(points, intersectPoint)
					logrus.Infoln("IntersectPoint", intersectPoint.String())
				}
			}
		}
	}

	return points
}

func randBool() bool {
	return rand.Float32() < 0.5
}
