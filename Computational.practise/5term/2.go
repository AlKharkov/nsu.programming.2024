/*
Computational practice. Task 2; variant 19
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

// Через характеристики
// tau/h	0.1			0.01		0.001		0.0001
// 0.1		1.17e-01	1.17e-01	1.17e-01	1.17e-01
// 0.01		3.55e-02	1.10e-02	1.10e-02	1.10e-02
// 0.001	3.25e-02	1.09e-03	1.09e-03	1.09e-03
// 0.0001	3.24e-02	3.56e-04	1.09e-04	1.09e-04

// Точное
// tau/h	0.1			0.01		0.001
// 0.1		4.28e-02	4.56e-04	4.56e-06
// 0.01		5.04e-02	4.77e-04	4.76e-06
// 0.001	5.06e-02	4.77e-04	4.77e-06

// Явная схема
// tau/h	0.1			0.01		0.001		0.0001
// 0.1		5.36e-02	6.75e-03	5.32e-03	5.54e-03
// 0.01		9.05e+06	9.24e-04	1.13e-04	6.23e-05
// 0.001	4.01e+17	3.05e+94	4.01e-04	1.22e-06
// 0.0001	5.15e+26	8.25e+210	NaN			7.71e+17

// Центральная явная + дополнительная
// tau/h	0.1			0.01		0.001
// 0.1		8.43e-02	7.09e-03	7.60e-04
// 0.01		4.93e+01	1.03e-03	1.18e-04
// 0.001	6.79e+04	7.52e+38	8.33e-06

func main() {
	p := 3
	eps := createEps(p)
	matrix := make([][]float64, p)
	for i := range p {
		matrix[i] = make([]float64, p)
		for j := range p {
			if i == p-1 && j == p-1 {
				matrix[i][j] = solve(eps[i], eps[j], false)
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
	ux0 := func(x float64) float64 { return (x+0.1)*(x+0.1) + x }
	u0t := func(t float64) float64 { return 0.1*0.1 - math.Sin(2*math.Pi*t)/2 - 3.5*t }
	cxt := func(x, t float64) float64 { return (math.Pi*math.Cos(2*math.Pi*t) + 3.5) / (2*x + 1.2) }
	uAns := func(x, t float64) float64 { return (x+0.1)*(x+0.1) - math.Sin(2*math.Pi*t)/2 + x - 3.5*t }
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
	}
	// Start calculating ujk
	var vector []float64
	matrix := make([][]float64, Nh-2)
	for j := range Nt - 1 {
		for i := range Nh - 2 {
			matrix[i] = make([]float64, Nh-2)
		}
		vector = make([]float64, Nh-2)
		//k := Nh - 1
		//c := (toGrid(cxt, k, j+1) + 2*toGrid(cxt, k, j) + toGrid(cxt, k-1, j)) / 4
		//uGrid[k][j+1] = tau*c*uGrid[k-1][j]/h + (h-tau*c)*uGrid[k][j]/h
		//uGrid[Nh-1][j+1] = toGrid(uAns, Nh-1, j+1)
		k := Nh - 2
		uGrid[k][j+1] = tau*toGrid(cxt, k, j)*(uGrid[k-1][j]-uGrid[k+1][j])/2/h + uGrid[k][j]
		k += 1
		c := cxt((float64(k)-0.5)*h, (float64(j)+0.5)*tau)
		c1, c2 := h/(h+tau*c), tau*c/(h+tau*c)
		uGrid[k][j+1] = (uGrid[k][j]-uGrid[k-1][j+1])*(c1-c2) + uGrid[k-1][j]*(c1+c2)
		for k := 1; k < Nh-1; k++ {
			switch k {
			case 1:
				p := tau * (toGrid(cxt, k, j) + toGrid(cxt, k, j+1)) / 8 / h
				vector[0] = uGrid[k][j] + p*(uGrid[k-1][j+1]+uGrid[k-1][j]-uGrid[k+1][j])
				matrix[0][0] = 1
				matrix[0][1] = p
			case Nh - 2:
				p := 8 * h / tau / (toGrid(cxt, k, j) + toGrid(cxt, k, j+1))
				vector[k-1] = uGrid[k+1][j+1] - uGrid[k-1][j] + uGrid[k+1][j] - p*uGrid[k][j]
				matrix[k-1][k-2] = 1
				matrix[k-1][k-1] = -p
			default:
				p := 8 * h / tau / (toGrid(cxt, k, j) + toGrid(cxt, k, j+1))
				vector[k-1] = uGrid[k+1][j] - uGrid[k-1][j] - p*uGrid[k][j]
				matrix[k-1][k-2] = 1
				matrix[k-1][k-1] = -p
				matrix[k-1][k] = -1
			}
		}
		// Now we have 3-diagonal matrix and vector so that Au=b
		for i, v := range runThroughMethod(matrix, vector) {
			uGrid[i+1][j+1] = v
		}
	}
	// Calculating maximum deviation between uGrid & u
	for j := range Nt {
		for k := range Nh {
			maxDeviation = max(maxDeviation, math.Abs(uGrid[k][j]-toGrid(uAns, k, j)))
		}
		//maxDeviation = max(maxDeviation, math.Abs(uGrid[Nh-1][j]-fToGrid(uAns, Nh-1, j))) // for border
	}
	if needPrint {
		//fmt.Println("--------Task 2--------")
		//showMatrix(uGrid)
		showMatrix(make([][]float64, 0))
		fmt.Println()

		fmt.Println("--------Task 3--------")
		step := 0.01
		iN := int(math.Round(1/step)) + 1
		args := make([]float64, iN)
		values := make([]opts.LineData, iN)
		//values1 := make([]opts.LineData, iN)
		//values2 := make([]opts.LineData, iN)
		for i := range iN {
			args[i] = round(step*float64(i), 6)
			values[i] = opts.LineData{Value: math.Abs(uGrid[Nh/2][int(math.Round(float64(i)*step/tau))] - uAns(0.5, float64(i)*step))}
			//values1[i] = opts.LineData{Value: uGrid[Nh/2][int(math.Round(float64(i)*step/tau))]}
			//values2[i] = opts.LineData{Value: gridFunc(u, Nh/2, int(math.Round(float64(i)*step/tau)))}
		}
		chart := charts.NewLine()
		chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Decision error"}))
		chart.SetXAxis(args).AddSeries("deviation", values)
		//chart.SetXAxis(args).AddSeries("uGrid", values1)
		//chart.SetXAxis(args).AddSeries("uExact", values2)
		file, _ := os.Create("task3.html")
		defer file.Close()
		chart.Render(file)
		fmt.Println("Successful! Look at the file <task3.html>")
	}
	return
}

func showMatrix(m [][]float64) {
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

// Ax = b. Find x
func runThroughMethod(a [][]float64, b []float64) (x []float64) {
	n := len(a)
	for i := range n - 1 {
		b[i+1] = b[i+1] - b[i]/a[i][i]
		a[i+1][i+1] = a[i+1][i+1] - a[i][i+1]/a[i][i]
		a[i+1][i] = 0
	}
	x = make([]float64, n)
	x[n-1] = b[n-1] / a[n-1][n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = (b[i] - x[i+1]*a[i][i+1]) / a[i][i]
	}
	return x
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
