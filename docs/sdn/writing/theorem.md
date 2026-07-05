---
outline: deep
---

# LaTeX 定理环境

`amsthm` 用于定义定理、引理、定义、命题和证明等环境。它负责标题样式、自动编号和证明结束符，通常与 `amsmath` 一起使用。

## 1 基础设置

### 1.1 加载宏包

```latex
\usepackage{amsmath}
\usepackage{amsthm}
```

同时使用两个宏包时，应先加载 `amsmath`，再加载 `amsthm`。AMS 文档类已经包含 `amsthm` 的功能，不需要重复加载。

### 1.2 定义环境

```latex
\newtheorem{theorem}{Theorem}
\newtheorem{lemma}{Lemma}
```

`\newtheorem{环境名}{显示名称}` 定义一个带编号的环境：

```latex
\begin{theorem}
Every finite tree has at least one leaf.
\end{theorem}
```

环境名用于 LaTeX 源码，显示名称出现在最终文档中。两者可以不同。

## 2 编号方式

### 2.1 独立编号

```latex
\newtheorem{theorem}{Theorem}
\newtheorem{lemma}{Lemma}
```

此时定理和引理分别从 1 开始编号。

### 2.2 共享编号

```latex
\newtheorem{theorem}{Theorem}
\newtheorem{lemma}[theorem]{Lemma}
\newtheorem{corollary}[theorem]{Corollary}
```

`[theorem]` 表示引理和推论共用 `theorem` 计数器，编号可能依次为 Theorem 1、Lemma 2、Corollary 3。

### 2.3 按节编号

```latex
\newtheorem{theorem}{Theorem}[section]
```

最后的 `[section]` 表示每节重新编号，得到 Theorem 1.1、Theorem 1.2、Theorem 2.1 等形式。共享该计数器的环境会沿用相同规则：

```latex
\newtheorem{lemma}[theorem]{Lemma}
```

### 2.4 不编号环境

```latex
\newtheorem*{remark}{Remark}
```

带星号的 `\newtheorem*` 不创建编号，适合备注、说明和只出现一次的结论。

## 3 定理样式

`amsthm` 提供三种默认样式：

| 样式         | 常见用途                       | 正文字体             |
| ------------ | ------------------------------ | -------------------- |
| `plain`      | 定理、引理、命题、推论         | 通常为斜体           |
| `definition` | 定义、条件、例子、问题         | 通常为正体           |
| `remark`     | 备注、注释、结论               | 正体且间距较紧       |

`\theoremstyle` 只影响它后面定义的环境：

```latex
\theoremstyle{plain}
\newtheorem{theorem}{Theorem}[section]
\newtheorem{lemma}[theorem]{Lemma}

\theoremstyle{definition}
\newtheorem{definition}[theorem]{Definition}
\newtheorem{example}[theorem]{Example}

\theoremstyle{remark}
\newtheorem*{remark}{Remark}
```

### 3.1 自定义样式

默认样式不能满足需求时，可以使用 `\newtheoremstyle`：

```latex
\newtheoremstyle{compact}
  {4pt}        % 上方间距
  {4pt}        % 下方间距
  {\itshape}   % 正文字体
  {}           % 缩进
  {\bfseries}  % 标题字体
  {.}          % 标题后的标点
  {.5em}       % 标题后的间距
  {}           % 标题格式

\theoremstyle{compact}
\newtheorem{claim}{Claim}
```

自定义样式应在对应的 `\newtheorem` 之前启用。

## 4 标题与交叉引用

### 4.1 补充标题

环境的可选参数会显示在标题后面：

```latex
\begin{theorem}[Shortest Path Existence]
Every finite graph with nonnegative edge weights has a shortest path
between any connected pair of vertices.
\end{theorem}
```

### 4.2 设置标签

`\label` 应紧跟在 `\begin{...}` 或其可选标题之后：

```latex
\begin{theorem}[Shortest Path Existence]\label{thm:shortest-path}
Every finite graph with nonnegative edge weights has a shortest path.
\end{theorem}

By Theorem~\ref{thm:shortest-path}, the minimum is attained.
```

建议统一使用 `thm:`、`lem:`、`def:`、`prop:` 等标签前缀。

## 5 证明环境

### 5.1 基本证明

```latex
\begin{proof}
The result follows from the finiteness of the graph.
\end{proof}
```

`proof` 会自动添加证明标题和末尾的 QED 符号，不需要手工输入方框。

### 5.2 自定义证明标题

```latex
\begin{proof}[Proof of Theorem~\ref{thm:shortest-path}]
The result follows from the finiteness of the graph.
\end{proof}
```

中文文档可以修改默认标题：

```latex
\renewcommand{\proofname}{证明}
```

### 5.3 调整 QED 位置

证明以独立公式或列表结束时，QED 符号可能另起一行。把 `\qedhere` 放在期望的位置：

```latex
\begin{proof}
Therefore,
\begin{equation*}
C(P') = C(P) - C(Q) \leq C(P). \qedhere
\end{equation*}
\end{proof}
```

不要在同一个证明中同时手工输入 `\qed` 和使用自动结束符。

## 6 完整示例

```latex
\documentclass[UTF8]{ctexart}
\usepackage{amsmath}
\usepackage{amsthm}

\theoremstyle{plain}
\newtheorem{theorem}{定理}[section]
\newtheorem{lemma}[theorem]{引理}

\theoremstyle{definition}
\newtheorem{definition}[theorem]{定义}

\theoremstyle{remark}
\newtheorem*{remark}{注}

\renewcommand{\proofname}{证明}

\begin{document}

\section{图上的路径}

\begin{definition}[路径代价]\label{def:path-cost}
设 $P=(v_0,v_1,\ldots,v_k)$ 是一条游走，其代价定义为
\[
C(P)=\sum_{i=1}^{k} w(v_{i-1},v_i).
\]
\end{definition}

\begin{lemma}\label{lem:remove-cycle}
若所有边权非负，从游走中删除一个环不会增加路径代价。
\end{lemma}

\begin{proof}
设删除的环为 $Q$，删除后得到 $P'$。由定义~\ref{def:path-cost}，
\begin{equation*}
C(P')=C(P)-C(Q)\leq C(P). \qedhere
\end{equation*}
\end{proof}

\begin{theorem}[最短路径存在性]\label{thm:shortest-path}
设 $G$ 是有限有向图，且所有边权非负。若顶点 $s$ 到 $t$ 可达，
则 $s$ 到 $t$ 的最短路径存在。
\end{theorem}

\begin{proof}
对任意一条从 $s$ 到 $t$ 的游走，反复应用引理~\ref{lem:remove-cycle}，
可以得到一条代价不更大的简单路径。有限图中只有有限条简单路径，
而候选集合非空，因此其中必有一条取得最小代价。
\end{proof}

\begin{remark}
若允许负权环，上述删除环的论证不再成立。
\end{remark}

\end{document}
```

## 7 常见问题

### 7.1 环境未定义

`amsthm` 不会自动创建名为 `theorem`、`lemma` 或 `definition` 的环境。使用前必须在导言区通过 `\newtheorem` 定义。

### 7.2 编号没有共享

检查共享计数器是否写在环境名之后：

```latex
\newtheorem{lemma}[theorem]{Lemma}
```

它不同于按节编号的写法 `\newtheorem{theorem}{Theorem}[section]`。

### 7.3 定义正文变成斜体

在定义环境的 `\newtheorem` 之前加入 `\theoremstyle{definition}`。样式切换不会追溯影响已经定义的环境。

### 7.4 定理以列表开头时排版异常

先进入水平模式，再开始列表：

```latex
\begin{theorem}
\leavevmode
\begin{enumerate}
  \item First condition.
  \item Second condition.
\end{enumerate}
\end{theorem}
```

### 7.5 QED 符号位置不正确

证明最后是独立公式或列表时使用 `\qedhere`。不要在 `$$...$$` 或 `eqnarray` 中使用它，应改用 `amsmath` 提供的公式环境。

### 7.6 与投稿模板冲突

部分文档类已经定义了定理或证明环境。遇到重复定义、间距或标题样式异常时，应优先使用模板提供的设置，不要强行覆盖。
