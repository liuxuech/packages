package cache

import "container/list"

type Cache struct {
	maxBytes  int64                    // 允许使用最大内存
	nBytes    int64                    // 当前已使用内存
	ll        *list.List               // 双向链表实现的队列
	cache     map[string]*list.Element // key：键，value：实际元素的地址
	OnEvicted func(string, Value)      // 记录被移除时执行
}

func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{maxBytes: maxBytes, OnEvicted: onEvicted}
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// 获取缓存值
//   双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if el, ok := c.cache[key]; ok {
		// 将节点移动到队列尾部
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		return kv.value, true
	}
	return
}

// 缓存淘汰
func (c *Cache) RemoveOldest() {
	// 获取队列的最后一个元素
	if el := c.ll.Back(); el != nil {
		// 从队列中删除
		c.ll.Remove(el)
		// 取出key，再从map中删除映射
		kv := el.Value.(*entry)
		delete(c.cache, kv.key)
		// 重新计算已使用内存,记录锁占的字节数：(len(key) + kv.value.Len())
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 如果有设置回调，执行
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// 先判断元素是否存在
	if el, ok := c.cache[key]; ok {
		// 将元素移动到队尾
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		// 重新计算已使用内存
		c.nBytes += int64(len(key)) + int64(value.Len())
		// 更新记录值
		kv.value = value
	} else {
		// 不存在，则在队尾添加新节点
		c.ll.PushFront(&entry{key, value})
		// 重新计算已使用内存
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 新添加元素后，计算内存是否超过最大内存使用
	// 注意：这里使用for不用if，是因为可能会 remove 多次。
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// 返回记录个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
