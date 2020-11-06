package snowflake

var defaultNode = &Snowflake{epoch: defaultEpoch, node: 1}

func SetDefault(sets ...func(node *Snowflake)) {
	for _, set := range sets {
		set(defaultNode)
	}
}

//32位秒时间戳+16位自增+5位机器
func NextID() ID {
	return defaultNode.Next()
}

//32位秒时间戳+16位自增+5位机器
func NextInt() int64 {
	return NextID().Int()
}

func NextFormat(radix int) string {
	return NextID().Format(radix)
}

func NextString() string {
	return NextID().String()
}
