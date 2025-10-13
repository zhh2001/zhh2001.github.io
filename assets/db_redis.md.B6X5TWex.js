import{_ as s,c as e,o as n,ag as t}from"./chunks/framework.BZemHgQ6.js";const h=JSON.parse('{"title":"Redis","description":"","frontmatter":{"outline":[2,3]},"headers":[],"relativePath":"db/redis.md","filePath":"db/redis.md"}'),l={name:"db/redis.md"};function i(p,a,o,c,d,u){return n(),e("div",null,a[0]||(a[0]=[t(`<h1 id="redis" tabindex="-1">Redis <a class="header-anchor" href="#redis" aria-label="Permalink to &quot;Redis&quot;">​</a></h1><p>Redis 诞生于 2009 年，全称是 <span style="color:red;">Re</span>mote <span style="color:red;">Di</span>ctionary <span style="color:red;">S</span>erver。</p><h2 id="_1-通用命令" tabindex="-1">1 通用命令 <a class="header-anchor" href="#_1-通用命令" aria-label="Permalink to &quot;1 通用命令&quot;">​</a></h2><h3 id="_1-1-keys" tabindex="-1">1.1 <code>KEYS</code> <a class="header-anchor" href="#_1-1-keys" aria-label="Permalink to &quot;1.1 \`KEYS\`&quot;">​</a></h3><ul><li>语法：<code>KEYS pattern</code></li><li>功能：查看复合模版的所有 <code>key</code>，不建议在生产环境使用</li><li>时间复杂度：<code>O(N)</code>，其中 <code>N</code> 是数据库中的键数</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; MSET firstname Jack lastname Stuntman age 35</span></span>
<span class="line"><span>&quot;OK&quot;</span></span>
<span class="line"><span>redis&gt; KEYS *name*</span></span>
<span class="line"><span>1) &quot;firstname&quot;</span></span>
<span class="line"><span>2) &quot;lastname&quot;</span></span>
<span class="line"><span>redis&gt; KEYS a??</span></span>
<span class="line"><span>1) &quot;age&quot;</span></span>
<span class="line"><span>redis&gt; KEYS *</span></span>
<span class="line"><span>1) &quot;age&quot;</span></span>
<span class="line"><span>2) &quot;firstname&quot;</span></span>
<span class="line"><span>3) &quot;lastname&quot;</span></span></code></pre></div><h3 id="_1-2-del" tabindex="-1">1.2 <code>DEL</code> <a class="header-anchor" href="#_1-2-del" aria-label="Permalink to &quot;1.2 \`DEL\`&quot;">​</a></h3><ul><li>语法：<code>DEL key [key ...]</code></li><li>功能：删除指定的 <code>key</code>，如果 <code>key</code> 不存在则忽略</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SET key1 &quot;Hello&quot;</span></span>
<span class="line"><span>&quot;OK&quot;</span></span>
<span class="line"><span>redis&gt; SET key2 &quot;World&quot;</span></span>
<span class="line"><span>&quot;OK&quot;</span></span>
<span class="line"><span>redis&gt; DEL key1 key2 key3</span></span>
<span class="line"><span>(integer) 2</span></span></code></pre></div><h3 id="_1-3-exists" tabindex="-1">1.3 <code>EXISTS</code> <a class="header-anchor" href="#_1-3-exists" aria-label="Permalink to &quot;1.3 \`EXISTS\`&quot;">​</a></h3><ul><li>语法：<code>EXISTS key [key ...]</code></li><li>功能：判断指定的 <code>key</code> 是否存在</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SET key1 &quot;Hello&quot;</span></span>
<span class="line"><span>&quot;OK&quot;</span></span>
<span class="line"><span>redis&gt; EXISTS key1</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; EXISTS nosuchkey</span></span>
<span class="line"><span>(integer) 0</span></span>
<span class="line"><span>redis&gt; SET key2 &quot;World&quot;</span></span>
<span class="line"><span>&quot;OK&quot;</span></span>
<span class="line"><span>redis&gt; EXISTS key1 key2 nosuchkey</span></span>
<span class="line"><span>(integer) 2</span></span></code></pre></div><h3 id="_1-4-expire" tabindex="-1">1.4 <code>EXPIRE</code> <a class="header-anchor" href="#_1-4-expire" aria-label="Permalink to &quot;1.4 \`EXPIRE\`&quot;">​</a></h3><ul><li>语法：<code>EXPIRE key seconds</code></li><li>功能：设置 <code>key</code> 的过期时长</li></ul><h3 id="_1-5-ttl" tabindex="-1">1.5 <code>TTL</code> <a class="header-anchor" href="#_1-5-ttl" aria-label="Permalink to &quot;1.5 \`TTL\`&quot;">​</a></h3><ul><li>语法：<code>TTL key</code></li><li>功能：查看指定 <code>key</code> 的剩余有效时长（秒）</li></ul><p>如果 <code>key</code> 存在但是没有设置过期时长，返回 <code>-1</code>。如果 <code>key</code> 不存在返回 <code>-2</code>。</p><h2 id="_2-string-类型" tabindex="-1">2 String 类型 <a class="header-anchor" href="#_2-string-类型" aria-label="Permalink to &quot;2 String 类型&quot;">​</a></h2><p>String 类型，也就是字符串类型，是 Redis 中最简单的存储类型。</p><p>其 value 是字符串，不过根据字符串的格式不同，又可以分为 3 类：</p><ul><li>string：普通字符串</li><li>int：整数类型，可以进行自增自减</li><li>float：浮点类型，可以进行自增自减</li></ul><p>不管哪种格式，底层都是字节数组形式存储，只不过是编码方式不同。字符串类型的最大空间不能超过 512m</p><h3 id="_2-1-set" tabindex="-1">2.1 <code>SET</code> <a class="header-anchor" href="#_2-1-set" aria-label="Permalink to &quot;2.1 \`SET\`&quot;">​</a></h3><ul><li>语法：<code>SET key value [NX | XX] [EX seconds | KEEPTTL]</code></li><li>功能：添加或修改一个 String 类型的键值对</li><li>可选项： <ul><li><code>NX</code>：只有 <code>key</code> 不存在时才设置</li><li><code>XX</code>：只有 <code>key</code> 已存在时才设置</li><li><code>EX</code>：设置过期时长</li><li><code>KEEPTTL</code>：保留 <code>key</code> 原有的过期时长</li></ul></li></ul><h3 id="_2-2-get" tabindex="-1">2.2 <code>GET</code> <a class="header-anchor" href="#_2-2-get" aria-label="Permalink to &quot;2.2 \`GET\`&quot;">​</a></h3><ul><li>语法：<code>GET key</code></li><li>功能：根据 <code>key</code> 获取 String 类型的 <code>value</code></li></ul><h3 id="_2-3-mset" tabindex="-1">2.3 <code>MSET</code> <a class="header-anchor" href="#_2-3-mset" aria-label="Permalink to &quot;2.3 \`MSET\`&quot;">​</a></h3><ul><li>语法：<code>MSET key value [key value ...]</code></li><li>功能：批量添加多个 String 类型的键值对</li></ul><h3 id="_2-4-mget" tabindex="-1">2.4 <code>MGET</code> <a class="header-anchor" href="#_2-4-mget" aria-label="Permalink to &quot;2.4 \`MGET\`&quot;">​</a></h3><ul><li>语法：<code>MGET key [key ...]</code></li><li>功能：批量获取多个 <code>key</code> 的 <code>value</code></li></ul><h3 id="_2-5-incr" tabindex="-1">2.5 <code>INCR</code> <a class="header-anchor" href="#_2-5-incr" aria-label="Permalink to &quot;2.5 \`INCR\`&quot;">​</a></h3><ul><li>语法：<code>INCR key</code></li><li>功能：整型自增 <code>1</code></li></ul><h3 id="_2-6-incrby" tabindex="-1">2.6 <code>INCRBY</code> <a class="header-anchor" href="#_2-6-incrby" aria-label="Permalink to &quot;2.6 \`INCRBY\`&quot;">​</a></h3><ul><li>语法：<code>INCRBY key increment</code></li><li>功能：整型自增 <code>increment</code></li></ul><h3 id="_2-7-incrbyfloat" tabindex="-1">2.7 <code>INCRBYFLOAT</code> <a class="header-anchor" href="#_2-7-incrbyfloat" aria-label="Permalink to &quot;2.7 \`INCRBYFLOAT\`&quot;">​</a></h3><ul><li>语法：<code>INCRBYFLOAT key increment</code></li><li>功能：浮点型自增 <code>increment</code></li></ul><h3 id="_2-8-setnx" tabindex="-1">2.8 <code>SETNX</code> <a class="header-anchor" href="#_2-8-setnx" aria-label="Permalink to &quot;2.8 \`SETNX\`&quot;">​</a></h3><blockquote><p>弃用，推荐采用 <code>SET key value NX</code></p></blockquote><ul><li>语法：<code>SETNX key value</code></li><li>功能：如果 <code>key</code> 不存在才新增。</li></ul><p><code>SETNX</code> 是 <strong>SET</strong> if <strong>N</strong>ot e<strong>X</strong>ists 的简写。</p><h3 id="_2-9-setex" tabindex="-1">2.9 <code>SETEX</code> <a class="header-anchor" href="#_2-9-setex" aria-label="Permalink to &quot;2.9 \`SETEX\`&quot;">​</a></h3><blockquote><p>弃用，推荐采用 <code>SET key value EX seconds</code></p></blockquote><ul><li>语法：<code>SETEX key seconds value</code></li><li>功能：新增 <code>key</code> 并设置有效时长</li></ul><h2 id="_3-key-的层级格式" tabindex="-1">3 Key 的层级格式 <a class="header-anchor" href="#_3-key-的层级格式" aria-label="Permalink to &quot;3 Key 的层级格式&quot;">​</a></h2><p>Redis 的 <code>key</code> 允许有多个单词形成层级结构，多个单词间用 <code>:</code> 隔开。</p><h2 id="_4-hash-类型" tabindex="-1">4 Hash 类型 <a class="header-anchor" href="#_4-hash-类型" aria-label="Permalink to &quot;4 Hash 类型&quot;">​</a></h2><p>Hash 类型，也叫散列，其 <code>value</code> 是一个无需字典，类似于 Java 中的 HashMap 结构。</p><p>String 结构将对象的所有字段保存为一整个字符串，如果要修改其中某个字段很不方便。</p><p>Hash 结构可以将对象中的每次字段独立存储，可以针对单个字段做 CRUD。</p><h3 id="_4-1-hset" tabindex="-1">4.1 <code>HSET</code> <a class="header-anchor" href="#_4-1-hset" aria-label="Permalink to &quot;4.1 \`HSET\`&quot;">​</a></h3><ul><li>语法：<code>HSET key field value [field value ...]</code></li><li>功能：添加或修改 hash 类型 <code>key</code> 的 <code>filed</code> 的值</li></ul><h3 id="_4-2-hget" tabindex="-1">4.2 <code>HGET</code> <a class="header-anchor" href="#_4-2-hget" aria-label="Permalink to &quot;4.2 \`HGET\`&quot;">​</a></h3><ul><li>语法：<code>HGET key field</code></li><li>功能：获取一个 hash 类型 <code>key</code> 的 <code>filed</code> 的值</li></ul><h3 id="_4-3-hmset" tabindex="-1">4.3 <code>HMSET</code> <a class="header-anchor" href="#_4-3-hmset" aria-label="Permalink to &quot;4.3 \`HMSET\`&quot;">​</a></h3><blockquote><p>弃用，采用 <code>HSET</code> 效果一样</p></blockquote><ul><li>语法：<code>HMSET key field value [field value ...]</code></li><li>功能：批量添加多个 hash 类型 <code>key</code> 的 <code>filed</code> 的值</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HMSET myhash field1 &quot;Hello&quot; field2 &quot;World&quot;</span></span>
<span class="line"><span>&quot;OK&quot;</span></span>
<span class="line"><span>redis&gt; HGET myhash field1</span></span>
<span class="line"><span>&quot;Hello&quot;</span></span>
<span class="line"><span>redis&gt; HGET myhash field2</span></span>
<span class="line"><span>&quot;World&quot;</span></span></code></pre></div><h3 id="_4-4-hmget" tabindex="-1">4.4 <code>HMGET</code> <a class="header-anchor" href="#_4-4-hmget" aria-label="Permalink to &quot;4.4 \`HMGET\`&quot;">​</a></h3><ul><li>语法：<code>HMGET key field [field ...]</code></li><li>功能：批量获取多个 hash 类型 <code>key</code> 的 <code>filed</code> 的值</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSET myhash field1 &quot;Hello&quot; field2 &quot;World&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; HMGET myhash field1 field2 nofield</span></span>
<span class="line"><span>1) &quot;Hello&quot;</span></span>
<span class="line"><span>2) &quot;World&quot;</span></span>
<span class="line"><span>3) (nil)</span></span></code></pre></div><h3 id="_4-5-hgetall" tabindex="-1">4.5 <code>HGETALL</code> <a class="header-anchor" href="#_4-5-hgetall" aria-label="Permalink to &quot;4.5 \`HGETALL\`&quot;">​</a></h3><ul><li>语法：<code>HGETALL key</code></li><li>功能：获取一个 hash 类型的 <code>key</code> 中的所有 <code>filed</code> 和对应 <code>value</code></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSET myhash field1 &quot;Hello&quot; field2 &quot;World&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; HGETALL myhash</span></span>
<span class="line"><span>1) &quot;field1&quot;</span></span>
<span class="line"><span>2) &quot;Hello&quot;</span></span>
<span class="line"><span>3) &quot;field2&quot;</span></span>
<span class="line"><span>4) &quot;World&quot;</span></span></code></pre></div><h3 id="_4-6-hkeys" tabindex="-1">4.6 <code>HKEYS</code> <a class="header-anchor" href="#_4-6-hkeys" aria-label="Permalink to &quot;4.6 \`HKEYS\`&quot;">​</a></h3><ul><li>语法：<code>HKEYS key</code></li><li>功能：获取一个 hash 类型的 <code>key</code> 中的所有的 <code>filed</code></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSET myhash field1 &quot;Hello&quot; field2 &quot;World&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; HKEYS myhash</span></span>
<span class="line"><span>1) &quot;field1&quot;</span></span>
<span class="line"><span>2) &quot;field2&quot;</span></span></code></pre></div><h3 id="_4-7-hvals" tabindex="-1">4.7 <code>HVALS</code> <a class="header-anchor" href="#_4-7-hvals" aria-label="Permalink to &quot;4.7 \`HVALS\`&quot;">​</a></h3><ul><li>语法：<code>HVALS key</code></li><li>功能：获取一个 hash 类型的 <code>key</code> 中的所有的 <code>value</code></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSET myhash field1 &quot;Hello&quot; field2 &quot;World&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; HVALS myhash</span></span>
<span class="line"><span>1) &quot;Hello&quot;</span></span>
<span class="line"><span>2) &quot;World&quot;</span></span></code></pre></div><h3 id="_4-8-hincrby" tabindex="-1">4.8 <code>HINCRBY</code> <a class="header-anchor" href="#_4-8-hincrby" aria-label="Permalink to &quot;4.8 \`HINCRBY\`&quot;">​</a></h3><ul><li>语法：<code>HINCRBY key field increment</code></li><li>功能：让一个 hash 类型的 <code>key</code> 的字段值增加指定步长</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSET myhash field 5</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; HINCRBY myhash field 1</span></span>
<span class="line"><span>(integer) 6</span></span>
<span class="line"><span>redis&gt; HINCRBY myhash field -10</span></span>
<span class="line"><span>(integer) -4</span></span></code></pre></div><h3 id="_4-9-hincrby" tabindex="-1">4.9 <code>HINCRBY</code> <a class="header-anchor" href="#_4-9-hincrby" aria-label="Permalink to &quot;4.9 \`HINCRBY\`&quot;">​</a></h3><ul><li>语法：<code>HINCRBY key field increment</code></li><li>功能：让一个 hash 类型的 <code>key</code> 的字段值增加指定步长</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSET myhash field 5</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; HINCRBY myhash field 1</span></span>
<span class="line"><span>(integer) 6</span></span>
<span class="line"><span>redis&gt; HINCRBY myhash field -10</span></span>
<span class="line"><span>(integer) -4</span></span></code></pre></div><h3 id="_4-9-hsetnx" tabindex="-1">4.9 <code>HSETNX</code> <a class="header-anchor" href="#_4-9-hsetnx" aria-label="Permalink to &quot;4.9 \`HSETNX\`&quot;">​</a></h3><ul><li>语法：<code>HSETNX key field value</code></li><li>功能：只有这个 <code>key</code> 的字段不存在才能设置</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; HSETNX myhash field &quot;Hello&quot;</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; HSETNX myhash field &quot;World&quot;</span></span>
<span class="line"><span>(integer) 0</span></span>
<span class="line"><span>redis&gt; HGET myhash field</span></span>
<span class="line"><span>&quot;Hello&quot;</span></span></code></pre></div><h2 id="_5-list-类型" tabindex="-1">5 List 类型 <a class="header-anchor" href="#_5-list-类型" aria-label="Permalink to &quot;5 List 类型&quot;">​</a></h2><p>Redis 的 List 类型与 Java 的 LinkedList 类似，可以看作双向链表。既可以正向检索也可以反向检索。</p><h3 id="_5-1-lpush" tabindex="-1">5.1 <code>LPUSH</code> <a class="header-anchor" href="#_5-1-lpush" aria-label="Permalink to &quot;5.1 \`LPUSH\`&quot;">​</a></h3><ul><li>语法：<code>LPUSH key element [element ...]</code></li><li>功能：向列表左侧插入元素，<code>key</code> 不存在则会创建</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; LPUSH mylist &quot;world&quot;</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; LPUSH mylist &quot;hello&quot;</span></span>
<span class="line"><span>(integer) 2</span></span></code></pre></div><h3 id="_5-2-rpush" tabindex="-1">5.2 <code>RPUSH</code> <a class="header-anchor" href="#_5-2-rpush" aria-label="Permalink to &quot;5.2 \`RPUSH\`&quot;">​</a></h3><ul><li>语法：<code>RPUSH key element [element ...]</code></li><li>功能：向列表右侧插入元素，<code>key</code> 不存在则会创建</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; RPUSH mylist &quot;one&quot; &quot;two&quot; &quot;three&quot; &quot;four&quot; &quot;five&quot;</span></span>
<span class="line"><span>(integer) 5</span></span></code></pre></div><h3 id="_5-3-lpop" tabindex="-1">5.3 <code>LPOP</code> <a class="header-anchor" href="#_5-3-lpop" aria-label="Permalink to &quot;5.3 \`LPOP\`&quot;">​</a></h3><ul><li>语法：<code>LPOP key [count]</code></li><li>功能：从列表左侧移除元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; RPUSH mylist &quot;one&quot; &quot;two&quot; &quot;three&quot; &quot;four&quot; &quot;five&quot;</span></span>
<span class="line"><span>(integer) 5</span></span>
<span class="line"><span>redis&gt; LPOP mylist</span></span>
<span class="line"><span>&quot;one&quot;</span></span>
<span class="line"><span>redis&gt; LPOP mylist 2</span></span>
<span class="line"><span>1) &quot;two&quot;</span></span>
<span class="line"><span>2) &quot;three&quot;</span></span></code></pre></div><h3 id="_5-4-rpop" tabindex="-1">5.4 <code>RPOP</code> <a class="header-anchor" href="#_5-4-rpop" aria-label="Permalink to &quot;5.4 \`RPOP\`&quot;">​</a></h3><ul><li>语法：<code>RPOP key [count]</code></li><li>功能：从列表右侧移除元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; RPUSH mylist &quot;one&quot; &quot;two&quot; &quot;three&quot; &quot;four&quot; &quot;five&quot;</span></span>
<span class="line"><span>(integer) 5</span></span>
<span class="line"><span>redis&gt; RPOP mylist</span></span>
<span class="line"><span>&quot;five&quot;</span></span>
<span class="line"><span>redis&gt; RPOP mylist 2</span></span>
<span class="line"><span>1) &quot;four&quot;</span></span>
<span class="line"><span>2) &quot;three&quot;</span></span></code></pre></div><h3 id="_5-5-lrange" tabindex="-1">5.5 <code>LRANGE</code> <a class="header-anchor" href="#_5-5-lrange" aria-label="Permalink to &quot;5.5 \`LRANGE\`&quot;">​</a></h3><ul><li>语法：<code>LRANGE key start stop</code></li><li>功能：返回索引在 <code>[start stop]</code> 内的所有元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; RPUSH mylist &quot;one&quot; &quot;two&quot; &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; LRANGE mylist 0 0</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>redis&gt; LRANGE mylist -3 2</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;two&quot;</span></span>
<span class="line"><span>3) &quot;three&quot;</span></span>
<span class="line"><span>redis&gt; LRANGE mylist -100 100</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;two&quot;</span></span>
<span class="line"><span>3) &quot;three&quot;</span></span>
<span class="line"><span>redis&gt; LRANGE mylist 5 10</span></span>
<span class="line"><span>(empty array)</span></span></code></pre></div><h2 id="_6-set-类型" tabindex="-1">6 Set 类型 <a class="header-anchor" href="#_6-set-类型" aria-label="Permalink to &quot;6 Set 类型&quot;">​</a></h2><p>Redis 的 Set 结构与 Java 中的 HashSet 类似，可以看作是一个 <code>value</code> 为 <code>null</code> 的 HashMap。因为也是一个 hash 表，因此具备与 HashSet 类似的特征：</p><ul><li>无序</li><li>元素不重复</li><li>查找快</li><li>支持交集、并集、差集等功能</li></ul><h3 id="_6-1-sadd" tabindex="-1">6.1 <code>SADD</code> <a class="header-anchor" href="#_6-1-sadd" aria-label="Permalink to &quot;6.1 \`SADD\`&quot;">​</a></h3><ul><li>语法：<code>SADD key member [member ...]</code></li><li>功能：往集合中添加元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD myset &quot;Hello&quot; &quot;World&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; SADD myset &quot;World&quot;</span></span>
<span class="line"><span>(integer) 0</span></span></code></pre></div><h3 id="_6-2-srem" tabindex="-1">6.2 <code>SREM</code> <a class="header-anchor" href="#_6-2-srem" aria-label="Permalink to &quot;6.2 \`SREM\`&quot;">​</a></h3><ul><li>语法：<code>SREM key member [member ...]</code></li><li>功能：移除集合中的指定元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD myset &quot;one&quot; &quot;two&quot; &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; SREM myset &quot;one&quot; &quot;three&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; SREM myset &quot;four&quot;</span></span>
<span class="line"><span>(integer) 0</span></span></code></pre></div><h3 id="_6-3-scard" tabindex="-1">6.3 <code>SCARD</code> <a class="header-anchor" href="#_6-3-scard" aria-label="Permalink to &quot;6.3 \`SCARD\`&quot;">​</a></h3><ul><li>语法：<code>SCARD key</code></li><li>功能：返回集合中的元素数量</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD myset &quot;one&quot; &quot;two&quot; &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; SCARD myset</span></span>
<span class="line"><span>(integer) 3</span></span></code></pre></div><h3 id="_6-4-sismember" tabindex="-1">6.4 <code>SISMEMBER</code> <a class="header-anchor" href="#_6-4-sismember" aria-label="Permalink to &quot;6.4 \`SISMEMBER\`&quot;">​</a></h3><ul><li>语法：<code>SISMEMBER key member</code></li><li>功能：判断元素是否在集合中</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD myset &quot;one&quot;</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; SISMEMBER myset &quot;one&quot;</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; SISMEMBER myset &quot;two&quot;</span></span>
<span class="line"><span>(integer) 0</span></span></code></pre></div><h3 id="_6-5-smembers" tabindex="-1">6.5 <code>SMEMBERS</code> <a class="header-anchor" href="#_6-5-smembers" aria-label="Permalink to &quot;6.5 \`SMEMBERS\`&quot;">​</a></h3><ul><li>语法：<code>SMEMBERS key</code></li><li>功能：获取集合中的全部元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD myset Hello World</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; SMEMBERS myset</span></span>
<span class="line"><span>1) &quot;Hello&quot;</span></span>
<span class="line"><span>2) &quot;World&quot;</span></span></code></pre></div><h3 id="_6-6-sinter" tabindex="-1">6.6 <code>SINTER</code> <a class="header-anchor" href="#_6-6-sinter" aria-label="Permalink to &quot;6.6 \`SINTER\`&quot;">​</a></h3><ul><li>语法：<code>SINTER key [key ...]</code></li><li>功能：求交集（intersection）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD s1 a b c d</span></span>
<span class="line"><span>(integer) 4</span></span>
<span class="line"><span>redis&gt; SADD s2 c</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; SADD s2 a c e</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; SINTER s1 s2 s3</span></span>
<span class="line"><span>1) &quot;c&quot;</span></span></code></pre></div><h3 id="_6-7-sdiff" tabindex="-1">6.7 <code>SDIFF</code> <a class="header-anchor" href="#_6-7-sdiff" aria-label="Permalink to &quot;6.7 \`SDIFF\`&quot;">​</a></h3><ul><li>语法：<code>SDIFF key [key ...]</code></li><li>功能：求差集（difference set）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD s1 a b c d</span></span>
<span class="line"><span>(integer) 4</span></span>
<span class="line"><span>redis&gt; SADD s2 c</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; SADD s2 a c e</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; SDIFF s1 s2 s3</span></span>
<span class="line"><span>1) &quot;d&quot;</span></span>
<span class="line"><span>2) &quot;b&quot;</span></span></code></pre></div><h3 id="_6-8-sunion" tabindex="-1">6.8 <code>SUNION</code> <a class="header-anchor" href="#_6-8-sunion" aria-label="Permalink to &quot;6.8 \`SUNION\`&quot;">​</a></h3><ul><li>语法：<code>SUNION key [key ...]</code></li><li>功能：求并集（union）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; SADD s1 a b c d</span></span>
<span class="line"><span>(integer) 4</span></span>
<span class="line"><span>redis&gt; SADD s2 c</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; SADD s2 a c e</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; SUNION s1 s2 s3</span></span>
<span class="line"><span>1) &quot;c&quot;</span></span>
<span class="line"><span>2) &quot;e&quot;</span></span>
<span class="line"><span>3) &quot;b&quot;</span></span>
<span class="line"><span>4) &quot;d&quot;</span></span>
<span class="line"><span>5) &quot;a&quot;</span></span></code></pre></div><h2 id="_7-sortedset-类型" tabindex="-1">7 SortedSet 类型 <a class="header-anchor" href="#_7-sortedset-类型" aria-label="Permalink to &quot;7 SortedSet 类型&quot;">​</a></h2><p>Redis 的 SortedSet 是一个可排序的集合，与 Java 中的 TreeSet 有些类似，但底层数据结构差别很大。SortedSet 中每个元素都带有一个 score 属性，可以基于 score 属性对元素排序，底层实现是一个跳表（SkipList）加 hash 表。</p><p>SortedSet 具备下列特性：</p><ul><li>可排序</li><li>元素不重复</li><li>查询速度快</li></ul><p>因为 SortedSet 的可排序特性，经常被用来实现排行榜这样的功能。</p><h3 id="_7-1-zadd" tabindex="-1">7.1 <code>ZADD</code> <a class="header-anchor" href="#_7-1-zadd" aria-label="Permalink to &quot;7.1 \`ZADD\`&quot;">​</a></h3><ul><li>语法：<code>ZADD key [NX | XX] score member [score member ...]</code></li><li>功能：添加元素到有序集合，如果已存在则更新 <code>score</code></li><li>可选项： <ul><li><code>XX</code>：仅更新已存在的元素。不添加新元素。</li><li><code>NX</code>：只添加新元素。不更新现有元素。</li></ul></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD myzset 1 one 1 uno 2 two 3 three</span></span>
<span class="line"><span>(integer) 4</span></span></code></pre></div><h3 id="_7-2-zrem" tabindex="-1">7.2 <code>ZREM</code> <a class="header-anchor" href="#_7-2-zrem" aria-label="Permalink to &quot;7.2 \`ZREM\`&quot;">​</a></h3><ul><li>语法：<code>ZREM key member [member ...]</code></li><li>功能：删除有序集合中的指定元素</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD myzset 1 one 1 uno 2 two 3 three</span></span>
<span class="line"><span>(integer) 4</span></span>
<span class="line"><span>redis&gt; ZREM myzset two</span></span>
<span class="line"><span>(integer) 1</span></span></code></pre></div><h3 id="_7-3-zscore" tabindex="-1">7.3 <code>ZSCORE</code> <a class="header-anchor" href="#_7-3-zscore" aria-label="Permalink to &quot;7.3 \`ZSCORE\`&quot;">​</a></h3><ul><li>语法：<code>ZSCORE key member</code></li><li>功能：查询指定元素的 <code>score</code></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD myzset 1 one</span></span>
<span class="line"><span>(integer) 1</span></span>
<span class="line"><span>redis&gt; ZSCORE myzset one</span></span>
<span class="line"><span>&quot;1&quot;</span></span></code></pre></div><h3 id="_7-4-zrank" tabindex="-1">7.4 <code>ZRANK</code> <a class="header-anchor" href="#_7-4-zrank" aria-label="Permalink to &quot;7.4 \`ZRANK\`&quot;">​</a></h3><ul><li>语法：<code>ZRANK key member</code></li><li>功能：获取指定元素在有序集合中的排名</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>127.0.0.1:6379&gt; ZADD z 1 one 2 two 3 three</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZRANK z three</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZRANK z four</span></span>
<span class="line"><span>(nil)</span></span></code></pre></div><h3 id="_7-5-zcard" tabindex="-1">7.5 <code>ZCARD</code> <a class="header-anchor" href="#_7-5-zcard" aria-label="Permalink to &quot;7.5 \`ZCARD\`&quot;">​</a></h3><ul><li>语法：<code>ZCARD key</code></li><li>功能：获取有序集合中的元素数量</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>127.0.0.1:6379&gt; ZADD z 1 one 2 two 3 three</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZCARD z</span></span>
<span class="line"><span>(integer) 3</span></span></code></pre></div><h3 id="_7-6-zcount" tabindex="-1">7.6 <code>ZCOUNT</code> <a class="header-anchor" href="#_7-6-zcount" aria-label="Permalink to &quot;7.6 \`ZCOUNT\`&quot;">​</a></h3><ul><li>语法：<code>ZCOUNT key min max</code></li><li>功能：获取有序集合中的 <code>score</code> 在 <code>[min, max]</code> 内的元素数量</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>127.0.0.1:6379&gt; ZADD z 1 one 2 two 3 three</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZCOUNT z -inf +inf</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZCOUNT z 2 3</span></span>
<span class="line"><span>(integer) 2</span></span></code></pre></div><h3 id="_7-7-zincrby" tabindex="-1">7.7 <code>ZINCRBY</code> <a class="header-anchor" href="#_7-7-zincrby" aria-label="Permalink to &quot;7.7 \`ZINCRBY\`&quot;">​</a></h3><ul><li>语法：<code>ZINCRBY key increment member</code></li><li>功能：让有序集合中指定元素的 <code>score</code> 自增，步长为 <code>increment</code></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD myzset 1 &quot;one&quot; 2 &quot;two&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; ZINCRBY myzset 2 &quot;one&quot;</span></span>
<span class="line"><span>&quot;3&quot;</span></span></code></pre></div><h3 id="_7-8-zrange" tabindex="-1">7.8 <code>ZRANGE</code> <a class="header-anchor" href="#_7-8-zrange" aria-label="Permalink to &quot;7.8 \`ZRANGE\`&quot;">​</a></h3><ul><li>语法：<code>ZRANGE key start stop [BYSCORE | BYLEX] [REV] [WITHSCORES]</code></li><li>功能：按照 <code>score</code> 升序排序后，获取指定排名范围内的元素，排名从 0 开始</li><li>参数： <ul><li><code>REV</code>：是否降序顺序</li><li><code>WITHSCORES</code>：是否返回 <code>score</code></li></ul></li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>127.0.0.1:6379&gt; ZADD z 1 one 2 two 3 three</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZRANGE z 1 99</span></span>
<span class="line"><span>1) &quot;two&quot;</span></span>
<span class="line"><span>2) &quot;three&quot;</span></span>
<span class="line"><span>127.0.0.1:6379&gt; ZRANGE z 0 1 WITHSCORES</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;1&quot;</span></span>
<span class="line"><span>3) &quot;two&quot;</span></span>
<span class="line"><span>4) &quot;2&quot;</span></span></code></pre></div><h3 id="_7-9-zrangebyscore" tabindex="-1">7.9 <code>ZRANGEBYSCORE</code> <a class="header-anchor" href="#_7-9-zrangebyscore" aria-label="Permalink to &quot;7.9 \`ZRANGEBYSCORE\`&quot;">​</a></h3><blockquote><p>弃用，推荐 <code>ZRANGE</code></p></blockquote><ul><li>语法：<code>ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]</code></li><li>功能：按照 <code>score</code> 排序后，获取指定 <code>score</code> 范围内的元素</li></ul><h3 id="_7-10-zdiff" tabindex="-1">7.10 <code>ZDIFF</code> <a class="header-anchor" href="#_7-10-zdiff" aria-label="Permalink to &quot;7.10 \`ZDIFF\`&quot;">​</a></h3><ul><li>语法：<code>ZDIFF numkeys key [key ...] [WITHSCORES]</code></li><li>功能：求差集（difference set）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD zset1 1 &quot;one&quot; 2 &quot;two&quot; 3 &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; ZADD zset2 1 &quot;one&quot; 2 &quot;two&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; ZDIFF 2 zset1 zset2</span></span>
<span class="line"><span>1) &quot;three&quot;</span></span>
<span class="line"><span>redis&gt; ZDIFF 2 zset1 zset2 WITHSCORES</span></span>
<span class="line"><span>1) &quot;three&quot;</span></span>
<span class="line"><span>2) &quot;3&quot;</span></span></code></pre></div><h3 id="_7-11-zdiff" tabindex="-1">7.11 <code>ZDIFF</code> <a class="header-anchor" href="#_7-11-zdiff" aria-label="Permalink to &quot;7.11 \`ZDIFF\`&quot;">​</a></h3><ul><li>语法：<code>ZDIFF numkeys key [key ...] [WITHSCORES]</code></li><li>功能：求差集（difference set）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD zset1 1 &quot;one&quot; 2 &quot;two&quot; 3 &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; ZADD zset2 1 &quot;one&quot; 2 &quot;two&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; ZDIFF 2 zset1 zset2</span></span>
<span class="line"><span>1) &quot;three&quot;</span></span>
<span class="line"><span>redis&gt; ZDIFF 2 zset1 zset2 WITHSCORES</span></span>
<span class="line"><span>1) &quot;three&quot;</span></span>
<span class="line"><span>2) &quot;3&quot;</span></span></code></pre></div><h3 id="_7-12-zinter" tabindex="-1">7.12 <code>ZINTER</code> <a class="header-anchor" href="#_7-12-zinter" aria-label="Permalink to &quot;7.12 \`ZINTER\`&quot;">​</a></h3><ul><li>语法：<code>ZINTER numkeys key [key ...] [WITHSCORES]</code></li><li>功能：求交集（intersection）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD zset1 1 &quot;one&quot; 2 &quot;two&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; ZADD zset2 1 &quot;one&quot; 2 &quot;two&quot; 3 &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; ZINTER 2 zset1 zset2</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;two&quot;</span></span>
<span class="line"><span>redis&gt; ZINTER 2 zset1 zset2 WITHSCORES</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;2&quot;</span></span>
<span class="line"><span>3) &quot;two&quot;</span></span>
<span class="line"><span>4) &quot;4&quot;</span></span></code></pre></div><h3 id="_7-13-zunion" tabindex="-1">7.13 <code>ZUNION</code> <a class="header-anchor" href="#_7-13-zunion" aria-label="Permalink to &quot;7.13 \`ZUNION\`&quot;">​</a></h3><ul><li>语法：<code>ZUNION numkeys key [key ...] [WITHSCORES]</code></li><li>功能：求并集（union）</li></ul><div class="language-txt vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">txt</span><pre class="shiki shiki-themes github-light github-dark vp-code" tabindex="0"><code><span class="line"><span>redis&gt; ZADD zset1 1 &quot;one&quot; 2 &quot;two&quot;</span></span>
<span class="line"><span>(integer) 2</span></span>
<span class="line"><span>redis&gt; ZADD zset2 1 &quot;one&quot; 2 &quot;two&quot; 3 &quot;three&quot;</span></span>
<span class="line"><span>(integer) 3</span></span>
<span class="line"><span>redis&gt; ZUNION 2 zset1 zset2</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;three&quot;</span></span>
<span class="line"><span>3) &quot;two&quot;</span></span>
<span class="line"><span>redis&gt; ZUNION 2 zset1 zset2 WITHSCORES</span></span>
<span class="line"><span>1) &quot;one&quot;</span></span>
<span class="line"><span>2) &quot;2&quot;</span></span>
<span class="line"><span>3) &quot;three&quot;</span></span>
<span class="line"><span>4) &quot;3&quot;</span></span>
<span class="line"><span>5) &quot;two&quot;</span></span>
<span class="line"><span>6) &quot;4&quot;</span></span></code></pre></div>`,166)]))}const q=s(l,[["render",i]]);export{h as __pageData,q as default};
