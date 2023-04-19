package nodes

// NodeInfos 存储所有node
type NodeInfos struct {
	NodeInfos   []*NodeInfo
	NodeInfoNum int
}

func NewNodeInfos() *NodeInfos {
	return &NodeInfos{
		NodeInfos:   make([]*NodeInfo, 0),
		NodeInfoNum: 0,
	}
}

// NodeInfo node信息
type NodeInfo struct {
	Node       map[string]interface{}
	NodeName   string
	NodeLen    int
	NodeLabels map[string]string
}

func NewNodeInfo(name string) *NodeInfo {
	n := &NodeInfo{
		NodeName:   name,
		Node:       map[string]interface{}{},
		NodeLen:    0,
		NodeLabels: map[string]string{},
	}

	return n
}

func (i *NodeInfo) SetNodeLabels(key, value string) {
	i.NodeLabels[key] = value
}

func (i *NodeInfo) SetNodeName(name string) {
	i.NodeName = name
}

func (ns *NodeInfos) AddNode(n *NodeInfo) {
	ns.NodeInfos = append(ns.NodeInfos, n)
	ns.NodeInfoNum++
}

func (ns *NodeInfos) DeleteNode(nName string) {
	for i, nodeInfo := range ns.NodeInfos {
		if nodeInfo.NodeName == nName {
			ns.NodeInfos = append(ns.NodeInfos[:i], ns.NodeInfos[i+1:]...)
			ns.NodeInfoNum--
		}
	}
}
