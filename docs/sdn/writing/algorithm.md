---
outline: deep
---

# LaTeX 算法

LaTeX 中的算法排版通常分为两层：`algorithm` 负责浮动体、标题和引用，`algpseudocode` 负责伪代码内容。下面的示例统一使用这组宏包。

## 1 基础设置

在导言区加载宏包：

```latex
\usepackage{algorithm}
\usepackage{algpseudocode}
```

`algpseudocode` 会自动加载 `algorithmicx`，不需要重复引入。若不想显示 `end if`、`end while` 等结束标记，可以使用：

```latex
\usepackage[noend]{algpseudocode}
```

::: warning 注意
`algorithmic`、`algpseudocode` 和 `algorithm2e` 是不同的命令体系。比如 `algorithmic` 使用 `\STATE`，`algpseudocode` 使用 `\State`，不要混写。
:::

## 2 使用 `algpseudocode`

### 2.1 基础语法

```latex
\begin{algorithm}[!t]
    \caption{Select the Least-Loaded Port}
    \label{alg:least-loaded-port}
    \begin{algorithmic}[1]
        \Require Nonempty port set $P$ and load function $L$
        \Ensure A port $p^\star$ with minimum load
        \State $p^\star \gets \Call{First}{P}$
        \ForAll{$p \in P$}
            \If{$L(p) < L(p^\star)$}
                \State $p^\star \gets p$
            \EndIf
        \EndFor
        \State \Return $p^\star$
    \end{algorithmic}
\end{algorithm}
```

`algorithmic` 后面的 `[1]` 表示每行都编号。数字改为 `2` 时，每两行显示一次行号；省略该参数则不显示行号。

### 2.2 常用命令

| 命令                              | 功能                     |
| --------------------------------- | ------------------------ |
| `\Require`                        | 输入或前置条件           |
| `\Ensure`                         | 输出或后置条件           |
| `\State`                          | 普通语句                 |
| `\Statex`                         | 不编号的普通语句         |
| `\If{条件}`                       | 条件判断                 |
| `\ElsIf{条件}`                    | 追加条件分支             |
| `\Else`                           | 否则分支                 |
| `\EndIf`                          | 结束条件判断             |
| `\For{循环条件}`                  | `for` 循环               |
| `\ForAll{循环条件}`               | `for all` 循环           |
| `\EndFor`                         | 结束 `for` 循环          |
| `\While{条件}`                    | `while` 循环             |
| `\EndWhile`                       | 结束 `while` 循环        |
| `\Repeat`、`\Until{条件}`         | 先执行、后判断的循环     |
| `\Procedure{名称}{参数}`          | 定义过程                 |
| `\Function{名称}{参数}`           | 定义函数                 |
| `\Call{名称}{参数}`               | 调用过程或函数           |
| `\Comment{注释}`                  | 添加行尾注释             |
| `\State \Return 返回值`           | 返回结果                 |

## 3 完整示例

下面的代码可以直接放进一个独立的 `.tex` 文件：

```latex
\documentclass{article}
\usepackage{algorithm}
\usepackage{algpseudocode}

\begin{document}

\begin{algorithm}[!t]
    \caption{Select the Least-Loaded Port}
    \label{alg:least-loaded-port}
    \begin{algorithmic}[1]
        \Require Nonempty port set $P$ and load function $L$
        \Ensure A port $p^\star$ with minimum load
        \State $p^\star \gets \Call{First}{P}$
        \ForAll{$p \in P$}
            \If{$L(p) < L(p^\star)$}
                \State $p^\star \gets p$
            \EndIf
        \EndFor
        \State \Return $p^\star$
    \end{algorithmic}
\end{algorithm}

See Algorithm~\ref{alg:least-loaded-port}.

\end{document}
```

`\label` 一般紧跟在 `\caption` 后面，引用时使用相同的标签名。

## 4 常用调整

### 4.1 添加注释

`algpseudocode` 本身提供 `\Comment`，不需要再加载其他伪代码宏包：

```latex
\begin{algorithmic}[1]
    \State $p^\star \gets \Call{First}{P}$ \Comment{Initial candidate}
\end{algorithmic}
```

### 4.2 自定义关键字

需要补充 `algpseudocode` 没有提供的关键字时，可以用 `\algnewcommand` 定义。例如添加 `break`：

```latex
\algnewcommand{\algorithmicbreak}{\textbf{break}}
\algnewcommand{\Break}{\State \algorithmicbreak}

\begin{algorithmic}[1]
    \While{condition}
        \If{stop condition}
            \Break
        \EndIf
    \EndWhile
\end{algorithmic}
```

### 4.3 隐藏结束标记

`noend` 适合版面紧张的论文，但嵌套较多时，保留结束标记通常更容易读：

```latex
\usepackage[noend]{algpseudocode}
```

### 4.4 使用 `algorithm2e`

`algorithm2e` 是另一套完整方案，不要与前面的 `algorithm`、`algpseudocode` 同时加载。

```latex
\usepackage[ruled,vlined,linesnumbered]{algorithm2e}

\begin{algorithm}[!t]
    \caption{An Example}
    \label{alg:algorithm2e}
    \KwIn{Input parameters}
    \KwOut{Output result}
    \If{condition}{
        perform operation\;
    }
    \KwRet{result}\;
\end{algorithm}
```

## 5 常见问题

### 5.1 算法位置不符合预期

`[!t]` 会优先把算法放在页顶。必须固定在当前位置时可以改用 `[H]`，但这样可能留下较大的页面空白，一般只在确有需要时使用。

### 5.2 行号不显示

检查 `\begin{algorithmic}[1]` 中的 `[1]` 是否存在。这个参数控制编号间隔，不是起始行号。

### 5.3 命令未定义

最常见的原因是混用了两套命令：

```latex
% algpseudocode
\usepackage{algpseudocode}
\State $x \gets 1$

% algorithmic（旧写法）
\usepackage{algorithmic}
\STATE $x \leftarrow 1$
```

新文档建议使用第一种写法。同一份文档里不要同时加载这两个宏包。

### 5.4 中文支持

需要在标题或伪代码中直接写中文时，可以使用 `ctexart`，并用 XeLaTeX 或 LuaLaTeX 编译：

```latex
\documentclass[UTF8]{ctexart}
\usepackage{algorithm}
\usepackage{algpseudocode}
```
