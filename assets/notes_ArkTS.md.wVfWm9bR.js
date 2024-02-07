import{_ as s,c as a,o as n,V as p}from"./chunks/framework.7l-q_CRe.js";const k=JSON.parse('{"title":"ArkTS","description":"","frontmatter":{},"headers":[],"relativePath":"notes/ArkTS.md","filePath":"notes/ArkTS.md"}'),e={name:"notes/ArkTS.md"},i=p(`<h1 id="arkts" tabindex="-1">ArkTS <a class="header-anchor" href="#arkts" aria-label="Permalink to &quot;ArkTS&quot;">​</a></h1><h2 id="自定义组件" tabindex="-1">自定义组件 <a class="header-anchor" href="#自定义组件" aria-label="Permalink to &quot;自定义组件&quot;">​</a></h2><p>定义名为 TitleComponent 的自定义组件：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Component  // 装饰器 @Component 表示这是个自定义组件</span></span>
<span class="line"><span>export struct TitleComponent {  // @Component 装饰的 struct 表示该结构体具有组件化能力，能够成为一个独立的组件</span></span>
<span class="line"><span>    build() {}  // 自定义组件必须定义 build 方法，在其中进行 UI 描述</span></span>
<span class="line"><span>}</span></span></code></pre></div><p>定义名为 RankPage 的自定义组件，在其中使用 TitleComponent 组件：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>import { TitleComponent } from &#39;../view/TitleComponent&#39;;</span></span>
<span class="line"><span>@Entry  // 装饰器 @Entry 表示这是个入口组件，加载页面时，将首先创建并呈现 @Entry 装饰的自定义组件，一个页面有且仅有一个 @Entry 装饰的组件</span></span>
<span class="line"><span>@Component</span></span>
<span class="line"><span>struct RankPage {</span></span>
<span class="line"><span>    build() {</span></span>
<span class="line"><span>        Column() {</span></span>
<span class="line"><span>            TitleComponent()</span></span>
<span class="line"><span>        }</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><h2 id="生命周期" tabindex="-1">生命周期 <a class="header-anchor" href="#生命周期" aria-label="Permalink to &quot;生命周期&quot;">​</a></h2><h3 id="自定义组件-1" tabindex="-1">自定义组件 <a class="header-anchor" href="#自定义组件-1" aria-label="Permalink to &quot;自定义组件&quot;">​</a></h3><ol><li>自定义组件创建</li><li>aboutToAppear()</li><li>aboutToDisappear()</li><li>自定义组件销毁</li></ol><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Component</span></span>
<span class="line"><span>struct RankPage {</span></span>
<span class="line highlighted"><span>    aboutToAppear() {}</span></span>
<span class="line highlighted"><span>    aboutToDisappear() {}</span></span>
<span class="line"><span>}</span></span></code></pre></div><h3 id="入口组件" tabindex="-1">入口组件 <a class="header-anchor" href="#入口组件" aria-label="Permalink to &quot;入口组件&quot;">​</a></h3><ol><li>页面入口组件创建</li><li>aboutToAppear()</li><li>onPageShow()，页面显示时自动调用</li><li>onBackPress()，用户点击返回按钮时自动调用</li><li>onPageHide()，页面消失时自动调用</li><li>aboutToDisappear()</li><li>页面入口组件销毁</li></ol><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Entry</span></span>
<span class="line"><span>@Component</span></span>
<span class="line"><span>struct RankPage {</span></span>
<span class="line highlighted"><span>    onPageShow() {}</span></span>
<span class="line highlighted"><span>    onPageHide() {}</span></span>
<span class="line highlighted"><span>    onBackPress() {</span></span>
<span class="line"><span>        // 返回值为 true 时，表示页面自己处理返回逻辑，不进行页面返回</span></span>
<span class="line"><span>        // 返回值为 false 时，表示由系统处理返回逻辑，默认为 false</span></span>
<span class="line"><span>        return true;</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><h2 id="渲染控制" tabindex="-1">渲染控制 <a class="header-anchor" href="#渲染控制" aria-label="Permalink to &quot;渲染控制&quot;">​</a></h2><p>条件渲染：<code>if ... else if ... else ...</code></p><p>循环渲染：使用 <code>ForEach</code> 迭代数组，并为每个数组项创建相应的组件</p><h2 id="状态管理" tabindex="-1">状态管理 <a class="header-anchor" href="#状态管理" aria-label="Permalink to &quot;状态管理&quot;">​</a></h2><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Component</span></span>
<span class="line"><span>export struct ListItemComponent {</span></span>
<span class="line"><span>    @State isChange: boolean = false;  // 装饰器 @State 表示组件中的状态变量，此状态变化会引起 UI 变更</span></span>
<span class="line"><span>    build() {</span></span>
<span class="line"><span>        Row() {</span></span>
<span class="line"><span>            Text(this.name).fontColor(this.isChange ? ItemStyle.COLOR_BLUE : ItemStyle.COLOR_BLACK)</span></span>
<span class="line"><span>        }.onClick(() =&gt; {</span></span>
<span class="line"><span>            this.isChange = !this.isChange;</span></span>
<span class="line"><span>        })</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Component</span></span>
<span class="line"><span>export struct TitleComponent {</span></span>
<span class="line"><span>    // @Link 装饰的变量可以和父组件的 @State 变量建立双向数据绑定。任何一方所做的修改都会反映给另一方。</span></span>
<span class="line"><span>    @Link isRefreshData: boolean;  // 没有对该变量初始化，需要父组件创建时该组件时赋值</span></span>
<span class="line"><span>}</span></span></code></pre></div><h2 id="应用程序入口——uiability" tabindex="-1">应用程序入口——UIAbility <a class="header-anchor" href="#应用程序入口——uiability" aria-label="Permalink to &quot;应用程序入口——UIAbility&quot;">​</a></h2><p>每一个UIAbility实例，都对应于一个最近任务列表中的任务。一个应用可以有一个UIAbility，也可以有多个UIAbility。</p><h3 id="uiability-内的页面跳转" tabindex="-1">UIAbility 内的页面跳转 <a class="header-anchor" href="#uiability-内的页面跳转" aria-label="Permalink to &quot;UIAbility 内的页面跳转&quot;">​</a></h3><p>先导入路由模块：<code>import router from &#39;@ohos.router&#39;</code></p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>// index页面</span></span>
<span class="line"><span>@Entry</span></span>
<span class="line"><span>@Component</span></span>
<span class="line"><span>struct Index {</span></span>
<span class="line"><span>    build() {</span></span>
<span class="line"><span>        Button(&#39;Next&#39;)</span></span>
<span class="line"><span>            .onClick(() =&gt; {</span></span>
<span class="line highlighted"><span>                router.pushUrl({</span></span>
<span class="line highlighted"><span>                    url: &#39;pages/second&#39;,  // 跳转</span></span>
<span class="line highlighted"><span>                    params: {  // 参数</span></span>
<span class="line highlighted"><span>                        &#39;msg&#39;: &#39;123&#39;</span></span>
<span class="line highlighted"><span>                    }</span></span>
<span class="line highlighted"><span>                })</span></span>
<span class="line"><span>            })</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>// second 页面</span></span>
<span class="line"><span>@Entry</span></span>
<span class="line"><span>@Component</span></span>
<span class="line"><span>struct Second {</span></span>
<span class="line highlighted"><span>    @State msg: string = router.getParams()[&#39;msg&#39;]  // 接收参数</span></span>
<span class="line"><span>    build() {</span></span>
<span class="line"><span>        Column() {</span></span>
<span class="line"><span>            Text(this.msg)</span></span>
<span class="line"><span>            Button(&#39;Back&#39;)</span></span>
<span class="line"><span>                .onClick(() =&gt; {</span></span>
<span class="line highlighted"><span>                    router.back()  // 返回上一页</span></span>
<span class="line"><span>                })</span></span>
<span class="line"><span>        }</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><h3 id="生命周期-1" tabindex="-1">生命周期 <a class="header-anchor" href="#生命周期-1" aria-label="Permalink to &quot;生命周期&quot;">​</a></h3><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>onCreate(want, launchParam) {}</span></span>
<span class="line"><span>onDestroy() {}</span></span>
<span class="line"><span>onWindowStageCreate(windowStage: window.WindowStage) {}</span></span>
<span class="line"><span>onWindowStageDestroy() {}</span></span>
<span class="line"><span>onForeground() {}</span></span>
<span class="line"><span>onBackground() {}</span></span></code></pre></div><h2 id="基础组件" tabindex="-1">基础组件 <a class="header-anchor" href="#基础组件" aria-label="Permalink to &quot;基础组件&quot;">​</a></h2><p>组件（Component）是界面搭建与显示的最小单位，ArkUI 提供了丰富多样的 UI 组件，我们可以使用这些组件轻松的编写出更加丰富、漂亮的界面。</p><p>组件根据功能可以分为以下五大类：基础组件、容器组件、媒体组件、绘制组件、画布组件。其中基础组件是视图层的基本组成单元，包括 <code>Text</code>、<code>Image</code>、<code>TextInput</code>、<code>Button</code>、<code>LoadingProgress</code> 等。</p><h3 id="image" tabindex="-1">Image <a class="header-anchor" href="#image" aria-label="Permalink to &quot;Image&quot;">​</a></h3><p>设置图片数据源：<code>Image(src: string|Resource)</code></p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>Image(&#39;https://xxxxxx&#39;)  // 使用string类型数据加载网络图片</span></span>
<span class="line"><span>Image($r(&#39;app.media.logo&#39;))  // 使用Resource加载图片</span></span></code></pre></div><p>其中加载网络图片需要在 <code>module.json5</code> 文件中添加网络访问权限：</p><div class="language-json5 vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">json5</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">{</span></span>
<span class="line"><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">    &quot;module&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> : {</span></span>
<span class="line"><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">        &quot;requestPermissions&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">:[</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">           {</span></span>
<span class="line"><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">             &quot;name&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">: </span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">&quot;ohos.permission.INTERNET&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">           }</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        ]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    }</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span></code></pre></div><p>设置图片大小：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>// 使用number数据设置图片大小</span></span>
<span class="line"><span>Image($r(&#39;app.media.icon&#39;))</span></span>
<span class="line"><span>    .width(80)</span></span>
<span class="line"><span>    .height(80)</span></span>
<span class="line"><span></span></span>
<span class="line"><span>// 使用string数据设置图片大小</span></span>
<span class="line"><span>Image($r(&#39;app.media.icon&#39;))</span></span>
<span class="line"><span>    .width(&#39;80vp&#39;)</span></span>
<span class="line"><span>    .height(&#39;80vp&#39;)</span></span>
<span class="line"><span>    </span></span>
<span class="line"><span>// 使用Resource数据设置图片大小</span></span>
<span class="line"><span>Image($r(&#39;app.media.icon&#39;))</span></span>
<span class="line"><span>    .width($r(&#39;app.float.logo_image_size&#39;))</span></span>
<span class="line"><span>    .height($r(&#39;app.float.logo_image_size&#39;))</span></span></code></pre></div><p>使用Resource数据设置图片大小，需要先在resources文件夹下定义：</p><div class="language-json vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">json</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">{</span></span>
<span class="line"><span style="--shiki-light:#005CC5;--shiki-dark:#79B8FF;">    &quot;float&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">: [</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        {</span></span>
<span class="line"><span style="--shiki-light:#005CC5;--shiki-dark:#79B8FF;">            &quot;name&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">: </span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">&quot;logo_image_size&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">,</span></span>
<span class="line"><span style="--shiki-light:#005CC5;--shiki-dark:#79B8FF;">            &quot;value&quot;</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">: </span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">&quot;80vp&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">        }</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">    ]</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span></code></pre></div><h3 id="text" tabindex="-1">Text <a class="header-anchor" href="#text" aria-label="Permalink to &quot;Text&quot;">​</a></h3><p>用于在节界面上展示一段文本信息</p><p>设置文本内容：<code>Text(content?: string|Resource)</code></p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>Text(&#39;登录界面&#39;)</span></span>
<span class="line"><span>Text($r(&#39;app.string.login_page&#39;))  // 需要在resources目录下的string.json文件中定义该资源</span></span></code></pre></div><p>设置文本大小：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>Text(&#39;abc&#39;).fontSize(24)</span></span>
<span class="line"><span>Text(&#39;abc&#39;).fontSize(&#39;24fp&#39;)</span></span>
<span class="line"><span>Text(&#39;abc&#39;).fontSize($r(&#39;app.float.page_title_size&#39;))</span></span></code></pre></div><p>设置文本粗细：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>Text(&#39;123&#39;).fontWeight(500)  // 默认400</span></span>
<span class="line"><span>Text(&#39;123&#39;).fontWeight(FontWeight.Medium)  // Lighter, Normal, Regular, Medium, Bold, Bolder</span></span></code></pre></div><p>设置文本颜色：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>Text(&#39;abc&#39;).fontColor(Color.Black)</span></span>
<span class="line"><span>Text(&#39;abc&#39;).fontColor(0x182431)</span></span>
<span class="line"><span>Text(&#39;abc&#39;).fontColor(&#39;#182431&#39;)</span></span>
<span class="line"><span>Text(&#39;abc&#39;).fontColor($r(&#39;app.color.my_color&#39;))</span></span></code></pre></div><h3 id="textinput" tabindex="-1">TextInput <a class="header-anchor" href="#textinput" aria-label="Permalink to &quot;TextInput&quot;">​</a></h3><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>TextInput({ placeholder: &quot;账号&quot; })  // 设置提示文本</span></span>
<span class="line"><span>    .placeholderColor(0x999999)     // 提示文本的颜色</span></span>
<span class="line"><span>    .fontColor(Color.Blue)          // 设置字体颜色</span></span>
<span class="line"><span>    .fontSize(20)                   // 设置字体大小</span></span>
<span class="line"><span>    .fontStyle(FontStyle.Italic)    // 设置斜体</span></span>
<span class="line"><span>    .fontWeight(FontWeight.Bold)    // 设置字体粗细</span></span>
<span class="line"><span>    .maxLength(8)                   // 设置最大输入字符数</span></span>
<span class="line"><span>    .type(InputType.Number)         // 设置文本框输入类型</span></span>
<span class="line"><span>    .onChange((value: string) =&gt; {  // 设置onChange事件</span></span>
<span class="line"><span>        console.log(value)</span></span>
<span class="line"><span>    })</span></span></code></pre></div><h3 id="button" tabindex="-1">Button <a class="header-anchor" href="#button" aria-label="Permalink to &quot;Button&quot;">​</a></h3><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>Button(&#39;登录&#39;, { type: ButtonType.Capsule, stateEffect: false })</span></span>
<span class="line"><span>    .width(&#39;90%&#39;)</span></span>
<span class="line"><span>    .height(40)</span></span>
<span class="line"><span>    .fontSize(16)</span></span>
<span class="line"><span>    .fontWeight(FontWeight.Medium)</span></span>
<span class="line"><span>    .backgroundColor(&#39;#007DFF&#39;)</span></span></code></pre></div><p>type用于定义按钮样式，stateEffect用于设置按钮按下时是否开启切换效果，当状态置为false时，点击效果关闭，默认值为true。</p><p>type的值为枚举类型 <code>ButtonType</code> 中的一种：<code>Capsule</code>(圆角按钮)、<code>Circle</code>(圆形按钮)、<code>Normal</code>(普通按钮、方形)</p><h2 id="容器组件" tabindex="-1">容器组件 <a class="header-anchor" href="#容器组件" aria-label="Permalink to &quot;容器组件&quot;">​</a></h2><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>// Column表示沿垂直方向布局的容器，space表示容器内组件间距</span></span>
<span class="line"><span>Column({space: 30}) {</span></span>
<span class="line"><span>    Text(&#39;abc&#39;)</span></span>
<span class="line"><span>    // Row表示沿水平方向布局的容器，space表示容器内组件间距</span></span>
<span class="line"><span>    Row({space: 20}) {</span></span>
<span class="line"><span>        Button(&#39;登录&#39;)</span></span>
<span class="line"><span>        Button(&#39;注册&#39;)</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><h3 id="对齐方式" tabindex="-1">对齐方式 <a class="header-anchor" href="#对齐方式" aria-label="Permalink to &quot;对齐方式&quot;">​</a></h3><table><thead><tr><th>属性</th><th>描述</th></tr></thead><tbody><tr><td><code>justifyContent</code></td><td>设置子组件在主轴方向上的对齐格式</td></tr><tr><td><code>alignItems</code></td><td>设置子组件在交叉轴方向上的对齐格式</td></tr></tbody></table><p>子组件在主轴方向上的对齐使用 <code>justifyContent()</code> 来设置，其参数为枚举类型 <code>FlexAlign</code>。<code>FlexAlign</code> 定义了以下几种类型：</p><ul><li><code>Start</code>：元素在主轴方向首端对齐，第一个元素与行首对齐，同时后续的元素与前一个对齐。</li><li><code>Center</code>：元素在主轴方向中心对齐，第一个元素与行首的距离以及最后一个元素与行尾距离相同。</li><li><code>End</code>：元素在主轴方向尾部对齐，最后一个元素与行尾对齐，其他元素与后一个对齐。</li><li><code>SpaceBetween</code>：元素在主轴方向均匀分配弹性元素，相邻元素之间距离相同。 第一个元素与行首对齐，最后一个元素与行尾对齐。</li><li><code>SpaceAround</code>：元素在主轴方向均匀分配弹性元素，相邻元素之间距离相同。 第一个元素到行首的距离和最后一个元素到行尾的距离是相邻元素之间距离的一半。</li><li><code>SpaceEvenly</code>：元素在主轴方向等间距布局，无论是相邻元素还是边界元素到容器的间距都一样。</li></ul><p>子组件在交叉轴方向上的对齐方式使用 <code>alignItems()</code> 来设置。</p><p><code>Column</code> 容器的主轴是垂直方向，交叉轴是水平方向，其参数为枚举类型 <code>HorizontalAlign</code>（水平对齐），<code>HorizontalAlign</code> 定义了以下几种类型：</p><ul><li><code>Start</code>：设置子组件在水平方向上按照起始端对齐。</li><li><code>Center</code>（默认值）：设置子组件在水平方向上居中对齐。</li><li><code>End</code>：设置子组件在水平方向上按照末端对齐。</li></ul><p><code>Row</code> 容器的主轴是水平方向，交叉轴是垂直方向，其参数为枚举类型 <code>VerticalAlign</code>（垂直对齐），<code>VerticalAlign</code> 定义了以下几种类型：</p><ul><li><code>Top</code>：设置子组件在垂直方向上居顶部对齐。</li><li><code>Center</code>（默认值）：设置子组件在竖直方向上居中对齐。</li><li><code>Bottom</code>：设置子组件在竖直方向上居底部对齐。</li></ul><h3 id="list" tabindex="-1"><code>List</code> <a class="header-anchor" href="#list" aria-label="Permalink to &quot;\`List\`&quot;">​</a></h3><p>List是很常用的滚动类容器组件，一般和子组件ListItem一起使用，List列表中的每一个列表项对应一个ListItem组件。</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Entry</span></span>
<span class="line"><span>@Component</span></span>
<span class="line"><span>struct ListDemo {</span></span>
<span class="line"><span>    private arr: number[] = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]</span></span>
<span class="line"><span>    build() {</span></span>
<span class="line"><span>        List({ space: 10 }) {</span></span>
<span class="line"><span>            ForEach(this.arr, (item: number) =&gt; {</span></span>
<span class="line"><span>                ListItem() {</span></span>
<span class="line"><span>                    Text(\`\${item}\`)</span></span>
<span class="line"><span>                        .width(&#39;100%&#39;)</span></span>
<span class="line"><span>                        .height(50)</span></span>
<span class="line"><span>                        .backgroundColor(0x007DFF)</span></span>
<span class="line"><span>                }</span></span>
<span class="line"><span>            })</span></span>
<span class="line"><span>        }</span></span>
<span class="line"><span>    }</span></span>
<span class="line"><span>}</span></span></code></pre></div><p><code>List</code> 组件子组件 <code>ListItem</code> 之间默认是没有分割线的，可以使用 <code>List</code> 组件的 <code>divider</code> 方法。<code>divider</code> 方法包含四个参数：</p><table><thead><tr><th>参数</th><th>说明</th></tr></thead><tbody><tr><td><code>strokeWidth</code></td><td>分割线的线宽</td></tr><tr><td><code>color</code></td><td>分割线的颜色</td></tr><tr><td><code>startMargin</code></td><td>分割线距离列表侧边起始端的距离</td></tr><tr><td><code>endMargin</code></td><td>分割线距离列表侧边结束端的距离</td></tr></tbody></table><p>List组件提供了一系列事件方法用来监听列表的滚动：</p><table><thead><tr><th>事件方法</th><th>说明</th></tr></thead><tbody><tr><td><code>onScroll</code></td><td>列表滑动时触发，返回值 <code>scrollOffset</code> 为滑动偏移量，<code>scrollState</code> 为当前滑动状态</td></tr><tr><td><code>onScrollIndex</code></td><td>列表滑动时触发，返回值分别为滑动起始位置索引值与滑动结束位置索引值</td></tr><tr><td><code>onReachStart</code></td><td>列表到达起始位置时触发</td></tr><tr><td><code>onReachEnd</code></td><td>列表到底末尾位置时触发</td></tr><tr><td><code>onScrollStop</code></td><td>列表滑动停止时触发</td></tr></tbody></table><p>List 组件里面的列表项默认是按垂直方向排列的，可以将 List 组件的 <code>listDirection</code> 属性设置为 <code>Axis.Horizontal</code> 实现列表项横向排列，默认为 <code>Axis.Vertical</code>。</p>`,74),t=[i];function l(o,c,d,h,r,g){return n(),a("div",null,t)}const b=s(e,[["render",l]]);export{k as __pageData,b as default};
