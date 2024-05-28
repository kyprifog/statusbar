package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"sort"
)

var defStyle tcell.Style
var cheatDays int

func bar_path() string {
	args := os.Args[1]
	return filepath.Join(args)
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func days_since(s string) int {
	const layout = "2006-01-02"
	t, _ := time.Parse(layout, s)
	return int((time.Since(t)).Hours() / (24))
}

func add_day(s string) string {
	const layout = "2006-01-02"
	t, _ := time.Parse(layout, s)
	newT := t.AddDate(0, 0, 1)
	return newT.Format(layout)
}

func sub_day(s string) string {
	const layout = "2006-01-02"
	t, _ := time.Parse(layout, s)
	newT := t.AddDate(0, 0, -1)
	return newT.Format(layout)
}

func insert(a []map[string]interface{}, index int, value map[string]interface{}) []map[string]interface{} {
		if len(a) == index { 
				return append(a, value)
		}
		a = append(a[:index+1], a[index:]...) 
		a[index] = value
		return a
}

func sort_bars(bars []map[string]interface{}) []map[string]interface{} {
	// Define a custom comparison function
	sort.Slice(bars, func(i, j int) bool {
		// Get the necessary values for the first element
		startDateI := bars[i]["start_date"].(string)
		lengthI := bars[i]["length"].(int)
		nameI := bars[i]["name"].(string)

		// Get the necessary values for the second element
		startDateJ := bars[j]["start_date"].(string)
		lengthJ := bars[j]["length"].(int)
		nameJ := bars[j]["name"].(string)

		// Calculate the sorting key
		diffI := days_since(startDateI) - lengthI
		diffJ := days_since(startDateJ) - lengthJ

		// Compare the sorting keys
		if diffI == diffJ {
			return nameI < nameJ
		}
		return diffI > diffJ
	})

	return bars
}

func render_bars(s tcell.Screen, max_bar_length int, bars []map[string]interface{}) {

	green := tcell.StyleDefault.Foreground(tcell.ColorLawnGreen)
	yellow := tcell.StyleDefault.Foreground(tcell.Color184)
	orange := tcell.StyleDefault.Foreground(tcell.ColorDarkOrange)
	red := tcell.StyleDefault.Foreground(tcell.ColorRed)
	blue := tcell.StyleDefault.Foreground(tcell.ColorBlue)

	theme := []string{"█", " ", "|", "▐", "▀", "▄" }

	index := 0
	maxBarLength := max_bar_length

	cheatString := fmt.Sprintf("Cheat + | -%s", strings.Repeat(" " + theme[4], cheatDays))
	emitStr(s, 1, index, blue, cheatString)

	for _, el := range bars {
		length := el["length"].(int)
		inc := el["inc"].(int)

		days_since := days_since(el["start_date"].(string))

		day_bar := days_since % maxBarLength
		overflow := length / maxBarLength
		bar_length := length % maxBarLength

		barString := fmt.Sprintf("\r %s%s%s",
			theme[3],
			strings.Repeat(theme[0], bar_length),
			" + | -",
		)
		errorBarString := fmt.Sprintf("\r %s%s",
			theme[3],
			strings.Repeat(theme[0], day_bar),
		)
		name_string := fmt.Sprintf("\r %s (%v) ", el["name"].(string), length)

		var bar_color = green
		if days_since <= length {
			bar_color = green
		} else {
			if (days_since - length) < 5 {
				bar_color = yellow
			} else if (days_since - length) < 10 {
				bar_color = orange
			} else {
				bar_color = red
			}
		}
		inc_string := fmt.Sprintf(" (+-%d)", inc)
		day_string := fmt.Sprintf("   %v day(s)", days_since)
		medal_string := strings.Repeat(" " + theme[4], overflow)
		emitStr(s, 2, index+2, bar_color, name_string)
		emitStr(s, len(name_string)+1, index+2, bar_color, inc_string)
		emitStr(s, len(inc_string)+len(name_string)+1, index+2, blue, day_string)
		emitStr(s, len(inc_string)+len(name_string)+len(day_string)+1, index+2, yellow, medal_string)
		emitStr(s, 2, index+3, bar_color, fmt.Sprintf(barString))
		emitStr(s, 2, index+4, blue, fmt.Sprintf(errorBarString))
		index += 4
	}
}

func save_bars(bars []map[string]interface{}) {
	data := make(map[string]interface{})
	data["bars"] = bars
	data["cheat_days"] = cheatDays
	d, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	f, err := os.Create(bar_path())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	f.Write(d)
}

func inc_dec_bars(max_bar_length int, x int, y int, bars []map[string]interface{}) []map[string]interface{} {
	new_bars := []map[string]interface{}{}
	if y == 0 {
		if x > 9 {
			for _, el := range bars {
				date := el["start_date"].(string)
				el["start_date"] = sub_day(date)
				new_bars = append(new_bars, el)
			}
			if cheatDays > 0 { // Ensure cheatDays doesn't go below 0
				cheatDays--
			}

		} else {
			for _, el := range bars {
				date := el["start_date"].(string)
				el["start_date"] = add_day(date)
				new_bars = append(new_bars, el)
			}
			cheatDays++
		}
	} else {
		for i, el := range bars {
			if 4*(i+1)-1 == y {
				length := el["length"].(int)
				bar_length := length % max_bar_length
				inc := el["inc"].(int)
				if x <= (bar_length + 6) {
					new_length := el["length"].(int) + inc
					el["length"] = new_length
					new_bars = append(new_bars, el)
				} else if x > (bar_length + 6) {
					new_length := el["length"].(int) - inc
					if new_length < 0 {
						el["length"] = 0
					} else {
						el["length"] = new_length
					}
					new_bars = append(new_bars, el)
				} else {
					new_bars = append(new_bars, el)
				}
			} else {
				new_bars = append(new_bars, el)
			}
		}
	}
	save_bars(new_bars)
	return new_bars
}

func get_bars() ([]map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(bar_path())

	type BarConfig struct {
		Bars      []map[string]interface{}
		CheatDays int `yaml:"cheat_days"`
	}

	bars := BarConfig{}
	if err == nil {
		err = yaml.Unmarshal(yamlFile, &bars)
	}

	cheatDays = bars.CheatDays

	return bars.Bars, err
}

func main() {
	s, e := tcell.NewScreen()

	encoding.Register()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)

	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	max_bar_length := 30

	bars, err := get_bars()

	ecnt := 0
	last_press := time.Now().AddDate(-1, 0, 0)

	if err != nil {
		s.Fini()
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(0)
	}

	render_bars(s, max_bar_length, bars)
	s.Show()

	go func() {
		for {
			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					ecnt++
					if ecnt > 1 {
						s.Fini()
						os.Exit(0)
					}
				}
			case *tcell.EventMouse:
				x, y := ev.Position()
				switch ev.Buttons() {
				case tcell.Button1, tcell.Button2, tcell.Button3:
					s.Clear()
					new_bars, _ := get_bars()
					if time.Now().Sub(last_press).Seconds() > 0.5 {
						new_bars = inc_dec_bars(max_bar_length, x, y, new_bars)
						last_press = time.Now()
					} else {
						new_bars, _ = get_bars()
					}
					render_bars(s, max_bar_length, new_bars)
					s.Sync()
					s.Show()
				}
			}
		}
	}()
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			new_bars, _ := get_bars()
			sorted_bars := sort_bars(new_bars)
			save_bars(sorted_bars)
			s.Clear()
			render_bars(s, max_bar_length, new_bars)
			s.Sync()
			s.Show()
		}
	}
}
