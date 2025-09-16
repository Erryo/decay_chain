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
		fmt.Println(parts)

		if len(parts) != 4 {
			fmt.Println("Parts too short")
			break
		}

		z := parts[0] // charge
		n := parts[1] // mass
		name := parts[2]
		decay_name := parts[3]

		charge, err := strconv.Atoi(z)
		if err != nil {
			fmt.Println("failed to convert charge", err)
			break
		}
		mass, err := strconv.Atoi(n)
		if err != nil {
			fmt.Println("failed to convert mass", err)
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
			mass:     mass,
			protons:  charge,
			neutrons: mass - charge,
			decays:   possible_decays,
		}

		isotope_data = append(isotope_data, isotope)
	}

	return isotope_data
}

func combine_dupes(isotopes []Element) []Element {
	new_list := []Element{}
	max_distance := 5
	for idx, isotope := range isotopes {
		for _, is := range isotopes[idx : idx+max_distance+1] {
			if is.name == isotope.name {
			}
		}
		new_list = append(new_list, isotope)
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

func react(isotopes []Element, neutrons int32, protons int32) Reaction {
	var decay_reaction Reaction
	return decay_reaction
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

	// He-5
	reaction := react(isotopes, 5, 2)
	fmt.Println(reaction)
}
