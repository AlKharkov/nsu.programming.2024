// v1.2 (final)
// 31/03/2026

package main

import (
	"fmt"
	. "math"
	"strings"
)

// func printM(m *[][]float64) {
// 	for _, a := range *m {
// 		for _, v := range a {
// 			if Abs(v) < Pow10(-6) {
// 				fmt.Print("    \t")
// 			} else {
// 				fmt.Printf("%.2f\t", v)
// 			}
// 		}
// 		fmt.Println()
// 	}
// 	fmt.Println()
// }

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
			A[i+1][j+1] = fmt.Sprintf("%d|%e", n, eps)
		}
	}
	fmt.Println("Формат вывода: <число итераций> | <Сеточная норма матрицы (U - P)>, где P - точное решение")
	printTable(&A)
}

func solve(dStep, N int) (int, float64) {
	var a, b float64 = 1.1, 0.8

	f := func(y, x float64) float64 {
		return 1.1*Sin(x) + (3.2*x*x+4.4*y*y)*Cos(2*x*y)
	}

	phi := func(y, x float64) float64 {
		return Sin(x) + Cos(2*x*y)
	}

	h := 1.0 / float64(N)

	// Определение внутренних точек
	var IntD [][2]int
	for i := 1; i < N; i++ {
		for j := 1; j < N; j++ {
			if i < j || i+j < N {
				IntD = append(IntD, [2]int{i, j})
			}
		}
	}

	// Оперделение индексов границы области D
	var dD [][2]int
	for i := 0; i <= N; i++ {
		for j := 0; j <= N; j++ {
			if (i == j && i >= N/2) || (i+j == N && i > N/2) || i == 0 || (j == 0 && i != 0 && i != N) || (j == N && i != 0 && i != N) {
				dD = append(dD, [2]int{i, j})
			}
		}
	}

	// Создаем матрицы для решений на итерациях U_n+1, U_n, и правой части
	U := make([][]float64, N+1)
	Uold := make([][]float64, N+1)
	F := make([][]float64, N+1)

	// Выделяем память
	for i := range N + 1 {
		U[i] = make([]float64, N+1)
		Uold[i] = make([]float64, N+1)
		F[i] = make([]float64, N+1)
	}

	// Заполняем точки границы точным решением
	for _, v := range dD {
		i, j := v[0], v[1]
		y, x := h*float64(i), h*float64(j)
		U[i][j] = phi(y, x)
		Uold[i][j] = phi(y, x)
	}

	// Заполняем правую часть для внутренних точек
	for _, v := range IntD {
		x := h * float64(v[1])
		y := h * float64(v[0])
		F[v[0]][v[1]] = f(y, x)
	}

	// Начальное приближение только для внутренних точек
	for _, v := range IntD {
		U[v[0]][v[1]] = 1.0
		Uold[v[0]][v[1]] = 1.0
	}

	// Скалярное произведение в сеточной норме
	scalar := func(A, B *[][]float64) float64 {
		result := 0.0
		for _, v := range IntD {
			result += (*A)[v[0]][v[1]] * (*B)[v[0]][v[1]]
		}
		return result * h * h
	}

	// Норма матрицы в сеточной норме
	scalarNorm := func(m *[][]float64) float64 {
		return Sqrt(scalar(m, m))
	}

	// Вычисляем ||U_n+1 - U_n||
	deltaCalc := func() float64 {
		deltaU := make([][]float64, N+1)
		for i := range N + 1 {
			deltaU[i] = make([]float64, N+1)
		}
		for _, v := range IntD {
			deltaU[v[0]][v[1]] = U[v[0]][v[1]] - Uold[v[0]][v[1]]
		}
		return scalarNorm(&deltaU)
	}

	// Оптимальный параметр
	spR := Cos(Pi / float64(N+1))
	omega := 2.0 / (1.0 + Sqrt(1.0-spR*spR))

	n := 0
	eps := Pow10(-dStep)
	cnt := 0

	// Основной цикл
	for n = 0; cnt < 3 && n < 1000; n++ {
		// Для выхода нужно набрать 3 подходящих итерации подряд
		if deltaCalc() < eps {
			cnt += 1
		} else {
			cnt = 0
		}

		// Сохраняем текущее решение
		for _, v := range IntD {
			Uold[v[0]][v[1]] = U[v[0]][v[1]]
		}

		// Итерация - обход в естественном порядке
		for _, v := range IntD {
			i, j := v[0], v[1]

			// Разностная схема: a*(u_{i,j-1} - 2u_{i,j} + u_{i,j+1})/h^2 +
			//                   b*(u_{i-1,j} - 2u_{i,j} + u_{i+1,j})/h^2 = -f_{i,j}
			//
			// Выражаем u_{i,j}: u_{i,j} = (a*(u_{i,j-1}+u_{i,j+1}) + b*(u_{i-1,j}+u_{i+1,j}) + h^2*f_{i,j}) / (2*(a+b))

			uWave := (a*(U[i][j-1]+Uold[i][j+1]) + b*(U[i-1][j]+Uold[i+1][j]) + h*h*F[i][j]) / (2.0 * (a + b))
			U[i][j] = (1.0-omega)*Uold[i][j] + omega*uWave
		}
	}

	// Вычисляем ошибку на внутренних точках, т.е ||U - Phi||
	err := 0.0
	for _, v := range IntD {
		x := h * float64(v[1])
		y := h * float64(v[0])
		diff := U[v[0]][v[1]] - phi(y, x)
		err += diff * diff
	}
	err = Sqrt(err * h * h)

	return n, err
}

// Формат вывода: <число итераций> | <Сеточная норма матрицы (U - P)>, где P - точное решение
// d\N    10                 20                 40
// 6      28|4.335256e-05    51|1.112073e-05    97|2.792101e-06
// 7      33|4.334247e-05    59|1.116293e-05    112|2.811497e-06
// 8      36|4.334415e-05    67|1.114860e-05    126|2.813621e-06
