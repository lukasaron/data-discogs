package model

type Company struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	EntityType     string `json:"entityType"`
	EntityTypeName string `json:"entityTypeName"`
	ResourceUrl    string `json:"resourceUrl"`
}

/*
<companies>
  <company>
	 <id>266169</id>
	 <name>JTS Studios</name>
	 <catno />
	 <entity_type>29</entity_type>
	 <entity_type_name>Mastered At</entity_type_name>
	 <resource_url>https://api.discogs.com/labels/266169</resource_url>
  </company>
  <company>
	 <id>56025</id>
	 <name>MPO</name>
	 <catno />
	 <entity_type>17</entity_type>
	 <entity_type_name>Pressed By</entity_type_name>
	 <resource_url>https://api.discogs.com/labels/56025</resource_url>
  </company>
</companies>
*/
