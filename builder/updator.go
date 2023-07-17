package builder

func (ug *BMFUpdator) AddNode(inputAlias []string, outputAlias []string, nodeAlias string, moduleName string,
	moduleType BMFModuleType, option interface{}, modulePath string, moduleEntry string, inputManager BMFInputManager,
	scheduler int, preModule *CBMFModule) {
	toAdd := &BMFUpdatorNode{
		id:     len(ug.relatedG.nodes),
		action: "add",
		alias:  nodeAlias,
		option: option,
		moduleInfo: &BMFModuleMeta{
			name:     moduleName,
			language: moduleType,
			path:     modulePath,
			entry:    moduleEntry,
		},
		metaInfo: &BMFNodeMeta{
			preModuleId:      -1,
			callbackBindings: nil,
		},
		inputManager: inputManager,
		scheduler:    scheduler,
		inStream:     []*BMFUpdatorStream{},
		outStream:    []*BMFUpdatorStream{},
	}
	ug.relatedG.nodes = append(ug.relatedG.nodes, &BMFNode{
		graph: ug.relatedG,
		id:    len(ug.relatedG.nodes),
	})
	if preModule != nil {
		uid := preModule.UID()
		toAdd.metaInfo.preModuleId = int(uid)
	}
	for _, s := range inputAlias {
		toAdd.inStream = append(toAdd.inStream, &BMFUpdatorStream{alias: s})
	}
	for _, s := range outputAlias {
		toAdd.outStream = append(toAdd.outStream, &BMFUpdatorStream{alias: s})
	}
	ug.nodes = append(ug.nodes, toAdd)
}

func (ug *BMFUpdator) RemoveNode(nodeAlias string) {
	ug.nodes = append(ug.nodes, &BMFUpdatorNode{
		action: "remove",
		alias:  nodeAlias,
	})
}

func (ug *BMFUpdator) ResetNode(nodeAlias string, option interface{}) {
	ug.nodes = append(ug.nodes, &BMFUpdatorNode{
		action: "reset",
		alias:  nodeAlias,
		option: option,
	})
}
