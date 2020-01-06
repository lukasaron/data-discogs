package decode

import (
	"encoding/xml"
	"errors"
	"github.com/lukasaron/discogs-parser/model"
	"github.com/lukasaron/discogs-parser/write"
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

func NewXmlDecoder(fileName string, options *Options) Decoder {
	d := &XMLDecoder{}

	d.f, d.err = os.Open(fileName)
	d.d = xml.NewDecoder(d.f)

	if options != nil {
		d.o = *options
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

func (x *XMLDecoder) Reset() error {
	_, x.err = x.f.Seek(0, 0)
	return x.err
}

// Decode data
func (x *XMLDecoder) Decode(w write.Writer) error {
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

func (x *XMLDecoder) Artists(limit int) (int, []model.Artist, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	artists := x.parseArtists(limit)
	if x.err == nil || x.err == io.EOF {
		artists = x.filterArtists(artists)
	}
	return len(artists), artists, x.err
}

func (x *XMLDecoder) Labels(limit int) (int, []model.Label, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	labels := x.parseLabels(limit)
	if x.err == nil || x.err == io.EOF {
		labels = x.filterLabels(labels)
	}
	return len(labels), labels, x.err
}

func (x *XMLDecoder) Masters(limit int) (int, []model.Master, error) {
	if x.err != nil {
		return 0, nil, x.err
	}

	masters := x.parseMasters(limit)
	if x.err == nil || x.err == io.EOF {
		masters = x.filterMasters(masters)
	}

	return len(masters), masters, x.err
}

func (x *XMLDecoder) Releases(limit int) (int, []model.Release, error) {
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

func (x *XMLDecoder) filterArtists(as []model.Artist) []model.Artist {
	fa := make([]model.Artist, 0, len(as))
	for _, a := range as {
		if x.o.QualityLevel.Includes(ToQualityLevel(a.DataQuality)) {
			fa = append(fa, a)
		}
	}

	return fa
}

func (x *XMLDecoder) filterLabels(ls []model.Label) []model.Label {
	fl := make([]model.Label, 0, len(ls))
	for _, l := range ls {
		if x.o.QualityLevel.Includes(ToQualityLevel(l.DataQuality)) {
			fl = append(fl, l)
		}
	}

	return fl
}

func (x *XMLDecoder) filterMasters(ms []model.Master) []model.Master {
	fm := make([]model.Master, 0, len(ms))
	for _, m := range ms {
		if x.o.QualityLevel.Includes(ToQualityLevel(m.DataQuality)) {
			fm = append(fm, m)
		}
	}

	return fm
}

func (x *XMLDecoder) filterReleases(rs []model.Release) []model.Release {
	fr := make([]model.Release, 0, len(rs))
	for _, r := range rs {
		if x.o.QualityLevel.Includes(ToQualityLevel(r.DataQuality)) {
			fr = append(fr, r)
		}
	}

	return fr
}

//--------------------------------------------------- Decoders ---------------------------------------------------

func (x *XMLDecoder) decodeFunction() (func(write.Writer, int, bool) (int, error), error) {
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

func (x *XMLDecoder) decodeArtists(w write.Writer, blockSize int, write bool) (int, error) {
	num, a, err := x.Artists(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteArtists(a)
	}

	return num, err
}

func (x *XMLDecoder) decodeLabels(w write.Writer, blockSize int, write bool) (int, error) {
	num, l, err := x.Labels(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteLabels(l)
	}

	return num, err
}

func (x *XMLDecoder) decodeMasters(w write.Writer, blockSize int, write bool) (int, error) {
	num, m, err := x.Masters(blockSize)
	if (err != nil && err != io.EOF) || num == 0 {
		return num, err
	}

	if write {
		return num, w.WriteMasters(m)
	}

	return num, err
}

func (x *XMLDecoder) decodeReleases(w write.Writer, blockSize int, write bool) (int, error) {
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

//=================================================== Parsers ===================================================

//--------------------------------------------------- Artist ---------------------------------------------------

func (x *XMLDecoder) parseArtists(limit int) (artists []model.Artist) {
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

func (x *XMLDecoder) parseArtist(se xml.StartElement) (artist model.Artist) {
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

func (x *XMLDecoder) parseAliases() (aliases []model.Alias) {
	if x.err != nil {
		return
	}
	var t xml.Token

	for t, x.err = x.d.Token(); x.err == nil && !x.endElementName(t, "aliases"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			alias := model.Alias{
				Id:   se.Attr[0].Value,
				Name: x.parseValue(),
			}
			aliases = append(aliases, alias)
		}
	}

	return aliases
}

func (x *XMLDecoder) parseMembers() (members []model.Member) {
	if x.err != nil {
		return
	}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil && !x.endElementName(t, "members"); t, x.err = x.d.Token() {
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "name" {
			member := model.Member{
				Id:   se.Attr[0].Value,
				Name: x.parseValue(),
			}
			members = append(members, member)
		}
	}

	return members

}

//--------------------------------------------------- Company ---------------------------------------------------

func (x *XMLDecoder) parseCompanies() (companies []model.Company) {
	if x.err != nil {
		return companies
	}

	company := model.Company{}
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
			company = model.Company{}
		}
	}
	return companies
}

//--------------------------------------------------- Format ---------------------------------------------------

func (x *XMLDecoder) parseFormats() (formats []model.Format) {
	if x.err != nil {
		return formats
	}
	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {

		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "formats" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "format" {
			format := model.Format{}
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

func (x *XMLDecoder) parseImages(se xml.StartElement) (images []model.Image) {
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

func (x *XMLDecoder) parseImage(se xml.StartElement) (img model.Image) {
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

func (x *XMLDecoder) parseLabels(limit int) (labels []model.Label) {
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

func (x *XMLDecoder) parseLabel(se xml.StartElement) (label model.Label) {
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
				label.ParentLabel = &model.LabelLabel{
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

func (x *XMLDecoder) parseSubLabels() (labels []model.LabelLabel) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "sublabels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := model.LabelLabel{}
			label.Id = se.Attr[0].Value
			label.Name = x.parseValue()
			labels = append(labels, label)
		}
	}

	return labels
}

//--------------------------------------------------- Master ---------------------------------------------------

func (x *XMLDecoder) parseMasters(limit int) (masters []model.Master) {
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

func (x *XMLDecoder) parseMaster(se xml.StartElement) (master model.Master) {
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

func (x *XMLDecoder) parseReleases(limit int) (releases []model.Release) {
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

func (x *XMLDecoder) parseRelease(se xml.StartElement) (release model.Release) {
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

func (x *XMLDecoder) parseReleaseArtists(wrapperName string) (artists []model.ReleaseArtist) {
	if x.err != nil {
		return artists
	}

	artist := model.ReleaseArtist{}
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
			artist = model.ReleaseArtist{}
		}
	}
	return artists
}

func (x *XMLDecoder) parseReleaseLabels() (labels []model.ReleaseLabel) {
	if x.err != nil {
		return labels
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "labels" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "label" {
			label := model.ReleaseLabel{}

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

func (x *XMLDecoder) parseIdentifiers() (identifiers []model.Identifier) {
	if x.err != nil {
		return identifiers
	}

	var t xml.Token
	for t, x.err = x.d.Token(); x.err == nil; t, x.err = x.d.Token() {
		if ee, ok := t.(xml.EndElement); ok && ee.Name.Local == "identifiers" {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Local == "identifier" {
			identifier := model.Identifier{}
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

func (x *XMLDecoder) parseTrackList() (trackList []model.Track) {
	if x.err != nil {
		return trackList
	}

	track := model.Track{}

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
			track = model.Track{}
		}
	}

	return trackList
}

//--------------------------------------------------- Video ---------------------------------------------------

func (x *XMLDecoder) parseVideos() (videos []model.Video) {
	if x.err != nil {
		return videos
	}

	video := model.Video{}

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
			video = model.Video{}
		}
	}

	return videos
}
