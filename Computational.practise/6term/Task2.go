// v2.0
// 23/02/2025

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
	A := make([][]string, len(dStep)+1)
	for i := range A {
		A[i] = make([]string, len(N)+1)
	}
	A[0][0] = "d\\N"
	for j := range len(N) {
		A[0][j+1] = fmt.Sprintf("%d", N[j])
	}
	for i := range len(dStep) {
		A[i+1][0] = fmt.Sprintf("%d", dStep[i])
		for j := range len(N) {
			n, eps := solve(dStep[i], N[j])
			A[i+1][j+1] = fmt.Sprintf("%d|%f", n, eps)
		}
	}
	printTable(&A)
}

func solve(dStep, N int) (int, float64) {
	var a, b float64 = 1.1, 0.8
	f := func(x, y float64) float64 { return 1.1*Sin(x) + (3.2*x*x+4.4*y*y)*Cos(2*x*y) }
	phi := func(x, y float64) float64 { return Sin(x) + Cos(2*x*y) }
	h := 1 / float64(N)

	// Создаём области для удобства циклов
	// Случай М из варианта
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
		// Случай квадрата
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

	// Рассматриваем IntD в качестве вектора
	scalar := func(A, B *[][]float64) (result float64) {
		for _, v := range IntD {
			result += (*A)[v[0]][v[1]] * (*B)[v[0]][v[1]]
		}
		return
	}
	scalarNorm := func(m *[][]float64) float64 { return scalar(m, m) }
	// norm := func(m *[][]float64, n float64) {
	// 	if n == -1 {
	// 		n = Sqrt(scalar(m, m))
	// 	}
	// 	for _, v := range IntD {
	// 		(*m)[v[0]][v[1]] /= n
	// 	}
	// }

	init := func(m *[][]float64) {
		for i := range N + 1 {
			(*m)[i] = make([]float64, N+1)
		}
	}

	X := make([][]float64, N+1)
	init(&X)
	// Заполняем начальное приближение
	for _, v := range IntD {
		X[v[0]][v[1]] = 1.0
	}
	// Вносим начальные данные
	for _, v := range dD {
		X[v[0]][v[1]] = phi(h*float64(v[0]), h*float64(v[1]))
	}

	A := func(i, j int, B *[][]float64) float64 {
		return -float64(N*N) * (a*((*B)[i-1][j]-2*(*B)[i][j]+(*B)[i+1][j]) + b*((*B)[i][j-1]-2*(*B)[i][j]+(*B)[i][j+1]))
	}

	F := make([][]float64, N+1)
	init(&F)
	for _, v := range IntD {
		F[v[0]][v[1]] = f(h*float64(v[0]), h*float64(v[1]))
	}

	KsiOld := make([][]float64, N+1)
	init(&KsiOld)

	tauCalc := func() float64 {
		AKsi := make([][]float64, N+1)
		init(&AKsi)
		for _, v := range IntD {
			AKsi[v[0]][v[1]] = A(v[0], v[1], &KsiOld)
		}
		normAKsi := Sqrt(scalarNorm(&AKsi))
		for _, v := range IntD {
			AKsi[v[0]][v[1]] /= normAKsi
		}
		KsiCopy := make([][]float64, N+1)
		init(&KsiCopy)
		for _, v := range IntD {
			KsiCopy[v[0]][v[1]] = KsiOld[v[0]][v[1]] / normAKsi
		}
		return scalar(&AKsi, &KsiCopy)
	}

	XOld := make([][]float64, N+1)
	init(&XOld)

	Ksi := make([][]float64, N+1)
	init(&Ksi)

	calcKsi := func(K, X *[][]float64) {
		for _, v := range IntD {
			(*K)[v[0]][v[1]] = A(v[0], v[1], X) - F[v[0]][v[1]]
		}
	}

	deltaCalc := func() float64 {
		normKsi := Sqrt(scalarNorm(&KsiOld))
		deltaKsi := make([][]float64, N+1)
		init(&deltaKsi)
		for _, v := range IntD {
			deltaKsi[v[0]][v[1]] = (Ksi[v[0]][v[1]] - KsiOld[v[0]][v[1]]) / normKsi
		}
		return scalarNorm(&deltaKsi)
	}

	n := 0
	for delta := Pow10(-dStep); n > 0 && deltaCalc() > delta && n < 1000; n++ {
		for _, v := range IntD {
			XOld[v[0]][v[1]] = X[v[0]][v[1]]
		}
		calcKsi(&KsiOld, &XOld)
		tau := tauCalc()
		for _, v := range IntD {
			X[v[0]][v[1]] = XOld[v[0]][v[1]] - tau*KsiOld[v[0]][v[1]]
		}
		calcKsi(&Ksi, &X)
	}
	if N == 10 && dStep == 6 {
		printM(&X)
		fmt.Println()
	}

	Eps := make([][]float64, N+1)
	init(&Eps)
	for _, v := range IntD {
		Eps[v[0]][v[1]] = phi(h*float64(v[0]), h*float64(v[1])) - X[v[0]][v[1]]
	}

	return n, scalarNorm(&Eps)
}
