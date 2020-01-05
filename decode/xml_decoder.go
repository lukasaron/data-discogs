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
