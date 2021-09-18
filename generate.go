// decoder/generator for hershey fonts
// reads hershey.dat and creates hershey.go amd hersheyheights.go
// uses fontids.go for font generation

package hershey

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func GenerateHershey() (out string, idToIndex map[int]int, err error) {
	idToIndex = make(map[int]int)
	name := "hershey"
	data, err := ioutil.ReadFile(name + ".dat")
	if err != nil {
		return
	}
	fmt.Println("Decoding: ", name)
	s := strings.ReplaceAll(string(data), "\n", "")
	out = "//Coordinates: \n\nvar " + name + " = []FontPath{\n"
	i := 0
	for {
		if len(s) < 5 {
			break
		}
		// hershey character id
		id := s[0:5]
		nid, err := strconv.Atoi(strings.Trim(id, " "))
		if err != nil {
			return "", map[int]int{}, err
		}
		idToIndex[nid] = i
		//number of coordinate pairs
		n, err := strconv.Atoi(strings.Trim(s[5:8], " "))
		if err != nil {
			return "", map[int]int{}, err
		}
		s = s[8:]
		//number of characters
		cn := 2 * n
		ss := s[:cn]
		//left position
		lt := int(ss[0]) - int('R')
		//right position
		rt := int(ss[1]) - int('R')
		out += fmt.Sprintf("// %d: %s\n", i, id)
		out += fmt.Sprintf("{%d, %d, [][][]int{", lt, rt)
		i++
		//paths & coordinates
		if n > 1 {
			out += "{"
		}
		for i := 2; i < len(ss); i++ {
			if ss[i] == ' ' { //new path
				out += fmt.Sprintf("},{")
				i += 2
			}
			//coords
			dx := int(ss[i]) - int('R')
			dy := int(ss[i+1]) - int('R')
			i++
			out += fmt.Sprintf("{%d, %d},", dx, dy)
		}
		out = strings.TrimRight(out, ",")
		if n > 1 {
			out += "}"
		}
		out += fmt.Sprintf("}},\n")
		s = s[cn:]
	}
	out = strings.TrimRight(out, ",\n")
	out += "}\n\n"
	return
}

func GenerateFonts(idToIndex map[int]int) (out string) {
	out = "//Fonts: \n\nvar Fonts = map[string][]int{"
	for name, _ := range Id_Fonts {
		out += "\"" + name + "\": " + name + ",\n"
	}
	out += "}\n\n"
	for name, font := range Id_Fonts {
		out += "var " + name + " = []int{\n"
		for i, v := range font {
			out += fmt.Sprintf("%d, ", idToIndex[v])
			if (i+1)%10 == 0 {
				out += "\n"
			}
		}
		out += "}\n\n"
	}
	return
}

func GenerateTranslators(idToIndex map[int]int) (out string) {
	out1 := "//Id translation\n\nvar IdToIndex = map[int]int {\n"
	out2 := "\n\nvar IndexToId = map[int]int {\n"
	i := 1
	for k, v := range IdToIndex {
		out1 += fmt.Sprintf("%d: %d,", k, v)
		out2 += fmt.Sprintf("%d: %d,", v, k)
		if i%10 == 0 {
			out1 += "\n"
			out2 += "\n"
		}
		i++
	}
	out1 += "}\n\n"
	out2 += "}\n\n"
	out = out1 + out2
	return
}

func Generate() (err error) {
	packageName := "hershey"
	out0 := "// generated file, do not edit - see generate.go\n"
	out0 += "package " + packageName + "\n\n" + "type FontPath struct {\n   Lt     int\n   Rt     int\n   Coords [][][]int\n}\n\n"

	out1, idToIndex, err := GenerateHershey()
	if err != nil {
		return
	}
	out2 := GenerateFonts(idToIndex)
	out3 := GenerateTranslators(idToIndex)

	err = ioutil.WriteFile(packageName+".go", []byte(out0+out1+out2+out3), 0644)
	// SymbolTables(IdToIndex)
	return
}

//helper function to print symbol tables for any unused chars, (already added to fontids.go)
func SymbolTables(IdToIndex map[int]int) {
	var symbols []int
	for id, _ := range IdToIndex {
		found := false
		for _, font := range Id_Fonts {
			for _, fid := range font {
				if fid == id {
					found = true
					break
				}
			}
		}
		if !found {
			symbols = append(symbols, id)
			// fmt.Println("Unique Id: ", id)
		}
	}
	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i] < symbols[j]
	})
	fmt.Println("Unused Symbols ", len(symbols))
	i := 0
	for _, v := range symbols {
		fmt.Printf("%d,", v)
		if (i+1)%10 == 0 {
			println()
		}
		if (i+1)%96 == 0 {
			println("\n===================")
			i = -1
		}
		i++
	}
}

//generate min/max heights for every font into hersheyheights.go, used by DrawStringLines()
func GenerateHeights() (err error) {
	name := "hersheyheights"
	packageTxt := "hershey"
	fmt.Println("Creating: ", name)
	out1 := "// generated file, do not edit - see generate.go\n"
	out1 += "package " + packageTxt + "\n\n"

	out1 += "//Height: \n\nvar Height = map[string][]int{"

	for fontname, f := range Fonts {
		fmt.Println("  Processing Font: ", fontname)
		var minX, minY, maxX, maxY int
		haveFirst := false
		for i, _ := range f {
			var x, y int
			err = DrawChar(rune(i+32), fontname, 1, &x, &y, Bounds, Bounds, &haveFirst, &minX, &minY, &maxX, &maxY)
			if err != nil {
				return
			}
		}
		out1 += fmt.Sprintf("\"%s\": {%d,%d},\n", fontname, minY, maxY)
	}

	out1 += "}\n\n"

	err = ioutil.WriteFile(name+".go", []byte(out1), 0644)
	return
}
