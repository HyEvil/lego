package tpl

var(
	 Repository = `
package repository
var (
	{{beanVar}} = {{beanStruct}}{}
)

type {{beanStruct}} struct {
}

func (self *{{beanStruct}}) session(session ...*xorm.Session) *xorm.Session {
	if len(session) > 0 {
		return session[0].Table(table.{{table}}.Model())
	} else {
		return table.{{table}}.Session()
	}
}

func (self *{{ beanStruct }}) OneById(id int64, session ...*xorm.Session) (*model.{{model}}, error) {
	s := self.session(session...)
	m := model.{{model}}{}
	ok, err := s.ID(id).Get(&m)
	if err != nil {
		return nil, apiec.DB.WithError(err)
	}
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (self *{{ beanStruct }}) UpdateById(id int64, bean interface{}, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.ID(id).Update(bean)
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

func (self *{{ beanStruct }}) DeleteById(id int64, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.ID(id).Delete(table.{{table}}.Model())
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

func (self *{{ beanStruct }}) Add(m *model.{{model}}, session ...*xorm.Session) *sql.InsertResult {
	s := self.session(session...)
	affected, err := s.Insert(m)
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewInsertResult(affected, err)
}

func (self *{{ beanStruct }}) Exist(cond builder.Cond, session ...*xorm.Session) (bool, error) {
	s := self.session(session...)
	ok, err := s.Exist(cond)
	if err != nil {
		return false, apiec.DB.WithError(err)
	}
	return ok,nil
}

func (self *{{ beanStruct }}) One(cond builder.Cond, session ...*xorm.Session) (*model.{{model}}, error) {
	s := self.session(session...)
	m := model.{{model}}{}
	ok, err := s.Where(cond).Get(&m)
	if err != nil {
		return nil, apiec.DB.WithError(err)
	}
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (self *{{ beanStruct }}) List(cond builder.Cond, session ...*xorm.Session) ([]*model.{{model}}, *sql.FindResult) {
	s := self.session(session...)
	list := []*model.{{model}}{}
	count, err := s.Where(cond).FindAndCount(&list)
	if err != nil {
		return nil, sql.NewFindResult(0, apiec.DB.WithError(err))
	}
	return list, sql.NewFindResult(count, nil)
}

func (self *{{ beanStruct }}) Find(cond builder.Cond,  extCond ...sql.Cond) ([]*model.{{model}}, *sql.FindResult) {
	s := self.session()
	list := []*model.{{model}}{}
	s = s.Where(cond)
	for _, c := range extCond {
		err := c.Apply(s)
		if err != nil {
			return nil, sql.NewFindResult(0, apiec.Internal.WithError(err))
		}
	}

	count, err := s.FindAndCount(&list)
	if err != nil {
		return nil, sql.NewFindResult(0, apiec.DB.WithError(err))
	}
	return list, sql.NewFindResult(count, nil)
}

func (self *{{ beanStruct }}) FindEx(cond builder.Cond, session *xorm.Session, extCond ...sql.Cond) ([]*model.{{model}}, *sql.FindResult) {
	s := self.session(session)
	list := []*model.{{model}}{}
	s = s.Where(cond)
	for _, c := range extCond {
		err := c.Apply(s)
		if err != nil {
			return nil, sql.NewFindResult(0, apiec.Internal.WithError(err))
		}
	}

	count, err := s.FindAndCount(&list)
	if err != nil {
		return nil, sql.NewFindResult(0, apiec.DB.WithError(err))
	}
	return list, sql.NewFindResult(count, nil)
}

func (self *{{ beanStruct }}) Update(cond builder.Cond, bean interface{}, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.Where(cond).Update(bean)
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}

func (self *{{ beanStruct }}) Delete(cond builder.Cond, session ...*xorm.Session) *sql.AffectedResult {
	s := self.session(session...)
	affected, err := s.Where(cond).Delete(table.{{table}}.Model())
	if err != nil {
		err = apiec.DB.WithError(err)
	}
	return sql.NewAffectedResult(affected, err)
}
`
)
