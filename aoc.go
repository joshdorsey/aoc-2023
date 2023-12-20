package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func MustReadFileLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines := make([]string, 0, 100)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

// Day 1

func Day1() {
	fmt.Println("Day 1")
	lines := MustReadFileLines("day1.input")

	isDigit := func(r byte) bool {
		if r >= '0' && r <= '9' {
			return true
		}
		return false
	}

	{ // Part 1

		firstLastDigit := func(s string) (int, int) {
			first, last := -1, 0
			for i := range s {
				if isDigit(s[i]) {
					digit := int(s[i] - '0')

					last = digit
					if first == -1 {
						first = digit
					}
				}
			}

			return first, last
		}

		total := 0
		for _, line := range lines {
			first, last := firstLastDigit(line)
			calibration := first*10 + last
			total += calibration
		}

		fmt.Printf("\tPart 1: %d\n", total)

	}

	{ // Part 2

		initialDigit := func(s string) int {
			if len(s) == 0 {
				return -1
			}

			for val, str := range []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"} {
				if strings.HasPrefix(s, str) {
					return val
				}
			}
			if isDigit(s[0]) {
				return int(s[0] - '0')
			}

			return -1
		}

		total := 0
		for _, line := range lines {
			first, last := -1, 0

			s := line
			for len(s) != 0 {
				digit := initialDigit(s)

				if digit != -1 {
					last = digit
					if first == -1 {
						first = digit
					}
				}

				s = s[1:]
			}

			calibration := first*10 + last
			total += calibration
		}

		fmt.Printf("\tPart 2: %d\n", total)
	}
}

// Day 2

type GameParser struct {
	Line string
}

func (p *GameParser) IsEol() bool {
	return len(p.Line) == 0
}

func (p *GameParser) SkipWs() {
	if p.IsEol() {
		return
	}

	for unicode.IsSpace(rune(p.Line[0])) {
		p.Line = p.Line[1:]
	}
}

func (p *GameParser) ReadStr(str string) {
	p.SkipWs()

	if strings.HasPrefix(p.Line, str) {
		p.Line = p.Line[len(str):]
	}
}

func (p *GameParser) ReadNum() int {
	p.SkipWs()

	if !unicode.IsDigit(rune(p.Line[0])) {
		return -1
	}

	i := 0
	for ; unicode.IsDigit(rune(p.Line[i])); i++ {
	}

	numStr := p.Line[:i]
	p.Line = p.Line[i:]
	num, _ := strconv.ParseInt(numStr, 10, 64)
	return int(num)
}

func (p *GameParser) ReadColor() string {
	p.SkipWs()

	for _, color := range []string{"red", "green", "blue"} {
		if strings.HasPrefix(p.Line, color) {
			p.Line = p.Line[len(color):]
			return color
		}
	}

	return "unknown"
}

func (p *GameParser) ReadSep() byte {
	p.SkipWs()

	if p.IsEol() {
		return '\000'
	}

	for _, sep := range []byte{':', ',', ';'} {
		if p.Line[0] == sep {
			p.Line = p.Line[1:]
			return sep
		}
	}

	return '\000'
}

func (p *GameParser) ReadGame() Game {
	game := Game{}
	p.ReadStr("Game")
	game.Number = p.ReadNum()
	p.ReadSep()

	for !p.IsEol() {
		handful := Handful{}
		sep := byte('0')
		for !p.IsEol() && sep != ';' {
			num := p.ReadNum()
			color := p.ReadColor()

			switch color {
			case "red":
				handful.Red += num
			case "green":
				handful.Green += num
			case "blue":
				handful.Blue += num
			}

			sep = p.ReadSep()
		}

		game.Handfuls = append(game.Handfuls, handful)
	}

	return game
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Handful struct {
	Blue, Red, Green int
}

func (h Handful) LessThan(other Handful) bool {
	return h.Blue <= other.Blue && h.Red <= other.Red && h.Green <= other.Green
}

func (h Handful) Power() int {
	return h.Red * h.Green * h.Blue
}

func (h Handful) Max(other Handful) Handful {
	return Handful{
		Red:   max(h.Red, other.Red),
		Green: max(h.Green, other.Green),
		Blue:  max(h.Blue, other.Blue),
	}
}

type Game struct {
	Number   int
	Handfuls []Handful
}

func Day2() {
	lines := MustReadFileLines("day2.input")

	fmt.Println("Day 2")

	games := make([]Game, 0, len(lines))
	for _, line := range lines {
		parser := GameParser{Line: line}
		games = append(games, parser.ReadGame())
	}

	{ // Part 1
		guess := Handful{Red: 12, Green: 13, Blue: 14}
		totalPossible := 0
		for _, game := range games {
			gamePossible := true

			for _, hand := range game.Handfuls {
				if !hand.LessThan(guess) {
					gamePossible = false
					break
				}
			}

			if gamePossible {
				totalPossible += game.Number
			}
		}

		fmt.Printf("\tPart 1: %d\n", totalPossible)
	}

	{ // Part 2
		totalPower := 0
		for _, game := range games {
			minimum := Handful{}

			for _, hand := range game.Handfuls {
				minimum = minimum.Max(hand)
			}

			totalPower += minimum.Power()
		}

		fmt.Printf("\tPart 2: %d\n", totalPower)
	}
}

// Day 3

type Vec2 struct {
	val uint32
}

func (v *Vec2) SetX(val int16) {
	v.val = (v.val & 0xffff) + (uint32(uint16(val)) << 16)
}

func (v *Vec2) SetY(val int16) {
	v.val = (v.val & (0xffff << 16)) + uint32(uint16(val))
}

func (v Vec2) X() int16 {
	return int16(v.val >> 16)
}

func (v Vec2) Y() int16 {
	return int16(v.val & 0xffff)
}

func (v Vec2) W() int16 {
	return v.X()
}

func (v Vec2) H() int16 {
	return v.Y()
}

func (v Vec2) Get() (x, y int16) {
	return v.X(), v.Y()
}

func (v Vec2) Equals(o Vec2) bool {
	return v.val == o.val
}

func Vec2Cmp(a, b Vec2) int {
	if a.val > b.val {
		return 1
	} else if a.val < b.val {
		return -1
	}

	return 0
}

func NewVec2(x, y int16) Vec2 {
	return Vec2{val: uint32(x)<<16 + uint32(uint16(y))}
}

type Number struct {
	Value  int
	Pos    Vec2
	Length int8
}

type Symbol struct {
	Value rune
	Pos   Vec2
}

func (s Symbol) IsAdjacentTo(n Number) bool {
	for _, vOff := range []int16{-1, 1} {
		for d := int8(-1); d < n.Length+1; d++ {
			loc := NewVec2(n.Pos.Y()+vOff, n.Pos.X()+int16(d))
			if loc.Equals(s.Pos) {
				return true
			}
		}
	}

	left := NewVec2(n.Pos.Y(), n.Pos.X()-1)
	right := NewVec2(n.Pos.Y(), n.Pos.X()+int16(n.Length))
	if left.Equals(s.Pos) || right.Equals(s.Pos) {
		return true
	}

	return false
}

func (n Number) HasAdjacentSymbols(s Schematic) bool {
	for _, vOff := range []int16{-1, 1} {
		for d := int8(-1); d < n.Length+1; d++ {
			surrounding := s.At(n.Pos.Y()+vOff, n.Pos.X()+int16(d))
			if surrounding != '.' && !unicode.IsDigit(rune(surrounding)) {
				return true
			}
		}
	}

	left := rune(s.At(n.Pos.Y(), n.Pos.X()-1))
	right := rune(s.At(n.Pos.Y(), n.Pos.X()+int16(n.Length)))
	if (left != '.' && !unicode.IsDigit(left)) ||
		(right != '.' && !unicode.IsDigit(right)) {
		return true
	}

	return false
}

type Schematic struct {
	Dims    Vec2
	Lines   []string
	Numbers []Number
	Symbols []Symbol
}

func (s *Schematic) Set(row, col int16, val byte) {
	if row < 0 || col < 0 || row >= s.Dims.H() || col >= s.Dims.W() {
		return
	}
	bytes := []byte(s.Lines[row])
	bytes[col] = val

	s.Lines[row] = string(bytes)
}

func (s *Schematic) At(row, col int16) byte {
	if row < 0 || col < 0 || row >= s.Dims.H() || col >= s.Dims.W() {
		return '.'
	}
	return s.Lines[row][col]
}

func (s *Schematic) build() {
	{ // Calculate dimensions

		width := 0
		for _, line := range s.Lines {
			width = max(width, len(line))
		}

		s.Dims = NewVec2(int16(width), int16(len(s.Lines)))
	}

	{ // Find numbers and symbols
		symbols := make([]Symbol, 0, len(s.Lines)*2)
		numbers := make([]Number, 0, len(s.Lines))
		for row, line := range s.Lines {
			numDigits := int8(0)
			for c := int16(0); c < int16(len(line)+1); c++ {
				r := rune(s.At(int16(row), c))

				if unicode.IsDigit(r) {
					numDigits++
				} else if numDigits != 0 {
					val, _ := strconv.Atoi(line[c-int16(numDigits) : c])
					numbers = append(numbers,
						Number{
							Pos:    NewVec2(c-int16(numDigits), int16(row)),
							Length: numDigits,
							Value:  val,
						},
					)
					numDigits = 0
				}

				if r != '.' {
					symbols = append(symbols, Symbol{
						Pos:   NewVec2(int16(c), int16(row)),
						Value: r,
					})
				}
			}
		}

		s.Numbers = numbers
		slices.SortFunc(s.Numbers, func(a, b Number) int {
			return Vec2Cmp(a.Pos, b.Pos)
		})

		s.Symbols = symbols
		slices.SortFunc(s.Symbols, func(a, b Symbol) int {
			return Vec2Cmp(a.Pos, b.Pos)
		})
	}
}

func NewSchematic(lines []string) Schematic {
	s := Schematic{Lines: lines}
	s.build()

	return s
}

func Day3() {
	lines := MustReadFileLines("day3.test")

	fmt.Println("Day 3")

	schematic := NewSchematic(lines)

	{ // Part 1
		total := 0
		for _, num := range schematic.Numbers {
			if num.HasAdjacentSymbols(schematic) {
				total += num.Value
			}
		}

		fmt.Println("\tPart 1:", strconv.Itoa(total))
	}

	{ // Part 2
		for _, sym := range schematic.Symbols {
			nearby := make([]Number, 0, 3)

			for _, num := range schematic.Numbers {
				off := num.Pos.Y() - sym.Pos.Y()
				if off >= -1 && off <= 1 {
					nearby = append(nearby, num)
				}

				if off > 1 {
					break
				}
			}
		}
	}
}

func main() {
	Day1()
	Day2()
	Day3()
}
