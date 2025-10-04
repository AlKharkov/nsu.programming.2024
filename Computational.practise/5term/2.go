/*
Computational practice. Task 2; variant 19
*/

package main

import (
	"fmt"
	"math"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func main() {
	solve(0.1, 0.1, false)
	solve(0.1, 0.01, false)
	solve(0.1, 0.001, false)
	solve(0.01, 0.1, false)
	solve(0.01, 0.01, false)
	solve(0.01, 0.001, false)
	solve(0.001, 0.1, false)
	solve(0.001, 0.01, false)
	solve(0.001, 0.001, true)
}

func solve(h, tau float64, needPrint bool) {
	Nh, Nt := int(math.Round(1/h))+1, int(math.Round(1/tau))+1
	// Input data
	ux0 := func(x float64) float64 { return (x+0.1)*(x+0.1) + x }
	u0t := func(t float64) float64 { return 0.1*0.1 - math.Sin(2*math.Pi*t)/2 - 3.5*t }
	cxt := func(x, t float64) float64 { return (math.Pi*math.Cos(2*math.Pi*t) + 3.5) / (2*x + 1.2) }
	uAns := func(x, t float64) float64 { return (x+0.1)*(x+0.1) - math.Sin(2*math.Pi*t)/2 + x - 3.5*t }
	kToX := func(k int) float64 { return float64(k) * h }
	jToT := func(j int) float64 { return float64(j) * tau }
	fToGrid := func(fn func(float64, float64) float64, k, j int) float64 { return fn(kToX(k), jToT(j)) }
	// Start solution
	uGrid := make([][]float64, Nh)
	for k := range Nh {
		uGrid[k] = make([]float64, Nt)
		uGrid[k][0] = ux0(kToX(k))
	}
	for j := range Nt {
		uGrid[0][j] = u0t(jToT(j))
		uGrid[Nh-1][j] = u0t(methodNewton(alongCharacteristicLine(jToT(j))))
		// fmt.Println(uGrid[Nh-1][j] - uAns(1, jToT(j)))  // Characteristic deviations
	}
	// Start calculating ujk
	var vector []float64
	matrix := make([][]float64, Nh-2)
	for j := range Nt - 1 {
		for i := range Nh - 2 {
			matrix[i] = make([]float64, Nh-2)
		}
		vector = make([]float64, Nh-2)
		for k := 1; k < Nh-1; k++ {
			switch k {
			case 1:
				p := tau * (fToGrid(cxt, k, j) + fToGrid(cxt, k, j+1)) / 8 / h
				vector[0] = uGrid[k][j] + p*(uGrid[k-1][j+1]+uGrid[k-1][j]-uGrid[k+1][j])
				matrix[0][0] = 1
				matrix[0][1] = p
			case Nh - 2:
				p := 8 * h / tau / (fToGrid(cxt, k, j) + fToGrid(cxt, k, j+1))
				vector[k-1] = uGrid[k+1][j+1] - uGrid[k-1][j] + uGrid[k+1][j] - p*uGrid[k][j]
				matrix[k-1][k-2] = 1
				matrix[k-1][k-1] = -p
			default:
				p := 8 * h / tau / (fToGrid(cxt, k, j) + fToGrid(cxt, k, j+1))
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
	var maximumDeviation float64
	for j := range Nt {
		for k := range Nh {
			maximumDeviation = max(maximumDeviation, math.Abs(uGrid[k][j]-fToGrid(uAns, k, j)))
		}
	}
	fmt.Printf("tau=%.0e, h=%.0e: %f / %.2e\n", round(tau, 6), round(h, 6), maximumDeviation, maximumDeviation)
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

func alongCharacteristicLine(x float64) float64 {
	return math.Sin(2*math.Pi*x) + 7*x - 4.4
}

// sin(2*pi*x) + 7*x - c = 0
func methodNewton(c float64) (x float64) {
	f := func(x float64) float64 { return math.Sin(2*math.Pi*x) + 7*x - c }
	df := func(x float64) float64 { return 2*math.Pi*math.Cos(2*math.Pi*x) + 7 }
	x = 0.5
	for math.Abs(f(x)) > 0.000001 {
		x -= f(x) / df(x)
	}
	return
}
