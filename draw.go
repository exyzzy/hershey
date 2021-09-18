package hershey

import (
	"errors"
	"fmt"
	"math"
	"regexp"
)

type Fn func(s ...interface{})

func DrawChar(c rune, font string, scale int, x, y *int, fmv, fln Fn, s ...interface{}) (err error) {
	index := int(c) - 32
	f, ok := Fonts[font]
	if !ok {
		err = errors.New("Unknown Font: " + font)
		return
	}
	if index < 0 || index >= len(f) {
		err = errors.New(fmt.Sprintf("Font: %s, Length: %d, Unprintable character: %d", font, len(f), index))
		return
	}
	if scale < 0 {
		err = errors.New(fmt.Sprintf("Negative scale: %d", scale))
		return
	}
	for _, path := range hershey[f[index]].Coords {
		dx := *x + (path[0][0]-hershey[f[index]].Lt)*scale
		dy := *y - (path[0][1])*scale //chars assume Y inversion
		//prepend dx, dy
		sm := append([]interface{}{&dy}, s...)
		sm = append([]interface{}{&dx}, sm...)
		// fmv(dx, dy, s...)
		if fmv != nil {
			fmv(sm...)
		}
		for coord := 1; coord < len(path); coord++ {
			dx = *x + (path[coord][0]-hershey[f[index]].Lt)*scale
			dy = *y - (path[coord][1])*scale
			//prepend dx, dy
			sl := append([]interface{}{&dy}, s...)
			sl = append([]interface{}{&dx}, sl...)
			// fln(dx, dy, s...)
			if fln != nil {
				fln(sl...)
			}
		}
	}
	//update x for char width
	*x += (hershey[f[index]].Rt - hershey[f[index]].Lt) * scale
	return
}

// s = (haveFirst, minX, minY, maxX, maxY)
func Bounds(s ...interface{}) {
	x := *(s[0].(*int))
	y := *(s[1].(*int))
	b := *(s[2].(*bool))
	if !b {
		*(s[2].(*bool)) = true
		*(s[3].(*int)) = x
		*(s[4].(*int)) = y
		*(s[5].(*int)) = x
		*(s[6].(*int)) = y
	} else {
		*(s[3].(*int)) = int(math.Min(float64(x), float64(*(s[3].(*int)))))
		*(s[4].(*int)) = int(math.Min(float64(y), float64(*(s[4].(*int)))))
		*(s[5].(*int)) = int(math.Max(float64(x), float64(*(s[5].(*int)))))
		*(s[6].(*int)) = int(math.Max(float64(y), float64(*(s[6].(*int)))))
	}
}

//get the bounds of a string to be drawn (do not draw)
func StringBounds(font string, scale int, x, y int, str string) (minX, minY, maxX, maxY int, err error) {
	haveFirst := false //to seed max/min
	tx := x
	ty := y
	for _, c := range str {
		err = DrawChar(c, font, scale, &tx, &ty, Bounds, Bounds, &haveFirst, &minX, &minY, &maxX, &maxY)
		if err != nil {
			return
		}
	}
	return
}

func DrawString(str string, font string, scale int, x, y *int, fmv, fln Fn, s ...interface{}) (err error) {
	for _, c := range str {
		err = DrawChar(c, font, scale, x, y, fmv, fln, s...)
		if err != nil {
			return
		}
	}
	return
}

func DrawStringLines(str string, font string, scale int, x, y *int, lineSpace, width int, fmv, fln Fn, s ...interface{}) (err error) {
	re := regexp.MustCompile(`[\n| +]`)
	id := re.FindAllIndex([]byte(str), -1)
	if id == nil {
		return
	}
	height := Height[font]
	lastIndex := 0
	x0 := *x
	for strId := 0; strId < len(id); strId++ {
		*x = x0
		_, _, xx, _, err := StringBounds(font, scale, *x, *y, str[lastIndex:id[strId][0]])
		if err != nil {
			return err
		}
		if (xx > (x0 + width)) || (str[id[strId][0]] == '\n') { //exceeded width or , print to last good break
			if xx > (x0 + width) {
				strId--
				if strId < 0 {
					return errors.New("No breaks and too long") //does not fit??
				}
			}
			*x = x0
			*y -= (height[1] * scale)
			err = DrawString(str[lastIndex:id[strId][0]], font, scale, x, y, fmv, fln, s...)
			*y += (height[0] * scale) - lineSpace
			lastIndex = id[strId][0] + 1
		}
	}
	*x = x0
	if lastIndex < len(str) {
		*y -= (height[1] * scale)
		err = DrawString(str[lastIndex:], font, scale, x, y, fmv, fln, s...)
	}
	return
}
