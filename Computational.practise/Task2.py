"""
version: 3.2.3

Task 1:
Interpolation polynomial in Lagrange form
Interpolation polynomial in Newton form
Interpolation polynomial in Lagrange form with Chebyshev`s nodes

Task 2:
Formula for the average rectangles, trapezoid, Simpson
Gaussian quadratures

Task 3:
The dichotomy method
The Newton method
"""
from math import exp, log10, cos, pi, sqrt, log
from numpy.polynomial import Polynomial as Poly
import matplotlib.pyplot as plt


def simpN(n: int | float, count=8) -> int | float:
    return (n * 10 ** (count + 4) + 5555) // 10 ** 4 / 10 ** count


def interpolationLagrange(args: list[int | float], values: list[int | float]) -> Poly:
    n = len(args) - 1
    omega = Poly.fromroots(args)
    Pn = Poly(0)
    for i in range(n + 1):
        t = omega // [-args[i], 1]
        Pn += values[i] * t / t(args[i])
    return Pn


def task11():
    # a, b = 0, 0.3
    x = 0.15
    argsA = [0, 0.1, 0.2, 0.3]
    valuesA = [1, 1.10517, 1.22140, 1.34986]
    Pa = interpolationLagrange(argsA, valuesA)
    argsB = [0.1, 0.2, 0.3]
    valuesB = [1.10517, 1.22140, 1.34986]
    Pb = interpolationLagrange(argsB, valuesB)
    print('Тест №1 для полинома Лагранжа | f(x) = exp(x)')
    print(f'Можно изменить: x = {x}')
    print(f'a) Pa({x}) = {simpN(Pa(x))} - значение 1 многочлена в точке х')
    print(f'b) Pb({x}) = {simpN(Pb(x))} - значение 2 многочлена в точке х')
    print(f'Pa({x}) - Pb({x}) = {simpN(Pa(x) - Pb(x))} - разность значений полученных многочленов в точке х')
    tabular = simpN(exp(x), 5)
    print(f'exp({x}) = {tabular} - табличное значение')
    print(f'Ra({x}) = exp({x}) - Pa({x}) = {simpN(tabular - Pa(x))} - погрешность вычисления в случае a)')
    print(f'Rb({x}) = exp({x}) - Pb({x}) = {simpN(tabular - Pb(x))} - погрешность вычисления в случае b)')


def interpolationNewton(args: list[float], values: list[float]) -> Poly:
    n = len(args) - 1
    div_dif = [[0 for _ in range(n + 1)] for i in range(n + 1)]  # Divided differences
    for i in range(n + 1):
        div_dif[i][0] = values[i]
    for j in range(1, n + 1):
        for i in range(j, n + 1):
            div_dif[i][j] = (div_dif[i][j - 1] - div_dif[i - 1][j - 1]) / (args[i] - args[abs(i - j)])
    Pn = Poly(values[0])
    t = Poly(1)
    for i in range(1, n + 1):
        t *= [-args[i - 1], 1]
        Pn += div_dif[i][i] * t
    return Pn


def get_chebyshev_nodes(a: float, b: float, n: int) -> list[float]:
    return [(a + b) / 2 + (b - a) * cos((2 * i - 1) * pi / 2 / n) / 2 for i in range(1, n + 1)]


def task12():
    a, b = 1, 10
    x = 5.25
    argsA = [1, 2, 4, 6, 8]
    valuesA = [0, 0.3010, 0.6021, 0.7782, 0.9031]
    Pa = interpolationNewton(argsA, valuesA)
    argsB = [2, 4, 6, 8, 10]
    valuesB = [0.3010, 0.6021, 0.7782, 0.9031, 1]
    Pb = interpolationNewton(argsB, valuesB)
    argsC = [3, 4, 5, 6, 7]
    valuesC = [0.4771, 0.6021, 0.6990, 0.7782, 0.8451]
    Pc = interpolationNewton(argsC, valuesC)
    print('----Тест №2 для полинома Ньютона | f(x) = lg(x)----')
    print(f'Можно изменить: x = {x}')
    print(f'a) Pa({x}) = {simpN(Pa(x))} - значение 1 многочлена в точке х')
    print(f'b) Pb({x}) = {simpN(Pb(x))}')
    print(f'c) Pc({x}) = {simpN(Pc(x))}')
    print(f'Pa({x}) - Pb({x}) = {simpN(Pa(x) - Pb(x))} - разность значений полученных многочленов в точке х')
    print(f'Pa({x}) - Pc({x}) = {simpN(Pa(x) - Pc(x))}')
    print(f'Pb({x}) - Pc({x}) = {simpN(Pb(x) - Pc(x))}')
    tabular = simpN(log10(x), 4)
    print(f'lg({x}) = {tabular} - табличное значение')
    print(f'Ra({x}) = lg({x}) - Pa({x}) = {simpN(tabular - Pa(x))} - погрешность вычисления в случае a)')
    print(f'Rb({x}) = lg({x}) - Pb({x}) = {simpN(tabular - Pb(x))} - b)')
    print(f'Rc({x}) = lg({x}) - Pc({x}) = {simpN(tabular - Pc(x))} - c)')
    print('----Интерполирование с узлами Чебышева----')
    n = len(argsA)
    argsChebyshev = get_chebyshev_nodes(a, b, n)
    valuesChebyshev = list(map(log10, argsChebyshev))
    PChebyshev = interpolationLagrange(argsChebyshev, valuesChebyshev)
    print('Узлы Чебышева:', list(map(simpN, argsChebyshev)))
    print(f'PChebyshev({x}) = {simpN(PChebyshev(x))} - значение многочлена, вычисленного по узлам Чебышева')
    print(f'lg({x}) - PChebyshev({x}) = {simpN(tabular - PChebyshev(x))}')
    argsEquidistant = [a + (b - a) * i / n for i in range(1, n + 1)]
    valuesEquidistant = list(map(log10, argsEquidistant))
    PEquidistant = interpolationLagrange(argsEquidistant, valuesEquidistant)
    print('Равноудаленные узлы:', list(map(simpN, argsEquidistant)))
    print(f'PEquidistant({x}) = {simpN(PEquidistant(x))} - значение многочлена, вычисленного по равноудаленным узлам')
    print(f'lg({x}) - PEquidistant({x}) = {simpN(tabular - simpN(PEquidistant(x)))}')


def task13():
    a, b = -1, 1
    argsEquidistant = [-1, -1 / 3, 1 / 3, 1]
    valuesEquidistant = list(map(exp, argsEquidistant))
    PEquidistant = interpolationLagrange(argsEquidistant, valuesEquidistant)
    n = len(argsEquidistant)
    argsChebyshev = get_chebyshev_nodes(a, b, n)
    valuesChebyshev = list(map(exp, argsChebyshev))
    PChebyshev = interpolationLagrange(argsChebyshev, valuesChebyshev)
    print('----Задание 1.3----')
    print('а) Для равноудаленных узлов:')
    for i in range(n):
        print(f'exp({simpN(argsEquidistant[i])}) = {simpN(valuesEquidistant[i])}')
    print('Для узлов Чебышева:')
    for i in range(n):
        print(f'exp({simpN(argsChebyshev[i])}) = {simpN(valuesChebyshev[i])}')
    args = list()
    valuesE = list()
    valuesC = list()
    point = a
    number_sections = 10 ** 4
    delta = (b - a) / number_sections
    xMaxE = 0
    xMaxC = 0
    maxMistakeE = 0
    maxMistakeC = 0
    for i in range(number_sections + 1):
        if point > b:
            break
        args.append(point)
        valuesE.append(PEquidistant(point))
        valuesC.append(PChebyshev(point))
        q = abs(exp(point) - valuesE[-1])
        if q > maxMistakeE:
            maxMistakeE = q
            xMaxE = point
        q = abs(exp(point) - valuesC[-1])
        if q > maxMistakeC:
            maxMistakeC = q
            xMaxC = point
        point += delta
    print(f'в) Для равноудаленных узлов максимальная ошибка...')
    print(f'достигается в точке {simpN(xMaxE)} и составляет {simpN(maxMistakeE)}')
    print(f'Для узлов Чебышева: {simpN(xMaxC)} и составляет {simpN(maxMistakeC)}')

    plt.plot(args, [abs(valuesE[i] - exp(args[i])) for i in range(len(args))])
    plt.plot(args, [abs(valuesC[i] - exp(args[i])) for i in range(len(args))])
    plt.title('Interpolation polynomials')
    plt.xlabel('Arg of function: x')
    plt.ylabel('Value of function: exp(x)')
    plt.grid()
    plt.legend(['Equidistant', 'Chebyshev'.format(63.0 * n)])
    plt.show()


def average_rectangles_formula(f, a: float, b: float, n: int) -> float:
    h = (b - a) / n
    xk = lambda k: a + k * h
    return sum([f((xk(i) + xk(i - 1)) / 2) for i in range(1, n + 1)]) * h


def trapezoid_formula(f, a: float, b: float, n: int) -> float:
    h = (b - a) / n
    xk = lambda k: a + k * h
    return sum([f(xk(i)) + f(xk(i - 1)) for i in range(1, n + 1)]) * h / 2


def simpson_formula(f, a: float, b: float, n: int) -> float:
    h = (b - a) / n
    xk = lambda k: a + k * h
    return sum([f(xk(i)) + 4 * f((xk(i) + xk(i - 1)) / 2) + f(xk(i - 1)) for i in range(1, n + 1)]) * h / 6


def task22():
    f = lambda x: sqrt(2 * x + 1)
    a, b = 0, 1
    I = sqrt(3) - 1 / 3
    print(f'"Идеальное"(табличное) значение интеграла I = {simpN(I)}')
    print('------------------------------------------------------------------------------------------------------')

    print('Значение квадратурной формулы центральных прямоугольников при:')
    I20 = average_rectangles_formula(f, a, b, 20)
    print(f'n = 20 равняется I_20={simpN(I20)}. Ошибка равняется I - I_20 = {simpN(I - I20)}')
    I50 = average_rectangles_formula(f, a, b, 50)
    print(f'n = 50 равняется I_50={simpN(I50)}. Ошибка равняется I - I_50 = {simpN(I - I50)}')
    I100 = average_rectangles_formula(f, a, b, 100)
    print(f'n = 100 равняется I_100={simpN(I100)}. Ошибка равняется I - I_100 = {simpN(I - I100)}')
    print('------------------------------------------------------------------------------------------------------')
    print('Значение квадратурной формулы трапеций при:')
    I20 = trapezoid_formula(f, a, b, 20)
    print(f'n = 20 равняется I_20={simpN(I20)}. Ошибка равняется I - I_20 = {simpN(I - I20)}')
    I50 = trapezoid_formula(f, a, b, 50)
    print(f'n = 50 равняется I_50={simpN(I50)}. Ошибка равняется I - I_50 = {simpN(I - I50)}')
    I100 = trapezoid_formula(f, a, b, 100)
    print(f'n = 100 равняется I_100={simpN(I100)}. Ошибка равняется I - I_100 = {simpN(I - I100)}')
    print('------------------------------------------------------------------------------------------------------')
    print('Значение квадратурной формулы Симпсона при:')
    I20 = simpson_formula(f, a, b, 20)
    print(f'n = 20 равняется I_20={simpN(I20)}. Ошибка равняется I - I_20 = {simpN(I - I20)}')
    I50 = simpson_formula(f, a, b, 50)
    print(f'n = 50 равняется I_50={simpN(I50)}. Ошибка равняется I - I_50 = {simpN(I - I50)}')
    I100 = simpson_formula(f, a, b, 100)
    print(f'n = 100 равняется I_100={simpN(I100)}. Ошибка равняется I - I_100 = {simpN(I - I100)}')


def gauss_2_nodes(f, a: float, b: float) -> float:
    c = (a + b) / 2
    print('Узлы квадратурной формулы:')
    x0, x1 = c - (b - a) / 2 / sqrt(3), c + (b - a) / 2 / sqrt(3)
    print(f'x0 = {simpN(x0)} | x1 = {simpN(x1)}')
    print('Значения функции в этих узлах:')
    print(f'f(x0) = {simpN(f(x0))} | f(x1) = {simpN(f(x1))}')
    print('Коэффициенты квадратурной формулы:')
    print(f'A0 = (b - a) / 2 = {(b - a) / 2} = A1')
    return (f(c - (b - a) / 2 / sqrt(3)) + f(c + (b - a) / 2 / sqrt(3))) * (b - a) / 2


def gauss_3_nodes(f, a: float, b: float) -> float:
    print('Узлы квадратурной формулы:')
    d1, d2 = 0, 0.7745966692
    d_to_x = lambda d: (d * (b - a) + a + b) / 2
    D1, D2 = 0.8888888888, 0.5555555556
    x0, x1, x2 = d_to_x(-d2), d_to_x(d1), d_to_x(d2)
    print(f'x0 = {simpN(x0)} | x1 = {simpN(x1)} | x2 = {simpN(x2)}')
    print('Значения функции в этих узлах:')
    print(f'f(x0) = {simpN(f(x0))} | f(x1) = {simpN(f(x1))} | f(x2) = {simpN(f(x2))}')
    print('Коэффициенты квадратурной формулы:')
    print(f'A0 = {D1 / 2} = A2 | A1 = {D2 / 2}')
    return (D2 * f(x0) + D1 * f(x1) + D2 * f(x2)) * (b - a) / 2


def gauss_5_nodes(f, a: float, b: float) -> float:
    print('Узлы квадратурной формулы:')
    d1, d2, d3 = 0, 0.5384693101, 0.9061798459
    d_to_x = lambda d: (d * (b - a) + a + b) / 2
    D1, D2, D3 = 0.5688888888, 0.4786286705, 0.2369268851
    x0, x1, x2, x3, x4 = d_to_x(-d3), d_to_x(-d2), d_to_x(d1), d_to_x(d2), d_to_x(d3)
    print(f'x0 = {simpN(x0)} | x1 = {simpN(x1)} | x2 = {simpN(x2)} | x3 = {simpN(x3)} | x4 = {simpN(x4)}')
    print('Значения функции в этих узлах:')
    print(f'f(x0) = {simpN(f(x0))} | f(x1) = {simpN(f(x1))} |'
          f' f(x2) = {simpN(f(x2))} | f(x3) = {simpN(f(x3))} | f(x4) = {simpN(x4)}')
    print('Коэффициенты квадратурной формулы:')
    print(f'A0 = {D3 / 2} = A4 | A1 = {D2 / 2} = A3 | A2 = {D1 / 2}')
    return (D3 * f(x0) + D2 * f(x1) + D1 * f(x2) + D2 * f(x3) + D3 * f(x4)) * (b - a) / 2


def task23():
    f = lambda x: sqrt(2 * x + 1)
    a, b = 0, 1
    I = sqrt(3) - 1 / 3
    print(f'"Идеальное"(табличное) значение интеграла I = {simpN(I)}')
    print('------------------------------------------------------------------------------------------------------')
    print(f'Квадратура Гаусса на 2 узлах:')
    I2 = gauss_2_nodes(f, a, b)
    print(f'Значение квадратуры I2 = {simpN(I2)}')
    print(f'Полученная ошибка вычисления составляет R2 = {simpN(I - I2)}')
    print('------------------------------------------------------------------------------------------------------')
    print('Квадратура Гаусса на 3 узлах:')
    I3 = gauss_3_nodes(f, a, b)
    print(f'Значение квадратуры I3 = {simpN(I3)}')
    print(f'Полученная ошибка вычисления составляет R3 = {simpN(I - I3)}')
    print('------------------------------------------------------------------------------------------------------')
    print('Квадратура Гаусса на 5 узлах:')
    I5 = gauss_5_nodes(f, a, b)
    print(f'Значение квадратуры I5 = {simpN(I5)}')
    print(f'Полученная ошибка вычисления составляет R5 = {simpN(I - I5)}')


def dichotomy(f, a: float, b: float, eps=10 ** -8) -> (float, int):
    x = (a + b) / 2
    num_iter = 0
    while b > a and b - a > eps:
        num_iter += 1
        if f(a) * f(x) < 0:
            b = x
        else:
            a = x
        x = (a + b) / 2
    return x, num_iter


def the_Newton_method(f, df, x0: float, eps=10 ** -8) -> (float, int):
    x_np1 = x0
    num_iter = 0
    x_n = x_np1 + 2
    x_nm1 = x_n + 8
    while abs((x_np1 - x_n) / (1 - (x_np1 - x_n) / (x_n - x_nm1))) > eps:
        num_iter += 1
        x_nm1 = x_n
        x_n = x_np1
        x_np1 -= f(x_np1) / df(x_np1)
    return x_np1, num_iter


def task31():
    f = lambda x: exp(x) + x
    df = lambda x: exp(x) + 1
    a, b = -1, 1
    x1, n1 = dichotomy(f, a, b)
    print(f'Метод дихотомии нашел корень: {x1} за {n1} итераций')
    x2, n2 = the_Newton_method(f, df, 0)
    print(f'Метод Ньютона нашел корень: {x2} за {n2} итераций')


def draw_graphic(a: float, b: float, n=10 ** 6):
    f = lambda x: x ** 3 - x
    h = (b - a) / n
    args = [a + i * h for i in range(n)]
    values = [f(arg) for arg in args]
    plt.plot(args, values)
    plt.title('Graphic')
    plt.xlabel('Arg x')
    plt.ylabel('The value of x ** 3 - x')
    plt.grid()
    plt.show()


def draw32(a: float, b: float, needIterations=False, n=10 ** 6):
    f = lambda x: x ** 3 - x
    df = lambda x: 3 * x ** 2 - 1
    h = (b - a) / n
    args = [a + i * h for i in range(n)]
    values = [the_Newton_method(f, df, arg) for arg in args]
    valuesX = [value[0] for value in values]
    valuesN = [value[1] for value in values]
    plt.plot(args, valuesX)
    plt.title('The Newton method values')
    plt.xlabel('First x: x0')
    plt.ylabel('Value of the Newton method: ')
    plt.grid()
    plt.show()
    if needIterations:
        plt.plot(args, valuesN)
        plt.title('The number of iterations')
        plt.xlabel('First x: x0')
        plt.ylabel('The number of iterations: ')
        plt.grid()
        plt.show()


def task32():
    draw32(-100, 100, True)
    draw32(-1, 1)


if __name__ == '__main__':
    task32()
