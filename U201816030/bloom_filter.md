# 基于 Rust 语言的多维 Bloom Filter 设计与实现

## 引言
Bloom Filter 是由 Burton Howard BLoom 在 1970 年提出的一种用于数据去重，空间效率高的概率型数据结构。它专门用来检测集合中是否存在特定的元素。  
Rust 语言是一门现代系统级编程语言，同时兼顾高性能和安全。和 C/C++ 相比，Rust 语言引入了所有权和生命周期机制，保证了系统运行时的安全性，与 Java/Go 相比，它没有 GC 机制，因此具有更高效的运行时系统。  
这里笔者将基于 Rust 语言来实现一款 Bloom Filter 支持库，并拓展到多维，分析相关性能。  

## Bloom Filter 的设计与实现
### 设计思想
Bloom Filter（下文统一称为 BF）是由一个长度为 m 比特的位数组与 k 个哈希函数组成的数据结构。位数组均初始化为 0，所有哈希函数都可以分别把输入数据尽量均匀地散列。  
当要插入一个元素时，将其数据分别输入 k 个哈希函数，产生 k 个哈希值，将 k 个哈希值对应的数据位置 1 。
当要查询一个元素的时候，同样将其数据输入哈希函数，然后检查对应的 k 个数据位，如果有任意一个数据位为 0，那么该元素一定不在集合中，否则该数据有较大可能性在集合中。  
### 数据结构设计
对于一个 BF，我们完全可以在编译期就知道需要分配的内存空间是多少，因此我们可以使用`常量泛型`来实现这个数据结构。  
2021 年 Rust 1.51 版本中稳定了常量泛型，它的一个作用是用于构建包含数组类型成员的结构体：  
```Rust
struct ArrayPair<T, const N: usize> {
    left: [T; N],
    right: [T; N]
}
```

从上面这个例子可以看到我们可以达到在编译期数组的长度是可变的，但是在运行时里面数组里面的数据是放到栈上的效果。  
如果使用可变数组(Vec<T>)去实现这种数据结构，那么数据将是放到堆上的，这样会损失一点运行时开销。  
基于这种考虑我们使用常量泛型来实现 BF，结构体如下：  
```Rust
pub struct Filter<BHK: BuildHashKernels, const W: usize, const M: usize, const D: u8> {
    buckets: Buckets<W, M, D>,      // filter data
    hash_kernels: BHK::HK, // hash kernels
    p: usize,              // number of buckets to decrement,
}
```
这个结构体定义除了三个常量泛型外还有一个 BHK 泛型，是和哈希函数有关的。  
首先这个结构体由一个 Buckets 和一系列哈希函数组成。Buckets 里面包含数组结构，用于存放 BF 位数据信息，插入和查询都会从这里进行数据检索。  
哈希函数将输入值散列，结果将会成为访问 Buckets 的下标。  

### Buckets 的设计与实现
BF 设计的核心是正是 Buckets 的设计，它的结构体实现如下：  
```Rust
pub struct Buckets<const WordCount: usize, const BucketCount: usize, const BucketSize: u8> {
    data: [Word; WordCount],
    max: u8,
}
```
可以看到这个结构体有三个常量泛型，WordCount 表示数组里面的表项数目，BucketCount 表示 Buckets 中 Bucket 的数量，BucketSize 表示一个 Bucket 占多少字节。  
数据检索的最小单位就是 Bucket，有时候我们需要大一点的 BucketSize。  
通过常量泛型，我们可以在让 data 数组长度可变的同时，它的数据还是放在栈上。(如果使用 Vec<T> 的话数据会放到堆上)  
然后我们为这个结构体实现一系列的方法，向上提供一些 API：  
```Rust
impl<const WordCount: usize, const BucketCount: usize, const BucketSize: u8> Buckets<WordCount, BucketCount, BucketSize> {
    /// 设置某个 bucket 的值
    pub fn set(&mut self, bucket: usize, byte: u8) {
        let offset = bucket * BucketSize as usize;
        let length = BucketSize as usize;
        let word = if byte > self.max as u8 { self.max } else { byte } as Word;
        self.set_word(offset, length, word);
    }

    /// 获取某个 bucket 的值
    pub fn get(&self, bucket: usize) -> u8 {
        self.get_word(bucket * BucketSize as usize, BucketSize as usize) as u8
    }
}
```
通过这两个方法函数，我们就可以做到对某个 bucket 置 1 和获取某个 bucket 的值了。  

### 哈希函数设计与实现
在 BF 中，我们需要的是一系列的哈希函数，而不是单个，因此我们借助 Rust 语言的迭代器语法来设计哈希函数：  
```Rust
/// A trait for creating hash iterator of item.
/// Rust 里面的 trait 相当于 Java 里面的 interface
pub trait HashKernels {
    type HI: Iterator<Item = usize>;

    fn hash_iter<T: Hash>(&self, item: &T) -> Self::HI;
}
```
这个 HashKernels trait 只有一个方法 hash_iter，它的语义是返回一个哈希函数的迭代器，这样就可以抽象出“一系列哈希函数”的概念了。  

### insert 方法和 contains 方法的实现
有了 Buckets 和 HashKernels 的基础，我们就可以实现数据的插入和查询方法，为应用场景提供 API 了。  
```Rust
impl<BHK: BuildHashKernels, const W: usize, const B: usize, const S: u8> BloomFilter for Filter<BHK, W, B, S> {
    /// 插入数据，更新所有哈希结果对应的 bucket
fn insert<T: Hash>(&mut self, item: &T) {
        self.decrement();
        let max = self.buckets.max_value();
        self.hash_kernels.hash_iter(item).for_each(|i| self.buckets.set(i, max))
    }
    /// 查询数据是否存在在集合中，只有所有哈希结果对应的 bucket 都被置一才返回 true
    /// 可能会误报，但不可能漏报
    fn contains<T: Hash>(&self, item: &T) -> bool {
        self.hash_kernels.hash_iter(item).all(|i| self.buckets.get(i) > 0)
    }
}
```
这样我们就实现了 insert 和 contains 方法，可以插入和查询数据了。  

### 正确性测试
这里基于 Rust 语言内置的单元测试系统，来测试上述结构实现的正确性：  
```Rust
fn _contains(items: &[usize]) {
        let mut filter = Filter::<_, {compute_word_num(730, 3)}, 730, 3>::new(0.03, DefaultBuildHashKernels::new(random(), RandomState::new()));
        assert!(items.iter().all(|i| !filter.contains(i)));
        items.iter().for_each(|i| filter.insert(i));
        assert!(items.iter().all(|i| filter.contains(i)));
    }

    proptest! {
        #[test]
        fn contains(ref items in any_with::<Vec<usize>>(size_range(7).lift())) {
            _contains(items)
        }
    }
```
测试结果：  
```
running 1 test
test stable::tests::contains ... ok

test result: ok. 1 passed; 0 failed; 0 ignored; 0 measured; 6 filtered out; finished in 0.02s
```

## 多维 Bloom Filter 的设计与实现

## 测试分析
这里借助开源项目[criterion](https://github.com/bheisler/criterion.rs)进行系统测试和分析。  
该项目可以帮助我们运行测试任务，输出运行时间，统计尾延迟，运行时间最佳估计等。  
### 延迟
### false positive
### 空间开销