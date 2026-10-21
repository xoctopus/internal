package weakptr

import "weak"

// InternTable 按 key 做规范化 (intern): 同一 key 在 *Symbol 仍存活时始终返回同一指针.
//
// 使用场景: 多处需要指向「逻辑上同一个」对象 (符号名, 类型元数据, 解析后的常量等),
// 希望合并成一份, 但 intern table 不应像 map[string]*Symbol 那样永久拽住已不用的对象.
//
// table 只存 weak.Pointer; Intern 返回的 *Symbol 是强引用, 由调用方持有.
// 当没有任何 *Symbol 存活时 GC 可回收, 弱条目在下次 Intern 时被替换为新分配.
// 字符串去重可优先用标准库 unique.Make; 此处演示「按 key 指向堆对象」的 weak intern 思路.
type InternTable struct {
	entries map[string]weak.Pointer[Symbol]
	nextGen int
}

// Symbol 表示被 intern 的规范实例 (示例用). Gen 仅用于示例里观察「是否重新分配」.
type Symbol struct {
	Name string
	Gen  int
}

func (t *InternTable) Intern(key string) *Symbol {
	if t.entries == nil {
		t.entries = make(map[string]weak.Pointer[Symbol])
	}
	if wp, ok := t.entries[key]; ok {
		if s := wp.Value(); s != nil {
			return s
		}
		delete(t.entries, key)
	}
	t.nextGen++
	s := &Symbol{Name: key, Gen: t.nextGen}
	t.entries[key] = weak.Make(s)
	return s
}
