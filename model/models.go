package model

//--------------------------------------------------- Artist ---------------------------------------------------

/*
Artist is one of the main structure from Discogs and the XML version looks like this:
	<artist>
	   <images>
		  <image height="337" type="primary" uri="" uri150="" width="600" />
		  <image height="554" type="secondary" uri="" uri150="" width="600" />
	   </images>
	   <id>132</id>
	   <name>Minimum Wage Brothers</name>
	   <realname>Terrence Parker</realname>
	   <profile />
	   <data_quality>Correct</data_quality>
	   <urls>
		  <url>http://www.terrenceparker.net</url>
		  <url>http://www.myspace.com/terrenceparker</url>
	   </urls>
	   <namevariations>
		  <name>Minimum Wage Bros.</name>
		  <name>The Minimum Wage Bros.</name>
	   </namevariations>
	   <aliases>
		  <name id="10678">2 Sweat Doctors</name>
		  <name id="121">Disco Revisited</name>
	   </aliases>
</artist>
*/
type Artist struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	RealName       string   `json:"realName"`
	Images         []Image  `json:"images,omitempty"`
	Profile        string   `json:"profile"`
	DataQuality    string   `json:"dataQuality"`
	NameVariations []string `json:"nameVariations"`
	Urls           []string `json:"urls"`
	Aliases        []Alias  `json:"aliases"`
	Members        []Member `json:"members,omitempty"`
}

// Alias is a sub structure of Artist.
type Alias struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Member is a sub structure of Artist.
type Member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//--------------------------------------------------- Company ---------------------------------------------------

/*
Discogs Company XML structure:
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
type Company struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	EntityType     string `json:"entityType"`
	EntityTypeName string `json:"entityTypeName"`
	ResourceUrl    string `json:"resourceUrl"`
}

//--------------------------------------------------- Format ---------------------------------------------------

/*
Discogs Format XML structure
 	<formats>
      <format name="Vinyl" qty="1" text="">
         <descriptions>
            <description>12"</description>
            <description>33 ⅓ RPM</description>
         </descriptions>
      </format>
   </formats>
*/
type Format struct {
	Name         string   `json:"name"`
	Quantity     string   `json:"quantity"`
	Text         string   `json:"text"`
	Descriptions []string `json:"description"`
}

//--------------------------------------------------- Image ---------------------------------------------------
/*
Discogs Image XML structure
	<image height="337" type="primary" uri="" uri150="" width="600" />
*/

type Image struct {
	Height string `json:"height"`
	Width  string `json:"width"`
	Type   string `json:"type"`
	Uri    string `json:"uri"`
	Uri150 string `json:"uri150"`
}

//--------------------------------------------------- Label ---------------------------------------------------

/*
Label is one of the main structure from Discogs and the XML version looks like this:
	<label>
		<images>
		   <image height="600" type="primary" uri="" uri150="" width="600" />
		   <image height="338" type="secondary" uri="" uri150="" width="300" />
		</images>
		<id>43</id>
		<name>Axis</name>
		<contactinfo>Axis Records&#xD;
	P.O. Box 416600&#xD;
	Miami Beach, FL 33141 USA&#xD;
	Tel: +1 786 953 4176&#xD;
		</contactinfo>
		<profile>American techno label established in 1991 in Chicago</profile>
		<data_quality>Needs Vote</data_quality>
		<parentLabel id="1175241">Axis Records (10)</parentLabel>
		<urls>
			<url>http://www.axisrecords.com</url>
			<url>http://twitter.com/AxisRecords</url>
		</urls>
		<sublabels>
			<label id="15681">6277</label>
			<label id="4504">Luxury Records</label>
		</sublabels>
	</label>
*/
type Label struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Images      []Image      `json:"images,omitempty"`
	ContactInfo string       `json:"contactInfo"`
	Profile     string       `json:"profile"`
	DataQuality string       `json:"dataQuality"`
	Urls        []string     `json:"urls"`
	ParentLabel *LabelLabel  `json:"parentLabel,omitempty"`
	SubLabels   []LabelLabel `json:"subLabels"`
}

type LabelLabel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//--------------------------------------------------- Master ---------------------------------------------------

/*
Master is one of the main structure from Discogs and the XML version looks like this:
	<master id="33228">
	   <main_release>341048</main_release>
	   <images>
		  <image height="602" type="primary" uri="" uri150="" width="600" />
		  <image height="600" type="secondary" uri="" uri150="" width="600" />
	   </images>
	   <artists>
		  <artist>
			 <id>20691</id>
			 <name>Philip Glass</name>
			 <anv />
			 <join />
			 <role />
			 <tracks />
		  </artist>
	   </artists>
	   <genres>
		  <genre>Electronic</genre>
		  <genre>Classical</genre>
	   </genres>
	   <styles>
		  <style>Soundtrack</style>
		  <style>Modern Classical</style>
	   </styles>
	   <year>1988</year>
	   <title>Powaqqatsi</title>
	   <data_quality>Correct</data_quality>
	   <videos>
		  <video duration="303" embed="true" src="https://www.youtube.com/watch?v=jRg3agJn1Mg">
			 <title>Philip Glass - Powaqqatsi - 01. Serra Pelada</title>
			 <description>Philip Glass - Powaqqatsi - 01. Serra Pelada</description>
		  </video>
		  <video duration="25" embed="true" src="https://www.youtube.com/watch?v=4E4gn9nNGYs">
			 <title>Philip Glass - Powaqqatsi - 02. The Title</title>
			 <description>Philip Glass - Powaqqatsi - 02. The Title</description>
		  </video>
	   </videos>
	</master>
*/
type Master struct {
	Id          string          `json:"id"`
	MainRelease string          `json:"mainRelease"`
	Images      []Image         `json:"images,omitempty"`
	Artists     []ReleaseArtist `json:"artists"`
	Genres      []string        `json:"genres"`
	Styles      []string        `json:"styles"`
	Year        string          `json:"year"`
	Title       string          `json:"title"`
	DataQuality string          `json:"dataQuality"`
	Videos      []Video         `json:"videos"`
}

//--------------------------------------------------- Release ---------------------------------------------------

/*
Release is one of the main structure from Discogs and the XML version looks like this:
	<release id="2" status="Accepted">
	   <images>
		  <image height="394" type="primary" uri="" uri150="" width="400" />
		  <image height="600" type="secondary" uri="" uri150="" width="600" />
	   </images>
	   <artists>
		  <artist>
			 <id>2</id>
			 <name>Mr. James Barth &amp; A.D.</name>
			 <anv />
			 <join />
			 <role />
			 <tracks />
		  </artist>
	   </artists>
	   <title>Knockin' Boots Vol 2 Of 2</title>
	   <labels>
		  <label catno="SK 026" id="5" name="Svek" />
		  <label catno="SK026" id="5" name="Svek" />
	   </labels>
	   <extraartists>
		  <artist>
			 <id>26</id>
			 <name>Alexi Delano</name>
			 <anv />
			 <join />
			 <role>Producer, Recorded By</role>
			 <tracks />
		  </artist>
		  <artist>
			 <id>27</id>
			 <name>Cari Lekebusch</name>
			 <anv />
			 <join />
			 <role>Producer, Recorded By</role>
			 <tracks />
		  </artist>
	   </extraartists>
	   <formats>
		  <format name="Vinyl" qty="1" text="">
			 <descriptions>
				<description>12"</description>
				<description>33 ⅓ RPM</description>
			 </descriptions>
		  </format>
	   </formats>
	   <genres>
		  <genre>Electronic</genre>
	   </genres>
	   <styles>
		  <style>Broken Beat</style>
		  <style>Techno</style>
	   </styles>
	   <country>Sweden</country>
	   <released>1998-06-00</released>
	   <notes>All joints recorded in NYC (Dec.97).</notes>
	   <data_quality>Correct</data_quality>
	   <master_id is_main_release="true">713738</master_id>
	   <tracklist>
		  <track>
			 <position>A1</position>
			 <title>A Sea Apart</title>
			 <duration>5:08</duration>
		  </track>
	   </tracklist>
	   <identifiers>
		  <identifier description="Side A Runout Etching" type="Matrix / Runout" value="MPO SK026-A -J.T.S.-" />
		  <identifier description="Side B Runout Etching" type="Matrix / Runout" value="MPO SK026-B -J.T.S.-" />
	   </identifiers>
	   <videos>
		  <video duration="296" embed="true" src="https://www.youtube.com/watch?v=2h0YM1ve6dE">
			 <title>Mr. James Barth &amp; A.D. - Yeah Kid!</title>
			 <description>Mr. James Barth &amp; A.D. - Yeah Kid!</description>
		  </video>
		  <video duration="266" embed="true" src="https://www.youtube.com/watch?v=wRzbbCgg_jY">
			 <title>Mr. James Barth &amp; A.D. - Dutchmaster</title>
			 <description>Mr. James Barth &amp; A.D. - Dutchmaster</description>
		  </video>
	   </videos>
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
	</release>
*/
type Release struct {
	Id           string          `json:"id"`
	Status       string          `json:"status"`
	Images       []Image         `json:"images,omitempty"`
	Artists      []ReleaseArtist `json:"artists"`
	ExtraArtists []ReleaseArtist `json:"extraArtists"`
	Title        string          `json:"title"`
	Formats      []Format        `json:"formats"`
	Genres       []string        `json:"genres"`
	Styles       []string        `json:"styles"`
	Country      string          `json:"country"`
	Released     string          `json:"released"`
	Notes        string          `json:"notes"`
	DataQuality  string          `json:"dataQuality"`
	MasterId     string          `json:"masterId"`
	MainRelease  string          `json:"mainRelease"`
	TrackList    []Track         `json:"trackList"`
	Identifiers  []Identifier    `json:"identifiers"`
	Videos       []Video         `json:"videos"`
	Labels       []ReleaseLabel  `json:"labels"`
	Companies    []Company       `json:"companies"`
}

type ReleaseArtist struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Join   string `json:"join"`
	Anv    string `json:"anv"`
	Role   string `json:"role"`
	Tracks string `json:"tracks"`
}

type ReleaseLabel struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Identifier struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

//--------------------------------------------------- TrackList ---------------------------------------------------

/*
Discogs tracklist XML structure
	<tracklist>
	  <track>
		 <position>A1</position>
		 <title>A Sea Apart</title>
		 <duration>5:08</duration>
	  </track>
	</tracklist>
*/
type Track struct {
	Position string `json:"position"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

//--------------------------------------------------- Video ---------------------------------------------------

/*
Discogs video XML structure
	<videos>
      <video duration="303" embed="true" src="https://www.youtube.com/watch?v=jRg3agJn1Mg">
         <title>Philip Glass - Powaqqatsi - 01. Serra Pelada</title>
         <description>Philip Glass - Powaqqatsi - 01. Serra Pelada</description>
      </video>
      <video duration="25" embed="true" src="https://www.youtube.com/watch?v=4E4gn9nNGYs">
         <title>Philip Glass - Powaqqatsi - 02. The Title</title>
         <description>Philip Glass - Powaqqatsi - 02. The Title</description>
      </video>
   </videos>
*/
type Video struct {
	Duration    string `json:"duration"`
	Embed       string `json:"embed"`
	Src         string `json:"src"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
