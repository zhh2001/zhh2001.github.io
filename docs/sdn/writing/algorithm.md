---
outline: deep
---

# LaTeX 算法

## 1 基础设置

首先，在导言区加载必要的宏包：

```latex
\usepackage{algorithm}    % 提供算法浮动体
\usepackage{algorithmic}  % 提供算法伪代码语法（或用 algorithm2e）
```

## 2 使用 `algorithmic` 环境 ​

### 2.1 基础语法

```latex
\begin{algorithm}[!t]
    \caption{Calculate $y = x^n$}  % 算法标题
    \label{algo1}                  % 标签（用于引用）
    \begin{algorithmic}[1]
    \Require $n \geq 0 \vee x \neq 0$
    \Ensure $y = x^n$
    \State $y \Leftarrow 1$
    \If{$n < 0$}
        \State $X \Leftarrow 1 / x$
        \State $N \Leftarrow -n$
    \Else
        \State $X \Leftarrow x$
        \State $N \Leftarrow n$
    \EndIf
    \While{$N \neq 0$}
        \If{$N$ is even}
            \State $X \Leftarrow X \times X$
            \State $N \Leftarrow N / 2$
        \Else[$N$ is odd]
            \State $y \Leftarrow y \times X$
            \State $N \Leftarrow N - 1$
        \EndIf
    \EndWhile
    \end{algorithmic}
\end{algorithm}
```

### 2.2 常用命令 ​

| 命令             | 功能                       |
| ---------------- | -------------------------- |
| `\State`         | 普通语句                   |
| `\If{条件}`      | 条件判断                   |
| `\Else`          | 否则分支                   |
| `\ElsIf`         | 另一个条件（可选）         |
| `\EndIf`         | 结束条件判断               |
| `\For{循环条件}` | `for` 循环                 |
| `\EndFor`        | 结束 `for` 循环            |
| `\While{条件}`   | `while` 循环               |
| `\EndWhile`      | 结束 `while` 循环          |
| `\Repeat`        | 直到循环（先执行后判断）   |
| `\Until{条件}`   | 结束 `repeat` 循环         |
| `\Return`        | 返回值                     |
| `\Comment{注释}` | 添加注释（需额外宏包支持） |

## 3 完整示例

```latex
\documentclass{article}
\usepackage{algorithm}
\usepackage{algorithmic}

\begin{document}

\begin{algorithm}[!t]
    \caption{Calculate $y = x^n$}
    \label{algo1}
    \begin{algorithmic}[1]
    \Require $n \geq 0 \vee x \neq 0$
    \Ensure $y = x^n$
    \State $y \Leftarrow 1$
    \If{$n < 0$}
        \State $X \Leftarrow 1 / x$
        \State $N \Leftarrow -n$
    \Else
        \State $X \Leftarrow x$
        \State $N \Leftarrow n$
    \EndIf
    \While{$N \neq 0$}
        \If{$N$ is even}
            \State $X \Leftarrow X \times X$
            \State $N \Leftarrow N / 2$
        \Else[$N$ is odd]
            \State $y \Leftarrow y \times X$
            \State $N \Leftarrow N - 1$
        \EndIf
    \EndWhile
    \end{algorithmic}
\end{algorithm}

引用算法：见算法~\ref{alg:gcd}。

\end{document}
```

## 4 高级功能

### 4.1 添加注释​

若需在代码行旁添加注释，需加载 `algorithmicx` 和 `algpseudocode` 宏包（更强大）：

```latex
\usepackage{algorithm}
\usepackage[noend]{algpseudocode}  % noend 可隐藏结束标记

\begin{algorithmic}[1]
    \State $x \gets 1$ \Comment{初始化 x}
\end{algorithmic}
```

### 4.2 自定义操作

用 `\Statex` 或 `\algorithmicindent` 调整缩进，或定义新命令：

```latex
\algnewcommand{\LineComment}[1]{\State $\triangleright$ #1}
\begin{algorithmic}[1]
    \State $x \gets 1$ \LineComment{初始化}
\end{algorithmic}
```

### 4.3 使用 `algorithm2e` 宏包​

`algorithm2e` 提供更灵活的排版（需调整参数）：

```latex
\usepackage[ruled,vlined]{algorithm2e}  % ruled 显示边框，vlined 显示竖线

\begin{algorithm}[H]
    \SetAlgoLined
    \KwIn{输入参数}
    \KwOut{输出结果}
    \If{条件}{
        执行操作\;
    }
    \Return{结果}
\end{algorithm}
```

## 5 常见问题​

### 5.1 算法不浮动

使用 `[H]` 强制定位（需 `float` 宏包），但建议优先让 LaTeX 自动调整位置。

### 5.2 行号不显示​

检查 `\begin{algorithmic}[1]` 中的 `[1]` 是否存在。

### 5.3 ​​中文支持

若需中文，加载 `ctex` 宏包，并确保字体支持：

```latex
\usepackage[UTF8]{ctex}
```
