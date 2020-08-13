# dycache

dynamic cache: cache size should be changed dynamicly according to the miss rate.

根据未命中缓存个数，推导出应设置的缓存数量。
也可以动态调整缓存数量，最终收敛于最佳缓存数。

## proof

`groupcache` use lru algorithe to purge the element. Once an element was purged, `OnEvicted` will be called if seted.
dycache replace the OnEvicted function, that is necessary for dycache, which can count the range of k/v.
If OnEvicted is used originally, it will be doing by the other solution later.


## 假设与证明

假设1 lru-cache存储大小不足以容纳所有元素，将会导致部分元素被驱逐。
假设2 在duration（可以设置为10s或1m），每个元素都会被驱逐至少一次。

那么，只需要评估当前被驱逐的元素（放入一个set）集合元素的个数，就可以推算出lru-cache应该设置的存储大小。

假设1是成立的，当缓存空间不足时，必有一部分数据进入缓存，从而导致部分元素被驱逐。
著名的2-8定理，80%的请求集中在20%的数据，这部分数据是热数据。但是冷数据也可能会被访问到。那么假设在duration的时间内，所有元素都会被访问到是成立的，取决于duration的设定。
需要给duration一个预估值，比如1m,10m。

cache应设置的容量应是当前cache的容量+当前被驱逐的元素数量。比较遗憾的是，无法评估当前被驱逐的元素数量，而只能评估一段时间内被驱逐的元素数量。

o(c) optimal : 理想缓存数量
c(c) current : 当前缓存数量
e(c) evicted : 一段时间内被驱逐元素数量

o(c) = c(c) + e(c), 
e(c) = s(d) d ~ 0, 当前时间（d趋向于0, 当前时刻）被驱逐元素数量

再者，在高并发场景下，e(c)将会极其不准确。因此，e(c)的评估，需要评估一段时间d的集合大小。

元素i会被不断驱逐，我们需要记录元素被驱逐的轨迹，统计当前时间范围内驱逐元素个数：

[Ae1 ,Be1 ,Ce1 , ... , Ae2, ...]

偶数次(odd)被驱逐，元素个数+1; 奇数次被驱逐，元素个数-1。

观察某一个元素在set中的出现次数，呈现0-1-0-1锯齿状的曲线。曲线映射到x轴的面积刚好是0.5，即出现在被驱逐集合里，该元素的平均个数是0.5，期望是1（最优个数）。
推演到所有元素，图像的面积（当前集合中元素个数）刚好是期望值（cache理想容量）的一半。

由此可知，cache的理想容量是被驱逐集合元素的2倍。
