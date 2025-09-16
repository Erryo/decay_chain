package main

import (
	"testing"
)

var (
	neutrino      Element = Element{name: "neutrino", neutrons: 0, protons: 0, decays: []Decay{}}
	anti_neutrino Element = Element{name: "anti neutrino", neutrons: 0, protons: 1, decays: []Decay{}}
	electron      Element = Element{name: "electron", mass: 0, charge: -1, decays: []Decay{}}
	positron      Element = Element{name: "positron", mass: 0, charge: 1, decays: []Decay{}}
	helium_std    Element = Element{name: "4He", neutrons: 2, protons: 2, mass: 4, charge: 2, decays: []Decay{}}
)

var (
	beta_minus Decay = Decay{name: "B-", delta_mass: 0, delta_charge: +1, by_products: []Element{electron, anti_neutrino}}
	beta_plus  Decay = Decay{name: "B+", delta_mass: 0, delta_charge: -1, by_products: []Element{positron, neutrino}}
	alpha_d    Decay = Decay{name: "a", delta_mass: -4, delta_charge: -2, by_products: []Element{helium_std}}
	gamma_d    Decay = Decay{name: "g", delta_mass: 0, delta_charge: 0, by_products: []Element{}}
)

func TestCombine(t *testing.T) {
	a := []Decay{beta_minus, beta_plus}
	b := []Decay{alpha_d, beta_plus}

	res := combine_decay(a, b)

	if !Equal(res, []Decay{alpha_d, beta_plus, beta_minus}) {
		t.Errorf("wrong  %v", res)
	}
}

func Equal(a, b []Decay) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.name != b[i].name {
			return false
		}
	}
	return true
}
