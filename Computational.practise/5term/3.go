/*
Computational practice. Task 3; variant 9
*/

package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// tau/h	0.1			0.01		0.001		0.0001
// 0.1		9.97e-02	9.75e-03	7.67e-04	1.32e-04
// 0.01		2.87e+11	3.82e+50	9.97e-04	9.76e-05
// 0.001	3.64e+29	5.49e+258	NaN			NaN
// 0.0001	3.65e+47	NaN			NaN			NaN

func main() {
	p := 4
	eps := createEps(p)
	matrix := make([][]float64, p)
	for i := range p {
		matrix[i] = make([]float64, p)
		for j := range p {
			if i == p-1 && j == p-1 {
				matrix[i][j] = solve(eps[i], eps[j], true) // true/false
			} else {
				matrix[i][j] = solve(eps[i], eps[j], false)
			}
		}
	}
	showEpsTable(&matrix)
}

func solve(h, tau float64, needPrint bool) (maxDeviation float64) {
	Nh, Nt := int(math.Round(1/h))+1, int(math.Round(1/tau))+1
	// Input data
	a := 0.012
	ux0 := func(x float64) float64 { return math.Pow(x, 4) - math.Pow(x, 3) }
	u0t := func(t float64) float64 { return t*t - t }
	u1t := func(t float64) float64 { return t*t + t - t*math.Exp(1) }
	f := func(x, t float64) float64 { return x + 2*t - math.Exp(x) - a*(12*x*x-6*x-t*math.Exp(x)) }
	uAns := func(x, t float64) float64 { return math.Pow(x, 4) - math.Pow(x, 3) + t*x + t*t - t*math.Exp(x) }
	// Auxilary functions
	kToX := func(k int) float64 { return float64(k) * h }
	jToT := func(j int) float64 { return float64(j) * tau }
	toGrid := func(fn func(float64, float64) float64, k, j int) float64 { return fn(kToX(k), jToT(j)) }
	// Start solution
	uGrid := make([][]float64, Nh)
	for k := range Nh {
		uGrid[k] = make([]float64, Nt)
		uGrid[k][0] = ux0(kToX(k))
	}
	for j := range Nt {
		uGrid[0][j] = u0t(jToT(j))
		uGrid[Nh-1][j] = u1t(jToT(j))
	}
	// Start calculating ujk
	for j := range Nt - 1 {
		for k := 1; k < Nh-1; k++ {
			uGrid[k][j+1] = a*tau*(uGrid[k+1][j]+uGrid[k-1][j])/h/h + uGrid[k][j]*(1-2*a*tau/h/h) + tau*toGrid(f, k, j)
		}
	}
	// Calculating maximum deviation between uGrid & u
	for j := range Nt {
		for k := range Nh {
			maxDeviation = max(maxDeviation, math.Abs(uGrid[k][j]-toGrid(uAns, k, j)))
		}
	}
	if needPrint {
		showMatrix(make([][]float64, 0))
		step := 0.01
		iN := int(math.Round(1/step)) + 1
		args := make([]float64, iN)
		values := make([]opts.LineData, iN)
		for i := range iN {
			args[i] = round(step*float64(i), 6)
			values[i] = opts.LineData{Value: math.Abs(uGrid[Nh/2][int(math.Round(float64(i)*step/tau))] - uAns(0.5, float64(i)*step))}

		}
		chart := charts.NewLine()
		chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Decision error"}))
		chart.SetXAxis(args).AddSeries("deviation", values)
		file, _ := os.Create("task3.html")
		defer file.Close()
		chart.Render(file)
		fmt.Printf("Task 3: Successful! Look at the file <task3.html>\n\n")
	}
	return
}

func showMatrix(m [][]float64) {
	if len(m) == 0 {
		return
	}
	for i := range m {
		v := m[len(m)-i-1]
		for _, u := range v {
			if u != 0 {
				fmt.Printf("%.0e\t", u)
			} else {
				fmt.Print("0\t")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func round(x float64, n int) float64 {
	return math.Round(x*math.Pow10(n)) / math.Pow10(n)
}

func createEps(n int) []float64 {
	eps := make([]float64, n)
	eps[0] = 0.1
	for i := 1; i < n; i++ {
		eps[i] = eps[i-1] / 10
	}
	return eps
}

func showEpsTable(m *[][]float64) {
	n := len(*m)
	// eps := createEps(n)
	fmt.Printf("tau/h")
	fmt.Print("\t0.1")
	for i := 1; i < n; i++ {
		fmt.Print("\t\t0." + strings.Repeat("0", i) + "1")
	}
	fmt.Println()
	for i := range n {
		fmt.Print("0." + strings.Repeat("0", i) + "1")
		for j := range n {
			fmt.Printf("\t%.2e", (*m)[i][j])
		}
		fmt.Println()
	}
}
