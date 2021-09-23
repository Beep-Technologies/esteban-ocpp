{{- $alias := .Aliases.Table .Table.Name -}}
{{- $orig_tbl_name := .Table.Name -}}

// {{$alias.UpSingular}} is an object representing the database table.
type {{$alias.UpSingular}} struct {
	{{- range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{- $orig_col_name := $column.Name -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}gorm:"column:{{$column.Name}}"`
	{{end -}}
}

func ({{$alias.UpSingular}}) TableName() string {
    _ = strconv.Quote("")
	return constants.SCHEMA + "{{$orig_tbl_name}}"
}