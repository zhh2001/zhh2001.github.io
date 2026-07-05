---
outline: deep
---

# LaTeX 图片

`graphicx` 负责插入、缩放、裁剪和旋转图片，`figure` 环境负责标题、编号和浮动位置。只需要显示图片时可以直接使用 `\includegraphics`。需要标题或交叉引用时，再把它放进 `figure`。

## 1 插入图片

### 1.1 加载宏包

```latex
\usepackage{graphicx}
```

最简单的写法如下：

```latex
\includegraphics{figures/topology.pdf}
```

图片路径区分大小写。扩展名可以省略，便于在不修改正文的情况下替换图片格式。

### 1.2 设置图片宽度

```latex
\includegraphics[width=0.8\linewidth]{figures/topology}
```

论文中优先使用 `\linewidth` 控制宽度。它会随当前环境变化，在单栏、双栏和子图中都比固定厘米值更容易复用。

同时限制宽度和高度时，应保留纵横比：

```latex
\includegraphics[
  width=\linewidth,
  height=0.35\textheight,
  keepaspectratio
]{figures/topology}
```

### 1.3 统一图片路径

图片较多时，可以在导言区设置搜索目录：

```latex
\graphicspath{{figures/}{images/}}
```

每个目录都需要单独放在一组花括号中。设置后可以简写为：

```latex
\includegraphics[width=0.8\linewidth]{topology}
```

## 2 浮动图片

### 2.1 基本结构

```latex
\begin{figure}[htbp]
  \centering
  \includegraphics[width=0.8\linewidth]{figures/topology}
  \caption{Testbed topology}
  \label{fig:topology}
\end{figure}
```

常用位置参数如下：

| 参数 | 含义             |
| ---- | ---------------- |
| `h`  | 当前文字附近     |
| `t`  | 页面顶部         |
| `b`  | 页面底部         |
| `p`  | 单独的浮动体页面 |
| `!`  | 放宽部分限制     |

`[htbp]` 表示放置偏好，并不保证图片固定在源码位置。LaTeX 会根据页面剩余空间、浮动体数量和排版参数选择位置。

### 2.2 交叉引用

`\label` 应放在 `\caption` 后面，确保记录的是当前图片编号：

```latex
As shown in Figure~\ref{fig:topology}, the switches form a tree.
```

`~` 是不可断行空格，可以避免 “Figure” 和编号分到两行。

### 2.3 限制浮动范围

确实不希望图片越过某个位置时，可以使用 `placeins`：

```latex
\usepackage{placeins}

% 正文和图片
\FloatBarrier
```

`\FloatBarrier` 会先排出此前尚未放置的浮动体。它适合章节边界，不宜在每张图片后都使用。

需要强制放在当前位置时，可以使用 `float` 提供的 `[H]`：

```latex
\usepackage{float}

\begin{figure}[H]
  % ...
\end{figure}
```

`[H]` 会减少 LaTeX 调整页面的空间，连续使用时容易留下大片空白，因此只适合少量确实不能浮动的图片。

## 3 标题设置

### 3.1 标题样式

`caption` 宏包可以统一设置图片标题：

```latex
\usepackage{caption}

\captionsetup[figure]{
  font=small,
  labelfont=bf,
  labelsep=period
}
```

局部设置只影响当前图片：

```latex
\begin{figure}[htbp]
  \captionsetup{justification=centering}
  % ...
\end{figure}
```

投稿模板已经规定标题格式时，不应再用 `caption` 强行覆盖。

### 3.2 目录中的短标题

标题较长时，可以为图片目录提供短标题：

```latex
\caption[Testbed topology]{
  Testbed topology and the connections between all programmable switches.
}
```

方括号中的内容写入图片目录，花括号中的内容显示在图片下方。

## 4 子图

### 4.1 基本写法

使用 `subcaption` 提供的 `subfigure` 环境：

```latex
\usepackage{subcaption}

\begin{figure}[htbp]
  \centering
  \begin{subfigure}[b]{0.47\linewidth}
    \centering
    \includegraphics[width=\linewidth]{figures/baseline}
    \caption{Baseline}
    \label{fig:baseline}
  \end{subfigure}
  \hfill
  \begin{subfigure}[b]{0.47\linewidth}
    \centering
    \includegraphics[width=\linewidth]{figures/proposed}
    \caption{Proposed method}
    \label{fig:proposed}
  \end{subfigure}
  \caption{Latency comparison}
  \label{fig:latency-comparison}
\end{figure}
```

两个子图宽度与中间间距之和不能超过当前行宽。`\hfill` 会把剩余空间放在两个子图之间。

### 4.2 引用子图

```latex
Figure~\ref{fig:latency-comparison} compares both methods.
The result in Figure~\ref{fig:proposed} has lower latency.
```

引用整个组合图时使用外层标签，引用单个子图时使用对应 `subfigure` 中的标签。

不要同时加载 `subcaption`、`subfig` 和已经过时的 `subfigure`。它们提供相近功能，混用容易造成命令或标题格式冲突。

## 5 裁剪与旋转

### 5.1 裁剪空白

```latex
\includegraphics[
  width=0.8\linewidth,
  trim={20 10 20 10},
  clip
]{figures/result.pdf}
```

`trim` 的顺序为左、下、右、上。只有同时加入 `clip`，超出裁剪区域的内容才会被隐藏。

### 5.2 旋转图片

```latex
\includegraphics[
  width=0.6\linewidth,
  angle=90,
  origin=c
]{figures/result.pdf}
```

`origin=c` 表示围绕图片中心旋转。整页横向图片较多时，使用专门的横向页面环境通常比逐张旋转更容易维护。

## 6 双栏与环绕图片

### 6.1 跨双栏图片

双栏文档中，普通 `figure` 只占一栏。跨越两栏使用 `figure*`：

```latex
\begin{figure*}[t]
  \centering
  \includegraphics[width=0.9\textwidth]{figures/architecture}
  \caption{System architecture}
  \label{fig:architecture}
\end{figure*}
```

双栏浮动体通常只能出现在页顶或页底，位置限制比单栏图片更严格。

### 6.2 文字环绕

使用 `wrapfig` 可以让正文环绕较小的图片：

```latex
\usepackage{wrapfig}

\begin{wrapfigure}{r}{0.4\textwidth}
  \centering
  \includegraphics[width=0.38\textwidth]{figures/device}
  \caption{Test device}
  \label{fig:device}
\end{wrapfigure}
```

环绕图片应放在相关段落之前，并避开章节标题、列表和其他浮动体。正文太短时，图片周围容易出现不自然的空白。

## 7 图片格式与清晰度

- 曲线图、拓扑图和架构图优先使用 PDF 等矢量格式，缩放后边缘仍然清晰。
- 截图和带透明区域的图像适合 PNG。
- 照片适合 JPEG，不要反复压缩后再用于论文。
- 位图应按最终排版尺寸导出，避免在 LaTeX 中大幅放大。
- 图中的字体、线宽和颜色应在导出阶段统一，不要依赖整体缩放修正比例。

`\includegraphics` 不能直接处理所有 SVG 工作流。可以预先转换为 PDF，或者在编译环境允许调用外部程序时使用 `svg` 宏包。

## 8 完整示例

下面的 `example-image-a` 和 `example-image-b` 来自 `mwe` 宏包，实际使用时替换为自己的图片路径。

```latex
\documentclass[UTF8]{ctexart}
\usepackage{graphicx}
\usepackage{subcaption}
\usepackage{mwe}

\captionsetup[figure]{
  font=small,
  labelfont=bf,
  labelsep=period
}

\begin{document}

图~\ref{fig:comparison} 对比了两组实验结果。

\begin{figure}[htbp]
  \centering
  \begin{subfigure}[b]{0.47\linewidth}
    \centering
    \includegraphics[width=\linewidth]{example-image-a}
    \caption{基线方法}
    \label{fig:example-baseline}
  \end{subfigure}
  \hfill
  \begin{subfigure}[b]{0.47\linewidth}
    \centering
    \includegraphics[width=\linewidth]{example-image-b}
    \caption{改进方法}
    \label{fig:example-proposed}
  \end{subfigure}
  \caption{两种方法的实验结果对比}
  \label{fig:comparison}
\end{figure}

\end{document}
```

## 9 常见问题

### 9.1 找不到图片

检查相对路径、文件名大小写和扩展名。项目从主文件所在目录编译时，图片路径通常应相对于主文件或通过 `\graphicspath` 配置的目录，而不是相对于被 `\input` 的章节文件。

### 9.2 图片超出版心

把宽度限制为当前环境的行宽：

```latex
\includegraphics[width=\linewidth]{figures/result}
```

子图内部也应使用 `\linewidth`，此时它表示子图自身的宽度。

### 9.3 图片没有出现在源码位置

这是浮动体的正常行为。先检查图片是否过大，再尝试 `[!htbp]` 或在章节边界使用 `\FloatBarrier`。不要把 `[H]` 作为所有图片的默认设置。

### 9.4 引用编号错误

确保 `\label` 位于 `\caption` 后，并至少编译两次，让 LaTeX 写入并读取交叉引用信息。

### 9.5 图片模糊

优先检查原始文件。LaTeX 只能缩放已有图像，不能恢复被压缩或低分辨率图片中已经丢失的细节。

### 9.6 图片周围有大块空白

空白可能来自图片文件本身，也可能由浮动位置造成。前者可在导出时去除画布边缘，或使用 `trim` 和 `clip`。后者应调整图片尺寸和浮动参数。
