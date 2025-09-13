package main

import (
	"fmt"
	"math"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func main() {
	// Input data
	a := 0.012
	mu := func(x float64) float64 { return x*x*x*x - x*x*x }
	mu1 := func(t float64) float64 { return t*t - t }
	mu2 := func(t float64) float64 { return t*t + t - t*math.Exp(1) }
	f := func(x, t float64) float64 { return x + 2*t - math.Exp(x) - a*(12*x*x-6*x-t*math.Exp(x)) }
	u := func(x, t float64) float64 { return x*x*x*x - x*x*x + t*x + t*t - t*math.Exp(x) }
	h, tau := 0.001, 0.001
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
	var vector []float64
	matrix := make([][]float64, Nh-2)
	for j := range Nt - 1 {
		for i := range Nh - 2 {
			matrix[i] = make([]float64, Nh-2)
		}
		vector = make([]float64, Nh-2)
		for k := 1; k < Nh-1; k++ {
			vector[k-1] = uGrid[k][j]/tau + gridFunc(f, k, j+1)
			if k == 1 {
				vector[k-1] = vector[k-1] + uGrid[k-1][j+1]*a/h/h
			} else { // k > 1
				matrix[k-1][k-2] = -a / h / h
			}
			matrix[k-1][k-1] = 1/tau + 2*a/h/h
			if k == Nh-2 {
				vector[k-1] = vector[k-1] + uGrid[k+1][j+1]*a/h/h
			} else { // k < Nh - 2
				matrix[k-1][k] = -a / h / h
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
			maximumDeviation = max(maximumDeviation, math.Abs(uGrid[k][j]-gridFunc(u, k, j)))
		}
	}
	fmt.Println("--------Task 1--------")
	fmt.Printf("Maximum deviation: %f\n\n", maximumDeviation)
	fmt.Printf("%.2e ", maximumDeviation)
	fmt.Println(maximumDeviation)

	fmt.Println("--------Task 2--------")
	// showMatrix(uGrid)
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
		values[i] = opts.LineData{Value: math.Abs(uGrid[Nh/2][int(math.Round(float64(i)*step/tau))] - u(0.5, float64(i)*step))}
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
func runThroughMethod(a [][]float64, b []float64) []float64 {
	n := len(a)
	for i := range n - 1 {
		b[i+1] = b[i+1] - b[i]*a[i+1][i]/a[i][i]
		a[i+1][i+1] = a[i+1][i+1] - a[i][i+1]*a[i+1][i]/a[i][i]
		a[i+1][i] = 0
	}
	x := make([]float64, n)
	x[n-1] = b[n-1] / a[n-1][n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = (b[i] - x[i+1]*a[i][i+1]) / a[i][i]
	}
	return x
}
