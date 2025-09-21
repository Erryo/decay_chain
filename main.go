package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const Data_Location = "data/nndc_nudat_DecayMode.csv"
const (
	b_minus = iota
	b_plus
	alpha
	gamma
)

/*
read csv
(ideally)show plot
ask for input
perform decay (cascading)
show result
*/

type ChainNode struct {
	reaction Reaction
	childred []ChainNode
}
type Element struct {
	name     string
	neutrons int
	protons  int
	mass     int
	charge   int
	decays   []Decay
}

type Decay struct {
	name         string
	delta_mass   int
	delta_charge int
	by_products  []Element
}

type Reaction struct {
	parent_el   Element
	child_el    Element
	by_products []Element
	decay_name  string
}

func (c *ChainNode) addReaction(reac Reaction) {
	newChild := ChainNode{reaction: reac, childred: []ChainNode{}}
	c.childred = append(c.childred, newChild)
}

func (c *ChainNode) convertArrayToChildren(re []Reaction) {
	for _, reaction := range re {
		c.addReaction(reaction)
	}
}

// breadth first Searc
func (c ChainNode) print_bfs() {
	queue := []ChainNode{c}
	fmt.Println(len(c.childred))
	for idx := 0; idx < len(queue); idx += 1 {
		node := queue[idx]
		print_reaction("Tree BFS", node.reaction)
		queue = append(queue, node.childred...)
		fmt.Println(len(node.childred))
	}
}

func (c ChainNode) print_dfs() {
	print_reaction("Tree DFS", c.reaction)
	for _, child := range c.childred {
		child.print_dfs()
	}
}

// aos
func read_csv(decays []Decay) []Element {
	file, err := os.Open(Data_Location)
	if err != nil {
		fmt.Println("could not open file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// skip the first line
	if scanner.Scan() {
	}

	var isotope_data []Element
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")

		if len(parts) != 4 {
			fmt.Println("Parts too short")
			break
		}

		z := parts[0] // charge
		n := parts[1] // neutrons
		name := parts[2]
		decay_name := parts[3]

		charge, err := strconv.Atoi(z)
		if err != nil {
			fmt.Println("failed to convert charge", err)
			break
		}
		neutrons, err := strconv.Atoi(n)
		if err != nil {
			fmt.Println("failed to convert neutrons", err)
			break
		}
		var possible_decays []Decay

		if strings.Contains(decay_name, "B-") {
			possible_decays = append(possible_decays, decays[b_minus])
		} else if strings.Contains(decay_name, "B+") {
			possible_decays = append(possible_decays, decays[b_plus])
		} else if strings.Contains(decay_name, "A") {
			possible_decays = append(possible_decays, decays[alpha])
		} else {
			continue
		}

		isotope := Element{
			name:     name,
			charge:   charge,
			mass:     neutrons + charge,
			protons:  charge,
			neutrons: neutrons,
			decays:   possible_decays,
		}

		isotope_data = append(isotope_data, isotope)
	}

	return isotope_data
}

func combine_dupes(isotopes []Element) []Element {
	new_list := []Element{}
	max_distance := 70
	for idx := 0; idx < len(isotopes); idx += 1 {
		isotope := isotopes[idx]
		last_dupe_at := idx

		loop_end := idx + max_distance + 1
		loop_end = min(loop_end, len(isotopes))

		for jd := idx; jd < loop_end; jd += 1 {
			is := isotopes[jd]
			if is.name == isotope.name {
				isotope.decays = combine_decay(isotope.decays, is.decays)
				last_dupe_at = jd
			}
		}
		new_list = append(new_list, isotope)

		// gonna be skipped bc of loop conditionn (; idx +=1)
		idx = last_dupe_at
	}
	return new_list
}

func combine_decay(a, b []Decay) []Decay {
	result := b

	for _, dc_a := range a {
		in_both := false
		for _, dc_b := range b {
			if dc_b.name == dc_a.name {
				in_both = true
			}
		}
		if !in_both {
			result = append(result, dc_a)
		}
	}
	return result
}

// Binary search won't work cus idx isn't the chcarge
// // Could work actually
//

type SeachNeturon struct{}

func (s SeachNeturon) bigger(isotope Element, needle int) bool {
	return needle < isotope.neutrons
}

func (s SeachNeturon) smaller(isotope Element, needle int) bool {
	return needle > isotope.neutrons
}

func (s SeachNeturon) equal(isotope Element, needle int) bool {
	return needle == isotope.neutrons
}

type SearchCharge struct{}

func (s SearchCharge) bigger(isotope Element, needle int) bool {
	return needle < isotope.charge
}

func (s SearchCharge) smaller(isotope Element, needle int) bool {
	return needle > isotope.charge
}

func (s SearchCharge) equal(isotope Element, needle int) bool {
	return needle == isotope.charge
}

type Searcheable interface {
	bigger(Element, int) bool
	smaller(Element, int) bool
	equal(Element, int) bool
}

func binary_all_isotopes(isotopes []Element, needle int, se Searcheable) []Element {
	results := []Element{}
	_, idx := binary_seach_single_isotope(isotopes, needle, se)

	if idx == -1 {
		return []Element{}
	}
	// if there are more isotopes of this element behind
	// go to the first isotope of this Element
	for true {
		if idx-1 > 0 {
			if se.equal(isotopes[idx-1], needle) {
				idx += -1
			} else {
				break
			}
		}
	}
	//	if there are more isotopes in front add them
	for true {
		if idx+1 < len(isotopes) {
			if se.equal(isotopes[idx+1], needle) {
				results = append(results, isotopes[idx+1])
				idx += 1
			} else {
				break
			}
		}
	}
	return results
}

func binary_seach_single_isotope(isotopes []Element, needle int, se Searcheable) (Element, int) {
	hi := len(isotopes) - 1
	lo := 0
	for lo <= hi {
		idx := (hi + lo) / 2

		// needle is to the left
		if se.bigger(isotopes[idx], needle) {
			hi = idx - 1
		}
		// needle is to the left
		if se.smaller(isotopes[idx], needle) {
			lo = idx + 1
		}
		// needle is to the left
		if se.equal(isotopes[idx], needle) {
			return isotopes[idx], idx
		}
	}
	return Element{}, -1
}

func get_isotopes_by_charge(isotopes []Element, charge int) []Element {
	found := []Element{}
	for _, isotope := range isotopes {
		if isotope.charge == charge {
			found = append(found, isotope)
		}
	}
	return found
}

func get_isotopes_by_neutron(isotopes []Element, neutron int) []Element {
	found := []Element{}
	for _, isotope := range isotopes {
		if isotope.neutrons == neutron {
			found = append(found, isotope)
		}
	}
	return found
}

func get_isotopes_by_name(isotopes []Element, name string) (Element, error) {
	for _, isotope := range isotopes {
		if isotope.name == name {
			return isotope, nil
		}
	}
	return Element{}, fmt.Errorf("Could not find isotope by name")
}

func react(isotopes []Element, parent Element) []Reaction {
	var result_react []Reaction
	for _, decay := range parent.decays {
		var decay_reaction Reaction

		decay_reaction.parent_el = parent
		new_charge := parent.charge + decay.delta_charge
		new_mass := parent.mass + decay.delta_mass
		new_neutrons := new_mass - new_charge

		child_els := binary_all_isotopes(isotopes, new_charge, SearchCharge{})
		child_isotope := get_isotopes_by_neutron(child_els, new_neutrons)
		if len(child_isotope) != 1 {
			//			print_elements("child_els", child_els)
			fmt.Println("Len not 1")
			continue
		}

		decay_reaction.child_el = child_isotope[0]
		decay_reaction.by_products = decay.by_products
		decay_reaction.decay_name = decay.name

		result_react = append(result_react, decay_reaction)
	}
	return result_react
}

// breadth first alike
func react_tree(isotopes []Element, parent ChainNode) ChainNode {
	for idx, child := range parent.childred {
		reaction_list := react(isotopes, child.reaction.child_el)
		if len(reaction_list) == 0 {
			continue
		}
		(&child).convertArrayToChildren(reaction_list)
		parent.childred[idx] = react_tree(isotopes, child)
	}
	return parent
}

func print_elements(title string, e []Element) {
	fmt.Println("")
	fmt.Println("_____", title, "_____")
	for _, element := range e {
		fmt.Print(element.name, ":", "p", element.charge, " n", element.neutrons, " A", element.mass, ":")
		for _, decay := range element.decays {
			fmt.Print(decay.name, ",")
		}
		fmt.Println()
	}
	fmt.Println("_____----_____")
}

func print_reaction(title string, r Reaction) {
	fmt.Println("====", title, "====")
	fmt.Printf("%s --%s--> %s ", r.parent_el.name, r.decay_name, r.child_el.name)
	for _, elem := range r.by_products {
		fmt.Print("+ ", elem.name)
	}
	fmt.Printf("\nP.p %d P.n %d ---> C.p %d C.n %d \n", r.parent_el.protons, r.parent_el.neutrons, r.child_el.protons, r.child_el.neutrons)
	fmt.Println("=======")
}

func get_input(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message)

	str, _, err := reader.ReadLine()
	if err != nil {
		return ""
	}
	return string(str)
}

func show_isotopes(isotopes []Element) Element {
	sc_height := rl.GetScreenHeight()
	selected_iso := Element{}
	const offset_y = 2
	for i, isoto := range isotopes {
		button_height := float32(sc_height / len(isotopes))
		if button_height > 70 {
			button_height = 70
		}
		split_idx := 0
		if button_height < 25 {
			button_height = 25
			split_idx = (sc_height)/int(button_height+offset_y) - 1
		}
		var x float32 = 100
		var y float32 = float32(i)*button_height + float32(offset_y)
		if split_idx != 0 {
			if i > split_idx {
				x = 330
				y -= float32(split_idx)*button_height + float32(offset_y)
			}
		}
		if y > float32(sc_height) {
			fmt.Println(y)
		}
		button := gui.Button(rl.Rectangle{X: x, Y: y, Width: 200, Height: button_height}, isoto.name)
		fmt.Println(button_height, split_idx)
		if button {
			fmt.Println(isoto.name)
			selected_iso = isoto
		}
	}
	return selected_iso
}

func draw_isotope(iso Element, x int32, y int32) {
	fmt.Println("Draw", iso.name, x, y)
	stripped_name := strings.TrimLeft(iso.name, "0123456789")
	rl.DrawText(stripped_name, x, y, 16, rl.Black)
	rl.DrawText(strconv.Itoa(iso.charge), x-20, y+10, 12, rl.Black)
	rl.DrawText(strconv.Itoa(iso.mass), x-20, y-10, 12, rl.Black)

	rl.DrawRing(
		rl.NewVector2(float32(x-2), float32(y+6)),
		27, 30, // inner and outer radius (30 radius, 4 thickness)
		0, 360, // full circle
		32, // segments
		rl.Blue,
	)
}

func show_react_tree(x, y int32, root ChainNode) {
	fmt.Println(len(root.childred))
	draw_isotope(root.reaction.child_el, x, y)
	x -= 50
	for i, child := range root.childred {
		var new_x int32 = x + int32((60+80/len(root.childred))*i)
		var new_y int32 = y + 100
		show_react_tree(new_x, new_y, child)
	}
}

func main() {
	var neutrino Element = Element{name: "neutrino", neutrons: 0, protons: 0, decays: []Decay{}}
	var anti_neutrino Element = Element{name: "anti neutrino", neutrons: 0, protons: 1, decays: []Decay{}}

	var electron Element = Element{name: "electron", mass: 0, charge: -1, decays: []Decay{}}
	var positron Element = Element{name: "positron", mass: 0, charge: 1, decays: []Decay{}}

	//	var proton Element = Element{name: "proton", mass: 1, charge: 1, decays: []Decay{}}
	//	var neutron Element = Element{name: "neutron", mass: 1, charge: 0, decays: []Decay{}}

	var helium_std Element = Element{name: "4He", neutrons: 2, protons: 2, mass: 4, charge: 2, decays: []Decay{}}

	var beta_minus Decay = Decay{name: "B-", delta_mass: 0, delta_charge: +1, by_products: []Element{electron, anti_neutrino}}
	var beta_plus Decay = Decay{name: "B+", delta_mass: 0, delta_charge: -1, by_products: []Element{positron, neutrino}}
	var alpha Decay = Decay{name: "a", delta_mass: -4, delta_charge: -2, by_products: []Element{helium_std}}
	var gamma Decay = Decay{name: "g", delta_mass: 0, delta_charge: 0, by_products: []Element{}}
	var decays []Decay = []Decay{beta_minus, beta_plus, alpha, gamma}

	isotopes := read_csv(decays)
	fmt.Println("No of isotopes:", len(isotopes))
	isotopes = combine_dupes(isotopes)
	fmt.Println("No of isotopes after combine_dupes:", len(isotopes))

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 20)

	rl.SetTargetFPS(60)

	var selected_iso Element
	var selected_element []Element
	charge_str := ""
	var charge int
	var react_root ChainNode
	title := "Enter Element Charge"

	for !rl.WindowShouldClose() {
		sc_h := rl.GetScreenHeight()
		sc_w := rl.GetScreenWidth()
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		if len(selected_element) == 0 {
			bound := rl.Rectangle{X: float32(sc_w-200) / 2, Y: float32(sc_h-100) / 2, Width: 200, Height: 100}
			rl.DrawText(title, int32(bound.X)-120, int32(bound.Y)-100, 40, rl.LightGray)
			gui.TextBox(bound, &charge_str, 4, true)

			if rl.IsKeyPressed(rl.KeyEnter) {
				var err error
				charge, err = strconv.Atoi(charge_str)
				if err != nil {
					fmt.Println("converting charge_str,", err)
					title = "Please insert only numbers"
					rl.EndDrawing()
					continue
				}
				selected_element = binary_all_isotopes(isotopes, charge, SearchCharge{})
				if len(selected_element) == 0 {
					charge_str = ""
					title = "Could not find"
				}
			}

		} else if selected_iso.name == "" {
			selected_iso = show_isotopes(selected_element)
			reactions := react(isotopes, selected_iso)
			react_root.reaction = Reaction{parent_el: Element{}, decay_name: "", child_el: selected_iso}
			react_root.convertArrayToChildren(reactions)
			react_root = react_tree(isotopes, react_root)
		} else if react_root.reaction.child_el.name != "" {
			show_react_tree(int32(sc_w)/2, int32(sc_h/2), react_root)
		}

		rl.EndDrawing()

	}
}
