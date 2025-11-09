---
outline: deep
---

# LaTeX 公式

## 1 数学环境

### 1.1 行内公式

使用 `$...$` 包裹公式。

示例：

```latex
$x = y + 1$
```

效果：

$x = y + 1$

### 1.2 块公式

使用 `$$...$$` 包裹公式。

示例：

```latex
$$
x = y + 1
$$
```

效果：

$$
x = y + 1
$$

## 2 数学符号

### 2.1 希腊字母

| 写法       | 效果       |     | 写法       | 效果       |
| ---------- | ---------- | --- | ---------- | ---------- |
| `\alpha`   | $\alpha$   |
| `\beta`    | $\beta$    |
| `\gamma`   | $\gamma$   |     | `\Gamma`   | $\Gamma$   |
| `\delta`   | $\delta$   |     | `\Delta`   | $\Delta$   |
| `\epsilon` | $\epsilon$ |
| `\zeta`    | $\zeta$    |
| `\eta`     | $\eta$     |
| `\theta`   | $\theta$   |     | `\Theta`   | $\Theta$   |
| `\iota`    | $\iota$    |
| `\kappa`   | $\kappa$   |
| `\lambda`  | $\lambda$  |     | `\Lambda`  | $\Lambda$  |
| `\mu`      | $\mu$      |
| `\nu`      | $\nu$      |
| `\xi`      | $\xi$      |     | `\Xi`      | $\Xi$      |
| `\pi`      | $\pi$      |     | `\Pi`      | $\Pi$      |
| `\rho`     | $\rho$     |
| `\sigma`   | $\sigma$   |     | `\Sigma`   | $\Sigma$   |
| `\tau`     | $\tau$     |
| `\upsilon` | $\upsilon$ |     | `\Upsilon` | $\Upsilon$ |
| `\phi`     | $\phi$     |     | `\Phi`     | $\Phi$     |
| `\chi`     | $\chi$     |
| `\psi`     | $\psi$     |     | `\Psi`     | $\Psi$     |
| `\omega`   | $\omega$   |     | `\Omega`   | $\Omega$   |

### 2.2 上下标

`_` 表示下标，`^` 表示上标。

| 示例      | 效果      |
| --------- | --------- |
| `x_1`     | $x_1$     |
| `x_{ij}`  | $x_{ij}$  |
| `x^2`     | $x^2$     |
| `x^{n+1}` | $x^{n+1}$ |
| `x_n^2`   | $x_n^2$   |

### 2.3 分数

示例：`\frac{a}{b}`

效果：$\frac{a}{b}$

### 2.4 根号

- 平方根：`\sqrt{x}`  
  效果：$\sqrt{x}$

- n 次根：`\sqrt[n]{x}`  
  效果：$\sqrt[n]{x}$

### 2.5 求和与积分

| 示例                      | 效果                      |
| ------------------------- | ------------------------- |
| `\sum_{i=1}^n a_i`        | $\sum_{i=1}^n a_i$        |
| `\prod_{i=1}^n a_i`       | $\prod_{i=1}^n a_i$       |
| `\int_a^b f(x)dx`         | $\int_a^b f(x)dx$         |
| `\iint_D f(x,y)dxdy`      | $\iint_D f(x,y)dxdy$      |
| `\iiint_V f(x,y,z)dxdydz` | $\iiint_V f(x,y,z)dxdydz$ |
| `\oint_C f(z)dz`          | $\oint_C f(z)dz$          |

### 2.6 极限

示例：`\lim_{x\to\infty} f(x)`

效果：$\lim_{x\to\infty} f(x)$

### 2.7 微分与偏导

| 示例                            | 效果                            |
| ------------------------------- | ------------------------------- |
| `\frac{dy}{dx}`                 | $\frac{dy}{dx}$                 |
| `\frac{\partial u}{\partial t}` | $\frac{\partial u}{\partial t}$ |
| `\nabla f`                      | $\nabla f$                      |

## 3 矩阵与数组

### 3.1 简单矩阵

示例：

```latex
$$
\begin{matrix}
a & b \\
c & d
\end{matrix}
$$
```

效果：

$$
\begin{matrix}
a & b \\
c & d
\end{matrix}
$$

### 3.2 带括号矩阵

示例：

```latex
$$
\begin{pmatrix}  % 圆括号
a & b \\
c & d
\end{pmatrix}
$$

$$
\begin{bmatrix}  % 方括号
a & b \\
c & d
\end{bmatrix}
$$

$$
\begin{vmatrix}  % 单竖线
a & b \\
c & d
\end{vmatrix}
$$

$$
\begin{Vmatrix}  % 双竖线
a & b \\
c & d
\end{Vmatrix}
$$
```

效果：

$$
\begin{pmatrix}  % 圆括号
a & b \\
c & d
\end{pmatrix}
$$

$$
\begin{bmatrix}  % 方括号
a & b \\
c & d
\end{bmatrix}
$$

$$
\begin{vmatrix}  % 单竖线
a & b \\
c & d
\end{vmatrix}
$$

$$
\begin{Vmatrix}  % 双竖线
a & b \\
c & d
\end{Vmatrix}
$$

### 3.3 大矩阵（跨多行）

示例：

```latex
$$
\begin{pmatrix}
a_{11} & a_{12} & \cdots & a_{1n} \\
a_{21} & a_{22} & \cdots & a_{2n} \\
\vdots & \vdots & \ddots & \vdots \\
a_{m1} & a_{m2} & \cdots & a_{mn}
\end{pmatrix}
$$
```

效果：

$$
\begin{pmatrix}
a_{11} & a_{12} & \cdots & a_{1n} \\
a_{21} & a_{22} & \cdots & a_{2n} \\
\vdots & \vdots & \ddots & \vdots \\
a_{m1} & a_{m2} & \cdots & a_{mn}
\end{pmatrix}
$$

### 3.4 数组环境

示例：

```latex
$$
\begin{array}{ccc}
a & b & c \\
d & e & f \\
g & h & i
\end{array}
$$
```

效果：

$$
\begin{array}{ccc}
a & b & c \\
d & e & f \\
g & h & i
\end{array}
$$

列对齐选项：`l`(左对齐)、`c`(居中)、`r`(右对齐)

## 4 分段函数

示例：

```latex
$$
f(x) =
\begin{cases}
x^2 & \text{if } x \geq 0 \\
-x & \text{if } x < 0
\end{cases}
$$
```

效果：

$$
f(x) =
\begin{cases}
x^2 & \text{if } x \geq 0 \\
-x & \text{if } x < 0
\end{cases}
$$

## 5 多行公式

示例：

```latex
$$
\begin{align}
a^2 + b^2 &= c^2 \\
e^{i\pi} + 1 &= 0
\end{align}
$$
```

效果：

$$
\begin{align}
a^2 + b^2 &= c^2 \\
e^{i\pi} + 1 &= 0
\end{align}
$$

## 6 特殊符号

### 6.1 集合符号

| 写法         | 效果         |
| ------------ | ------------ |
| `\in`        | $\in$        |
| `\subset`    | $\subset$    |
| `\subseteq`  | $\subseteq$  |
| `\supset`    | $\supset$    |
| `\supseteq`  | $\supseteq$  |
| `\cup`       | $\cup$       |
| `\cap`       | $\cap$       |
| `\setminus`  | $\setminus$  |
| `\emptyset`  | $\emptyset$  |
| `\mathbb{N}` | $\mathbb{N}$ |
| `\mathbb{Z}` | $\mathbb{Z}$ |
| `\mathbb{Q}` | $\mathbb{Q}$ |
| `\mathbb{R}` | $\mathbb{R}$ |
| `\mathbb{C}` | $\mathbb{C}$ |

### 6.2 逻辑符号

| 写法              | 效果              |
| ----------------- | ----------------- |
| `\land`           | $\land$           |
| `\lor`            | $\lor$            |
| `\lnot`           | $\lnot$           |
| `\Rightarrow`     | $\Rightarrow$     |
| `\Leftrightarrow` | $\Leftrightarrow$ |
| `\forall`         | $\forall$         |
| `\exists`         | $\exists$         |

### 6.3 关系符号

| 写法      | 效果      |
| --------- | --------- |
| `\lt`     | $\lt$     |
| `\gt`     | $\gt$     |
| `\leq`    | $\leq$    |
| `\geq`    | $\geq$    |
| `\ll`     | $\ll$     |
| `\gg`     | $\gg$     |
| `\approx` | $\approx$ |
| `\equiv`  | $\equiv$  |
| `\sim`    | $\sim$    |
| `\simeq`  | $\simeq$  |
| `\propto` | $\propto$ |

### 6.4 箭头符号

| 写法              | 效果              |     | 写法              | 效果              |
| ----------------- | ----------------- | --- | ----------------- | ----------------- |
| `\rightarrow`     | $\rightarrow$     |     | `\Rightarrow`     | $\Rightarrow$     |
| `\leftarrow`      | $\leftarrow$      |     | `\Leftarrow`      | $\Leftarrow$      |
| `\leftrightarrow` | $\leftrightarrow$ |     | `\Leftrightarrow` | $\Leftrightarrow$ |
| `\mapsto`         | $\mapsto$         |
| `\hookrightarrow` | $\hookrightarrow$ |

### 6.5 特殊运算符

| 写法       | 效果       |
| ---------- | ---------- |
| `\pm`      | $\pm$      |
| `\mp`      | $\mp$      |
| `\times`   | $\times$   |
| `\div`     | $\div$     |
| `\cdot`    | $\cdot$    |
| `\ast`     | $\ast$     |
| `\star`    | $\star$    |
| `\circ`    | $\circ$    |
| `\bullet`  | $\bullet$  |
| `\diamond` | $\diamond$ |
| `\oplus`   | $\oplus$   |
| `\otimes`  | $\otimes$  |
| `\odot`    | $\odot$    |

## 7 高级技巧

### 7.1 自定义命令

在导言区定义常用命令：

```latex
\newcommand{\abs}[1]{\left| #1 \right|}
\newcommand{\norm}[1]{\left\| #1 \right\|}
\newcommand{\pfrac}[2]{\left( \frac{#1}{#2} \right)}
```

使用时：

```latex
\abs{x}, \norm{v}, \pfrac{a}{b}
```

### 7.2 空间控制

```latex
\,     % 小空格
\:     % 中空格
\;     % 大空格
\!     % 负空格
\quad  % 较大空格
\qquad % 很大空格
```

### 7.3 文本插入

| 示例                  | 效果                  |
| --------------------- | --------------------- |
| `\text{some text}`    | $\text{some text}$    |
| `\mathrm{text}`       | $\mathrm{text}$       |
| `\mathbf{bold}`       | $\mathbf{bold}$       |
| `\mathit{italic}`     | $\mathit{italic}$     |
| `\mathsf{sans-serif}` | $\mathsf{sans-serif}$ |
| `\mathtt{typewriter}` | $\mathtt{typewriter}$ |

## 8 常见问题

### 8.1 分数显示过小

使用 `\dfrac` 代替 `\frac`：

```latex
\dfrac{a}{b}  % 显示更大的分数
```

效果：

$$
\dfrac{a}{b}
$$

### 8.2 矩阵元素对齐

使用 `array` 环境并指定对齐方式：

```latex
\begin{array}{rcl}
left & center & right \\
a & b & c
\end{array}
```

效果：

$$
\begin{array}{rcl}
left & center & right \\
a & b & c
\end{array}
$$

### 8.3 长公式换行

使用 `multline` 环境：

```latex
\begin{multline}
a + b + c + d + e + f \\
+ g + h + i + j + k + l
\end{multline}
```
