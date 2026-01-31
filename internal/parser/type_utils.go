package parser

func UsesRawInModels(schema Schema) bool {
	for _, model := range schema.Models {
		for _, field := range model.Fields {
			if HasRawType(field.Type) {
				return true
			}
		}
	}
	return false
}

func UsesRawInRPCs(schema Schema) bool {
	for _, rpc := range schema.RPCs {
		for _, param := range rpc.Parameters {
			if HasRawType(param.Type) {
				return true
			}
		}
		if rpc.HasReturn && HasRawType(rpc.Returns) {
			return true
		}
	}
	return false
}

func HasRawType(t TypeRef) bool {
	switch t.Kind {
	case TypeList:
		if t.Elem == nil {
			return false
		}
		return HasRawType(*t.Elem)
	case TypeMap:
		if t.Value == nil {
			return false
		}
		return HasRawType(*t.Value)
	default:
		return t.Name == "raw"
	}
}
