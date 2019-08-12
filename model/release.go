package model

type Release struct {
	Id           string
	Status       string
	Images       []*Image
	Artists      []*ReleaseArtist
	ExtraArtists []*ReleaseArtist
	Title        string
	Formats      []*ReleaseFormat
	Genres       []string
	Styles       []string
	Country      string
	Released     string
	Notes        string
	DataQuality  string
	MasterId     string
	MainRelease  string
	TrackList    []*ReleaseTrack
	Identifiers  []*ReleaseIdentifier
	Videos       []*ReleaseVideo
	Labels       []*ReleaseLabel
	Companies    []*ReleaseCompany
}

type ReleaseArtist struct {
	Id     string
	Name   string
	Join   string
	Anv    string
	Role   string
	Tracks string
}

type ReleaseLabel struct {
	Id       string
	Name     string
	Category string
}

type ReleaseFormat struct {
	Name         string
	Quantity     string
	Text         string
	Descriptions []string
}

type ReleaseTrack struct {
	Position string
	Title    string
	Duration string
}

type ReleaseIdentifier struct {
	Description string
	Type        string
	Value       string
}
type ReleaseVideo struct {
	Duration    string
	Embed       string
	Src         string
	Title       string
	Description string
}

type ReleaseCompany struct {
	Id             string
	Name           string
	Category       string
	EntityType     string
	EntityTypeName string
	ResourceUrl    string
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<release id="2" status="Accepted">
   <images>
      <image height="394" type="primary" uri="" uri150="" width="400" />
      <image height="600" type="secondary" uri="" uri150="" width="600" />
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
      <artist>
         <id>26</id>
         <name>Alexi Delano</name>
         <anv>A. Delano</anv>
         <join />
         <role>Written-By</role>
         <tracks />
      </artist>
      <artist>
         <id>27</id>
         <name>Cari Lekebusch</name>
         <anv>C. Lekebusch</anv>
         <join />
         <role>Written-By</role>
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
      <style>Tech House</style>
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
      <track>
         <position>A2</position>
         <title>Dutchmaster</title>
         <duration>4:21</duration>
      </track>
      <track>
         <position>B1</position>
         <title>Inner City Lullaby</title>
         <duration>4:22</duration>
      </track>
      <track>
         <position>B2</position>
         <title>Yeah Kid!</title>
         <duration>4:46</duration>
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
      <video duration="316" embed="true" src="https://www.youtube.com/watch?v=6iZu93js2oM">
         <title>Mr. James Barth &amp; A.D. - A Sea Apart</title>
         <description>Mr. James Barth &amp; A.D. - A Sea Apart</description>
      </video>
      <video duration="271" embed="true" src="https://www.youtube.com/watch?v=yvwrhp0hrv8">
         <title>Mr. James Barth &amp; A.D. - Inner City Lullaby</title>
         <description>Mr. James Barth &amp; A.D. - Inner City Lullaby</description>
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
