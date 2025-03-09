package models

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Hotel represents the structure of a hotel document
type Hotel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	WikiLink    string `json:"wikiLink"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Image       string `json:"image"`
	Coordinates string `json:"coordinates"`
	Website     string `json:"website"`
	Address     string `json:"address"`
}

func GetAll(db *mongo.Collection) ([]Hotel, error) {
	var hotels []Hotel

	ctx := context.TODO()

	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Couldn't get ALL, cursor Find", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &hotels); err != nil {
		log.Fatal("Couldn't get ALL", err)
		return nil, err
	}

	return hotels, nil
}

// GetByFilter searches for hotels by city name, then by province, and finally by coordinate distance
// FILTER FORMAT: "cityName, provinceName, coordinates"
/*func GetByFilter(filter string, db *mongo.Collection) ([]Hotel, error) {
	filters := strings.Split(filter, ",")
	if len(filters) != 3 {
		log.WithFields(log.Fields{
			"error":    "Invalid filter format",
			"expected": "cityName, provinceName, coordinates",
			"received": filter,
		}).Error("Invalid filter format")
		return nil, fmt.Errorf("invalid filter format")
	}

	cityName := strings.TrimSpace(filters[0])
	provinceName := strings.TrimSpace(filters[1])
	coordinates := strings.TrimSpace(filters[2])

	var hotels []Hotel
	ctx := context.TODO()

	// Step 1: Search by city name
	log.WithFields(log.Fields{
		"city": cityName,
	}).Info("Searching hotels by city")
	cursor, err := db.Find(ctx, bson.M{"city": cityName})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"city":  cityName,
		}).Error("Failed to search hotels by city")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &hotels); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"city":  cityName,
		}).Error("Failed to decode hotels by city")
		return nil, err
	}

	log.WithFields(log.Fields{
		"city":  cityName,
		"count": len(hotels),
	}).Debug("Hotels found by city")

	// If we have 3 or more results, return them
	if len(hotels) >= 3 {
		log.WithFields(log.Fields{
			"city":   cityName,
			"count":  len(hotels),
			"action": "Returning results",
		}).Info("Found sufficient hotels by city")
		return hotels, nil
	}

	// Step 2: If fewer than 3 results, search by province
	log.WithFields(log.Fields{
		"province": provinceName,
	}).Info("Searching hotels by province")
	cursor, err = db.Find(ctx, bson.M{"province": provinceName})
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"province": provinceName,
		}).Error("Failed to search hotels by province")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &hotels); err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"province": provinceName,
		}).Error("Failed to decode hotels by province")
		return nil, err
	}

	log.WithFields(log.Fields{
		"province": provinceName,
		"count":    len(hotels),
	}).Debug("Hotels found by province")

	// If we have 3 or more results, return them
	if len(hotels) >= 3 {
		log.WithFields(log.Fields{
			"province": provinceName,
			"count":    len(hotels),
			"action":   "Returning results",
		}).Info("Found sufficient hotels by province")
		return hotels, nil
	}

	// Step 3: If still fewer than 3 results, search by "close enough" coordinates
	log.WithFields(log.Fields{
		"coordinates": coordinates,
	}).Info("Searching hotels by coordinates")
	coords := strings.Split(coordinates, ";")
	if len(coords) != 2 {
		log.WithFields(log.Fields{
			"error":    "Invalid coordinates format",
			"expected": "latitude; longitude",
			"received": coordinates,
		}).Error("Invalid coordinates format")
		return nil, fmt.Errorf("invalid coordinates format")
	}

	latitude := parseFloat(strings.TrimSpace(coords[0]))
	longitude := parseFloat(strings.TrimSpace(coords[1]))

	// Define a range for "close enough" coordinates
	latitudeRange := 20.0  // ±20 degrees latitude
	longitudeRange := 20.0 // ±20 degrees longitude

	// Calculate the latitude and longitude bounds
	minLatitude := latitude - latitudeRange
	maxLatitude := latitude + latitudeRange
	minLongitude := longitude - longitudeRange
	maxLongitude := longitude + longitudeRange

	log.WithFields(log.Fields{
		"minLatitude":  minLatitude,
		"maxLatitude":  maxLatitude,
		"minLongitude": minLongitude,
		"maxLongitude": maxLongitude,
	}).Info("Querying hotels by coordinate range")

	// Perform a range-based query for "close enough" coordinates
	cursor, err = db.Find(ctx, bson.M{
		"$expr": bson.M{
			"$and": []bson.M{
				// Ensure the coordinates field is not empty and contains exactly two parts
				{"$eq": []interface{}{bson.M{"$size": bson.M{"$split": []string{"$coordinates", ";"}}}, 2}},
				// Ensure latitude and longitude are not empty
				{"$ne": []interface{}{bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 0}}, ""}},
				{"$ne": []interface{}{bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 1}}, ""}},
				// Parse latitude and longitude, and check bounds
				{"$gte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 0}}}}}, minLatitude}},
				{"$lte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 0}}}}}, maxLatitude}},
				{"$gte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 1}}}}}, minLongitude}},
				{"$lte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 1}}}}}, maxLongitude}},
			},
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"coordinates": coordinates,
		}).Error("Failed to search hotels by coordinates")
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &hotels); err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"coordinates": coordinates,
		}).Error("Failed to decode hotels by coordinates")
		return nil, err
	}

	log.WithFields(log.Fields{
		"coordinates": coordinates,
		"count":       len(hotels),
	}).Debug("Hotels found by coordinates")

	// Return the results, even if fewer than 3
	log.WithFields(log.Fields{
		"coordinates": coordinates,
		"count":       len(hotels),
		"action":      "Returning results",
	}).Info("Found hotels by coordinates")
	return hotels, nil
}*/

// GetOneByFilter searches for one hotel for a given filter string.
// FILTER FORMAT: "cityName, provinceName, coordinates"
func GetOneByFilter(filter string, db *mongo.Collection) (Hotel, error) {
	filters := strings.Split(filter, ",")
	if len(filters) != 3 {
		log.WithFields(log.Fields{
			"error":    "Invalid filter format",
			"expected": "cityName, provinceName, coordinates",
			"received": filter,
		}).Error("Invalid filter format")
		return Hotel{}, fmt.Errorf("invalid filter format")
	}

	cityName := strings.TrimSpace(filters[0])
	provinceName := strings.TrimSpace(filters[1])
	coordinates := strings.TrimSpace(filters[2])

	ctx := context.TODO()

	// Step 1: Search by city name
	log.WithFields(log.Fields{
		"city": cityName,
	}).Info("Searching hotels by city")
	cursor, err := db.Find(ctx, bson.M{"city": cityName})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"city":  cityName,
		}).Error("Failed to search hotels by city")
		return Hotel{}, err
	}
	defer cursor.Close(ctx)

	var hotels []Hotel
	if err = cursor.All(ctx, &hotels); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"city":  cityName,
		}).Error("Failed to decode hotels by city")
		return Hotel{}, err
	}

	if len(hotels) >= 3 {
		log.WithFields(log.Fields{
			"city":   cityName,
			"count":  len(hotels),
			"action": "Returning first hotel from city search",
		}).Info("Found sufficient hotels by city")
		return hotels[0], nil
	}

	// Step 2: Search by province
	log.WithFields(log.Fields{
		"province": provinceName,
	}).Info("Searching hotels by province")
	cursorProv, err := db.Find(ctx, bson.M{"province": provinceName})
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"province": provinceName,
		}).Error("Failed to search hotels by province")
		return Hotel{}, err
	}
	defer cursorProv.Close(ctx)

	var hotelsProv []Hotel
	if err = cursorProv.All(ctx, &hotelsProv); err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"province": provinceName,
		}).Error("Failed to decode hotels by province")
		return Hotel{}, err
	}

	if len(hotelsProv) >= 3 {
		log.WithFields(log.Fields{
			"province": provinceName,
			"count":    len(hotelsProv),
			"action":   "Returning first hotel from province search",
		}).Info("Found sufficient hotels by province")
		return hotelsProv[0], nil
	}

	// Step 3: Search by coordinates
	log.WithFields(log.Fields{
		"coordinates": coordinates,
	}).Info("Searching hotels by coordinates")
	coords := strings.Split(coordinates, ";")
	if len(coords) != 2 {
		log.WithFields(log.Fields{
			"error":    "Invalid coordinates format",
			"expected": "latitude; longitude",
			"received": coordinates,
		}).Error("Invalid coordinates format")
		return Hotel{}, fmt.Errorf("invalid coordinates format")
	}

	latitude := parseFloat(strings.TrimSpace(coords[0]))
	longitude := parseFloat(strings.TrimSpace(coords[1]))

	// Define a range for "close enough" coordinates
	latitudeRange := 20.0  // ±20 degrees latitude
	longitudeRange := 20.0 // ±20 degrees longitude

	minLatitude := latitude - latitudeRange
	maxLatitude := latitude + latitudeRange
	minLongitude := longitude - longitudeRange
	maxLongitude := longitude + longitudeRange

	log.WithFields(log.Fields{
		"minLatitude":  minLatitude,
		"maxLatitude":  maxLatitude,
		"minLongitude": minLongitude,
		"maxLongitude": maxLongitude,
	}).Info("Querying hotels by coordinate range")

	cursorCoord, err := db.Find(ctx, bson.M{
		"$expr": bson.M{
			"$and": []bson.M{
				{"$eq": []interface{}{bson.M{"$size": bson.M{"$split": []string{"$coordinates", ";"}}}, 2}},
				{"$ne": []interface{}{bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 0}}, ""}},
				{"$ne": []interface{}{bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 1}}, ""}},
				{"$gte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 0}}}}}, minLatitude}},
				{"$lte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 0}}}}}, maxLatitude}},
				{"$gte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 1}}}}}, minLongitude}},
				{"$lte": []interface{}{bson.M{"$toDouble": bson.M{"$trim": bson.M{"input": bson.M{"$arrayElemAt": []interface{}{bson.M{"$split": []string{"$coordinates", ";"}}, 1}}}}}, maxLongitude}},
			},
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"coordinates": coordinates,
		}).Error("Failed to search hotels by coordinates")
		return Hotel{}, err
	}
	defer cursorCoord.Close(ctx)

	var hotelsCoord []Hotel
	if err = cursorCoord.All(ctx, &hotelsCoord); err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"coordinates": coordinates,
		}).Error("Failed to decode hotels by coordinates")
		return Hotel{}, err
	}

	if len(hotelsCoord) > 0 {
		log.WithFields(log.Fields{
			"coordinates": coordinates,
			"count":       len(hotelsCoord),
			"action":      "Returning first hotel from coordinates search",
		}).Info("Found hotels by coordinates")
		return hotelsCoord[0], nil
	}

	return Hotel{}, fmt.Errorf("no hotel found for filter: %s", filter)
}

func parseFloat(s string) float64 {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"value": s,
		}).Error("Failed to parse coordinate to float")
	}
	return f
}
