---
title: LaTeX 算法伪代码与 algpseudocode 用法
description: 使用 algorithm 与 algpseudocode 排版算法，编写输入输出、循环和条件语句，设置行号与交叉引用，并说明较长算法的拆分方式及宏包兼容性。
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

`algpseudocode` 会自动加载 `algorithmicx`。需要隐藏 `end if`、`end while` 等结束标记时，可以使用：

```latex
\usepackage[noend]{algpseudocode}
```

::: info 宏包与命令体系
`algorithmic`、`algpseudocode` 和 `algorithm2e` 各有自己的命令体系。比如 `algorithmic` 使用 `\STATE`，`algpseudocode` 使用 `\State`，命令需与所选宏包对应。
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

示例假定端口集合有限且非空，`First(P)` 表示按遍历顺序选取第一个端口，负载函数在一次遍历中保持不变。负载相同的候选不会替换当前结果。

`algorithmic` 后面的 `[1]` 表示每行都编号。数字改为 `2` 时，每两行显示一次行号，省略该参数则不显示行号。`\Require`、`\Ensure` 和 `\Statex` 不参与普通语句编号。

### 2.2 常用命令

| 命令                      | 功能                 |
| ------------------------- | -------------------- |
| `\Require`                | 输入或前置条件       |
| `\Ensure`                 | 输出或后置条件       |
| `\State`                  | 普通语句             |
| `\Statex`                 | 不编号的普通语句     |
| `\If{条件}`               | 条件判断             |
| `\ElsIf{条件}`            | 追加条件分支         |
| `\Else`                   | 否则分支             |
| `\EndIf`                  | 结束条件判断         |
| `\For{循环条件}`          | `for` 循环           |
| `\ForAll{循环条件}`       | `for all` 循环       |
| `\EndFor`                 | 结束 `for` 循环      |
| `\While{条件}`            | `while` 循环         |
| `\EndWhile`               | 结束 `while` 循环    |
| `\Repeat`、`\Until{条件}` | 先执行、后判断的循环 |
| `\Procedure{名称}{参数}`  | 定义过程             |
| `\EndProcedure`           | 结束过程             |
| `\Function{名称}{参数}`   | 定义函数             |
| `\EndFunction`            | 结束函数             |
| `\Call{名称}{参数}`       | 调用过程或函数       |
| `\Comment{注释}`          | 添加行尾注释         |
| `\State \Return 返回值`   | 返回结果             |

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

`algpseudocode` 提供 `\Comment`，用于添加行尾注释：

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

### 4.4 长算法与分页

`algorithm` 浮动体不能跨页。长算法可将 `algorithmic` 直接放在正文中，随正文分页。标题和算法编号需另用模板认可的非浮动方案设置，也可将算法分段放入浮动体。

需要分段并继续行号、缩进和未结束的代码块时，`algorithmicx` 提供 `\algstore` 和 `\algrestore`：

```latex
\begin{algorithmic}[1]
    \State $i \gets 0$
    \While{$i < n$}
        \State Process item $i$
        \algstore{loop-state}
\end{algorithmic}

% 后续段落中继续，也可以放入另一个 algorithm 浮动体
\begin{algorithmic}[1]
    \algrestore{loop-state}
        \State $i \gets i + 1$
    \EndWhile
\end{algorithmic}
```

`\algstore` 必须是前一段结束前的最后一条命令，`\algrestore` 必须是后一段的第一条命令。外层 `algorithm` 的标题和编号需在各段浮动体中另行设置。

### 4.5 使用 `algorithm2e`

采用 `algorithm2e` 时，用它替换前面的 `algorithm`、`algpseudocode` 组合，算法环境与伪代码命令都由 `algorithm2e` 提供。

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

`[!t]` 请求页顶放置，并放宽部分浮动限制。需要固定在当前位置时，可以改用 `[H]`。算法放不进当前页时，固定位置可能留下较大的空白。

### 5.2 行号不显示

检查 `\begin{algorithmic}[1]` 中的编号间隔参数，`[1]` 表示每行显示行号，省略时不显示。

### 5.3 命令未定义

核对已加载的宏包与命令写法，例如：

```latex
% algpseudocode
\usepackage{algpseudocode}
\State $x \gets 1$

% algorithmic（旧写法）
\usepackage{algorithmic}
\STATE $x \leftarrow 1$
```

本页采用 `algpseudocode`，对应 `\State` 等命令。

### 5.4 中文支持

需要在标题或伪代码中直接写中文时，可以使用 `ctexart`，并用 XeLaTeX 或 LuaLaTeX 编译：

```latex
\documentclass[UTF8]{ctexart}
\usepackage{algorithm}
\usepackage{algpseudocode}
```
