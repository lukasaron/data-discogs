package decode

import (
	"encoding/xml"
)

func (x *XMLDecoder) parseCompanies() (companies []Company) {
	if x.err != nil {
		return companies
	}

	company := Company{}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "companies" {
			break
		}

		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "id":
				company.Id = x.parseValue()
			case "name":
				company.Name = x.parseValue()
			case "catno":
				company.Category = x.parseValue()
			case "entity_type":
				company.EntityType = x.parseValue()
			case "entity_type_name":
				company.EntityTypeName = x.parseValue()
			case "resource_url":
				company.ResourceUrl = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "company" {
			companies = append(companies, company)
			company = Company{}
		}
	}
	return companies
}
