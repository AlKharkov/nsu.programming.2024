/*
Computing workshop: task 1; variant 9(the implicit scheme)

Version 1.0
Last updated: 12.09.2025
*/

package main

import (
	"fmt"
	"math"
)

func main() {
	// Input data
	a := 0.012
	mu := func(x float64) float64 { return x*x*x*x - x*x*x }
	mu1 := func(x float64) float64 { return x*x - x }
	mu2 := func(x float64) float64 { return x*x + x - x*math.Exp(1) }
	f := func(x, t float64) float64 { return x + 2*t - math.Exp(x) - a*(12*x*x-6*x-t*math.Exp(x)) }
	u := func(x, t float64) float64 { return x*x*x*x - x*x*x + t*x + t*t - t*math.Exp(x) }
	h, tau := 0.1, 0.1
	gridFunc := func(fn func(float64, float64) float64, k, j int) float64 { return fn(float64(k)*h, float64(j)*tau) }
	Nh, Nt := int(1/h)+1, int(1/tau)+1
	// Start solution
	uGrid := make([][]float64, Nh)
	for k := range Nh {
		uGrid[k] = make([]float64, Nt)
		uGrid[k][0] = mu(float64(k) * h)
	}
	for j := range Nt {
		uGrid[0][j] = mu1(float64(j) * tau)
		uGrid[Nh-1][j] = mu2(float64(j) * tau)
	}
	// Start calculating ujk
	for q := Nt - 2; q > 0; q-- {
		for p := 1; p < Nh-1; p++ {
			uGrid[p][q] = tau * (-a/h/h*uGrid[p-1][q+1] + (1/tau+2*a/h/h)*uGrid[p][q+1] - a/h/h*uGrid[p+1][q+1] - gridFunc(f, p, q+1))
		}
	}
	// Calculating maximum deviation between uGrid & u
	var maximumDeviation float64
	for j := range Nt {
		for k := range Nh {
			maximumDeviation = max(maximumDeviation, math.Abs(uGrid[k][j]-gridFunc(u, k, j)))
		}
	}
	fmt.Println("--------Task 1--------")
	fmt.Printf("Maximum deviation: %f\n\n", maximumDeviation)

	fmt.Println("--------Task 2--------")
	showMatrix(uGrid)
	fmt.Println()

	step := 0.1
	var args [11]float64
	var values [11]float64
	for i := range 11 {
		args[i] = step * float64(i)
		values[i] = uGrid[int(Nh/2)][int(float64(i)*step/tau)]
	}
	fmt.Println("--------Task 3--------")
	fmt.Println(args)
	fmt.Println(values)

}

func showMatrix(m [][]float64) {
	for _, v := range m {
		for _, u := range v {
			if u != 0 {
				fmt.Printf("%.0e\t", u)
			} else {
				fmt.Print("0\t")
			}
		}
		fmt.Println()
	}
}
