# Вычисление базиса Грёбнера и решение системы

Программа позволяет вычислить базис Грёбнера для заданной системы уравнений и решить эту систему.

## Формат входных данных
На вход программе подается файл, который имеет следующий формат:
```
N
[# x1 > x2 > x3 > ... > xN]
P1(x1, ..., xN)
P2(x1, ..., xN)
...
Pk(x1, ..., xN)
```
* На первой строке располагается число `N` - это размерность системы.
* На второй строке опционально расмещается измененный лексикографический порядок системы, он обязательно должен начинаться с символа `#` (даже пробел не должен быть первым символом!). При задании порядка нужно перечислить все N переменных системы в желаемом порядке. *Если не вписывать эту строку, будет использоваться обычный лексикографический порядок*.
* Далее на k строках перечисляются k многочленов. Многочлен задается как сумма одночленов. Все переменные обязательно должны иметь имя `x`. Например: `x1 + x2^3 * x3^16 + x4*x9^3`. Индексы переменных и степени должны быть положительны. При задании многочленов можно как угодно использовать пробелы. Перед парсингом они удаляются.

Пример входного файла `input.txt`:
```
3
# x2 > x3 > x1
x3 + x1*x2
x2 + x1*x3
```

## Запуск программы
Чтобы запустить программу для конкретного входного файла, нужно из директории проекта выполнить следующее:
```
$ go run main.go input.txt
```
* `input.txt` - это входной файл, формат которого был описан в предыдущем разделе.
* В `stdout` выводится информация о размерности системы, оригинальной системе, приведенной системе и их равенстве.
* В `stderr` выводится прочая вспомогательная информация о выполнении программы.

Чтобы запустить скрипт, который генерирует системы и запускает решатель, нужно выполнить следующее:
```
$ python3 checker.py
```
В файл `output.txt` будет выведена информация о решении каждой из сгенерированных систем. Чтобы проверить, есть ли системы, у которых оригинальный и приведенный вариант имеют разные решения, нужно выполнить:
```
$ grep "false" output.txt
```

## Описание кода
Проект состоит из двух составных частей:
* Парсер (`parser/parser.go`). Парсер разбирает поданный на вход файл и создает структуры данных для одночленов и многочленов.
* Решатель (`f2/f2.go`). Решатель вычисляет базис Гребнера и минимизирует его.

### Парсер
Рассмотрим, что делает каждая функция в парсере.
* `func lookahead(line *string, index int) byte` - возвращает символ, следующий за символом с индексом `idx` в строке `line`. Если строка исчерпана, возвращает 0.
* `func isEndOfMonomial(char byte) bool` - проверяет, является ли символ символом конца одночлена (стоит ли после одночлена `+` или кончилась ли строка).
* `func readNumber(line *string, index int) (number, lastIndex int, err error)` - прочитать число. `index` - это индекс символа, находящегося перед первой цифрой текущего числа в строке `line`. Функция возвращает считанное число `number`, индекс последней цифры числа `lastIndex` и ошибку считывания `err`.
* `func ReadMonomial(line *string, index int) (monom f2.Monomial, lastIndex int, err error)` - считывает новый одночлен из строки `line`, где `index` - это индекс символа, находящегося перед одночленом. Возвращает считанный одночлен, индекс последнего символа одночлена `lastIndex` и ошибку считывания `err`.
* `func ParseOrder(line *string, n int) (monom []int, err error)` - разбирает во входном файле строку `line`, в которой задан порядок переменных для системы размерности `n`. Возвращает массив положительных чисел `monom` - индексов переменных в правильном порядке - и ошибку считывания `err`. Например, для входного файла из примера выше массив будет такой: `[2 3 1]`.
* `func ParsePolynomial(line *string) (f2.Polynomial, error)` - разбирает строку `line` и превращает ее в полином. Возвращает считанный полином и ошибку считывания `err`.

### Решатель
Парсер генерирует для решателя следующие объекты:
* полиномы, состоящие из мономов;
* порядок переменных в виде массива целых чисел - индексов переменных.

**Мономы** (`type Monomial map[int]int`) представлены в виде карты (map) - это хеш-таблица в Go. map можно рассматривать как множество пар *ключ-значение*. Мономы удобно представлять в таком виде, задавая отображение индекса переменной на ее степень. Например, моном `x2^3 * x4` будет представлен в виде следующей карты: `[ 2 -> 3 , 4 -> 1 ]`. *Чтобы задать константу 1, используется пустая карта `[]`*.

**Полиномы** (`type Polynomial []Monomial`) представлены в виде массивов мономов. *Если полином - это константа 0, то задается пустой массив*.

**Базис** (`type Basis []Polynomial`) представлен в виде массива полиномов.

**Система** (`type System struct`) представлена в виде структуры, имеющей следующие поля:
* `N int` - размерность;
* `Polynomials []Polynomial` - массив полиномов.

Решатель имеет следующие функции:
* `func SetOrder(orderlist []int, n int) (err error)` - устанавливает порядок переменных, который хранится в переменной `f2.Order`, используя массив `orderlist`. Возвращает ошибку `err`. Возможны две ошибки: не хватает какой-то переменной или какая-то переменная задана дважды.
* `func (m Monomial) String() string` - функция, преобразующая одночлен в строку.
* `func (p Polynomial) String() string` - функция, преобразующая многочлен в строку.
* `func (sys System) String() string` - функция, преобразующая систему в строку.
* `func MultMono(m1, m2 *Monomial) (m3 Monomial)` - умножает одночлены `m1` и `m2` и возвращает результат - новый одночлен.
* `func MultPoly(p1, p2 *Polynomial) (p3 Polynomial)` - умножает многочлены `p1` и `p2` и возвращает результат - новый многочлен.
* `func MultMonoPoly(m *Monomial, p *Polynomial) Polynomial` - умножает одночлен `m` на многочлен `p` и возвращает результат - новый многочлен.
* `func Simplify(src *Polynomial) (dst Polynomial)` - упрощает многочлен. Если в составе многочлена есть такие одночлены, которые сами с собой суммируются четное количество раз, то такие одночлены удаляются. Которые суммируются нечетное количество раз - превращается в один многочлен.
* `func AddPoly(p1, p2 *Polynomial) (p3 Polynomial)` - суммирует полиномы, возвращает новый полином.
* `func SubPoly(p1, p2 *Polynomial) (p3 Polynomial)` - вычитает полиномы. Делает то же самое, что `AddPoly()`.
* `func CompareByOrder(idx1, idx2 int) int` - сравнивает переменные xi и xj, то есть вычисляет разность между индексами индексов i (представлен как idx1) и j (представлен как idx2) в массиве `Order`. Например, для нашего примера вызовем `CompareByOrder(2, 1)`. Разность между переменными x2 и x1 в массиве `Order = [2 3 1]` будет равна 2, то есть x2 > x1.
* `func CompareMono(m1, m2 *Monomial) int` - сравнивает одночлены `m1`, `m2`. Возвращает:
	* 1, если `m1 > m2`;
	* -1, если `m1 < m2`;
	* 0, если `m1 == m2`.
* `func (p *Polynomial) GetTopMonomial() (topMonomial *Monomial)` - возвращает одночлен, который лексикографически больше всех остальных одночленов в многочлене `p` ("старший" одночлен).
* `func (p *Polynomial) DiscardTopMonomial() (pNew *Polynomial)` - выкинуть старший одночлен из многочлена `p`, вернуть новый многочлен `pNew`.
* `func (m1 *Monomial) Divide(m2 *Monomial) *Monomial` - разделить одночлен `m1` на `m2`. Если делится, то вернуть частное. Если нет, то вернуть `nil`.
* `func (m1 *Monomial) Gcd(m2 *Monomial) Monomial` - найти `НОД(m1, m2)` двух одночленов.
* `func (m *Monomial) IsConstant() bool` - проверить, является ли многочлен `m` константой 1.
* `func (h *Polynomial) Reduce(f []Polynomial) *Polynomial` - провести редукцию многочлена `h` по набору многочленов `f`, вернуть получившийся многочлен. Если `h` не редуцируется, то вернуть `nil`.
* `func (p1 *Polynomial) HasCommonChain(p2 *Polynomial) bool` - проверить, имеют ли многочлены `p1` и `p2` зацепление.
* `func GetGroebnerBasis(ideal []Polynomial) (basis Basis)` - построить базис Гребнера по заданному идеалу (массиву полиномов).
* `func (b *Basis) Minimize()` - провести минимизацию + редукцию для базиса `b`.
* `func (s *System) Solve() [][]int` - решить систему уравнений. Функция перебирает все возможные векторы длины N и подставляет их в систему. Если вектор подходит всем уравнениям в системе, то он записывается в выходной массив. Функция возвращает массив, состоящий из массивов из 0 и 1 - решений системы.

### Разбор файла main.go
Файл `main.go` содержит код, вызывающий парсер и решатель.
* Для начала вручную парсится N - размер системы. 
* Затем парсер вызывается для строки, задающей порядок переменных. На этом этапе обязательно надо вызвать функцию `f2.SetOrder()`, чтобы установить считанный порядок.
* Далее построчно считываются строки полиномов и вызывается парсер. Парсер возвращает разобранный полином, полиномы добавляются к системе.
* Для системы уравнений ищется базис Гребнера.
* Базис минимизируется + редуцируется.
* Затем ищем решения начальной системы и приведенной системы.
* Результат сравнения решений лежит  в `areEqual`.

## Заключительные замечания к коду

* Если `err == nil`, это означает, что ошибки в ходе работы функции не случилось.
* Чтобы напечатать что-то в `stderr`, надо вызывать `log.Println()`, `log.Printf()` и прочие функции из `log`.
* Чтобы напечатать что-то в `stdout`, нужно вызвать `fmt.Println()`, `fmt.Printf()` и прочие функции из `fmt`.
