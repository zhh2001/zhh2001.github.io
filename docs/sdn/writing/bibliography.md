---
title: LaTeX 参考文献与 BibTeX、Biber 编译流程
description: 维护 LaTeX 文献数据库，区分 BibTeX 与 biblatex、Biber 的编译流程，说明引用命令、文献样式、DOI 字段和条目整理中的常见问题。
outline: deep
---

# LaTeX 参考文献

LaTeX 通常把文献数据放在 `.bib` 文件中，再由 BibTeX 或 Biber 生成参考文献列表。投稿时优先遵循期刊、会议模板指定的工具和样式，不要在同一份文档中混用两套流程。

## 1 文献数据库

### 1.1 `.bib` 文件

新建 `references.bib`，每条文献由类型、引用键和字段组成：

```bibtex
@article{zhang2026ttint,
  author  = {Zhang, Henghua and Chen, Jue and Wu, Yuhang and Xiong, Yujie},
  title   = {{TT-INT}: A Time-Threshold-Based Lightweight In-Band Network
             Telemetry Scheme for {P4}-Enabled Programmable Networks},
  journal = {IEEE Transactions on Network and Service Management},
  year    = {2026},
  volume  = {23},
  pages   = {4335--4354},
  doi     = {10.1109/TNSM.2026.3688086}
}
```

示例使用已发表的 [TT-INT 论文](https://doi.org/10.1109/TNSM.2026.3688086)。`zhang2026ttint` 是引用键，应在文档加载的所有文献库中保持唯一。引用键一般采用“第一作者姓氏 + 年份 + 标题关键词”的形式。

### 1.2 常见文献类型

| 类型             | 用途         | 常用字段                                           |
| ---------------- | ------------ | -------------------------------------------------- |
| `@article`       | 期刊论文     | `author`、`title`、`journal`、`year`、`doi`        |
| `@inproceedings` | 会议论文     | `author`、`title`、`booktitle`、`year`             |
| `@book`          | 书籍         | `author` 或 `editor`、`title`、`publisher`、`year` |
| `@mastersthesis` | 硕士学位论文 | `author`、`title`、`school`、`year`                |
| `@phdthesis`     | 博士学位论文 | `author`、`title`、`school`、`year`                |
| `@misc`          | 其他材料     | 根据文献类型填写                                   |
| `@online`        | 在线资料     | `author`、`title`、`url`、`urldate`                |

这张表不是必填字段的完整清单。`@online` 属于 `biblatex` 数据模型，需要兼容传统 BibTeX 时，可改用 `@misc`。传统 BibTeX 的学位论文使用 `school`，biblatex 的对应字段通常使用 `institution`，并以 `@thesis` 及 `type` 区分论文类型。

### 1.3 字段写法

- 多位作者使用 `and` 分隔，不能用逗号代替。
- 姓名写成 `Zhang, Henghua`，可以明确区分姓和名。
- 缩写和专有名词用额外的 `{}` 保护，例如 `{P4}`、`{SDN}` 和 `{TT-INT}`。
- 页码范围使用两个连字符，例如 `pages = {123--135}`。
- `doi` 字段只填写 DOI，不要添加 `https://doi.org/`。
- 不确定的信息宁可暂时省略，不要用占位字符串填充。

机构作为作者时，用额外的一组花括号保护名称，避免被按人名拆分：

```bibtex
author = {{P4 Language Consortium}},
```

## 2 BibTeX

BibTeX 适合已经提供 `.bst` 样式文件的投稿模板。IEEE 论文通常使用 `IEEEtran.bst`。

### 2.1 文档写法

假设主文件为 `main.tex`，文献库为 `references.bib`：

```latex
\documentclass{article}

\begin{document}

In-band network telemetry can collect hop-level network state
\cite{zhang2026ttint}.

\bibliographystyle{IEEEtran}
\bibliography{references}

\end{document}
```

`\bibliographystyle` 指定 `.bst` 样式，`\bibliography` 填写文献库文件名，不带 `.bib` 后缀。

### 2.2 编译顺序

```sh
xelatex main
bibtex main
xelatex main
xelatex main
```

第一次 XeLaTeX 生成 `.aux`，BibTeX 读取引用、文献库和 `.bst`，生成 `.bbl`，后两次 XeLaTeX 更新参考文献和交叉引用。这里的 `main` 是主文档名，不是 `.bib` 文件名。

### 2.3 常用样式

| 样式       | 排序方式                             |
| ---------- | ------------------------------------ |
| `plain`    | 按作者姓名排序                       |
| `unsrt`    | 按正文中的引用顺序排列               |
| `alpha`    | 按作者等信息排序，使用作者与年份标签 |
| `IEEEtran` | 按 IEEE 格式和引用顺序排列           |

投稿模板已经指定样式时，不要自行替换。

## 3 `biblatex` 与 Biber

`biblatex` 在 LaTeX 层控制引用和参考文献格式，通常搭配 Biber。它更适合需要 UTF-8、多语言、筛选或多份参考文献列表的文档。

### 3.1 文档写法

```latex
\documentclass{article}
\usepackage{csquotes}
\usepackage[
  backend=biber,
  style=ieee,
  sorting=none
]{biblatex}

\addbibresource{references.bib}

\begin{document}

In-band network telemetry can collect hop-level network state
\cite{zhang2026ttint}.

\printbibliography

\end{document}
```

`style=ieee` 需要 `biblatex-ieee`。没有安装该样式时，可以先使用 `style=numeric`。

### 3.2 编译顺序

```sh
xelatex main
biber main
xelatex main
xelatex main
```

第一次 XeLaTeX 为 biblatex 生成 `.bcf`，Biber 根据它读取文献库并生成 `.bbl`。使用 `backend=biber` 时必须运行 `biber`，不能改成 `bibtex`。biblatex 的 `.bbl` 格式与传统 BibTeX 样式生成的格式不同，切换流程时不能直接沿用。

### 3.3 引用命令

| 命令               | 用途                         |
| ------------------ | ---------------------------- |
| `\cite{key}`       | 普通引用                     |
| `\cite{key1,key2}` | 同时引用多条文献             |
| `\parencite{key}`  | 括号式引用                   |
| `\textcite{key}`   | 把作者作为正文的一部分       |
| `\nocite{key}`     | 不在正文引用，但加入文献列表 |
| `\nocite{*}`       | 输出文献库中的所有条目       |

具体输出形式由所选样式决定。

## 4 两套流程如何选择

| 场景                                | 建议                    |
| ----------------------------------- | ----------------------- |
| 模板提供 `.bst` 或明确要求 BibTeX   | 使用 BibTeX             |
| IEEE 投稿并使用官方 `IEEEtran` 模板 | 使用 `IEEEtran.bst`     |
| 自己编写报告、学位论文或个人文档    | 优先考虑 `biblatex`     |
| 文献包含大量 Unicode 或多语言内容   | 使用 `biblatex` + Biber |

两套流程的对应关系如下：

| BibTeX                      | `biblatex` + Biber                |
| --------------------------- | --------------------------------- |
| `\bibliographystyle{...}`   | `style=...`                       |
| `\bibliography{references}` | `\addbibresource{references.bib}` |
| 文末自动插入参考文献        | `\printbibliography`              |
| 运行 `bibtex main`          | 运行 `biber main`                 |

## 5 常见问题

### 5.1 引用显示为 `[?]`

依次检查：

1. 引用键是否与 `.bib` 文件完全一致。
2. `.bib` 条目是否缺少逗号或右花括号。
3. 是否按正确顺序运行了 BibTeX 或 Biber。
4. 后端运行后是否再次执行了两遍 LaTeX。

### 5.2 参考文献列表为空

默认情况下，被引用的条目或通过 `\nocite` 指定的条目会进入文献列表，biblatex 的筛选选项还可能进一步限制结果。调试时可以临时加入 `\nocite{*}`，确认文献库能否被正常读取。

### 5.3 标题中的缩写变成小写

部分 BibTeX 样式会调整标题大小写。需要保留的内容应写成 `{P4}`、`{IPv6}` 或 `{OpenFlow}`。

### 5.4 Biber 找不到数据源

确认 `\addbibresource{references.bib}` 包含文件后缀，并从主文件所在目录运行 `biber main`。

### 5.5 修改后结果没有更新

先检查 BibTeX 或 Biber 的 `.blg` 日志，确认后端是否成功运行。切换文献工具后，旧的 `.aux`、`.bbl`、`.bcf` 或 `.run.xml` 可能仍被使用，可清理这些中间文件并重新执行完整流程。清理前确认 `.bbl` 是可重新生成的文件，有些投稿项目只提供 `.bbl`，不能直接删除。

### 5.6 DOI 或 URL 没有显示

参考文献样式决定哪些字段会出现在最终结果中。即使 `.bib` 已填写 `doi` 或 `url`，当前 `.bst` 或 `biblatex` 样式也可能省略它们。投稿时应服从模板样式，不要把 DOI 手工拼进标题字段。

### 5.7 投稿模板与 `biblatex` 冲突

部分期刊类文件已经固定了参考文献方案。遇到这种情况应使用模板提供的命令和 `.bst`，不要强行加载 `biblatex`。
