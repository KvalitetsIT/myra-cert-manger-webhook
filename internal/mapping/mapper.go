package mapping

type Mapper[Internal, External any] interface {
	ToInternal(External) Internal // maps external → internal
	ToExternal(Internal) External // maps internal → external
}
