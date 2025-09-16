package main

import "fmt"

/*
read csv
(ideally)show plot
ask for input
perform decay (cascading)
show result
*/
type Element struct {
	name     string
	neutrons int32
	protons  int32
	mass     int32
	charge   int32
	decays   []Decay
}

type Decay struct {
	name         string
	delta_mass   int32
	delta_charge int32
	by_products  []Element
}

type Reaction struct {
	parent_el   Element
	child_el    Element
	by_products []Element
	decay_name  string
}

// aos
func read_csv(b_minus, b_plus, alpha, gamma Decay) []Element {
	var isotope_data []Element
	return isotope_data
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

	isotopes := read_csv(beta_minus, beta_plus, alpha, gamma)

	// He-5
	reaction := react(isotopes, 5, 2)
	fmt.Println(reaction)
}
