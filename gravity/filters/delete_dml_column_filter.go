package filters

import (
	"github.com/juju/errors"

	"github.com/moiot/gravity/gravity/registry"
	"github.com/moiot/gravity/pkg/core"
)

// [[filters]]
// type = "delete-dml-column"
// match-schema = "test"
// match-table = "test_table"
// columns = ["e", "f"]

const DeleteDMLColumnFilterName = "delete-dml-column"

type deleteDmlColumnFilter struct {
	BaseFilter
	columns []string
}

func (f *deleteDmlColumnFilter) Configure(data map[string]interface{}) error {
	err := f.ConfigureMatchers(data)
	if err != nil {
		return errors.Trace(err)
	}

	columns, ok := data["columns"]
	if !ok {
		return errors.Errorf("\"column\" is not configured")
	}

	c, ok := columns.([]string)
	if !ok {
		return errors.Errorf("\"column\" should be an array of string")
	}

	f.columns = c
	return nil
}

func (f *deleteDmlColumnFilter) Filter(msg *core.Msg) (continueNext bool, err error) {
	if !f.matchers.Match(msg) {
		return true, nil
	}

	if msg.DmlMsg == nil {
		return false, errors.Errorf("DmlMsg is null")
	}

	for _, name := range f.columns {
		delete(msg.DmlMsg.Data, name)

		if msg.DmlMsg.Old != nil {
			delete(msg.DmlMsg.Old, name)
		}

		if msg.DmlMsg.Pks != nil {
			delete(msg.DmlMsg.Pks, name)
		}
	}

	return true, nil
}

type deleteDMLColumnFilterFactoryType struct{}

func (factory *deleteDMLColumnFilterFactoryType) Configure(_ string, _ map[string]interface{}) error {
	return nil
}

func (factory *deleteDMLColumnFilterFactoryType) NewFilter() core.IFilter {
	return &deleteDmlColumnFilter{}
}

func init() {
	registry.RegisterPlugin(registry.FilterPlugin, DeleteDMLColumnFilterName, &deleteDMLColumnFilterFactoryType{}, true)
}
