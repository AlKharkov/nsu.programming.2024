package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type matrix struct {
	n     int
	value [][]byte
}

func (m *matrix) Read() {
	file, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	m.n = len(strings.Split(scanner.Text(), " "))
	m.value = make([][]byte, m.n)
	m.value[0] = make([]byte, m.n)
	for i, v := range strings.Split(scanner.Text(), " ") {
		var val byte = 0
		if v != "0" {
			val = 1
		}
		m.value[0][i] = val
	}
	for i := 1; i < m.n; i++ {
		m.value[i] = make([]byte, m.n)
		scanner.Scan()
		for j, v := range strings.Split(scanner.Text(), " ") {
			var val byte = 0
			if v != "0" {
				val = 1
			}
			m.value[i][j] = val
		}
	}
}

// func (m *matrix) Print() {
// 	fmt.Println(m.n)
// 	for _, v := range m.value {
// 		for j := range v {
// 			fmt.Print(v[j], " ")
// 		}
// 		fmt.Println()
// 	}
// }

func main() {
	var m matrix
	m.Read()
	fmt.Println(BronKerboschAlgorithm(&m))
}

func BronKerboschAlgorithm(m *matrix) [][]int {
	results := make([][]int, 0)
	compsub := make([]int, 0)
	candidates := make([]int, m.n)
	not := make([]int, 0)
	for i := range m.n {
		candidates[i] = i
	}
	extend(compsub, candidates, not, m, &results)
	return results
}

func extend(compsub, candidates, not []int, m *matrix, results *[][]int) {
	for len(candidates) != 0 && !containsVerticeConnectedToAllfrom(not, candidates, m) {
		v := candidates[0]
		compsub = append(compsub, v)
		newCandidates := connectedWithV(v, candidates, m)
		newNot := connectedWithV(v, not, m)
		if len(newCandidates) == 0 && len(newNot) == 0 {
			*results = append(*results, compsub)
		} else {
			extend(compsub, newCandidates, newNot, m, results)
		}
		candidates = candidates[1:]
		compsub = compsub[:len(compsub)-1]
		not = append(not, v)
	}
}

func containsVerticeConnectedToAllfrom(a, b []int, m *matrix) bool {
	for _, v := range a {
		connectedToAll := true
		for _, w := range b {
			if m.value[v][w] == 0 {
				connectedToAll = false
				break
			}
		}
		if connectedToAll {
			return true
		}
	}
	return false
}

func connectedWithV(v int, a []int, m *matrix) []int {
	b := make([]int, 0)
	for _, w := range a {
		if m.value[v][w] == 1 {
			b = append(b, w)
		}
	}
	return b
}
