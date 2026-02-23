// v1.1
// 14/02/2025

package main

import (
	"fmt"
	. "math"
	"strings"
)

func printM(m *[][]float64) {
	for _, a := range *m {
		for _, v := range a {
			if Abs(v) < Pow10(-6) {
				fmt.Print("    \t")
			} else {
				fmt.Printf("%.2f\t", v)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func printTable(m *[][]string) {
	l := make([]int, len((*m)[0]))
	for _, a := range *m {
		for j, v := range a {
			if len(v) > l[j] {
				l[j] = len(v)
			}
		}
	}
	var add int = 4
	for _, a := range *m {
		for j, v := range a {
			fmt.Print(v, strings.Repeat(" ", l[j]+add-len(v)))
		}
		fmt.Println()
	}
}

func main() {
	dStep := [...]int{6, 7, 8}
	N := [...]int{10, 20, 40}
	A := make([][]string, 4)
	for i := range A {
		A[i] = make([]string, 7)
	}
	A[0][0] = "d\\N"
	for j := 0; j < len(N); j++ {
		A[0][j+1] = fmt.Sprintf("%d", N[j])
	}
	for i := 0; i < len(dStep); i++ {
		A[i+1][0] = fmt.Sprintf("%d", dStep[i])
		for j := 0; j < len(N); j++ {
			nMax, lMax, nMin, lMin := solve(dStep[i], N[j])
			A[i+1][j+1] = fmt.Sprintf("%d|%f| |%d|%f", nMax, lMax, nMin, lMin)
		}
	}
	printTable(&A)
}

func solve(dStep, N int) (int, float64, int, float64) {
	var a, b float64 = 1.1, 0.8
	// f := func(x, y float64) float64 { return 1.1*Sin(x) + (3.2*x*x+4.4*y*y)*Cos(2*x*y) }
	phi := func(x, y float64) float64 { return 0 } //Sin(x) + Cos(2*x*y) }
	h := 1 / float64(N)

	// Создаём области для удобства циклов
	var dD [][2]int
	for j := 0; j <= N; j++ {
		dD = append(dD, [2]int{0, j})
	}
	for i := 1; i <= N; i++ {
		dD = append(dD, [2]int{i, 0})
	}
	for i := 1; i <= N; i++ {
		dD = append(dD, [2]int{i, N})
	}
	for k := (N + 1) / 2; k < N; k++ {
		dD = append(dD, [2]int{k, k})
	}
	for k := (N+1)/2 + 1; k < N; k++ {
		dD = append(dD, [2]int{k, N - k})
	}
	var IntD [][2]int
	for i := 1; i < N/2; i++ {
		for j := 1; j < N; j++ {
			IntD = append(IntD, [2]int{i, j})
		}
	}
	for i := N / 2; i < N; i++ {
		for j := 1; j < N-i; j++ {
			IntD = append(IntD, [2]int{i, j})
		}
		for j := i + 1; j < N; j++ {
			IntD = append(IntD, [2]int{i, j})
		}
	}
	var D [][2]int
	for i := 0; i <= N; i++ {
		for j := 0; j <= N; j++ {
			if j >= i || i+j <= N {
				D = append(D, [2]int{i, j})
			}
		}
	}
	/*
		var dD [][2]int
		for j := 0; j <= N; j++ {
			dD = append(dD, [2]int{0, j})
		}
		for i := 1; i <= N; i++ {
			dD = append(dD, [2]int{i, 0})
		}
		for i := 1; i <= N; i++ {
			dD = append(dD, [2]int{i, N})
		}
		for j := 1; j < N; j++ {
			dD = append(dD, [2]int{N, j})
		}
		var IntD [][2]int
		for i := 1; i < N; i++ {
			for j := 1; j < N; j++ {
				IntD = append(IntD, [2]int{i, j})
			}
		}
		var D [][2]int
		for i := 0; i <= N; i++ {
			for j := 0; j <= N; j++ {
				D = append(D, [2]int{i, j})
			}
		}
	*/

	U := make([][]float64, N+1)
	for i := 0; i <= N; i++ {
		U[i] = make([]float64, N+1)
	}

	// Вносим начальные данные
	for _, v := range dD {
		U[v[0]][v[1]] = phi(h*float64(v[0]), h*float64(v[1]))
	}
	// Рассматриваем IntD в качестве вектора
	scalar := func(A, B *[][]float64) (result float64) {
		for _, v := range IntD {
			result += (*A)[v[0]][v[1]] * (*B)[v[0]][v[1]]
		}
		return
	}
	norm := func(m *[][]float64) {
		var n float64 = Sqrt(scalar(m, m))
		for _, v := range IntD {
			(*m)[v[0]][v[1]] /= n
		}
	}
	// Заполняем начальное приближение
	for _, v := range IntD {
		U[v[0]][v[1]] = 1
	}
	norm(&U)

	U_old := make([][]float64, N+1)
	for i := 0; i <= N; i++ {
		U_old[i] = make([]float64, N+1)
	}
	A := func(i, j int) float64 {
		return -float64(N*N) * (a*(U_old[i-1][j]-2*U_old[i][j]+U_old[i+1][j]) + b*(U_old[i][j-1]-2*U_old[i][j]+U_old[i][j+1]))
	}
	delta, nA := Pow10(-dStep), 0
	var lA1, lA float64 = 1, 0
	deltaCalc := func(lm, lm1 float64) float64 {
		return Abs(lm1-lm) / Abs(lm)
	}
	for nA = 0; deltaCalc(lA, lA1) > delta && nA < 1000; nA++ {
		for _, v := range D {
			U_old[v[0]][v[1]] = U[v[0]][v[1]]
		}
		for _, v := range IntD {
			U[v[0]][v[1]] = A(v[0], v[1])
		}
		lA = lA1
		lA1 = scalar(&U, &U_old) / scalar(&U_old, &U_old)
		norm(&U)
	}
	if N == 10 && dStep == 6 {
		printM(&U)
		fmt.Println()
	}

	// Приступим к вычислению минимального с.з.
	V := make([][]float64, N+1)
	for i := 0; i <= N; i++ {
		V[i] = make([]float64, N+1)
	}
	// Вносим начальные данные
	for _, v := range dD {
		V[v[0]][v[1]] = phi(h*float64(v[0]), h*float64(v[1]))
	}
	// Заполняем начальное приближение
	for _, v := range IntD {
		V[v[0]][v[1]] = 1
	}
	norm(&V)
	V_old := make([][]float64, N+1)
	for i := 0; i <= N; i++ {
		V_old[i] = make([]float64, N+1)
	}
	B := func(i, j int) float64 {
		return lA1*V_old[i][j] + float64(N*N)*(a*(V_old[i-1][j]-2*V_old[i][j]+V_old[i+1][j])+b*(V_old[i][j-1]-2*V_old[i][j]+V_old[i][j+1]))
	}

	nB := 0
	var lB1, lB float64 = 1, 0
	for nB = 0; deltaCalc(lB, lB1) > delta && nB < 1000; nB++ {
		for _, v := range D {
			V_old[v[0]][v[1]] = V[v[0]][v[1]]
		}
		for _, v := range IntD {
			V[v[0]][v[1]] = B(v[0], v[1])
		}
		lB = lB1
		lB1 = scalar(&V, &V_old) / scalar(&V_old, &V_old)
		norm(&V)
	}

	return nA, lA1, nB, lA1 - lB1
}
