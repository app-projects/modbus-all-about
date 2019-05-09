package trans

import (
	"net/source/proto/trans/interfaces"
	"plugin"
	"sync"
)




func GetEncoderPlugin(soName string) (interfaces.Decoder, error) {
	plg, err1 := plugin.Open(soName)
	if err1 == nil && plg != nil {
		fn, err1 := plg.Lookup("CreateIntance")
		if err1 == nil && fn != nil {
			ins := fn.(func() interfaces.Decoder)()
			return ins, nil
		}
	}
	return nil, err1
}

var encoderPluginsMap sync.Map
func init() {
	//PutEncoder()
}

func PutEncoder(code int32,soPath string)  {

}

