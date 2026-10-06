---
title: LaTeX 学术写作与排版笔记
description: 学术写作中常用的 LaTeX 排版方法，涵盖公式、数值单位、定理、图片、表格、算法、代码和参考文献，附各主题入口与编译环境说明。
outline: false
---

# LaTeX 学术写作与排版 {#latex-写作}

整理学术写作中常用的 LaTeX 排版方法，包括数学内容、实验结果、代码和参考文献。

## 数学排版

- [公式](/sdn/writing/formula)：行内公式、多行公式、矩阵和常用数学符号。
- [数值与单位](/sdn/writing/unit)：数值格式、物理量、复合单位和网络常用单位。
- [定理环境](/sdn/writing/theorem)：定理、引理、定义、证明及编号设置。

## 内容呈现

- [图片](/sdn/writing/figure)：尺寸控制、浮动位置、子图、标题和交叉引用。
- [表格](/sdn/writing/table)：三线表、对齐、跨行跨列和宽表处理。
- [算法](/sdn/writing/algorithm)：伪代码、行号、输入输出和跨页处理。
- [代码排版](/sdn/writing/code)：使用 `listings` 或 `minted` 排版源代码。

## 文献管理

- [参考文献](/sdn/writing/bibliography)：BibTeX、biblatex、引用命令和条目维护。

## 编译与阅读

公式页的效果由本站的 MathJax 渲染。图片、表格、算法等页面中的代码则用于独立的 `.tex` 文档，不会在网页中执行。只有带 `\documentclass` 和 `\begin{document}` 的示例才是完整文档，其余片段需要放到对应的导言区或正文中。

含中文的完整示例采用 `ctexart`，可保存为 `main.tex` 后用 XeLaTeX 编译。交叉引用通常需要编译两遍：

```sh
xelatex main.tex
xelatex main.tex
```

参考文献还需运行 BibTeX 或 Biber，具体顺序见参考文献页。本文示例已在 TeX Live 2026 下核验，`siunitx` 为 3.5.5，`minted` 为 3.8.0。较新的命令在各页注明了版本要求。

表格中的数值和图片中的方法名称用于演示排版，不代表本仓库的实验结果。投稿时以目标期刊、会议或学校模板为准，先确认宏包和编译方式是否受支持。
