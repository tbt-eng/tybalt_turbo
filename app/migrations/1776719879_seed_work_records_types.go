package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

var workRecordConsumableSeeds = []struct {
	id             string
	name           string
	inputKind      string
	unitLabel      string
	workRecordType string
	sortOrder      int
	relationID     string
}{
	{id: "egef5amyx680t0k", name: "Standpipe", inputKind: "number", unitLabel: "m", workRecordType: "4k4nxilem8k10cu", sortOrder: 10, relationID: "1btvyuf420z6y1m"},
	{id: "e8dgdcy1ggbitb2", name: "Bentonite", inputKind: "number", unitLabel: "bags", workRecordType: "4k4nxilem8k10cu", sortOrder: 20, relationID: "pl5zn8abmbfjm55"},
	{id: "cum6foslqon9px7", name: "Cement", inputKind: "number", unitLabel: "bags", workRecordType: "4k4nxilem8k10cu", sortOrder: 30, relationID: "ro4dzuq2a2k5e11"},
	{id: "3epvus1t4vrfg7s", name: "Sand", inputKind: "number", unitLabel: "bags", workRecordType: "4k4nxilem8k10cu", sortOrder: 40, relationID: "fv737qi1s3b2els"},
	{id: "jdr504icbp7u1q9", name: "Grout", inputKind: "number", unitLabel: "bags", workRecordType: "4k4nxilem8k10cu", sortOrder: 50, relationID: "msk65izjb34xftf"},
	{id: "1atxe3uvy6q380z", name: "Well Protectors", inputKind: "number", unitLabel: "each", workRecordType: "4k4nxilem8k10cu", sortOrder: 60, relationID: "p47fxlh8u4352k5"},
	{id: "o078ulmrklkqbyf", name: "Well Riser", inputKind: "number", unitLabel: "m", workRecordType: "4k4nxilem8k10cu", sortOrder: 70, relationID: "i81vqr3dlca7qyz"},
	{id: "y3fkrpjlc52lbvk", name: "Well Screen", inputKind: "number", unitLabel: "m", workRecordType: "4k4nxilem8k10cu", sortOrder: 80, relationID: "16fqwc20snnvnbe"},
	{id: "12irhin2knp7y1f", name: "End Caps", inputKind: "number", unitLabel: "each", workRecordType: "4k4nxilem8k10cu", sortOrder: 90, relationID: "3a5y8q42tmxwobr"},
	{id: "ti4gn1llghc6ecl", name: "J-Plugs", inputKind: "number", unitLabel: "each", workRecordType: "4k4nxilem8k10cu", sortOrder: 100, relationID: "5aom5dya7v6jc1w"},
	{id: "pchicif088d5zj2", name: "SIBs", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 10, relationID: "kfz4rk947j9evbz"},
	{id: "r5eahts6k9gnpph", name: "SSIBs", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 20, relationID: "at5rs5ub2c3n2ri"},
	{id: "tsrkmrd3186ya08", name: "IBs", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 30, relationID: "zr5dap0m9w7rbcx"},
	{id: "511vjmuhczl1c3r", name: "Plastic Bars", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 40, relationID: "1d1a51t6b2lujxs"},
	{id: "rn9n9xx9ov4jxo6", name: "Magnails", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 50, relationID: "0hz38ez6whb1gbh"},
	{id: "ihk7pa0s5n48ykf", name: "Rock Plug", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 60, relationID: "l7qio6d4w0dr06r"},
	{id: "8fdde0cwslewp9k", name: "Wood Stakes", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 70, relationID: "5rry77wyn0bx2ce"},
	{id: "d3b192m6wawdvqg", name: "Magnetic Nails", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 80, relationID: "ed11tjqtb1h38eo"},
	{id: "379ul4dtkevuikz", name: "Monument Marker Posts", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 90, relationID: "5birek34mzdm292"},
	{id: "je6uuu8p0caa49e", name: "Monument Marker Plaques", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 100, relationID: "ckqognq3phc58ov"},
	{id: "om4tyzt1pnjvhtd", name: "Painted Wood Stakes", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 110, relationID: "i3nvif9efp4kg5q"},
	{id: "qklevs2r24aqdji", name: "Spray Paint", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 120, relationID: "pj50k4fxex3z9ei"},
	{id: "rhlo38cctuix0u6", name: "Flagging Tape", inputKind: "number", unitLabel: "rolls", workRecordType: "oask5u50nka0km8", sortOrder: 130, relationID: "1sk6wi1vh0x3a4a"},
	{id: "va117idnpj6uuyd", name: "Eco Stakes", inputKind: "number", unitLabel: "each", workRecordType: "oask5u50nka0km8", sortOrder: 140, relationID: "o1i3y1wjpqc9fbm"},
	{id: "3e4kq37ckzmnha3", name: "Nitrile Gloves", inputKind: "number", unitLabel: "pair", workRecordType: "up9ndzdnyzv8snx", sortOrder: 10, relationID: "5ax4ldfwa8hg50e"},
	{id: "zs65q5wpwsfagte", name: "Bailers", inputKind: "number", unitLabel: "each", workRecordType: "up9ndzdnyzv8snx", sortOrder: 20, relationID: "lem8svi314vksqh"},
	{id: "fqjzoj8vg8d0omg", name: "Foot Valves", inputKind: "number", unitLabel: "each", workRecordType: "up9ndzdnyzv8snx", sortOrder: 30, relationID: "g2ux0v2wpril2dh"},
	{id: "bd963vl9ogt23rs", name: "Tubing (Waterra)", inputKind: "number", unitLabel: "m", workRecordType: "up9ndzdnyzv8snx", sortOrder: 40, relationID: "59denu0x18lhjz1"},
	{id: "djxxxntwwtvtt1n", name: "Tubing (VOC)", inputKind: "number", unitLabel: "m", workRecordType: "up9ndzdnyzv8snx", sortOrder: 45, relationID: "hwpmvaw90eytljr"},
	{id: "6d4p7ddaxn6y1il", name: "Sample Bags", inputKind: "number", unitLabel: "each", workRecordType: "up9ndzdnyzv8snx", sortOrder: 50, relationID: "fkvq0vbis9m0zl3"},
	{id: "f4yup3botoz09tz", name: "Field Filters", inputKind: "number", unitLabel: "each", workRecordType: "up9ndzdnyzv8snx", sortOrder: 60, relationID: "vspu5hb7r2240mn"},
	{id: "50we91ti05mszpg", name: "Waterra pump", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 70, relationID: "q5ct4aeki6efby5"},
	{id: "ot64hhzr4cddmru", name: "Interface Meter", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 80, relationID: "ltpb33ba9zwevc7"},
	{id: "gr31f2iskxkqzh0", name: "Water Level Meter", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 90, relationID: "9wa2xm2hxsonkor"},
	{id: "nuk1z6cf89r8iom", name: "PID", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 100, relationID: "wju9awb05iy5tig"},
	{id: "jo6w4ablwlibt9t", name: "Water Quality Meter", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 110, relationID: "oitaubag9svooyk"},
	{id: "wtuvd7cy89hx8wu", name: "Peristaltic Pump", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 120, relationID: "jnqiasl9gng0y5m"},
	{id: "z3t2oyfwmn420f1", name: "Field Pen (Multi-param)", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 130, relationID: "2tmi13dja6n7otm"},
	{id: "k5tpb6lirpru96k", name: "Field Meter (Multi-param meter)", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 135, relationID: "1ia3x1m2mv0zk8j"},
	{id: "yblnpz4nljl5l44", name: "Electrofisher", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 140, relationID: "uiim3is8rx92f2n"},
	{id: "d9kefa7pfb7feef", name: "Turbidity Meter", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 150, relationID: "vlt4aaszndlnng8"},
	{id: "sj3owvspw4xc2pd", name: "Fyke Net", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 160, relationID: "9yjkg9ou524eawq"},
	{id: "pq3prleew4hia0y", name: "Seine Net", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 170, relationID: "urc9sul703dp44k"},
	{id: "6fm14hl4594hgne", name: "Camera", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 180, relationID: "8bnvr3rkc51bxmz"},
	{id: "ra7usaffjfwbrvd", name: "Other GPS", inputKind: "boolean", unitLabel: "", workRecordType: "up9ndzdnyzv8snx", sortOrder: 190, relationID: "hmjpen9mqgkqu8f"},
}

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("work_records_types")
		if err != nil {
			return err
		}

		// Seed the work_records_types collection
		for _, seed := range []struct {
			id        string
			name      string
			sortOrder int
		}{
			{id: "rqayp9yoz1qe6ui", name: "General", sortOrder: 10},
			{id: "4k4nxilem8k10cu", name: "Drilling & Field Services", sortOrder: 20},
			{id: "oask5u50nka0km8", name: "Survey", sortOrder: 30},
			{id: "up9ndzdnyzv8snx", name: "Environmental", sortOrder: 40},
		} {
			record := core.NewRecord(collection)
			record.Set("id", seed.id)
			record.Set("name", seed.name)
			record.Set("max_subjects", 1)
			record.Set("allow_lane_kms", false)
			record.Set("sort_order", seed.sortOrder)
			record.Set("active", true)

			if err := app.Save(record); err != nil {
				return err
			}
		}

		consumablesCollection, err := app.FindCollectionByNameOrId("work_records_consumables")
		if err != nil {
			return err
		}

		typeConsumablesCollection, err := app.FindCollectionByNameOrId("work_records_type_consumables")
		if err != nil {
			return err
		}

		// Seed the work_records_consumables and work_records_type_consumables collections.
		for _, seed := range workRecordConsumableSeeds {
			record := core.NewRecord(consumablesCollection)
			record.Set("id", seed.id)
			record.Set("name", seed.name)
			record.Set("input_kind", seed.inputKind)
			record.Set("unit_label", seed.unitLabel)
			record.Set("active", true)

			if err := app.Save(record); err != nil {
				return err
			}

			typeConsumableRecord := core.NewRecord(typeConsumablesCollection)
			typeConsumableRecord.Set("id", seed.relationID)
			typeConsumableRecord.Set("type", seed.workRecordType)
			typeConsumableRecord.Set("consumable", seed.id)
			typeConsumableRecord.Set("sort_order", seed.sortOrder)
			typeConsumableRecord.Set("active", true)

			if err := app.Save(typeConsumableRecord); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		for _, seed := range workRecordConsumableSeeds {
			typeConsumableRecord, err := app.FindRecordById("work_records_type_consumables", seed.relationID)
			if err != nil {
				return err
			}

			if err := app.Delete(typeConsumableRecord); err != nil {
				return err
			}

			record, err := app.FindRecordById("work_records_consumables", seed.id)
			if err != nil {
				return err
			}

			if err := app.Delete(record); err != nil {
				return err
			}
		}

		for _, id := range []string{
			"rqayp9yoz1qe6ui",
			"4k4nxilem8k10cu",
			"oask5u50nka0km8",
			"up9ndzdnyzv8snx",
		} {
			record, err := app.FindRecordById("work_records_types", id)
			if err != nil {
				return err
			}

			if err := app.Delete(record); err != nil {
				return err
			}
		}

		return nil
	})
}
