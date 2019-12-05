package parser

import (
	"encoding/xml"
	"github.com/Twyer/discogs/model"
)

func ParseCompanies(tr xml.TokenReader) []model.Company {
	companies := make([]model.Company, 0, 0)
	company := model.Company{}
	for {
		t, _ := tr.Token()
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "companies" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "id":
				company.Id = parseValue(tr)
			case "name":
				company.Name = parseValue(tr)
			case "catno":
				company.Category = parseValue(tr)
			case "entity_type":
				company.EntityType = parseValue(tr)
			case "entity_type_name":
				company.EntityTypeName = parseValue(tr)
			case "resource_url":
				company.ResourceUrl = parseValue(tr)
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "company" {
			companies = append(companies, company)
			company = model.Company{}
		}
	}
	return companies
}
