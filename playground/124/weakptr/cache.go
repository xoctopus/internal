package weakptr

import "weak"

// Cache 按 key 缓存 *Big, 但不延长其生命周期.
//
// 使用场景: 加速重复访问 (解析结果, 编译产物等), 缓存不应像 map[string]*Big 那样
// 一直强引用对象; 仅当调用方仍持有 *Big 时 Get 才能 Value() 命中, 否则 GC 可回收后自然失效.
// 因为不控制大对象的生命周期, 所以可以避免 OOM 的发生.
// 和 LRU 等缓存思路不同的是这里不需要考虑内存容量和淘汰思路. 心智负担全部转移到 GC 机制.
// 代价是重复的 Big 可能会被多次创建或者构建, 但是避免了强引用造成的内存占用导致的 OOM.
type Cache struct {
	entries map[string]weak.Pointer[Big]
}

// Big 表示体积较大或构建代价高的缓存值 (示例用).
type Big struct {
	Label string
}

func (c *Cache) Put(key string, v *Big) {
	if c.entries == nil {
		c.entries = make(map[string]weak.Pointer[Big])
	}
	c.entries[key] = weak.Make(v)
}

// Get 返回仍存活的 *Big; 弱引用已失效时删除条目并返回 nil.
func (c *Cache) Get(key string) *Big {
	wp, ok := c.entries[key]
	if !ok {
		return nil
	}
	v := wp.Value()
	if v == nil {
		delete(c.entries, key)
		return nil
	}
	return v
}
