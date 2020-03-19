package gowe

// 场景信息模型
type SceneInfoModel struct {
	ID       string `json:"id"`        // 门店唯一标识
	Name     string `json:"name"`      // 门店名称
	AreaCode string `json:"area_code"` // 门店所在地行政区划码，详细见《最新县及县以上行政区划代码》
	Address  string `json:"address"`   // 门店详细地址
}
