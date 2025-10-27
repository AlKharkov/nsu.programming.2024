package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type graph struct {
	n int      // Количество вершин
	A [][]byte // Матрица смежности
	I [][]int  // Множество из список инцидентных вершин
}

// Считываем граф по матрице смежности из входного файла "input.txt"
func (g *graph) Read() {
	file, _ := os.Open("input.txt") // "input.txt" -> os.Stdin
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	g.n = len(strings.Split(scanner.Text(), " "))
	// Calculating m.value
	g.A = make([][]byte, g.n)
	g.A[0] = make([]byte, g.n)
	for i, v := range strings.Split(scanner.Text(), " ") {
		var val byte = 0
		if v != "0" {
			val = 1
		}
		g.A[0][i] = val
	}
	for i := 1; i < g.n; i++ {
		g.A[i] = make([]byte, g.n)
		scanner.Scan()
		for j, v := range strings.Split(scanner.Text(), " ") {
			var val byte = 0
			if v != "0" {
				val = 1
			}
			g.A[i][j] = val
		}
	}
	g.Build() // в таком представлении сокращается количество кода
}

// Строит для графа сопоставление каждой вершине множество инцидентных ей рёбер
func (g *graph) Build() {
	g.I = make([][]int, g.n)
	for i, u := range g.A {
		for j, v := range u {
			if v == 1 {
				g.I[i] = append(g.I[i], j)
			}
		}
	}
}

// Создаёт подграф, полученный из g, удалением вершины node
func (g *graph) Delete(node int) (g1 graph) {
	g1.n = g.n - 1
	g1.A = make([][]byte, g1.n)
	k := 0
	for i, v := range g.A {
		if i == node {
			k = 1
			continue
		}
		g1.A[i-k] = make([]byte, g1.n)
		l := 0
		for j := range v {
			if j == node {
				l = 1
				continue
			}
			g1.A[i-k][j-l] = g.A[i][j]
		}
	}
	g1.Build()
	return g1
}

// Обход в глубину - проверяет связность графа
func (g *graph) Dfs(node int, visited *[]bool) {
	(*visited)[node] = true
	for _, v := range g.I[node] {
		if !(*visited)[v] {
			g.Dfs(v, visited)
		}
	}
}

// Проверяет, состоит ли массив только из `true`?
func isFull(a *[]bool) bool {
	for _, v := range *a {
		if !v {
			return false
		}
	}
	return true
}

// При запуске программы вызывается эта функция
func main() {
	var g graph
	g.Read()
	var results []int
	for v := range g.n { // Проверяем, является ли шарниром вершина v
		g1 := g.Delete(v)
		visited := make([]bool, g1.n)
		g1.Dfs(0, &visited)
		if !isFull(&visited) { // Если граф перестал быть связным, то v - шарнир
			results = append(results, v+1)
		}
	}
	fmt.Println(len(results))
	fmt.Println(results) // Нумерация вершин 1...n - в соответсвии с математикой, а не программированием
}
