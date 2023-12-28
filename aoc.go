package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var AocOut io.Writer = os.Stdout

func Printf(format string, args ...any) (n int, err error) {
	return fmt.Fprintf(AocOut, format, args...)
}

func Println(args ...any) (n int, err error) {
	return fmt.Fprintln(AocOut, args...)
}

func MustReadFileLines(path string) []string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := make([]string, 0, 1024)

	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lines
}

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// Day 1

func Day1() {
	Println("Day 1")

	lines := MustReadFileLines("day1.input")

	isDigit := func(b byte) bool {
		return '0' <= b && b <= '9'
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

		Printf("\tPart 1: %d\n", total)

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

		Printf("\tPart 2: %d\n", total)
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

	if !IsDigit(p.Line[0]) {
		return -1
	}

	i := 0
	for ; IsDigit(p.Line[i]); i++ {
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
	Println("Day 2")

	lines := MustReadFileLines("day2.input")

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

		Printf("\tPart 1: %d\n", totalPossible)
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

		Printf("\tPart 2: %d\n", totalPower)
	}
}

// Day 3

type Vec2 struct {
	val uint32
}

func (v *Vec2) SetX(val int16) {
	v.val = (v.val & (0xffff << 16)) + uint32(uint16(val))
}

func (v *Vec2) SetY(val int16) {
	v.val = (v.val & 0xffff) + (uint32(uint16(val)) << 16)
}

func (v Vec2) X() int16 {
	return int16(v.val & 0xffff)
}

func (v Vec2) Y() int16 {
	return int16(v.val >> 16)
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
	return Vec2{val: uint32(y)<<16 + uint32(uint16(x))}
}

type Number struct {
	Value  int
	Pos    Vec2
	Length int16
}

type Symbol struct {
	Value byte
	Pos   Vec2
}

func (s Symbol) IsAdjacentTo(n Number) bool {
	vOff := n.Pos.Y() - s.Pos.Y()
	if vOff == 1 || vOff == -1 {
		return s.Pos.X() >= n.Pos.X()-1 &&
			s.Pos.X() <= n.Pos.X()+n.Length
	} else if vOff == 0 {
		return s.Pos.X() == n.Pos.X()-1 ||
			s.Pos.X() == n.Pos.X()+n.Length
	}

	return false
}

func (n Number) HasAdjacentSymbols(s Schematic) bool {
	for _, vOff := range []int16{-1, 1} {
		for d := int16(-1); d < n.Length+1; d++ {
			surrounding := s.At(n.Pos.Y()+vOff, n.Pos.X()+d)
			if surrounding != '.' && !unicode.IsDigit(rune(surrounding)) {
				return true
			}
		}
	}

	left := s.At(n.Pos.Y(), n.Pos.X()-1)
	right := s.At(n.Pos.Y(), n.Pos.X()+int16(n.Length))
	if (left != '.' && !IsDigit(left)) ||
		(right != '.' && !IsDigit(right)) {
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
			numDigits := int16(0)
			for c := int16(0); c < int16(len(line)+1); c++ {
				r := s.At(int16(row), c)

				if IsDigit(r) {
					numDigits++
				} else {
					if r == '*' {
						symbols = append(symbols, Symbol{
							Pos:   NewVec2(int16(c), int16(row)),
							Value: r,
						})
					}

					if numDigits != 0 {
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
	Println("Day 3")

	lines := MustReadFileLines("day3.input")

	schematic := NewSchematic(lines)

	{ // Part 1
		total := 0

		for _, num := range schematic.Numbers {
			if num.HasAdjacentSymbols(schematic) {
				total += num.Value
			}
		}

		Println("\tPart 1:", strconv.Itoa(total))
	}

	{ // Part 2
		total := 0

		for _, sym := range schematic.Symbols {
			nearby := make([]Number, 0, 3)
			adjacent := make([]Number, 0, 3)

			for _, num := range schematic.Numbers {
				off := num.Pos.Y() - sym.Pos.Y()
				if off >= -1 && off <= 1 {
					nearby = append(nearby, num)
				}

				if off > 1 {
					break
				}
			}

			for _, near := range nearby {
				if sym.IsAdjacentTo(near) {
					adjacent = append(adjacent, near)
				}
			}

			if len(adjacent) == 2 {
				total += adjacent[0].Value * adjacent[1].Value
			}
		}

		Println("\tPart 2:", strconv.Itoa(total))
	}
}

// Day 4

type CardParser struct {
	Line string
}

func (p *CardParser) IsEol() bool {
	return len(p.Line) == 0
}

func (p *CardParser) SkipWs() {
	if p.IsEol() {
		return
	}

	for unicode.IsSpace(rune(p.Line[0])) {
		p.Line = p.Line[1:]
	}
}

func (p *CardParser) ReadStr(str string) {
	p.SkipWs()

	if strings.HasPrefix(p.Line, str) {
		p.Line = p.Line[len(str):]
	}
}

func (p *CardParser) ReadChar() byte {
	p.SkipWs()

	if p.IsEol() {
		return ' '
	}

	result := p.Line[0]
	p.Line = p.Line[1:]

	return result
}

func (p *CardParser) SkipNum() {
	p.SkipWs()

	for IsDigit(p.Line[0]) {
		p.Line = p.Line[1:]
	}
}

func (p *CardParser) ReadNum() int {
	p.SkipWs()

	if p.IsEol() || !IsDigit(p.Line[0]) {
		return -1
	}

	i := 0
	for ; i < len(p.Line) && IsDigit(p.Line[i]); i++ {
	}

	numStr := p.Line[:i]
	p.Line = p.Line[i:]
	num, _ := strconv.ParseUint(numStr, 10, 64)
	return int(num)
}

func (p *CardParser) ReadCard() Card {
	card := Card{
		Winning: make([]int, 0, 16),
		Numbers: make([]int, 0, 16),
	}

	p.ReadStr("Card")
	p.SkipNum()
	p.ReadChar()

	num := p.ReadNum()
	for ; num != -1; num = p.ReadNum() {
		card.Winning = append(card.Winning, num)
	}

	p.ReadChar()

	num = p.ReadNum()
	for ; num != -1; num = p.ReadNum() {
		card.Numbers = append(card.Numbers, num)
	}

	slices.Sort(card.Winning)
	slices.Sort(card.Numbers)

	return card
}

type Card struct {
	Winning, Numbers []int
}

func Day4() {
	Println("Day 4")

	lines := MustReadFileLines("day4.input")

	cards := make([]Card, 0, len(lines))

	for _, line := range lines {
		parser := CardParser{Line: line}
		cards = append(cards, parser.ReadCard())
	}

	numMatches := make([]int, len(cards))

	{ // Part 1
		for i, card := range cards {
			wins, nums := card.Winning, card.Numbers
			iWin, iNum := 0, 0
			matches := 0

			for iWin < len(wins) && iNum < len(nums) {
				diff := wins[iWin] - nums[iNum]
				if diff == 0 {
					matches++
					iNum++
					iWin++
				} else if diff < 0 {
					iWin++
				} else if diff > 0 {
					iNum++
				}
			}

			numMatches[i] = matches
		}

		total := 0
		for _, matches := range numMatches {
			total += (1 << matches) >> 1
		}

		Printf("\tPart 1: %d\n", total)
	}

	{ // Part 2
		copies := make([]int, len(cards))
		for i := range copies {
			copies[i] = 1
		}

		for i := range cards {
			for j := 0; j < numMatches[i]; j++ {
				copies[i+j+1] += copies[i]
			}
		}

		total := 0
		for _, copies := range copies {
			total += copies
		}

		Printf("\tPart 2: %d\n", total)
	}
}

// Day 5

type IntMap struct {
	DstStarts, SrcStarts, Lengths []int64
}

func (m *IntMap) Insert(dstStart, srcStart, length int64) {
	m.DstStarts = append(m.DstStarts, dstStart)
	m.SrcStarts = append(m.SrcStarts, srcStart)
	m.Lengths = append(m.Lengths, length)
}

func (m *IntMap) Len() int {
	return len(m.SrcStarts)
}

func (m *IntMap) Swap(i, j int) {
	m.SrcStarts[i], m.SrcStarts[j] = m.SrcStarts[j], m.SrcStarts[i]
	m.DstStarts[i], m.DstStarts[j] = m.DstStarts[j], m.DstStarts[i]
	m.Lengths[i], m.Lengths[j] = m.Lengths[j], m.Lengths[i]
}

func (m *IntMap) Less(i, j int) bool {
	return m.SrcStarts[i] < m.SrcStarts[j]
}

func (m *IntMap) Map(v int64) int64 {
	i, found := slices.BinarySearch(m.SrcStarts, v)
	if !found {
		// Get the next-lowest element (or -1)
		i -= 1
	}

	if i < 0 {
		return v
	}

	off := v - m.SrcStarts[i]
	if off >= 0 && off <= m.Lengths[i] {
		return m.DstStarts[i] + off
	}

	return v
}

func IntMapCombine(a, b IntMap) IntMap {
	//
	return IntMap{}
}

type Almanac struct {
	SeedToSoil     IntMap
	SoilToFert     IntMap
	FertToWater    IntMap
	WaterToLight   IntMap
	LightToTemp    IntMap
	TempToHumidity IntMap
	HumidityToLoc  IntMap
}

func (a Almanac) SeedToLocation(seed int64) int64 {
	soil := a.SeedToSoil.Map(seed)
	fert := a.SoilToFert.Map(soil)
	water := a.FertToWater.Map(fert)
	light := a.WaterToLight.Map(water)
	temp := a.LightToTemp.Map(light)
	humidity := a.TempToHumidity.Map(temp)
	loc := a.HumidityToLoc.Map(humidity)
	return loc
}

func ReadAlmanac(lines []string, a *Almanac, s *[]int64) {
	readNumList := func(line string, numValues int) []int64 {
		result := make([]int64, 0, numValues)
		numStrs := strings.Fields(line)

		for _, str := range numStrs {
			val, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				panic(err)
			}

			result = append(result, int64(val))
		}

		return result
	}

	*s = readNumList(lines[0][len("seeds: "):], 100)

	lines = lines[2:]

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		dst := (*IntMap)(nil)

		switch line {
		case "seed-to-soil map:":
			dst = &a.SeedToSoil
		case "soil-to-fertilizer map:":
			dst = &a.SoilToFert
		case "fertilizer-to-water map:":
			dst = &a.FertToWater
		case "water-to-light map:":
			dst = &a.WaterToLight
		case "light-to-temperature map:":
			dst = &a.LightToTemp
		case "temperature-to-humidity map:":
			dst = &a.TempToHumidity
		case "humidity-to-location map:":
			dst = &a.HumidityToLoc
		default:
			break
		}

		i++
		for ; i < len(lines) && lines[i] != ""; i++ {
			nums := readNumList(lines[i], 3)
			dst.Insert(nums[0], nums[1], nums[2])
		}
		sort.Sort(dst)
	}
}

func Day5() {
	lines := MustReadFileLines("day5.input")

	Println("Day 5")

	almanac := Almanac{}
	seeds := make([]int64, 0, 100)
	ReadAlmanac(lines, &almanac, &seeds)

	{ // Part 1
		min := int64(math.MaxInt64)
		for _, seed := range seeds {
			loc := almanac.SeedToLocation(seed)
			if loc < min {
				min = loc
			}
		}

		Printf("\tPart 1: %d\n", min)
	}

	{ // Part 2
		min := int64(math.MaxInt64)
		for i := 0; i < len(seeds); i += 2 {
			start, length := seeds[i], seeds[i+1]
			for s := start; s < start+length; s++ {
				loc := almanac.SeedToLocation(s)
				if loc < min {
					min = loc
				}
			}
		}

		Printf("\tPart 2: %d\n", min)
	}
}

func main() {
	Day1()
	Day2()
	Day3()
	Day4()
	Day5()
}
