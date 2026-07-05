---
outline: deep
---

# LaTeX 表格

`tabular` 负责表格内容和列格式，外层的 `table` 环境负责标题、编号和浮动位置。简单表格只用 `tabular` 即可。需要引用时，再把它放进 `table`。

## 1 基本表格结构

### 1.1 最简单的表格

使用 `tabular` 环境创建基本表格：

```latex
\begin{tabular}{|c|c|c|}
\hline
列1 & 列2 & 列3 \\ \hline
A & B & C \\ \hline
D & E & F \\ \hline
\end{tabular}
```

- `|` 表示竖线
- `c` 表示居中对齐
- `&` 分隔列
- `\\` 换行
- `\hline` 添加水平线

### 1.2 对齐方式

- `l` 左对齐
- `c` 居中对齐
- `r` 右对齐
- `p{宽度}` 固定宽度并自动换行

示例：

```latex
\begin{tabular}{|l|c|r|}
\hline
左对齐 & 居中对齐 & 右对齐 \\ \hline
A & B & C \\ \hline
D & E & F \\ \hline
\end{tabular}
```

## 2 表格边框控制

### 2.1 无边框表格

```latex
\begin{tabular}{ccc}
A & B & C \\
D & E & F \\
\end{tabular}
```

### 2.2 部分边框

```latex
\begin{tabular}{|l|c|r|}
\hline
左 & 中 & 右 \\ \hline
A & B & C \\
D & E & F \\ \hline
\end{tabular}
```

上面的例子只在首尾添加水平线。

### 2.3 更精细的边框控制

`\cline{i-j}` 只绘制第 $i$ 列到第 $j$ 列之间的横线：

```latex
\begin{tabular}{|l|c|r|}
\hline
左 & 中 & 右 \\ \hline
A & B & C \\ \cline{2-3}
D & E & F \\ \hline
\end{tabular}
```

## 3 合并单元格

### 3.1 横向合并

语法为 `\multicolumn{列数}{列格式}{内容}`：

```latex
\begin{tabular}{|c|c|c|}
\hline
\multicolumn{2}{|c|}{合并两列} & 列3 \\ \hline
A & B & C \\ \hline
D & E & F \\ \hline
\end{tabular}
```

### 3.2 纵向合并

`\multirow{行数}{宽度}{内容}`

需要 `multirow` 包：

```latex
\usepackage{multirow}
% ...
\begin{tabular}{|c|c|c|}
\hline
\multirow{2}{*}{合并两行} & 列2 & 列3 \\ \cline{2-3}
 & B & C \\ \hline
D & E & F \\ \hline
\end{tabular}
```

## 4 表格尺寸控制

### 4.1 固定列宽

`p{宽度}` 会让内容在固定宽度内自动换行，并按段落顶部对齐：

```latex
\begin{tabular}{|p{2cm}|p{3cm}|p{1.5cm}|}
\hline
短 & 较长文本会自动换行 & 很短 \\ \hline
A & B & C \\ \hline
\end{tabular}
```

### 4.2 表格整体宽度

使用 `tabularx` 包：

```latex
\usepackage{tabularx}
% ...
\begin{tabularx}{\textwidth}{|X|X|X|}
\hline
自动调整宽度 & 自动调整宽度 & 自动调整宽度 \\ \hline
A & B & C \\ \hline
\end{tabularx}
```

`X` 表示自动分配剩余宽度的列。

## 5 长表格处理

对于跨页表格，使用 `longtable` 包：

```latex
\usepackage{longtable}

\begin{longtable}{|c|c|c|}
\caption{跨页表格示例}\label{tab:longtable} \\
\hline
\textbf{列1} & \textbf{列2} & \textbf{列3} \\ \hline
\endfirsthead
\multicolumn{3}{c}{续表~\thetable} \\
\hline
\textbf{列1} & \textbf{列2} & \textbf{列3} \\ \hline
\endhead
\hline
\multicolumn{3}{|r|}{\small 续下页} \\ \hline
\endfoot
\hline
\endlastfoot
A & B & C \\ \hline
D & E & F \\ \hline
% 继续添加数据行
\end{longtable}
```

`longtable` 本身会分页，不要再套在 `table` 浮动体中。

## 6 表格样式美化

### 6.1 使用 `booktabs` 包

```latex
\usepackage{booktabs}
% ...
\begin{tabular}{ccc}
\toprule
列1 & 列2 & 列3 \\
\midrule
A & B & C \\
D & E & F \\
\botrule
\end{tabular}
```

- `\toprule`：顶部横线
- `\midrule`：表头与数据之间的横线
- `\botrule`：底部横线

`booktabs` 表格通常不使用竖线，也不需要在每一行后重复画横线。

### 6.2 添加颜色

使用 `xcolor` 的 `table` 选项：

```latex
\usepackage[table]{xcolor}

\begin{tabular}{|l|c|r|}
\hline
\rowcolor{gray!30} 左 & 中 & 右 \\ \hline
A & \cellcolor{red!30}B & C \\ \hline
D & E & F \\ \hline
\end{tabular}
```

## 7 复杂表格示例

```latex
\documentclass[UTF8]{ctexart}
\usepackage{array}
\usepackage{multirow}
\usepackage[table]{xcolor}

\begin{document}

\begin{table}[htbp]
\centering
\caption{复杂表格示例}
\label{tab:complex}
\begin{tabular}{|>{\columncolor{gray!20}}l|c|r|}
\hline
\rowcolor{gray!40}\textbf{项目} & \textbf{数值} & \textbf{说明} \\ \hline
\multirow{2}{*}{第一组} & 10 & 描述1 \\ \cline{2-3}
 & 20 & 描述2 \\ \hline
\multirow{3}{*}{第二组} & \cellcolor{yellow!30}30 & 描述3 \\ \cline{2-3}
 & 40 & 描述4 \\ \cline{2-3}
 & 50 & 描述5 \\ \hline
\end{tabular}
\end{table}

\end{document}
```

## 8 表格与文本环绕

使用 `wrapfig` 宏包提供的 `wraptable` 环境实现文字环绕：

```latex
\usepackage{wrapfig}
% ...
\begin{wraptable}{r}{0.5\textwidth}
\centering
\begin{tabular}{|c|c|}
\hline
A & B \\ \hline
C & D \\ \hline
\end{tabular}
\caption{环绕表格}
\label{tab:wrapped}
\end{wraptable}
这里是与表格相关的正文内容。
```

## 9 对齐与浮动体

### 9.1 表格浮动

`tabular` 不是浮动体；把它放进 `table` 后，LaTeX 才会根据可用空间调整位置。常用的位置参数如下：

```latex
\begin{table}[htbp]
% h：当前位置附近
% t：页顶
% b：页底
% p：浮动体专页
\centering
% 表格内容
\end{table}
```

这些参数表示放置偏好，并不保证绝对位置。确实需要固定时可以使用 `float` 宏包提供的 `[H]`，但不宜滥用。

### 9.2 精确对齐

定宽列中使用 `\centering`、`\raggedright`、`\raggedleft` 时，要用 `array` 宏包的 `\arraybackslash` 恢复行尾的 `\\`：

```latex
\usepackage{array}

\begin{tabular}{
    |>{\raggedright\arraybackslash}p{3cm}
    |>{\centering\arraybackslash}p{2cm}
    |>{\raggedleft\arraybackslash}p{3cm}|
}
\hline
左对齐文本 & 居中文本 & 右对齐文本 \\ \hline
左侧内容 & 居中内容 & 右侧内容 \\ \hline
\end{tabular}
```

## 10 常见问题解决

### 10.1 表格线不显示

检查列格式中的 `|` 和表格行之间的 `\hline`。使用 `booktabs` 时，则应改用 `\toprule`、`\midrule` 和 `\bottomrule`。

### 10.2 单元格内容过长

单列内容过长时使用 `p{宽度}`；需要表格适应版心宽度时使用 `tabularx`。不建议一开始就用 `\resizebox` 缩放整张表，因为文字也会一起变小。

### 10.3 表格跨页

使用 `longtable`，并避免在外层嵌套 `table`。

### 10.4 中文支持

使用 `ctexart` 文档类，并用 XeLaTeX 或 LuaLaTeX 编译：

```latex
\documentclass[UTF8]{ctexart}
```

## 11 完整示例

```latex
\documentclass[UTF8]{ctexart}
\usepackage{array}
\usepackage{multirow}
\usepackage[table]{xcolor}
\usepackage{wrapfig}

\begin{document}

\begin{table}[htbp]
\centering
\caption{示例表格}
\label{tab:example}
\begin{tabular}{|>{\columncolor{gray!10}}l|c|r|}
\hline
\rowcolor{gray!30}\textbf{项目} & \textbf{数值} & \textbf{说明} \\ \hline
A & 10 & 描述1 \\ \hline
B & 20 & 描述2 \\ \hline
\multirow{2}{*}{C} & 30 & \cellcolor{yellow!20}描述3 \\ \cline{2-3}
 & 40 & 描述4 \\ \hline
\end{tabular}
\end{table}

\begin{wraptable}{r}{0.4\textwidth}
\centering
\begin{tabular}{|c|c|}
\hline
X & Y \\ \hline
1 & 2 \\ \hline
3 & 4 \\ \hline
\end{tabular}
\caption{环绕表格}
\label{tab:wrapped-example}
\end{wraptable}
环绕表格应放在相关段落之前。正文足够长时，文字会沿着表格另一侧继续排版；段落太短则可能出现较大的空白。

\end{document}
```
