package gormAdapter

import (
	"errors"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jinzhu/gorm"
)

var tablePrefix string

type CasbinRule struct {
	TablePrefix string `gorm:"-"`
	PType       string `gorm:"type:varchar"`
	V0          string `gorm:"type:varchar"`
	V1          string `gorm:"type:varchar"`
	V2          string `gorm:"type:varchar"`
	V3          string `gorm:"type:varchar"`
	V4          string `gorm:"type:varchar"`
	V5          string `gorm:"type:varchar"`
}

type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

func (c *CasbinRule) TableName() string {
	return c.TablePrefix + "cb_rules" //as Gorm keeps table names are plural, and we love consistency
}

// Adapter represents the Gorm adapter for policy storage.
type Adapter struct {
	tablePrefix    string
	driverName     string
	dataSourceName string
	dbSpecified    bool
	db             *gorm.DB
	isFiltered     bool
}

// NewAdapterByDB obtained through an existing Gorm instance get  a adapter, specify the table prefix
// Example: gormadapter.NewAdapterByDBUsePrefix(&db, "cms_") Automatically generate table name like this "cms_casbin_rule"
func NewAdapterByDBUsePrefix(db *gorm.DB, prefix string) (*Adapter, error) {
	a := &Adapter{
		tablePrefix: prefix,
		db:          db,
	}
	tablePrefix = prefix
	err := a.createTable()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func NewAdapterByDB(db *gorm.DB) (*Adapter, error) {
	a := &Adapter{
		db: db,
	}

	err := a.createTable()
	if err != nil {
		return nil, err
	}

	return a, nil
}

// getTableInstance return the dynamic table name
func (a *Adapter) getTableInstance() *CasbinRule {
	return &CasbinRule{TablePrefix: a.tablePrefix}
}

func (a *Adapter) createTable() error {
	if a.db.HasTable(a.getTableInstance()) {
		return nil
	}

	return a.db.CreateTable(a.getTableInstance()).Error
}

func (a *Adapter) dropTable() error {
	return a.db.DropTable(a.getTableInstance()).Error
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	if err := a.db.Find(&lines).Error; err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

// LoadFilteredPolicy loads only policy rules that match the filter.
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	var lines []CasbinRule

	filterValue, ok := filter.(Filter)
	if !ok {
		return errors.New("invalid filter type")
	}

	if err := a.db.Scopes(a.filterQuery(a.db, filterValue)).Find(&lines).Error; err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}
	a.isFiltered = true

	return nil
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *Adapter) IsFiltered() bool {
	return a.isFiltered
}

// filterQuery builds the gorm query to match the rule filter to use within a scope.
func (a *Adapter) filterQuery(db *gorm.DB, filter Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(filter.PType) > 0 {
			db = db.Where("p_type in (?)", filter.PType)
		}
		if len(filter.V0) > 0 {
			db = db.Where("v0 in (?)", filter.V0)
		}
		if len(filter.V1) > 0 {
			db = db.Where("v1 in (?)", filter.V1)
		}
		if len(filter.V2) > 0 {
			db = db.Where("v2 in (?)", filter.V2)
		}
		if len(filter.V3) > 0 {
			db = db.Where("v3 in (?)", filter.V3)
		}
		if len(filter.V4) > 0 {
			db = db.Where("v4 in (?)", filter.V4)
		}
		if len(filter.V5) > 0 {
			db = db.Where("v5 in (?)", filter.V5)
		}
		return db
	}
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{TablePrefix: tablePrefix}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	err := a.dropTable()
	if err != nil {
		return err
	}
	err = a.createTable()
	if err != nil {
		return err
	}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			err := a.db.Create(&line).Error
			if err != nil {
				return err
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			err := a.db.Create(&line).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	err := a.db.Create(&line).Error
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	err := rawDelete(a.db, line) //can't use db.Delete as we're not using primary key http://jinzhu.me/gorm/crud.html#delete
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{TablePrefix: tablePrefix}

	line.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err := rawDelete(a.db, line)
	return err
}

func rawDelete(db *gorm.DB, line CasbinRule) error {
	queryArgs := []interface{}{line.PType}

	queryStr := "p_type = ?"
	if line.V0 != "" {
		queryStr += " and v0 = ?"
		queryArgs = append(queryArgs, line.V0)
	}
	if line.V1 != "" {
		queryStr += " and v1 = ?"
		queryArgs = append(queryArgs, line.V1)
	}
	if line.V2 != "" {
		queryStr += " and v2 = ?"
		queryArgs = append(queryArgs, line.V2)
	}
	if line.V3 != "" {
		queryStr += " and v3 = ?"
		queryArgs = append(queryArgs, line.V3)
	}
	if line.V4 != "" {
		queryStr += " and v4 = ?"
		queryArgs = append(queryArgs, line.V4)
	}
	if line.V5 != "" {
		queryStr += " and v5 = ?"
		queryArgs = append(queryArgs, line.V5)
	}
	args := append([]interface{}{queryStr}, queryArgs...)
	err := db.Delete(CasbinRule{TablePrefix: tablePrefix}, args...).Error
	return err
}
