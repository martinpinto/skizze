package cml

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/count-min-log/count-min-log"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

const defaultEpsilon = 0.00000543657
const defaultDelta = 0.99

/*
Domain is the toplevel domain to control the count-min-log implementation
*/
type Domain struct {
	*abstract.Info
	impl *cml.Sketch16
	lock sync.RWMutex
}

/*
NewDomain ...
*/
func NewDomain(info *abstract.Info) (*Domain, error) {
	manager = storage.GetManager()
	manager.Create(info.ID)
	sketch16, _ := cml.NewSketch16ForEpsilonDelta(info.ID, defaultEpsilon, defaultDelta)
	d := Domain{info, sketch16, sync.RWMutex{}}
	err := d.Save()
	if err != nil {
		logger.Error.Println("an error has occurred while saving domain: " + err.Error())
	}
	return &d, nil
}

/*
NewDomainFromData ...
*/
func NewDomainFromData(info *abstract.Info) (*Domain, error) {
	sketch16, _ := cml.NewSketch16ForEpsilonDelta(info.ID, defaultEpsilon, defaultDelta)
	// FIXME: create domain from new data
	return &Domain{info, sketch16, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Domain) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.impl.IncreaseCount(value)
	d.Save()
	return true, nil
}

/*
AddMultiple ...
*/
func (d *Domain) AddMultiple(values [][]byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	for _, value := range values {
		d.impl.IncreaseCount(value)
	}
	d.Save()
	return true, nil
}

/*
Remove ...
*/
func (d *Domain) Remove(value []byte) (bool, error) {
	logger.Error.Println("This domain type does not support deletion")
	return false, errors.New("This domain type does not support deletion")
}

/*
RemoveMultiple ...
*/
func (d *Domain) RemoveMultiple(values [][]byte) (bool, error) {
	logger.Error.Println("This domain type does not support deletion")
	return false, errors.New("This domain type does not support deletion")
}

/*
GetCount ...
*/
func (d *Domain) GetCount() uint {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return uint(d.impl.Count())
}

/*
Clear ...
*/
func (d *Domain) Clear() (bool, error) {
	d.impl.Reset()
	return true, nil
}

/*
Save ...
*/
func (d *Domain) Save() error {
	count := d.impl.Count()
	d.Info.State["count"] = uint64(count)
	infoData, err := json.Marshal(d.Info)
	if err != nil {
		return err
	}
	return storage.GetManager().SaveInfo(d.Info.ID, infoData)
}

/*
GetType ...
*/
func (d *Domain) GetType() string {
	return d.Type
}

/*
GetFrequency ...
*/
func (d *Domain) GetFrequency(values [][]byte) interface{} {
	res := make(map[string]uint)
	for _, value := range values {
		res[string(value)] = uint(d.impl.GetCount(value))
	}
	return res
}
