"""
version: 1.10
Gauss method
Determinant(with upper triangular form)
Inverse matrix(with extended matrix form)
Run-throw method
Reflection method(Householder's)
Richardson method with Chebyshev`s params
Conjugate gradient method
Power method
"""


def simplify_number(n: float, accuracy=8, sign_space=False) -> str:
    m = str(round(n * 10 ** accuracy) / 10 ** accuracy)
    if '.' in m and '9' in m:
        a = m.rindex('9')
        b = a
        while b > 1 and m[b - 1] == '9':
            b -= 1
        if 2 * (a - b + 1) >= accuracy:
            if m[b - 1] == '.':
                m = str(int(m[:b - 1]) + 1)
            else:
                m = m[:b - 1] + str(int(m[b - 1]) + 1)
    if abs(round(float(m)) - float(m)) < 10 ** (-accuracy // 4):
        m = str(round(float(m)))
    if sign_space:
        return m if m[0] == '-' else ' ' + m
    return m


def print_matrix(A: list[list[float]], accuracy=8) -> None:  # Print beautiful matrix
    n = len(A)
    m = len(A[0])
    A = [[simplify_number(A[i][j], accuracy, sign_space=any([A[i][j] < 0 for i in range(n)])) for j in range(m)]
         for i in range(n)]
    len_of_element = [max([len(A[i][j]) for i in range(n)]) for j in range(m)]
    for i in range(n):
        print('⌈' if i == 0 else '⌊' if i + 1 == n else '|', end='')
        for j in range(m):
            print(A[i][j] + ' ' * (len_of_element[j] - len(A[i][j])), end='')
            if j + 1 != m:
                print('  ', end='')
        print('⌉' if i == 0 else '⌋' if i + 1 == n else '|')


def print_vector(m: list[float], ends='\n') -> None:
    print('x = (' + '; '.join(map(simplify_number, m)) + ')', end=ends)


def is_upper_triangular(m: list[list[float]]) -> bool:
    for i in range(len(m)):
        for j in range(i):
            if m[i][j] != 0:
                return False
    return True


def to_upper_triangular(A: list[list[float]]) -> list[list[float]]:  # Straight running method Gauss
    n = len(A)
    m = len(A[0])
    A = [A[i][:] for i in range(n)]
    for k in range(n - 1):
        ind_max = k + A[k:].index(max(A[k:], key=lambda x: x[k]))  # Swap rows with max element and first not zero
        c = A[ind_max][:]
        A[ind_max] = A[k]
        A[k] = c
        if A[ind_max] != A[k]:
            for j in range(m):
                A[ind_max][j] *= -1
        t = A[k][k]
        for i in range(k + 1, n):
            for j in range(m - 1, k - 1, -1):
                A[i][j] -= A[k][j] * A[i][k] / t
    return A


def det(m: list[list[float]]) -> str:  # Return determinant of matrix m
    if is_upper_triangular(m):
        res = 1
        for i in range(len(m)):
            res *= m[i][i]
        return simplify_number(res)
    return det(to_upper_triangular(m))


def inverse_matrix(m: list[list[float]]) -> list[list[float]]:
    n = len(m)
    m = [[m[i][j] if j < n else 1 if i + n == j else 0 for j in range(2 * n)] for i in range(n)]  # m = [A|E]
    for k in range(n):
        ind_max = k + m[k:].index(max(m[k:], key=lambda x: x[k]))  # Swap rows with max element and first not zero
        c = m[ind_max][:]
        m[ind_max] = m[k]
        m[k] = c
        t = m[k][k]
        for j in range(2 * n):
            m[k][j] /= t
        for i in range(n):
            if i != k:
                for j in range(2 * n - 1, k - 1, -1):
                    m[i][j] -= m[k][j] * m[i][k]
    return [line[n:] for line in m]


def check_correct_ax_b(m: list[list[float]], b: list[float]) -> None:
    n = len(m)
    if n != len(b) or any([len(line) != n for line in m]):  # Incorrect size of A or b
        raise ValueError('Incorrect size')
    if det(m) == 0:  # Incorrect type of matrix
        raise ValueError('This matrix is degenerate')


def the_gauss_method(A: list[list[float]], b: list[float]) -> list[float]:  # Gauss method with column maximum selection
    check_correct_ax_b(A, b)
    n = len(A)
    A = [[A[i][j] if j < n else b[i] for j in range(n + 1)] for i in range(n)]  # m = [A|b]
    A = to_upper_triangular(A)  # Straight of gauss method
    # Calculating x
    x = [0.0 for _ in range(n)]
    for i in range(n - 1, -1, -1):
        x[i] = (A[i][n] - sum([A[i][j] * x[j] for j in range(i + 1, n)])) / A[i][i]
    return x


def run_throw_method(m: list[list[float]], b: list[float]) -> list[float]:  # A - triple diagonal matrix
    n = len(m)
    m = [[m[i][j] if j < n else b[i] for j in range(n + 1)] for i in range(n)]  # m = [A|b]
    for k in range(n - 1):
        m[k][n] /= m[k][k]
        for j in range(min(n - 1, k + 1), k - 1, -1):
            m[k][j] /= m[k][k]
        m[k + 1][n] -= m[k][n] * m[k + 1][k]
        for j in range(min(n - 1, k + 1), k - 1, -1):
            m[k + 1][j] -= m[k][j] * m[k + 1][k]
    for j in range(n, n - 3, -1):
        m[n - 1][j] /= m[n - 1][n - 1]
    # Calculating x
    x = [0.0 for _ in range(n)]
    x[n - 1] = m[n - 1][n]
    for i in range(n - 2, -1, -1):
        x[i] = (m[i][n] - m[i][i + 1] * x[i + 1]) / m[i][i]
    return x


def norm_vector(v: list[float]) -> float:  # The norm 2
    return sum(map(lambda q: q * q, v)) ** 0.5


def reflection_method(m: list[list[float]], b: list[float]) -> list[float]:
    def matrices_mul(q: list[list[float]], p: list[list[float]]):
        return [[sum([q[i][z] * p[z][j] for z in range(len(p))]) for j in range(len(p[0]))] for i in range(len(q))]

    def vvt_multiplication(q: list[float], p: list[float]) -> list[list[float]]:
        if len(q) != len(p):
            raise ValueError('Incorrect size of vector by v*vT multiplication')
        return [[q[i] * p[j] for j in range(len(p))] for i in range(len(q))]

    def sign(q: float) -> int:
        return 1 if q > 0 else 0 if q == 0 else -1

    def matrices_div(q: list[list[float]], p: list[list[float]]) -> list[list[float]]:
        return [[q[i][j] - p[i][j] for j in range(len(q))] for i in range(len(q))]

    check_correct_ax_b(m, b)
    n = len(m)
    m = [[m[i][j] if j < n else b[i] for j in range(n + 1)] for i in range(n)]  # m = [A|b]
    for k in range(n - 1):
        Enk = [[1 if i == j else 0 for j in range(n - k)] for i in range(n - k)]
        rk = norm_vector([m[i][k] for i in range(k, n)])
        if rk != 0:
            t = [m[i][k] for i in range(k, n)]
            t[0] -= rk * sign(m[k][k])
            norm_t = norm_vector(t)
            w = [t[i] / norm_t for i in range(n - k)]
            W = vvt_multiplication([2 * w[j] for j in range(len(w))], w)
            Ank = matrices_div(Enk, W)
            Hk = [[1 if i == j and i < k else Ank[i - k][j - k] if i >= k and j >= k else 0 for j in range(n)]
                  for i in range(n)]
            m = matrices_mul(Hk, m)
    # Calculating x
    x = [0.0 for _ in range(n)]
    for i in range(n - 1, -1, -1):
        x[i] = (m[i][n] - sum([m[i][j] * x[j] for j in range(i + 1, n)])) / m[i][i]
    return x


def matrix_vector_mul(m: list[list[float]], v: list[float]) -> list[float]:
    return [sum([m[k][p] * v[p] for p in range(len(v))]) for k in range(len(v))]


def vectors_sub(v: list[float], u: list[float]) -> list[float]:
    return [v[k] - u[k] for k in range(len(v))]


def num_vector_mul(p: float, v: list[float]) -> list[float]:
    return [p * v[k] for k in range(len(v))]


def vectors_add(v: list[float], u: list[float]) -> list[float]:
    return [v[k] + u[k] for k in range(len(v))]


def richardson_with_chebyshev_params_method(
        A: list[list[float]], b: list[float], l_min: float, l_max: float, eps=10 ** -5) -> (int, list[float]):
    def mistake(m: list[list[float]], v: list[float], current_x: list[float]) -> list[float]:
        return vectors_sub(matrix_vector_mul(m, current_x), v)

    n = len(b)
    x0 = [0 for _ in range(n)]
    xk = x0
    t_opt = 2 / (l_min + l_max)
    xkp1 = vectors_sub(x0, num_vector_mul(t_opt, vectors_sub(matrix_vector_mul(A, x0), b)))
    tk = 1
    t1 = (l_min + l_max) / (l_min - l_max)
    tkp1 = t1
    number_iteration = 1
    while norm_vector(mistake(A, b, xkp1)) > eps:
        number_iteration += 1
        tkm1 = tk
        tk = tkp1
        tkp1 = 2 * t1 * tk - tkm1
        xkm1 = xk
        xk = xkp1
        xkp1 = vectors_sub(vectors_add(xk, num_vector_mul(tkm1 / tkp1, vectors_sub(xk, xkm1))),
                           num_vector_mul(t_opt * (1 + tkm1 / tkp1), vectors_sub(matrix_vector_mul(A, xk), b)))
    return number_iteration, xkp1


def print_iterative(pair: (int, list[float])) -> None:
    print('The number of last iteration:', pair[0], end=' | ')
    print_vector(pair[1])


def vtv_mul(u: list[float], v: list[float]) -> float:
    return sum([u[i] * v[i] for i in range(len(u))])


def conjugate_gradient_method(A: list[list[float]], b: list[float], x_prev=None, eps=10 ** -5) -> (int, list[float]):
    n = len(A)
    if x_prev is None:
        x_prev = [0 for _ in range(n)]
    r = vectors_sub(matrix_vector_mul(A, x_prev), b)
    g = r
    alpha = vtv_mul(r, g) / vtv_mul(matrix_vector_mul(A, g), g)
    x = vectors_sub(x_prev, num_vector_mul(alpha, g))
    r = vectors_sub(matrix_vector_mul(A, x), b)
    num_iteration = 1
    while norm_vector(r) > eps:
        num_iteration += 1
        gamma = vtv_mul(matrix_vector_mul(A, r), g) / vtv_mul((matrix_vector_mul(A, g)), g)  # (A * r, g) / (A * g, g)
        g = vectors_sub(r, num_vector_mul(gamma, g))  # g = r - gamma * g
        alpha = vtv_mul(r, g) / (vtv_mul(matrix_vector_mul(A, g), g))
        x = vectors_sub(x, num_vector_mul(alpha, g))
        r = vectors_sub(matrix_vector_mul(A, x), b)
    return num_iteration, x


''' # Three-layer formulas
def conjugate_gradient_method(A: list[list[float]], b: list[float], x_prev=None) -> (int, list[float]):
    n = len(A)
    if x_prev is None:
        x_prev = [0 for _ in range(n)]
    r = vectors_sub(matrix_vector_mul(A, x_prev), b)  # r0 = A * x0 - b
    alpha = vtv_mul(r, r) / vtv_mul(matrix_vector_mul(A, r), r)
    x = vectors_sub(x_prev, num_vector_mul(alpha, r))
    r_prev = r
    r = vectors_sub(matrix_vector_mul(A, x), b)
    for i in range(n):
        v = the_gauss_method(
            [[vtv_mul(vectors_sub(r, r_prev), vectors_sub(x, x_prev)), vtv_mul(vectors_sub(r, r_prev), r)],
             [vtv_mul(vectors_sub(r, r_prev), r), vtv_mul(matrix_vector_mul(A, r), r)]],
            [vtv_mul(r, vectors_sub(x, x_prev)), vtv_mul(r, r)])
        x_new = vectors_sub(vectors_sub(x, num_vector_mul(v[0], vectors_sub(x, x_prev))), num_vector_mul(v[1], r))
        x_prev = x
        x = x_new
        r_prev = r
        r = vectors_sub(matrix_vector_mul(A, x), b)
    return x'''


def power_method(A: list[list[float]], eps=10 ** -5) -> (int, float):
    n = len(A)
    x_new = [1 for _ in range(n)]
    num_iteration = 1
    # for num_iteration in range(100):
    while norm_vector(vectors_sub(matrix_vector_mul(A, x_new), num_vector_mul(norm_vector(x_new), x_new))) > eps:
        num_iteration += 1
        x = x_new
        x_new = num_vector_mul(1 / norm_vector(x), matrix_vector_mul(A, x))
    print_vector(num_vector_mul(1 / norm_vector(x_new), x_new))
    return num_iteration, norm_vector(x_new)


def norm_1_matrix(A: list[list[float]]) -> float:
    n = len(A)
    return max([sum([A[i][j] for i in range(n)]) for j in range(n)])


def matrices_sub(A: list[list[float]], B: list[list[float]]) -> list[list[float]]:
    n = len(A)
    return [[A[i][j] - B[i][j] for j in range(n)] for i in range(n)]


def create_identity(n: int) -> list[list[float]]:
    return [[1 if i == j else 0 for j in range(n)] for i in range(n)]


def num_matrices_mul(k: float, A: list[list[float]]) -> list[list[float]]:
    n = len(A)
    return [[k * A[i][j] for j in range(n)] for i in range(n)]


def print_power_method(A: list[list[float]], eps: float) -> None:
    n = len(A)
    v = power_method(A, eps)
    print('The number of last iteration:', v[0], '| lambda_max =', simplify_number(v[1]))
    nA = norm_1_matrix(A)
    v = power_method(matrices_sub(num_matrices_mul(nA, create_identity(n)), A), eps)
    print('The number of last iteration:', v[0], '| lambda_min =', simplify_number(nA - v[1]))


def main() -> None:
    try:
        with open('input.txt', 'r') as f:
            n = int(f.readline())  # Input size of matrix nxn
            A = [list(map(float, f.readline().split())) for _ in range(n)]  # Input matrix
            # b = list(map(float, f.readline().split()))  # Input vector b
            '''# l_min, l_max = map(float, f.readline().split()) # for richardson`s method'''
            eps = f.readline()
            if eps == '\n':
                eps = 10 ** -5
            else:
                eps = 10 ** int(eps)

            # print('Gauss`s method result: ', end='') print_vector(the_gauss_method(A, b))
            # print_matrix(inverse_matrix(A))
            # print(det(A))
            # print_vector(run_throw_method(A, b))
            # print_vector(reflection_method(A, b))
            # print_iterative(richardson_with_chebyshev_params_method(A, b, l_min, l_max, eps))
            # print('Conjugate gradient:    ', end="") print_iterative(conjugate_gradient_method(A, b))
            print_power_method(A, eps)
    except Exception as e:
        print(str(e.__class__)[8:-2] + ': ' + str(e))


if __name__ == '__main__':
    main()
