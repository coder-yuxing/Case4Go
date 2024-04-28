# 一个简单的Goroutine池实现

## 为什么要使用Goroutine池？
goroutine的开销很小，一个goroutine的起始栈大小为2KB，并且创建、切换和销毁的代价很低，我们可以创建成千上万甚至更多的goroutine。但是，
goroutine开销虽然“廉价”，但并不是免费的。一旦规模化后，这种非零成本也会成为瓶颈。我们以一个 Goroutine 分配 2KB 执行栈为例，100w Goroutine 就是 2GB 的内存消耗。

TODO: Go1. 连续栈 & Go1.3 分段栈对比？ 

Goroutine 池就是一种常见的解决方案。这个方案的核心思想是对 Goroutine 的重用，也就是把 M 个计算任务调度到 N 个 Goroutine 上，而不是为每个
计算任务分配一个独享的 Goroutine，从而提高计算资源的利用率。

## goroutine池的实现原理
此goroutine池，workerpool 采用完全基于channel+select 的实现方案。主要分为三个部分：
1. pool的创建和销毁
2. pool中worker（goroutine）的管理
3. task的提交与调度

![1](https://github.com/coder-yuxing/Case4Go/blob/main/projects/workerpool/docs/img/1.jpg)
