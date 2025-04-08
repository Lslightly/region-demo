# Detailed Design解读

...

## region和arena的比较


回答 http://222.195.92.204:1480/vm/golang/work/go1.23.1/-/issues/2#note_2165

[Comparison with the arenas](https://go.googlesource.com/proposal/+/refs/heads/master/design/70257-memory-regions.md#comparison-with-the-arenas)

### implicit allocation 而不是显式的 allocation

只需要通过 region.Do 就可以在运行时自动做区域内存分配替换。而 arena 需要显式的 API 调用。

对于库函数，如果要支持 arena 的话，就需要一个新版本的库来支持 arena。

有人可能会想说让 arena implicit 以及 goroutine-bound，但是隐式分配内存到 arena 的问题在于不知道什么时候可以安全地释放 arena。

相反，region 的设计将特定的内存分配在 region 中，且防止不知道什么时候释放内存的问题。新的 write barrier 设计将无法安全释放的内存给 pin 住不释放。而其他内存就可以在 region 结束的时候释放。

region.Ignore 则提供了一种细粒度控制内存分配方式的机制。

arean 引入的对逃逸分析的挑战在于：
- 如果采用显式 arena 的话，逃逸分析会出问题。逃逸分析可以判断一个变量是栈分配的，但是如果这个变量实际是分配在 arena 中的，那么就会出问题。
- 特殊的逃逸分析支持可以让 arena-allocated variables 是栈分配的候选者，但是就会比较复杂。

> 这里还是存疑的。如果 arena 总是返回一个堆地址，或者让变量逃逸就可以了，在运行时实际对象分配在哪里由 arena 操作。

### 细粒度的内存复用

arena 实验的一个特征是 arena 非常大(8MiB)。这么大是为了摊还 syscall 来强制内存错误的代价。

为了能够用好 arena，需要将 arena 在多个逻辑操作之间共享（比如说多个 http 情况）。但是在实际中，arena 实验只复用 partially-filled arena chunks，在未来将他们填满。但是，这些 partially-filled chunks 仍然对 live-heap 有贡献，虽然其内部的大部分内存已经死了。

> 这个想表达的含义应该是说 arena 的大块内存利用率不高，程序员手动管理 arena 的效率也不高，（静态分析的效率当然估计也不高了，这个还是需要很多运行时数据的）。所以这样的话实际上就还是在运行时进行细粒度管理的效率可能会好一点。这应该也是为什么 region 中会采用 block line 进行细粒度管理的原因。

region blocks 相比于 arena chunks 要小三个数量级。这意味着没有必要复用 partially-filled region blocks。因此 region.Do 的基线 live heap contribution 应该会很小。

### 偶尔回收不可达内存会有更少的风险

Failing to free an arena, or allocating too much into an arena, may result in running out of memory,

arena 中所有对象的生命期都是绑定在一起的。总而言之就是不够灵活。

