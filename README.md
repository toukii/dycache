# dycache

dynamic cache: cache size should be changed dynamicly according to the miss rate.

根据未命中缓存个数，推导出应设置的缓存数量。
也可以动态调整缓存数量，最终收敛于最佳缓存数。

## proof

`groupcache` use lru algorithe to purge the element. Once an element was purged, `OnEvicted` will be called if seted.
dycache replace the OnEvicted function, that is necessary for dycache, which can count the range of k/v.
If OnEvicted is used originally, it will be doing by the other solution later.


