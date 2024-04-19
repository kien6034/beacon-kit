package tree

type SSZType interface{}
type BasicValue interface {
	ByteLen() int
}

type Container interface{}

func ItemLength(typ SSZType) int {
	// Return the number of bytes in a basic type, or 32 (a full hash) for compound types.
	if typ, ok := typ.(BasicValue); ok {
		return typ.ByteLen()
	}
	return 32
}

// func getElemType(typ interface{}, indexOrVariableName interface{}) SSZType {
// 	// Return the type of the element of an object of the given type with the given index
// 	// or member variable name (eg. `7` for `x[7]`, `"foo"` for `x.foo`)
// 	switch t := typ.(type) {
// 	case Container:
// 		fields := reflect.TypeOf(t).Fields()
// 		switch idx := indexOrVariableName.(type) {
// 		case int:
// 			return fields[idx]
// 		case SSZVariableName:
// 			return fields[idx]
// 		}
// 	case BaseBytes, BaseList:
// 		return t.ElemType()
// 	}
// 	return nil
// }

// func chunkCount(typ SSZType) int {
//     // Return the number of hashes needed to represent the top-level elements in the given type
//     switch t := typ.(type) {
//     case BasicValue:
//         return 1
//     case Bits:
//         return (t.Length() + 255) / 256
//     case Elements:
//         return (t.Length() * itemLength(t.ElemType()) + 31) / 32
//     case Container:
//         return len(t.GetFields())
//     default:
//         panic(fmt.Sprintf("Type not supported: %v", typ))
//     }
// }

// func getItemPosition(typ SSZType, indexOrVariableName interface{}) (int, int, int) {
//     // Return three variables:
//     // (i) the index of the chunk in which the given element of the item is represented;
//     // (ii) the starting byte position within the chunk;
//     // (iii) the ending byte position within the chunk.
//     switch t := typ.(type) {
//     case Elements:
//         index := indexOrVariableName.(int)
//         start := index * itemLength(t.ElemType())
//         return start / 32, start % 32, start%32 + itemLength(t.ElemType())
//     case Container:
//         variableName := indexOrVariableName.(SSZVariableName)
//         return t.GetFieldNames().Index(variableName), 0, itemLength(getElemType(typ, variableName))
//     default:
//         panic("Only lists/vectors/containers supported")
//     }
// }

// func getGeneralizedIndex(typ SSZType, path ...interface{}) GeneralizedIndex {
//     // Converts a path into the generalized index representing its position in the Merkle tree.
//     root := GeneralizedIndex(1)
//     for _, p := range path {
//         if _, ok := typ.(BasicValue); ok {
//             panic("Path cannot continue further into a basic type")
//         }
//         switch p := p.(type) {
//         case string:
//             if p == "__len__" {
//                 typ = uint64(0) // Assuming uint64 is a type that represents List or ByteList
//                 root = GeneralizedIndex(root*2 + 1)
//             } else {
//                 pos, _, _ := getItemPosition(typ, p)
//                 baseIndex := GeneralizedIndex(1)
//                 if _, ok := typ.(List, ByteList); ok {
//                     baseIndex = GeneralizedIndex(2)
//                 }
//                 root = GeneralizedIndex(root * baseIndex * getPowerOfTwoCeil(chunkCount(typ)) + pos)
//                 typ = getElemType(typ, p)
//             }
//         case int:
//             pos, _, _ := getItemPosition(typ, p)
//             baseIndex := GeneralizedIndex(1)
//             if _, ok := typ.(List, ByteList); ok {
//                 baseIndex = GeneralizedIndex(2)
//             }
//             root = GeneralizedIndex(root * baseIndex * getPowerOfTwoCeil(chunkCount(typ)) + pos)
//             typ = getElemType(typ, p)
//         }
//     }
//     return root
// }
