package mappert

// NameConvertorFunc is a function which converts names between structs and maps
//
type NameConvertorFunc func(string) string

// IgnoreFieldsFunc is a function that is used to see if certain fields should not be mapped
//
type IgnoreFieldsFunc func(string) bool

// MapConfiguration holds mapping specific configuration parameters
//
type MapConfiguration struct {
	nameConvertor NameConvertorFunc
	ignoreFields  IgnoreFieldsFunc
}

// And combines two configurations into one
//
func (config *MapConfiguration) And(other *MapConfiguration) *MapConfiguration {
	if other.nameConvertor != nil {
		config.nameConvertor = other.nameConvertor
	}
	if other.ignoreFields != nil {
		config.ignoreFields = other.ignoreFields
	}
	return config
}

// Name applies configuration's name convertor to given name
//
func (config *MapConfiguration) Name(name string) string {
	if config.nameConvertor == nil {
		return name
	}
	return config.nameConvertor(name)
}

// ShouldIgnore applies configuration's field ignoration function to given name
//
func (config *MapConfiguration) ShouldIgnore(name string) bool {
	if config.ignoreFields == nil {
		return false
	}
	return config.ignoreFields(name)
}

// ConvertNamesUsing creates a configuration with a name convertor
//
func ConvertNamesUsing(f NameConvertorFunc) *MapConfiguration {
	return &MapConfiguration{
		nameConvertor: f,
	}
}

// IgnoreFields creates a configuration with a field ignoration
//
func IgnoreFields(names ...string) *MapConfiguration {
	fieldMap := make(map[string]bool, len(names))
	for _, name := range names {
		fieldMap[name] = true
	}

	return &MapConfiguration{
		ignoreFields: func(name string) bool {
			_, found := fieldMap[name]
			return found
		},
	}
}

func combineConfiguration(configs ...*MapConfiguration) *MapConfiguration {

	if len(configs) == 1 {
		return configs[0]
	}

	first := &MapConfiguration{}
	if len(configs) == 0 {
		return first
	}

	for _, next := range configs {
		first = first.And(next)
	}

	return first
}
