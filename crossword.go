package crossword

import (
	"sort"
	"fmt"
	"math/rand"
	"time"
	"encoding/json"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Crossword struct {
	Words     []string    `json:"words"`
	Placed    []string    `json:"placed"`
	NotPlaced []string    `json:"not_placed"`
	Grid      *[][]string `json:"-"`
}

type Coordinate struct {
	RowIndex int `json:"row"`
	ColIndex int `json:"col"`

	Word      string `json:"word,omitempty"`
	Letter    string `json:"letter,omitempty"`
	WordIndex int    `json:"word_index,omitempty"`

	Direction Direction     `json:"direction"`
	Placement []*Coordinate `json:"placement,omitempty"`
}

type Direction int

const (
	Vertical   = Direction(iota)
	Horizontal
)

func New(size int) *Crossword {

	var grid [][]string

	for i := 0; i <= size; i++ {
		grid = append(grid, make([]string, size))
	}

	crossword := &Crossword{
		Words: []string{"hello", "goodbye"},
		Grid:  &grid,
	}

	sort.Slice(crossword.Words, func(i, j int) bool {
		return len(crossword.Words[i]) > len(crossword.Words[j])
	})

	crossword.Placed = []string{}
	crossword.NotPlaced = crossword.Words

	crossword.Generate()

	return crossword
}

func (crossword *Crossword) String() string {

	for ri, row := range *crossword.Grid {
		for ci, col := range row {
			if col == "" {

				col = fmt.Sprintf(" %02d,%02d ", ri, ci)
				//col = fmt.Sprint("   ")
			} else {
				col = "   " + col + "   "
			}
			fmt.Print(col)
		}
		fmt.Println()
	}
	j, _ := json.Marshal(crossword)
	return string(j)
}

func (crossword *Crossword) Generate() {
	crossword.InsertWords()
}

func (crossword *Crossword) InsertWords() {
	for _, word := range crossword.NotPlaced {
		crossword.InsertWord(word)
	}
}

func (crossword *Crossword) InsertWord(word string) {
	if len(crossword.Placed) == 0 {
		crossword.InsertFirst(word)
		crossword.UpdatePlacedWords(word)
		return
	}

	intersectionPoints := crossword.FindIntersectPoints(word)
	verticalPoints := crossword.ValidVerticalPoints(intersectionPoints)
	horizontalPoints := crossword.ValidHorizontalPoints(intersectionPoints)

	points := append(verticalPoints, horizontalPoints...)

	fmt.Println(toString(points))
}

func (crossword *Crossword) UpdatePlacedWords(word string) {
	crossword.Placed = append(crossword.Placed, word)
	i := SliceIndex(len(crossword.NotPlaced), func(i int) bool { return crossword.NotPlaced[i] == word })
	crossword.NotPlaced = append(crossword.NotPlaced [:i], crossword.NotPlaced [i+1:]...)
}

func (crossword *Crossword) InsertFirst(word string) {
	grid := *crossword.Grid
	centerRow := len(grid) / 2
	centerCol := len(grid[centerRow]) / 2
	startPos := centerCol - len(word)/2

	var coordinates []*Coordinate
	for i, char := range word {
		coordinates = append(coordinates, &Coordinate{
			RowIndex: centerRow,
			ColIndex: startPos + i,
			Letter:   string(char),
		})
	}

	crossword.Insert(coordinates)
}

func (crossword *Crossword) Insert(coordinates []*Coordinate) {
	grid := *crossword.Grid
	for _, coordinate := range coordinates {
		grid[coordinate.RowIndex][coordinate.ColIndex] = coordinate.Letter
	}
	crossword.Grid = &grid
}

func (crossword *Crossword) FindIntersectPoints(word string) []*Coordinate {

	grid := *crossword.Grid
	var coordinates []*Coordinate
	for letterIndex, char := range word {
		for rowIndex, row := range grid {
			for colIndex, col := range row {
				if string(char) == col {
					coordinates = append(coordinates, &Coordinate{
						RowIndex:  rowIndex,
						ColIndex:  colIndex,
						Letter:    string(char),
						WordIndex: letterIndex,
						Word:      word,
					})
				}
			}
		}
	}
	return coordinates
}

func (crossword *Crossword) ValidVerticalPoints(points []*Coordinate) []*Coordinate {
	var validPoints []*Coordinate
	for _, point := range points {
		if validPoint := crossword.ValidVerticalPoint(point); validPoint != nil {
			validPoint.Direction = Vertical
			validPoints = append(validPoints, validPoint)
			return validPoints
		}
	}
	return validPoints
}

func (crossword *Crossword) ValidHorizontalPoints(points []*Coordinate) []*Coordinate {
	var validPoints []*Coordinate
	for _, point := range points {
		if validPoint := crossword.ValidHorizontalPoint(point); validPoint != nil {
			validPoint.Direction = Horizontal
			validPoints = append(validPoints, validPoint)
			return validPoints
		}
	}
	return validPoints
}

func (crossword *Crossword) ValidHorizontalPoint(point *Coordinate) *Coordinate {
	return nil
}

func (crossword *Crossword) ValidVerticalPoint(point *Coordinate) *Coordinate {

	grid := *crossword.Grid
	var placements []*Coordinate

	for charIndex, char := range point.Word {
		fmt.Println(charIndex, string(char))
	}

	for i := -point.WordIndex; i < len(point.Word)-point.WordIndex; i++ {
		j := i + point.WordIndex

		point.Placement = make([]*Coordinate, 0)
		row := point.RowIndex + i
		col := point.ColIndex

		fmt.Println(row, col, j)

		// the exact line
		currentChar := grid[row][col]
		if currentChar != "" {
			if currentChar != string(point.Word[j]) {
				return nil
			}
		}

		// either side of the line
		leftChar := grid[row][col-1]
		if leftChar != "" {
			if leftChar != string(point.Word[j]) {
				return nil
			}
		}

		newC := &Coordinate{
			RowIndex: row,
			ColIndex: col,
			Letter:   string(point.Word[j]),
		}
		fmt.Println(newC)
		placements = append(placements, newC)

	}

	point.Placement = placements
	crossword.Grid = &grid
	return point
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func toString(i interface{}) string {
	j, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(j)
}
