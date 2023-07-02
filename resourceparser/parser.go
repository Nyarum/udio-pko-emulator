package resourceparser

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/gelsrc/go-charset"
	"github.com/jszwec/csvutil"
)

type ItemPre struct {
	ID   string `csv:"id"`
	Name string `csv:"name"`
}

type ItemInfo struct {
	ID                              string `csv:"ID"`
	Name                            string `csv:"Name"`
	ICON                            string `csv:"ICON"`
	Model1                          string `csv:"Model 1"`
	Model2                          string `csv:"Model 2"`
	Model3                          string `csv:"Model 3"`
	Model4                          string `csv:"Model 4"`
	Model5                          string `csv:"Model 5"`
	ShipSymbol                      string `csv:"Ship Symbol"`
	ShipSizNumber                   string `csv:"ship siz number"`
	Type                            string `csv:"Type"`
	BbtainPrefixRate                string `csv:"obtain prefix rate"`
	SetID                           string `csv:"set ID"`
	ForgingLevel                    string `csv:"Forging Level"`
	StableValue                     string `csv:"Stable value"`
	OnlyID                          string `csv:"only ID"`
	Trade                           string `csv:"Trade"`
	Picked                          string `csv:"picked"`
	Discard                         string `csv:"Discard"`
	ConfirmToDelete                 string `csv:"Confirm to delete"`
	Stackable                       string `csv:"stackable"`
	IsItInstantiation               string `csv:"is it instantiation"`
	Price                           string `csv:"Price"`
	Size                            string `csv:"Size"`
	CharacterLevel                  string `csv:"Character Level"`
	Class                           string `csv:"Class"`
	CharacterNick                   string `csv:"Character Nick"`
	CharacterReputation             string `csv:"Character Reputation"`
	ItemCanequipLocation            string `csv:"item can equip location"`
	ItemSwitchLocation              string `csv:"item switch location"`
	ItemObtainIntoLocationDetermine string `csv:"item obtain into location determine"`
	StrModulusBonus                 string `csv:"Str modulus bonus"`
	AgiModulusBonus                 string `csv:"Agi modulus bonus"`
	DexModulusBonus                 string `csv:"Dex modulus bonus"`
	ConModulusBonus                 string `csv:"Con modulus bonus"`
	SprModulusBonus                 string `csv:"Spr modulus bonus"`
	LukModulusBonus                 string `csv:"Luk modulus bonus"`
	HitRateModulusBonus             string `csv:"Hit rate modulus bonus"`
	HuiZnaet3                       string `csv:"-"`
	MinAttackModulusBonus           string `csv:"Min attack modulus bonus"`
	MaxAttackModulusBonus           string `csv:"Max Attack modulus bonus"`
	DefenseModulusBonus             string `csv:"Defense modulus bonus     "`
	MaxHPModulusBonus               string `csv:"Max HP modulus bonus"`
	MxspModulusBonus                string `csv:"Mxsp modulus bonus"`
	FleeModulusBonus                string `csv:"flee modulus bonus"`
	HitModulusBonus                 string `csv:"Hit modulus bonus"`
	CrtModulusBonus                 string `csv:"crt modulus bonus"`
	MfModulusBonus                  string `csv:"mf modulus bonus"`
	HrecModulusBonus                string `csv:"hrec modulus bonus"`
	SrecModulusBonus                string `csv:"srec modulus bonus"`
	MspdModulusBonus                string `csv:"mspd modulus bonus"`
	ColModulusBonus                 string `csv:"col modulus bonus"`
	StrconstantBonus                string `csv:"Strconstant bonus"`
	Agiconstantbonus                string `csv:"Agiconstant bonus"`
	Dexconstantbonus                string `csv:"Dexconstant bonus"`
	Conconstantbonus                string `csv:"Conconstant bonus"`
	Staconstantbonus                string `csv:"Staconstant bonus"`
	Lukconstantbonus                string `csv:"Lukconstant bonus"`
	HuiZnaetChtoETO                 string `csv:"-"`
	AttackRangeConstantBonus        string `csv:"Attack range constant bonus  "`
	MinAttackConstantBonus          string `csv:"Min Attack constant bonus   "`
	MaxAttackConstantBonus          string `csv:"Max Attack constant bonus"`
	HuiZnaetChtoETO2                string `csv:"-"`
	MaxHPConstantBonus              string `csv:"Max HP constant bonus"`
	MxspconstantBonus               string `csv:"Mxspconstant bonus"`
	FleeconstantBonus               string `csv:"fleeconstant bonus"`
	HitconstantBonus                string `csv:"hitconstant bonus"`
	CrtconstantBonus                string `csv:"crtconstant bonus"`
	MfconstantBonus                 string `csv:"mfconstant bonus"`
	HrecconstantBonus               string `csv:"hrecconstant bonus"`
	SrecconstantBonus               string `csv:"srecconstant bonus"`
	MspdconstantBonus               string `csv:"mspdconstant bonus"`
	ColconstantBonus                string `csv:"colconstant bonus"`
	PhysicalResist                  string `csv:"Physical Resist"`
	ItemLeftHandExertIdentifier     string `csv:"item left hand exert identifier"`
	ItemEnergy                      string `csv:"Item Energy"`
	Durability                      string `csv:"Durability"`
	MaxInstantiationHoleValue       string `csv:"Max instantiation hole value"`
	ShipDurabilityRecovered         string `csv:"Ship durability recovered"`
	CanContainCannonQuantity        string `csv:"can contain cannon quantity"`
	ShipMemberCount                 string `csv:"ship member count"`
	MemberLabel                     string `csv:"member label"`
	CargoCapacity                   string `csv:"Cargo Capacity"`
	FuelConsumption                 string `csv:"Fuel consumption"`
	CannonballPathOfFlightSpeed     string `csv:"Cannonball Path of Flight speed"`
	ShipMovementSpeed               string `csv:"ship movement speed"`
	UsageEffect                     string `csv:"usage effect"`
	DisplayEffect                   string `csv:"display effect"`
	ItemBindEffect                  string `csv:"item bind effect"`
	ItemBindEffectDummy             string `csv:"item bind effect dummy"`
	DisplayItemEffectItem           string `csv:"Display item effect (item put at object slot 1 to show effect)"`
	ItemDropModelEffect             string `csv:"item drop model effect"`
	ItemUsageEffect                 string `csv:"item usage effect"`
	DescriptionItemLevel            string `csv:"Description (Item level)"`
	Remark                          string `csv:"Remark"`
}

type Parser struct {
	ItemPre  []ItemPre
	ItemInfo []ItemInfo
}

func NewParser() (*Parser, error) {
	p := &Parser{}
	err := p.LoadAndParse()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) LoadAndParse() error {
	files, err := ioutil.ReadDir("./resource")
	if err != nil {
		return err
	}

	for _, f := range files {
		switch f.Name() {
		case "ItemPre.txt":
			itemPre, err := os.Open("./resource/" + f.Name())
			if err != nil {
				return err
			}

			csvReader := csv.NewReader(itemPre)
			csvReader.Comma = '\t'

			dec, err := csvutil.NewDecoder(csvReader)
			if err != nil {
				log.Fatal(err)
			}

			for {
				var itemPre ItemPre

				if err := dec.Decode(&itemPre); err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				itemPre.Name = string(charset.Cp1251BytesToRunes([]byte(itemPre.Name)))

				p.ItemPre = append(p.ItemPre, itemPre)
			}
		case "ItemInfo.txt":
			itemInfo, err := os.Open("./resource/" + f.Name())
			if err != nil {
				return err
			}

			csvReader := csv.NewReader(itemInfo)
			csvReader.Comma = '\t'
			csvReader.LazyQuotes = true

			dec, err := csvutil.NewDecoder(csvReader)
			if err != nil {
				log.Fatal(err)
			}

			for {
				var itemInfo ItemInfo

				if err := dec.Decode(&itemInfo); err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				itemInfo.Name = string(charset.Cp1251BytesToRunes([]byte(itemInfo.Name)))
				itemInfo.DescriptionItemLevel = string(charset.Cp1251BytesToRunes([]byte(itemInfo.DescriptionItemLevel)))
				itemInfo.Remark = string(charset.Cp1251BytesToRunes([]byte(itemInfo.Remark)))

				p.ItemInfo = append(p.ItemInfo, itemInfo)
			}
		}
	}

	return nil
}
