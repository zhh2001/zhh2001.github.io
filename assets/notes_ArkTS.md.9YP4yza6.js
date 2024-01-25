import{_ as a,c as s,o as n,V as p}from"./chunks/framework.7l-q_CRe.js";const b=JSON.parse('{"title":"ArkTS","description":"","frontmatter":{},"headers":[],"relativePath":"notes/ArkTS.md","filePath":"notes/ArkTS.md"}'),e={name:"notes/ArkTS.md"},l=p(`<h1 id="arkts" tabindex="-1">ArkTS <a class="header-anchor" href="#arkts" aria-label="Permalink to &quot;ArkTS&quot;">​</a></h1><h2 id="自定义组件" tabindex="-1">自定义组件 <a class="header-anchor" href="#自定义组件" aria-label="Permalink to &quot;自定义组件&quot;">​</a></h2><p>定义名为 TitleComponent 的自定义组件：</p><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>@Component  // 装饰器 @Component 表示这是个自定义组件</span></span>
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
<span class="line"><span>}</span></span></code></pre></div><h3 id="声明周期" tabindex="-1">声明周期 <a class="header-anchor" href="#声明周期" aria-label="Permalink to &quot;声明周期&quot;">​</a></h3><div class="language-ArkTS vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">ArkTS</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span>onCreate(want, launchParam) {}</span></span>
<span class="line"><span>onDestroy() {}</span></span>
<span class="line"><span>onWindowStageCreate(windowStage: window.WindowStage) {}</span></span>
<span class="line"><span>onWindowStageDestroy() {}</span></span>
<span class="line"><span>onForeground() {}</span></span>
<span class="line"><span>onBackground() {}</span></span></code></pre></div><h2 id="基础组件" tabindex="-1">基础组件 <a class="header-anchor" href="#基础组件" aria-label="Permalink to &quot;基础组件&quot;">​</a></h2><p>组件（Component）是界面搭建与显示的最小单位，HarmonyOS ArkUI声明式开发范式为开发者提供了丰富多样的UI组件，我们可以使用这些组件轻松的编写出更加丰富、漂亮的界面。</p><p>组件根据功能可以分为以下五大类：基础组件、容器组件、媒体组件、绘制组件、画布组件。其中基础组件是视图层的基本组成单元，包括<code>Text</code>、<code>Image</code>、<code>TextInput</code>、<code>Button</code>、<code>LoadingProgress</code> 等。</p>`,30),i=[l];function t(o,c,r,h,d,u){return n(),s("div",null,i)}const k=a(e,[["render",t]]);export{b as __pageData,k as default};
