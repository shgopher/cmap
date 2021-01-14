# concurrent-map
go实现一个分片的哈希表

使用方法：

````go
// 线程安全
cs := NewCamp()
cs.Set("shgopher",shgopher)
cs.Get("shgopher")
````

本方法就是一种使用了，一致性哈希，然后将问题分布式的一种算法。通过哈希函数，将这些写入和读取的过程分为不同的组去处理，然后锁的颗粒度自然就降低了
