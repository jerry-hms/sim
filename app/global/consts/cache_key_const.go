package consts

const (
	ImSessionKeyFormat          string = "im:session:%d"      // im的会话缓存key，hash类型，格式：im:session:{user_id}
	ImSessionFieldFormat        string = "rel_id:%d"          // im的会话缓存field，hash类型，格式：rel_id:{rel_id}
	ImSessionTopicSortKeyFormat string = "im:session:sort:%d" // im会话的排序key, sort类型
)
