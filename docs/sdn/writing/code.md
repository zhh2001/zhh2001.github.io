---
outline: deep
---

# LaTeX 代码排版

LaTeX 中常用 `listings` 或 `minted` 排版源代码。前者完全由 TeX 处理，配置简单。后者调用 Pygments，支持更多语言和更细致的语法高亮。

## 1 选择宏包

| 对比项       | `listings`                 | `minted`                         |
| ------------ | -------------------------- | -------------------------------- |
| 外部程序     | 不需要                     | 需要 `latexminted` 和 Pygments   |
| 语言支持     | 常见语言                   | Pygments 支持的语言              |
| 高亮效果     | 基于关键词规则             | 基于词法分析器                   |
| UTF-8 代码   | 多字节字符支持有限         | 默认使用 UTF-8                   |
| 投稿兼容性   | 通常更稳妥                 | 取决于编译环境是否允许外部程序   |
| 配置复杂度   | 较低                       | 较低，但依赖编译环境             |

期刊模板没有特殊要求时，两者选一个即可。同一份文档不建议同时加载，以免代码清单名称和相关命令冲突。

## 2 使用 `listings`

### 2.1 完整示例

```latex
\documentclass{article}
\usepackage{xcolor}
\usepackage{listings}

\definecolor{codebg}{RGB}{248,248,248}

\lstdefinestyle{code}{
  backgroundcolor=\color{codebg},
  basicstyle=\ttfamily\small,
  keywordstyle=\color{blue!70!black},
  commentstyle=\color{green!40!black},
  stringstyle=\color{orange!60!black},
  numbers=left,
  numberstyle=\tiny\color{gray},
  numbersep=8pt,
  frame=single,
  rulecolor=\color{black!15},
  breaklines=true,
  breakatwhitespace=false,
  showstringspaces=false,
  keepspaces=true,
  columns=fullflexible,
  tabsize=4,
  captionpos=b
}

\begin{document}

\begin{lstlisting}[
  style=code,
  language=Go,
  caption={A worker function},
  label={lst:worker}
]
package main

func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * job
    }
}
\end{lstlisting}

Listing~\ref{lst:worker} receives jobs and sends results.

\end{document}
```

`\lstdefinestyle` 用于保存一组配置，后续代码块通过 `style=code` 复用。局部参数会覆盖同名的全局参数。

### 2.2 内联代码

```latex
Use \lstinline[language=Go]|make(chan int)| to create a channel.
```

分隔符不一定使用 `|`，只要该字符没有出现在代码内容中即可。

### 2.3 引入外部文件

```latex
\lstinputlisting[
  style=code,
  language=Go,
  firstline=5,
  lastline=12,
  firstnumber=5
]{worker.go}
```

`firstline` 和 `lastline` 选择源文件范围，`firstnumber` 控制显示的起始行号。论文中的代码来自真实项目时，引用外部文件比复制代码更容易维护。

### 2.4 常用参数

| 参数                       | 作用                         |
| -------------------------- | ---------------------------- |
| `language=Go`              | 指定语言                     |
| `numbers=left`             | 在左侧显示行号               |
| `firstnumber=10`           | 设置起始行号                 |
| `breaklines=true`          | 自动折行                     |
| `showstringspaces=false`   | 不标记字符串中的空格         |
| `keepspaces=true`          | 保留代码中的空格             |
| `columns=fullflexible`     | 按等宽字体的自然宽度排版     |
| `frame=single`             | 添加边框                     |
| `caption={标题}`           | 添加标题                     |
| `label={lst:key}`          | 设置交叉引用标签             |
| `float=tbp`                | 把代码块作为浮动体           |

较长代码不宜设置 `float`，否则代码块不能自然跨页。

## 3 使用 `minted`

### 3.1 完整示例

```latex
\documentclass{article}
\usepackage{xcolor}
\usepackage{minted}

\definecolor{codebg}{RGB}{248,248,248}

\setminted{
  fontsize=\small,
  linenos,
  breaklines,
  breakanywhere,
  autogobble,
  frame=single,
  framesep=2mm,
  bgcolor=codebg
}

\begin{document}

\begin{listing}[tbp]
\begin{minted}{go}
package main

func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * job
    }
}
\end{minted}
\caption{A worker function}
\label{lst:minted-worker}
\end{listing}

Listing~\ref{lst:minted-worker} is highlighted by Pygments.

\end{document}
```

`\setminted` 设置全局参数。也可以写成 `\setminted[go]{...}`，只对 Go 代码生效。

### 3.2 内联代码

```latex
Use \mintinline{go}|make(chan int)| to create a channel.
```

### 3.3 引入外部文件

```latex
\inputminted[
  firstline=5,
  lastline=12,
  firstnumber=5,
  linenos
]{go}{worker.go}
```

语法为 `\inputminted[参数]{语言}{文件名}`。外部文件默认按 UTF-8 解码。

### 3.4 浮动代码清单

`listing` 环境提供标题、编号和交叉引用，其行为与 `figure`、`table` 类似：

```latex
\begin{listing}[tbp]
\begin{minted}{go}
func add(a, b int) int {
    return a + b
}
\end{minted}
\caption{Add two integers}
\label{lst:add}
\end{listing}
```

浮动代码清单不能跨页。长代码应直接使用 `minted` 环境，不要放进 `listing`。

### 3.5 编译权限

minted 3 使用 `latexminted` 生成高亮结果。TeX Live 2024 及以上版本已将它加入受信任程序列表，通常可以直接编译：

```sh
xelatex main.tex
```

旧版 TeX Live 或部分 MiKTeX 环境可能仍需显式授权：

```sh
xelatex -shell-escape main.tex
```

`-shell-escape` 允许 LaTeX 执行外部命令，只应对可信文档启用。

## 4 排版细节

### 4.1 长代码行

`listings` 使用 `breaklines=true`。`minted` 使用 `breaklines`，必要时再加 `breakanywhere`。自动折行只改变排版，不会修改源文件。

### 4.2 缩进与制表符

论文代码最好统一使用空格缩进。制表符宽度在不同工具中可能不一致：`listings` 使用 `tabsize`，`minted` 使用 `tabsize`。

### 4.3 中文注释

`listings` 对多字节 UTF-8 源文件支持有限，中文注释可能需要额外转换或转义。需要直接排版含中文的 UTF-8 源代码时，优先使用 minted，并用 XeLaTeX 或 LuaLaTeX 编译。

### 4.4 在代码中插入 LaTeX

仅在确有需要时启用转义，避免源代码被误当作 LaTeX：

```latex
% listings
\begin{lstlisting}[escapeinside={(*@}{@*)}]
x := 1 // (*@\textit{initialize}@*)
\end{lstlisting}

% minted
\begin{minted}[escapeinside=||]{go}
x := 1 // |\textit{initialize}|
\end{minted}
```

## 5 常见问题

### 5.1 找不到语言

先确认语言名称是否被当前宏包支持。`listings` 与 Pygments 的语言名称不完全相同，例如同一种语言可能使用不同的别名。

### 5.2 minted 无法执行

依次检查 `minted.sty`、`latexminted` 和 Python 是否已安装并位于 `PATH`，再检查当前 TeX 发行版是否允许执行 `latexminted`。不要一遇到报错就直接开启全局 shell escape。

### 5.3 代码超出页边距

确认已启用自动折行。若某个很长的标识符仍无法换行，minted 可以增加 `breakanywhere`，`listings` 可以调整 `breakatwhitespace`。

### 5.4 行号与源文件不一致

截取外部文件时，同时设置 `firstline` 和 `firstnumber`。前者控制读取位置，后者只控制打印出来的数字。

### 5.5 特殊字符报错

代码必须放在 `lstlisting`、`minted`、`\lstinline` 或 `\mintinline` 中。直接写进普通正文时，`#`、`%`、`_`、`{`、`}` 等字符仍会按 LaTeX 语法解析。
