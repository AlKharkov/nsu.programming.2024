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
	h, tau := 0.01, 0.01
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
	for p := Nh - 2; p > 0; p-- {
		for q := 1; q < Nt-1; q++ {
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
	// showMatrix(uGrid)
	showMatrix(make([][]float64, 0))
	fmt.Println()

	fmt.Println("--------Task 3--------")
	step := 0.01
	iN := int(math.Round(1/step)) + 1
	args := make([]float64, iN)
	values1 := make([]opts.LineData, iN)
	values2 := make([]opts.LineData, iN)
	for i := range iN {
		args[i] = round(step*float64(i), 6)
		// values[i] = opts.LineData{Value: math.Abs(uGrid[Nh/2][int(math.Round(float64(i)*step/tau))] - u(0.5, float64(i)*step))}
		values1[i] = opts.LineData{Value: uGrid[Nh/2][int(math.Round(float64(i)*step/tau))]}
		values2[i] = opts.LineData{Value: gridFunc(u, Nh/2, int(math.Round(float64(i)*step/tau)))}
	}
	chart := charts.NewLine()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Decision error graph"}))
	chart.SetXAxis(args).AddSeries("uGrid", values1)
	chart.SetXAxis(args).AddSeries("uExact", values2)
	file, _ := os.Create("task3.html")
	defer file.Close()
	chart.Render(file)
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
