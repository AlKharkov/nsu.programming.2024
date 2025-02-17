"""
version: 1.2.0

Interpolation polynomial in Lagrange form
Interpolation polynomial in Newton form
"""

from math import exp, log10
from numpy.polynomial import Polynomial as Poly


def simpN(n: int | float, count=8) -> int | float:
    if abs(n) <= 10 ** -count:
        cnt = 0
        while n < 1:
            n *= 10
            cnt += 1
        return simpN(n, 0) / 10 ** cnt
    t = n * 10 ** count
    if t % 1 >= 0.5:
        return int(t + 0.5) / 10 ** count
    return int(t) / 10 ** count


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
    div_dif = [[0 for _ in range(i + 1)] for i in range(n + 1)]  # Divided differences
    for i in range(n + 1):
        div_dif[i][0] = values[i]
    for j in range(1, n + 1):
        for i in range(j, n + 1):
            div_dif[i][j] = (div_dif[i][j - 1] - div_dif[i - 1][j - 1]) / (args[i] - args[i - 1])
    Pn = Poly(values[0])
    t = Poly(1)
    for i in range(1, n + 1):
        t *= [-args[i - 1], 1]
        Pn += div_dif[i][i] * t
    return Pn


def task12():
    # a, b = 1, 10
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
    print('Тест №2 для полинома Ньютона | f(x) = lg(x)')
    print(f'Можно изменить: x = {x}')
    print(f'a) Pa({x}) = {simpN(Pa(x))} - значение 1 многочлена в точке х')
    print(f'b) Pb({x}) = {simpN(Pb(x))} - значение 2 многочлена в точке х')
    print(f'c) Pc({x}) = {simpN(Pc(x))} - значение 3 многочлена в точке х')
    print(f'Pa({x}) - Pb({x}) = {simpN(Pa(x) - Pb(x))} - разность значений полученных многочленов в точке х')
    print(f'Pa({x}) - Pc({x}) = {simpN(Pa(x) - Pc(x))}')
    print(f'Pb({x}) - Pc({x}) = {simpN(Pb(x) - Pc(x))}')
    tabular = simpN(log10(x), 4)
    print(f'lg({x}) = {tabular} - табличное значение')
    print(f'Ra({x}) = lg({x}) - Pa({x}) = {simpN(tabular - Pa(x))} - погрешность вычисления в случае a)')
    print(f'Rb({x}) = lg({x}) - Pb({x}) = {simpN(tabular - Pb(x))} - погрешность вычисления в случае b)')
    print(f'Rc({x}) = lg({x}) - Pc({x}) = {simpN(tabular - Pc(x))} - погрешность вычисления в случае c)')


if __name__ == '__main__':
    task12()
