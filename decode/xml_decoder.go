package decode

import (
	"encoding/xml"
	"errors"
	"github.com/lukasaron/discogs-parser/writer"
	"io"
	"log"
	"os"
	"strings"
)

var (
	// Errors returned when failure occurs
	wrongTypeSpecified    = errors.New("wrong file type specified")
	notCorrectStarElement = errors.New("token is not a correct start element")
)

// XML Decoder type is behaviour structure that implements Decoder interface and supports
// the Discogs XML dump data decoding.
type XMLDecoder struct {
	f   *os.File
	d   *xml.Decoder
	o   Options
	err error
}

func NewXmlDecoder(fileName string, options ...Options) Decoder {
	d := &XMLDecoder{}

	d.f, d.err = os.Open(fileName)
	d.d = xml.NewDecoder(d.f)

	if options != nil && len(options) > 0 {
		d.o = options[0]
	}

	return d
}

func (x *XMLDecoder) Close() error {
	return x.f.Close()
}

func (x *XMLDecoder) Options() Options {
	return x.o
}

func (x *XMLDecoder) SetOptions(opt Options) {
	x.o = opt
}

// Decode data
func (x *XMLDecoder) Decode(w writer.Writer) error {
	if x.err != nil {
		return x.err
	}

	if x.o.Block.Limit < 1 {
		x.o.Block.Limit = int(^uint(0) >> 1)
	}

	// get decode function based on the file type
	fn, err := x.decodeFunction()
	if err != nil {
		return err
	}

	for blockCount := 1; blockCount <= x.o.Block.Limit; blockCount++ {
		// call appropriate decoder function
		num, err := fn(w, x.o.Block.Size, blockCount > x.o.Block.Skip)
		if err != nil && err != io.EOF {
			log.Printf("Block %d failed [%d]\n", blockCount, num)
			return err
		}

		if num == 0 && err == io.EOF {
			break
		}

		if blockCount > x.o.Block.Skip {
			log.Printf("Block %d written [%d]\n", blockCount, num)
		} else {
			log.Printf("Block %d skipped [%d]\n", blockCount, num)
		}
	}

	return nil
}

func (x *XMLDecoder) Artists(limit int) (int, []Artist, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	artists := x.parseArtists(limit)
	if x.err == nil || x.err == io.EOF {
		artists = x.filterArtists(artists)
	}
	return len(artists), artists, x.err
}

func (x *XMLDecoder) Labels(limit int) (int, []Label, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	labels := x.parseLabels(limit)
	if x.err == nil || x.err == io.EOF {
		labels = x.filterLabels(labels)
	}
	return len(labels), labels, x.err
}

func (x *XMLDecoder) Masters(limit int) (int, []Master, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	masters := x.parseMasters(limit)
	if x.err == nil || x.err == io.EOF {
		masters = x.filterMasters(masters)
	}

	return len(masters), masters, x.err
}

func (x *XMLDecoder) Releases(limit int) (int, []Release, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	releases := x.parseReleases(limit)
	if x.err == nil || x.err == io.EOF {
		releases = x.filterReleases(releases)
	}
	return len(releases), releases, x.err
}

//--------------------------------------------------- FILTERS ---------------------------------------------------

func (x *XMLDecoder) filterArtists(as []Artist) []Artist {
	fa := make([]Artist, 0, len(as))
	for _, a := range as {
		if x.o.QualityLevel.Includes(ToQualityLevel(a.DataQuality)) {
			fa = append(fa, a)
		}
	}

	return fa
}

func (x *XMLDecoder) filterLabels(ls []Label) []Label {
	fl := make([]Label, 0, len(ls))
	for _, l := range ls {
		if x.o.QualityLevel.Includes(ToQualityLevel(l.DataQuality)) {
			fl = append(fl, l)
		}
	}

	return fl
}

func (x *XMLDecoder) filterMasters(ms []Master) []Master {
	fm := make([]Master, 0, len(ms))
	for _, m := range ms {
		if x.o.QualityLevel.Includes(ToQualityLevel(m.DataQuality)) {
			fm = append(fm, m)
		}
	}

	return fm
}

func (x *XMLDecoder) filterReleases(rs []Release) []Release {
	fr := make([]Release, 0, len(rs))
	for _, r := range rs {
		if x.o.QualityLevel.Includes(ToQualityLevel(r.DataQuality)) {
			fr = append(fr, r)
		}
	}

	return fr
}

//--------------------------------------------------- Decoders ---------------------------------------------------

func (x *XMLDecoder) decodeFunction() (func(writer.Writer, int, bool) (int, error), error) {
	switch x.o.FileType {
	case Artists:
		return x.decodeArtists, nil
	case Labels:
		return x.decodeLabels, nil
	case Masters:
		return x.decodeMasters, nil
	case Releases:
		return x.decodeReleases, nil
	case Unknown:
		fallthrough
	default:
		return nil, wrongTypeSpecified
	}
}

func (x *XMLDecoder) decodeArtists(w writer.Writer, blockSize int, write bool) (int, error) {
	num, a, err := x.Artists(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteArtists(a)
	}

	return num, err
}

func (x *XMLDecoder) decodeLabels(w writer.Writer, blockSize int, write bool) (int, error) {
	num, l, err := x.Labels(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteLabels(l)
	}

	return num, err
}

func (x *XMLDecoder) decodeMasters(w writer.Writer, blockSize int, write bool) (int, error) {
	num, m, err := x.Masters(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteMasters(m)
	}

	return num, err
}

func (x *XMLDecoder) decodeReleases(w writer.Writer, blockSize int, write bool) (int, error) {
	num, r, err := x.Releases(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteReleases(r)
	}

	return num, err
}

//--------------------------------------------------- HELPERS ---------------------------------------------------

func (x *XMLDecoder) startElement(token xml.Token) bool {
	_, ok := token.(xml.StartElement)
	return ok
}

func (x *XMLDecoder) endElement(token xml.Token) bool {
	_, ok := token.(xml.EndElement)
	return ok
}

func (x *XMLDecoder) startElementName(token xml.Token, name string) bool {
	se, ok := token.(xml.StartElement)
	return ok && se.Name.Local == name
}

func (x *XMLDecoder) endElementName(token xml.Token, name string) bool {
	ee, ok := token.(xml.EndElement)
	return ok && ee.Name.Local == name
}

func (x *XMLDecoder) parseValue() string {
	sb := strings.Builder{}
	for {
		t, _ := x.d.Token()
		if x.endElement(t) {
			break
		}

		if cr, ok := t.(xml.CharData); ok {
			sb.Write(cr)
		}
	}
	return sb.String()
}

func (x *XMLDecoder) parseChildValues(parentName, childName string) (children []string) {
	for {
		t, _ := x.d.Token()
		if x.startElementName(t, childName) {
			children = append(children, x.parseValue())
		}
		if x.endElementName(t, parentName) {
			break
		}
	}
	return children
}

//=================================================== Models ===================================================

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

//=================================================== Parsers ===================================================

//--------------------------------------------------- Artist ---------------------------------------------------

func (x *XMLDecoder) parseArtists(limit int) (artists []Artist) {
	if x.err != nil {
		return artists
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.startElementName(t, "artist") {
			artist := x.parseArtist(t.(xml.StartElement))
			if x.err != nil {
				return artists
			}
			artists = append(artists, artist)
			cnt++
		}
	}

	return artists
}

func (x *XMLDecoder) parseArtist(se xml.StartElement) (artist Artist) {
	if x.err != nil {
		return artist
	}

	if se.Name.Local != "artist" {
		x.err = notCorrectStarElement
		return artist
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil && !x.endElementName(t, "artist"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return artist
				}

				artist.Images = imgs
			case "id":
				artist.Id = x.parseValue()
			case "name":
				artist.Name = x.parseValue()
			case "realname":
				artist.RealName = x.parseValue()
			case "namevariations":
				artist.NameVariations = x.parseChildValues("namevariations", "name")
			case "members":
				artist.Members = x.parseMembers()
			case "aliases":
				artist.Aliases = x.parseAliases()
			case "profile":
				artist.Profile = x.parseValue()
			case "data_quality":
				artist.DataQuality = x.parseValue()
			case "urls":
				artist.Urls = x.parseChildValues("urls", "url")
			}
		}
	}

	return artist
}

func (x *XMLDecoder) parseAliases() (aliases []Alias) {
	if x.err != nil {
		return
	}
	var t xml.Token

	for t, x.err = x.d.Token(); x.err == nil && !x.endElementName(t, "aliases"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			alias := Alias{
				Id:   se.Attr[0].Value,
				Name: x.parseValue(),
			}
			aliases = append(aliases, alias)
		}
	}

	return aliases
}

func (x *XMLDecoder) parseMembers() (members []Member) {
	if x.err != nil {
		return
	}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil && !x.endElementName(t, "members"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			member := Member{
				Id:   se.Attr[0].Value,
				Name: x.parseValue(),
			}
			members = append(members, member)
		}
	}

	return members

}

//--------------------------------------------------- Company ---------------------------------------------------

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

//--------------------------------------------------- Format ---------------------------------------------------

func (x *XMLDecoder) parseFormats() (formats []Format) {
	if x.err != nil {
		return formats
	}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "formats" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "format" {
			format := Format{}
			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "qty":
					format.Quantity = attr.Value
				case "name":
					format.Name = attr.Value
				case "text":
					format.Text = attr.Value
				}
			}

			format.Descriptions = x.parseChildValues("descriptions", "description")
			formats = append(formats, format)
		}
	}
	return formats
}

//--------------------------------------------------- Image ---------------------------------------------------

func (x *XMLDecoder) parseImages(se xml.StartElement) (images []Image) {
	if x.err != nil {
		return images
	}

	if se.Name.Local != "images" {
		x.err = notCorrectStarElement
		return images
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "image" {
			img := x.parseImage(se)
			if x.err != nil {
				return images
			}

			images = append(images, img)
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "images" {
			break
		}
	}

	return images
}

func (x *XMLDecoder) parseImage(se xml.StartElement) (img Image) {
	if x.err != nil {
		return img
	}

	if se.Name.Local != "image" {
		x.err = notCorrectStarElement
		return img
	}

	for _, attr := range se.Attr {
		switch attr.Name.Local {
		case "height":
			img.Height = attr.Value
		case "width":
			img.Width = attr.Value
		case "type":
			img.Type = attr.Value
		case "uri":
			img.Uri = attr.Value
		case "uri150":
			img.Uri150 = attr.Value
		}
	}

	return img
}

//--------------------------------------------------- Label ---------------------------------------------------

func (x *XMLDecoder) parseLabels(limit int) (labels []Label) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.startElementName(t, "label") {
			l := x.parseLabel(t.(xml.StartElement))
			if x.err != nil {
				return labels
			}

			labels = append(labels, l)
			cnt++
		}
	}

	return labels
}

func (x *XMLDecoder) parseLabel(se xml.StartElement) (label Label) {
	if x.err != nil {
		return label
	}

	if se.Name.Local != "label" {
		x.err = notCorrectStarElement
		return label
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return label
				}

				label.Images = imgs
			case "id":
				label.Id = x.parseValue()
			case "name":
				label.Name = x.parseValue()
			case "contactinfo":
				label.ContactInfo = x.parseValue()
			case "profile":
				label.Profile = x.parseValue()
			case "urls":
				label.Urls = x.parseChildValues("urls", "url")
			case "sublabels":
				label.SubLabels = x.parseSubLabels()
			case "data_quality":
				label.DataQuality = x.parseValue()
			case "parentLabel":
				label.ParentLabel = &LabelLabel{
					Id:   se.Attr[0].Value,
					Name: x.parseValue(),
				}
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "label" {
			break
		}
	}

	return label
}

func (x *XMLDecoder) parseSubLabels() (labels []LabelLabel) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "sublabels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := LabelLabel{}
			label.Id = se.Attr[0].Value
			label.Name = x.parseValue()
			labels = append(labels, label)
		}
	}

	return labels
}

//--------------------------------------------------- Master ---------------------------------------------------

func (x *XMLDecoder) parseMasters(limit int) (masters []Master) {
	if x.err != nil {
		return masters
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.startElementName(t, "master") {
			m := x.parseMaster(t.(xml.StartElement))
			if x.err != nil {
				return masters
			}

			masters = append(masters, m)
			cnt++
		}
	}

	return masters
}

func (x *XMLDecoder) parseMaster(se xml.StartElement) (master Master) {
	if x.err != nil {
		return master
	}

	if se.Name.Local != "master" {
		x.err = notCorrectStarElement
		return master
	}

	master.Id = se.Attr[0].Value
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return master
				}
				master.Images = imgs
			case "main_release":
				master.MainRelease = x.parseValue()
			case "artists":
				master.Artists = x.parseReleaseArtists("artists")
			case "genres":
				master.Genres = x.parseChildValues("genres", "genre")
			case "styles":
				master.Styles = x.parseChildValues("styles", "style")
			case "year":
				master.Year = x.parseValue()
			case "title":
				master.Title = x.parseValue()
			case "data_quality":
				master.DataQuality = x.parseValue()
			case "videos":
				master.Videos = x.parseVideos()
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "master" {
			break
		}
	}

	return master
}

//--------------------------------------------------- Release ---------------------------------------------------

func (x *XMLDecoder) parseReleases(limit int) (releases []Release) {
	if x.err != nil {
		return releases
	}

	var t xml.Token
	cnt := 0
	for t, x.err = x.d.Token(); t != nil && x.err == nil && cnt != limit; t, x.err = x.d.Token() {
		if x.startElementName(t, "release") {
			rls := x.parseRelease(t.(xml.StartElement))
			if x.err != nil {
				return releases
			}

			releases = append(releases, rls)
			cnt++
		}
	}

	return releases
}

func (x *XMLDecoder) parseRelease(se xml.StartElement) (release Release) {
	if x.err != nil {
		return release
	}

	if se.Name.Local != "release" {
		x.err = notCorrectStarElement
		return release
	}

	for _, attr := range se.Attr {
		switch attr.Name.Local {
		case "id":
			release.Id = attr.Value
		case "status":
			release.Status = attr.Value
		}
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "images":
				imgs := x.parseImages(se)
				if x.err != nil {
					return release
				}
				release.Images = imgs
			case "artists":
				release.Artists = x.parseReleaseArtists("artists")
			case "extraartists":
				release.ExtraArtists = x.parseReleaseArtists("extraartists")
			case "title":
				release.Title = x.parseValue()
			case "labels":
				release.Labels = x.parseReleaseLabels()
			case "formats":
				release.Formats = x.parseFormats()
			case "genres":
				release.Genres = x.parseChildValues("genres", "genre")
			case "styles":
				release.Styles = x.parseChildValues("styles", "style")
			case "country":
				release.Country = x.parseValue()
			case "released":
				release.Released = x.parseValue()
			case "notes":
				release.Notes = x.parseValue()
			case "data_quality":
				release.DataQuality = x.parseValue()
			case "master_id":
				release.MainRelease = se.Attr[0].Value
				release.MasterId = x.parseValue()
			case "tracklist":
				release.TrackList = x.parseTrackList()
			case "identifiers":
				release.Identifiers = x.parseIdentifiers()
			case "videos":
				release.Videos = x.parseVideos()
			case "companies":
				release.Companies = x.parseCompanies()
			}
		}
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "release" {
			break
		}
	}

	return release
}

func (x *XMLDecoder) parseReleaseArtists(wrapperName string) (artists []ReleaseArtist) {
	if x.err != nil {
		return artists
	}

	artist := ReleaseArtist{}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == wrapperName {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "id":
				artist.Id = x.parseValue()
			case "name":
				artist.Name = x.parseValue()
			case "anv":
				artist.Anv = x.parseValue()
			case "join":
				artist.Join = x.parseValue()
			case "role":
				artist.Role = x.parseValue()
			case "tracks":
				artist.Tracks = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "artist" {
			artists = append(artists, artist)
			artist = ReleaseArtist{}
		}
	}
	return artists
}

func (x *XMLDecoder) parseReleaseLabels() (labels []ReleaseLabel) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "labels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := ReleaseLabel{}

			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "id":
					label.Id = attr.Value
				case "name":
					label.Name = attr.Value
				case "catno":
					label.Category = attr.Value
				}
			}

			labels = append(labels, label)
		}
	}
	return labels
}

func (x *XMLDecoder) parseIdentifiers() (identifiers []Identifier) {
	if x.err != nil {
		return identifiers
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "identifiers" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "identifier" {
			identifier := Identifier{}
			for _, attr := range se.Attr {
				switch attr.Name.Local {
				case "description":
					identifier.Description = attr.Value
				case "type":
					identifier.Type = attr.Value
				case "value":
					identifier.Value = attr.Value
				}
			}

			identifiers = append(identifiers, identifier)
		}
	}
	return identifiers
}

//--------------------------------------------------- TrackList ---------------------------------------------------

func (x *XMLDecoder) parseTrackList() (trackList []Track) {
	if x.err != nil {
		return trackList
	}

	track := Track{}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "tracklist" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "position":
				track.Position = x.parseValue()
			case "title":
				track.Title = x.parseValue()
			case "duration":
				track.Duration = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "track" {
			trackList = append(trackList, track)
			track = Track{}
		}
	}

	return trackList
}

//--------------------------------------------------- Video ---------------------------------------------------

func (x *XMLDecoder) parseVideos() (videos []Video) {
	if x.err != nil {
		return videos
	}

	video := Video{}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "videos" {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			switch se.Name.Local {
			case "video":
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "duration":
						video.Duration = attr.Value
					case "embed":
						video.Embed = attr.Value
					case "src":
						video.Src = attr.Value
					}
				}
			case "title":
				video.Title = x.parseValue()
			case "description":
				video.Description = x.parseValue()
			}
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "video" {
			videos = append(videos, video)
			video = Video{}
		}
	}

	return videos
}
