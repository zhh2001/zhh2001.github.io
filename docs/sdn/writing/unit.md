---
outline: deep
---

# LaTeX 数值与单位

技术文档中的数值和单位可以使用 `siunitx` 统一排版。它会处理数字分组、科学计数法、数值与单位之间的间距，以及复合单位的字体和指数。

## 1 基础用法

### 1.1 加载宏包

```latex
\usepackage{siunitx}
```

`siunitx` 3 的三个基础命令分别处理数值、单位和物理量：

```latex
\num{12345.678}
\unit{\metre\per\second}
\qty{12.5}{\metre\per\second}
```

其中，`\qty` 同时接收数值和单位。一般不要写成 `12.5 m/s`，也不必手工添加 `\,` 或 `\mathrm`。

### 1.2 常用命令

| 命令                             | 用途                 |
| -------------------------------- | -------------------- |
| `\num{3.14}`                     | 排版数值             |
| `\unit{\kilo\metre}`             | 只排版单位           |
| `\qty{3.14}{\kilo\metre}`        | 排版数值和单位       |
| `\numrange{1}{5}`                | 数值范围             |
| `\qtyrange{1}{5}{\milli\second}` | 带单位的范围         |
| `\numlist{1;2;3}`                | 数值列表             |
| `\qtylist{1;2;3}{\milli\second}` | 带单位的数值列表     |
| `\ang{30;15;20}`                 | 角度                 |
| `\duration{1;20;30}`             | 时、分、秒形式的时长 |

## 2 数值排版

### 2.1 数字分组

```latex
\num{1234567.89}
```

较长的整数和小数会自动分组。源文件中可以始终使用小数点，输出形式再通过选项统一控制。

### 2.2 科学计数法

```latex
\num{6.022e23}
\qty{1.25e9}{\bit\per\second}
```

输入中的 `e` 表示十的幂，输出时会排成乘方形式。

### 2.3 不确定度

括号中的数字表示末尾若干位的不确定度：

```latex
\num{1.234(5)}
\qty{12.30(12)}{\milli\second}
```

需要输出为正负形式时，可以局部设置：

```latex
\num[uncertainty-mode = separate]{1.234(5)}
```

### 2.4 数值舍入

```latex
\num[
  round-mode = places,
  round-precision = 2
]{3.1415926}
```

舍入选项适合统一展示精度，但原始数据仍应保留完整数值。表格中的精度要求不一致时，优先在单个命令或单列上设置。

## 3 单位与物理量

### 3.1 SI 单位和前缀

```latex
\qty{2.4}{\giga\hertz}
\qty{15}{\micro\second}
\qty{25}{\degreeCelsius}
\qty{80}{\percent}
```

前缀和单位分开书写，例如 `\milli\second`、`\kilo\metre` 和 `\mega\hertz`。这样可以在不修改数据的情况下统一调整格式。

### 3.2 复合单位

```latex
\qty{9.81}{\metre\per\square\second}
\qty{1.2}{\joule\per\mole\per\kelvin}
\qty{25}{\kilo\bit\per\second}
```

`\per` 只作用于它后面的单位。平方和立方可以使用 `\square`、`\cubic`，也可以把 `\squared`、`\cubed` 写在单位后面。

### 3.3 除号形式

默认输出通常使用负指数。需要斜线时，可以全局或局部设置：

```latex
\sisetup{per-mode = symbol}

\qty{100}{\mega\bit\per\second}
\qty[per-mode = fraction]{8}{\joule\per\mole}
```

一条单位中有多个分母时，`siunitx` 会按所选模式添加必要的括号或分式结构。

## 4 范围、列表与角度

### 4.1 数值范围

```latex
\numrange{10}{20}
\qtyrange{5}{10}{\milli\second}
```

使用范围命令可以保证两个端点采用相同格式，也便于统一修改连接词和单位显示方式。

### 4.2 数值列表

```latex
\numlist{10;20;30}
\qtylist{1.2;1.5;1.8}{\giga\bit\per\second}
```

各项使用分号分隔。不要用逗号分隔列表，因为逗号也可能是数值的小数标记。

### 4.3 角度与时长

```latex
\ang{45}
\ang{30;15;20}
\duration{1;20;30}
```

`\ang{30;15;20}` 依次表示度、分、秒，`\duration{1;20;30}` 依次表示小时、分钟、秒。

## 5 网络相关单位

### 5.1 bit、byte 与速率

```latex
\qty{100}{\mega\bit\per\second}
\qty{12.5}{\mega\byte\per\second}
\qty{2.4}{\giga\hertz}
\qty{850}{\micro\second}
```

bit 和 byte 的符号不同，不能只靠大小写习惯在正文中手工输入。使用 `\bit` 和 `\byte` 可以避免同一篇文档中出现多套写法。

### 5.2 十进制与二进制前缀

```latex
\qty{1}{\giga\byte}
\qty{1}{\gibi\byte}
```

`\giga` 表示十进制前缀 G，`\gibi` 表示二进制前缀 Gi。链路速率通常使用十进制前缀，内存或文件大小则应根据数据来源明确选择，不能把 GB 和 GiB 混用。

### 5.3 自定义 packet 单位

`siunitx` 没有预定义所有领域单位。可以在导言区声明 packet：

```latex
\DeclareSIUnit\packet{packet}

\qty{1488095}{\packet\per\second}
\qty{64}{\byte\per\packet}
```

自定义单位集中放在导言区，正文只使用命令，后续修改缩写时不需要逐处替换。

## 6 全局与局部设置

### 6.1 全局设置

```latex
\sisetup{
  per-mode = symbol,
  group-minimum-digits = 4
}
```

`\sisetup` 适合放置整篇文档共用的规则。期刊或学校模板已有明确要求时，应在这里统一配置。

### 6.2 局部覆盖

```latex
\qty{10}{\metre\per\second}
\qty[per-mode = fraction]{10}{\metre\per\second}
```

方括号中的选项只影响当前命令，适合少量例外，不宜用来替代全局规范。

### 6.3 小数标记

```latex
\num[output-decimal-marker = {,}]{12.5}
```

输入和输出可以分开控制。源文件中的数据保持统一，展示时再按出版规范选择小数点或小数逗号。

## 7 表格中的数值

`siunitx` 提供 `S` 列，用于按小数点对齐：

```latex
\begin{tabular}{l S[table-format = 4.2]}
\toprule
Method & {Delay (\unit{\micro\second})} \\
\midrule
A & 8.5 \\
B & 125.37 \\
C & 1024.00 \\
\bottomrule
\end{tabular}
```

表头不是纯数值，需要用花括号保护。`table-format = 4.2` 表示整数部分最多四位、小数部分两位，合理声明格式可以让列宽更稳定。

## 8 完整示例

```latex
\documentclass[UTF8]{ctexart}
\usepackage{booktabs}
\usepackage{siunitx}

\DeclareSIUnit\packet{packet}

\sisetup{
  per-mode = symbol,
  group-minimum-digits = 4
}

\begin{document}

测试链路带宽为 \qty{10}{\giga\bit\per\second}，
测得平均时延为 \qty{23.48(16)}{\micro\second}。
实验持续 \duration{0;15;30}，负载范围为
\qtyrange{1}{8}{\giga\bit\per\second}。

\begin{table}[htbp]
  \centering
  \caption{不同负载下的转发性能}
  \label{tab:forwarding-performance}
  \begin{tabular}{
    S[table-format = 1.0]
    S[table-format = 4.2]
    S[table-format = 8.0]
  }
    \toprule
    {Load (\unit{\giga\bit\per\second})} &
    {Delay (\unit{\micro\second})} &
    {Throughput (\unit{\packet\per\second})} \\
    \midrule
    1 & 8.42  & 1488100 \\
    4 & 12.67 & 5952400 \\
    8 & 23.48 & 11904800 \\
    \bottomrule
  \end{tabular}
\end{table}

\end{document}
```

## 9 常见问题

### 9.1 手工控制间距

不要在不同位置混用 `10m`、`10~m` 和 `$10\,\mathrm{m}$`。把数值和单位交给 `\qty`，才能保证正文、公式和表格中的格式一致。

### 9.2 旧版命令

旧文档中常见 `\SI`、`\si`、`\SIrange` 和 `\SIlist`。`siunitx` 3 仍提供兼容接口，新文档应分别使用 `\qty`、`\unit`、`\qtyrange` 和 `\qtylist`。

### 9.3 摄氏度写法

不要把摄氏度拆成度符号和字母 C：

```latex
\qty{25}{\degreeCelsius}
```

### 9.4 单位写进数值参数

下面的写法会绕过 `siunitx` 对单位的处理：

```latex
% 不推荐
\num{10 Mbps}

% 推荐
\qty{10}{\mega\bit\per\second}
```

### 9.5 模板兼容

投稿模板可能预设数值格式或使用旧版 `siunitx`。出现未知选项时，先确认编译环境中的宏包版本，不要同时加载多个提供相似单位命令的宏包。
