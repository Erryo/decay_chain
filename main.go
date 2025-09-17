package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		}

		if strings.Contains(decay_name, "B+") {
			possible_decays = append(possible_decays, decays[b_plus])
		}

		if strings.Contains(decay_name, "A") {
			possible_decays = append(possible_decays, decays[alpha])
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
	for {
		idx := (hi + lo) / 2
		// needle is to the left
		if se.bigger(isotopes[idx], needle) {
			hi = idx
		}
		// needle is to the left
		if se.smaller(isotopes[idx], needle) {
			lo = idx
		}
		// needle is to the left
		if se.equal(isotopes[idx], needle) {
			return isotopes[idx], idx
		}
	}
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
		new_mass := parent.mass + decay.delta_charge
		new_neutrons := new_mass - new_charge

		child_els := binary_all_isotopes(isotopes, new_charge, SearchCharge{})
		print_elements("child_els", child_els)
		child_isotope := get_isotopes_by_neutron(child_els, new_neutrons)
		print_elements("child_isotope", child_isotope)
		if len(child_isotope) != 1 {
			fmt.Println("Len not 1")
		}

		result_react = append(result_react, decay_reaction)
	}
	return result_react
}

func print_elements(title string, e []Element) {
	fmt.Println("")
	fmt.Println("_____", title, "_____")
	for _, element := range e {
		fmt.Print(element.name, ":", element.charge, element.neutrons, ":")
		for _, decay := range element.decays {
			fmt.Print(decay.name, ",")
		}
		fmt.Println()
	}
	fmt.Println("_____----_____")
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
	reduced_isotopes := combine_dupes(isotopes)
	fmt.Println("No of isotopes after combine_dupes:", len(reduced_isotopes))

	th_isotopes := binary_all_isotopes(isotopes, 90, SearchCharge{})
	print_elements("THorium isotopes", th_isotopes)
	th_230, err := get_isotopes_by_name(th_isotopes, "230Th")
	print_elements("Thorium 230", []Element{th_230})

	if err != nil {
		fmt.Println(err)
		return
	}
	reaction := react(isotopes, th_230)
	fmt.Println("Reachion", reaction)
}
