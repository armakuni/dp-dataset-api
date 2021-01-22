package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type CenStructure struct {
	Structure	KeyfamiliesStruct  `bson:"structure,omitempty"      json:"structure,omitempty"`
}

type KeyfamiliesStruct struct {
	Keyfamilies	 *KeyfamilyArray   `bson:"keyfamilies,omitempty"    json:"keyfamilies,omitempty"`
}

type KeyfamilyArray struct {
	Keyfamily   []KeyfamilyDetails `bson:"keyfamily,omitempty"      json:"keyfamily,omitempty"`
}

type KeyfamilyDetails struct {
	ID  		string             `bson:"id,omitempty"              json:"id,omitempty"`
	Annotations *AnnotationDetails `bson:"annotations,omitempty"     json:"annotations,omitempty"`
	Name 	    *NameDetails       `bson:"name,omitempty"            json:"name,omitempty"`
}

type AnnotationDetails struct {
	Annotation  []AnnoTextTitle    `bson:"annotation,omitempty"       json:"annotation,omitempty"`
}
type AnnoTextTitle struct {
	Text        interface{}       `bson:"annotationtext,omitempty"    json:"annotationtext,omitempty"`
	Title       string            `bson:"annotationtitle,omitempty"   json:"annotationtitle,omitempty"`
}
type NameDetails struct {
	Value 		string 			  `bson:"value,omitempty"             json:"value,omitempty"`
}

var (
	censusNationalStatistic = false
	fileName    string
	fullURLFile string
	title string
)

const (
	censusYear string = "2011"
	censusVersion string = "1"
)

func CensusContactDetails() models.ContactDetails{
	return models.ContactDetails{
		Email : "support@nomisweb.co.uk",
		Name :"Nomis",
		Telephone: "+44(0) 191 3342680",
	}
}

func main() {
    downloadFile()
	ctx := context.Background()
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Event(ctx, "failed to initialise mongo", log.FATAL, log.Error(err))
		os.Exit(1)
	}
	defer session.Close()
	fileLocation := "./NOMIS/def.sdmx.json"
	f, err := os.Open(fileLocation)
	if err != nil {
		log.Event(ctx, "failed to open file", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Event(ctx, "failed to read json file as a byte array", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	//var censusData CenStructure
	res := CenStructure{}
	if err := json.Unmarshal(b, &res); err != nil {
		logData := log.Data{"json file": res}
		log.Event(ctx, "failed to unmarshal json", log.ERROR, log.Error(err), logData)
		fmt.Println("error")
		return
	}

	for index0,_:=range res.Structure.Keyfamilies.Keyfamily {

		censusEditionData :=models.EditionUpdate{}
		mapData := models.Dataset{}
		cenId := res.Structure.Keyfamilies.Keyfamily[index0].ID
		mapData.Title= res.Structure.Keyfamilies.Keyfamily[index0].Name.Value
		mapData.ID = cenId


		datasetUrl:="http://127.0.0.1:12345/datasets/"
		instanceUrl:="http://127.0.0.1:12345/instances/"
		editionUrl:="/editions"
		versionUrl:="/versions"

		createEditionLink :=fmt.Sprintf("%s%s%s",datasetUrl,cenId,editionUrl)
		createLatestVersion:=fmt.Sprintf("%s%s%s%s%s%s",datasetUrl,cenId,editionUrl,"/"+censusYear,versionUrl,"/"+censusVersion)
		mapData.Links = &models.DatasetLinks{
			Editions: &models.LinkObject{HRef: createEditionLink},
			LatestVersion: &models.LinkObject{HRef: createLatestVersion},
			Self:&models.LinkObject{HRef: fmt.Sprintf("%s%s",datasetUrl,cenId)},
		}
		mapData.Contacts=[]models.ContactDetails{
			CensusContactDetails(),
		}

		mapData.License="Open Government Licence v3.0"
		mapData.NationalStatistic= &censusNationalStatistic
		mapData.NextRelease="To Be Confirmed"
		mapData.ReleaseFrequency ="Decennially"
		mapData.State="published"
		mapData.Type="nomis"

		//Model to generate editions document in mongodb
		generalModel:= &models.Edition{
			Edition: censusYear,
			Links: &models.EditionUpdateLinks{
				Dataset: &models.LinkObject{HRef: fmt.Sprintf("%s%s",datasetUrl,cenId), ID: cenId},
				LatestVersion: &models.LinkObject{HRef: fmt.Sprintf("%s%s%s%s%s%s%s",datasetUrl,cenId,editionUrl,censusYear,versionUrl,"/",censusVersion),ID: censusVersion},
				Self:&models.LinkObject{HRef: fmt.Sprintf("%s%s%s%s",datasetUrl,cenId,"/editions/",censusYear)},
				Versions: &models.LinkObject{HRef: fmt.Sprintf("%s%s%s%s%s",datasetUrl,cenId,editionUrl,"/"+censusYear,versionUrl)},
			},
			State:"published",
		}

		censusEditionData.ID=uuid.NewV4().String()
		censusEditionData.Next=generalModel
		censusEditionData.Current=generalModel

		//Model to generate instances documents in mongodb
		generateId:=uuid.NewV4().String()
		censusInstances := models.Version{
			Edition: censusYear,
			ID: generateId,
			LastUpdated: mapData.LastUpdated,
			Links: &models.VersionLinks{
				Dataset: &models.LinkObject{HRef: fmt.Sprintf("%s%s",datasetUrl,cenId), ID: cenId},
				Edition: &models.LinkObject{HRef: fmt.Sprintf("%s%s%s%s%s",datasetUrl,cenId,editionUrl,"/",censusYear),ID: censusYear},
				Self: &models.LinkObject{HRef: fmt.Sprintf("%s%s",instanceUrl,generateId)},
				Version: &models.LinkObject{HRef: fmt.Sprintf("%s%s%s%s%s%s",datasetUrl,cenId,"/editions/",censusYear,versionUrl,"/"+censusVersion),ID: censusVersion},
			},
			State: "published",
			Version: 1,
			UsageNotes:&[]models.UsageNote{},
		}


		for indx, _ := range res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation {

			str := res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Title
			switch str {
			case "MetadataText0":
				mapData.Description = res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)

			case "Keywords":
				keywrd := res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)
				var split = strings.Split(keywrd, ",")
				mapData.Keywords = split

			case "LastUpdated":
				tt := res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)
				t, _ := time.Parse("2006-01-02 15:04:05", tt)
				mapData.LastUpdated = t
				generalModel.LastUpdated = mapData.LastUpdated

			case "Units":
				mapData.UnitOfMeasure = res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)

			case "Mnemonic":
				ref := res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)
				param := strings.Split(ref, "c2011")
				mapData.NomisReferenceURL = "https://www.nomisweb.co.uk/census/2011/" + param[1]

			case"FirstReleased":
				releaseDt := res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)
				rd,_ := time.Parse("2006-01-02 15:04:05",releaseDt)
				censusInstances.ReleaseDate=rd.String()
			}

			if strings.HasPrefix(str, "MetadataTitle") {
				title= res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string)
			}
			if strings.HasPrefix(str, "MetadataText") {
				*censusInstances.UsageNotes=append(*censusInstances.UsageNotes,models.UsageNote{
					Note: res.Structure.Keyfamilies.Keyfamily[index0].Annotations.Annotation[indx].Text.(string),
					Title: title,
				})
			}
		}
		createDocument(ctx,mapData,session,"datasets")
		createDocument(ctx,censusEditionData,session,"editions")
		createDocument(ctx,censusInstances,session,"instances")
	}
}


//Inserts a document in the specific collection
func createDocument( ctx context.Context,class interface{}, session *mgo.Session,document string){
	var err error
	logData:= log.Data{"data": class}
	if err = session.DB("datasets").C(document).Insert(class); err != nil {
		log.Event(ctx, "failed to insert data in collection", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}
}


//Download a file from nomis website for census 2011 data
func downloadFile(){
	fullURLFile = "https://www.nomisweb.co.uk/api/v01/dataset/def.sdmx.json?search=*c2011*"

	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		fmt.Println(("error Parsing"))
		os.Exit(1)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]
	newFileName :="./NOMIS/"+fileName

	// Create blank file
	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Println(("error creating the file"))
		os.Exit(1)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Put content on file
	resp, err := client.Get(fullURLFile)

	if err != nil {
		fmt.Println(("error writing the file"))
		os.Exit(1)
	}
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)
}