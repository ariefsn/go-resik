package helper

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilterOperator string

const (
	FoContains  FilterOperator = "$contains"
	FoStartWith FilterOperator = "$startWith"
	FoEndWith   FilterOperator = "$endWith"
	FoEq        FilterOperator = "$eq"
	FoNe        FilterOperator = "$ne"
	FoIn        FilterOperator = "$in"
	FoNin       FilterOperator = "$nin"
	FoGt        FilterOperator = "$gt"
	FoGte       FilterOperator = "$gte"
	FoLt        FilterOperator = "$lt"
	FoLte       FilterOperator = "$lte"
)

type MongoSortBy string

const (
	SortByAsc  MongoSortBy = "ASC"
	SortByDesc MongoSortBy = "DESC"
)

type MongoSort struct {
	SortField string
	SortBy    MongoSortBy
}

type MongoAggregate struct {
	Skip  *int64
	Limit *int64
	Sort  []MongoSort
	Match bson.D
}

func NewMongoAggregate() *MongoAggregate {
	a := new(MongoAggregate)

	return a
}

func (a *MongoAggregate) BuildPipe() mongo.Pipeline {
	pipe := mongo.Pipeline{}

	// build match
	if a.Match != nil && len(a.Match) > 0 {
		pipe = append(pipe, MongoMatch(a.Match))
	}

	// build sort
	if len(a.Sort) > 0 {
		pipe = append(pipe, bson.D{
			bson.E{
				Key:   "$sort",
				Value: MongoSorting(a.Sort...),
			},
		})
	}

	// build skip
	if a.Skip != nil && *a.Skip > -1 {
		pipe = append(pipe, MongoSkip(*a.Skip))
	}

	// build limit
	if a.Limit != nil && *a.Limit > -1 {
		if *a.Skip < 0 {
			pipe = append(pipe, MongoLimit(*a.Limit))
		} else {
			pipe = append(pipe, MongoLimit(*a.Limit))
		}
	}

	return pipe
}

func MongoPipe(aggregate MongoAggregate) mongo.Pipeline {
	return aggregate.BuildPipe()
}

func MongoMatch(filter bson.D) bson.D {
	return bson.D{
		bson.E{
			Key:   "$match",
			Value: filter,
		},
	}
}

func MongoSorting(sort ...MongoSort) bson.D {
	s := bson.D{}

	for _, v := range sort {
		sortBy := 1

		if v.SortBy == SortByDesc {
			sortBy = -1
		}

		s = append(s, bson.E{
			Key:   v.SortField,
			Value: sortBy,
		})
	}

	return s
}

func MongoSkip(skip int64) bson.D {
	return bson.D{
		bson.E{
			Key:   "$skip",
			Value: skip,
		},
	}
}

func MongoLimit(limit int64) bson.D {
	return bson.D{
		bson.E{
			Key:   "$limit",
			Value: limit,
		},
	}
}

func MongoFilter(operator FilterOperator, field string, value interface{}) bson.D {
	switch operator {
	case FoEq, FoNe, FoIn, FoNin, FoGt, FoGte, FoLt, FoLte:
		return bson.D{
			bson.E{
				Key: field,
				Value: bson.D{
					bson.E{
						Key:   string(operator),
						Value: value,
					},
				},
			},
		}
	case FoContains, FoStartWith, FoEndWith:
		var regexPattern string
		regexOpt := "i"

		switch operator {
		case FoContains:
			regexPattern = fmt.Sprintf(".*%s.*", value)
		case FoStartWith:
			regexPattern = fmt.Sprintf("^%s", value)
		case FoEndWith:
			regexPattern = fmt.Sprintf("%s$", value)
		}

		return bson.D{
			bson.E{
				Key: field,
				Value: primitive.Regex{
					Pattern: regexPattern,
					Options: regexOpt,
				},
			},
		}
	default:
		return nil
	}
}

func MongoFilterM(operator FilterOperator, field string, value interface{}) bson.M {
	m, _ := ToBsonM(MongoFilter(operator, field, value))

	return m
}

func BuildMongoOrders(orders string, separators ...string) []MongoSort {
	ordersSlice := []MongoSort{}

	separatorOrder := ","
	separatorOrderBy := "*"

	for k, v := range separators {
		switch k {
		case 0:
			separatorOrder = v
		case 1:
			separatorOrderBy = v
		}
	}

	orderSplit := strings.Split(orders, separatorOrder)

	if orders == "" || len(orderSplit) == 0 {
		return ordersSlice
	}

	for _, os := range orderSplit {
		sortSplit := strings.Split(os, separatorOrderBy)

		sortBy := SortByAsc

		if len(sortSplit) == 2 {

			if sortSplit[1] == "-1" || strings.ToLower(sortSplit[1]) == "desc" {
				sortBy = SortByDesc
			}

			ordersSlice = append(ordersSlice, MongoSort{
				SortField: sortSplit[0],
				SortBy:    sortBy,
			})
		}
	}

	return ordersSlice
}

type MongoLookupOptions struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
}

func MongoLookup(opt MongoLookupOptions) bson.D {
	return bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				bson.E{
					Key:   "from",
					Value: opt.From,
				},
				bson.E{
					Key:   "localField",
					Value: opt.LocalField,
				},
				bson.E{
					Key:   "foreignField",
					Value: opt.ForeignField,
				},
				bson.E{
					Key:   "as",
					Value: opt.As,
				},
			},
		},
	}
}

type MongoUnwindOptions struct {
	Path              string
	IncludeArrayIndex string
	PreserveEmpty     bool
}

func MongoUnwind(opt MongoUnwindOptions) bson.D {
	return bson.D{
		bson.E{
			Key: "$unwind",
			Value: bson.D{
				bson.E{
					Key:   "path",
					Value: opt.Path,
				},
				// bson.E{
				// 	Key:   "includeArrayIndex",
				// 	Value: opt.IncludeArrayIndex,
				// },
				bson.E{
					Key:   "preserveNullAndEmptyArrays",
					Value: opt.PreserveEmpty,
				},
			},
		},
	}
}

func MongoIn(field string, inValues ...interface{}) bson.D {
	return bson.D{
		bson.E{
			Key: field,
			Value: bson.D{
				bson.E{
					Key:   "$in",
					Value: inValues,
				},
			},
		},
	}
}

func MongoSet(field string, value interface{}) bson.D {
	return bson.D{
		bson.E{
			Key: "$set",
			Value: bson.D{
				bson.E{
					Key:   field,
					Value: value,
				},
			},
		},
	}
}

func MongoUnionWith(collection string, pipelines []bson.D) bson.D {
	return bson.D{
		bson.E{
			Key: "$unionWith",
			Value: bson.D{
				bson.E{
					Key:   "coll",
					Value: collection,
				},
				bson.E{
					Key:   "pipeline",
					Value: pipelines,
				},
			},
		},
	}
}

type MongoGraphLookupOptions struct {
	From             string
	StartWith        string
	ConnectFromField string
	ConnectToField   string
	DepthField       string
	As               string
}

func MongoGraphLookup(opt MongoGraphLookupOptions) bson.D {
	return bson.D{
		{
			Key: "$graphLookup",
			Value: bson.D{
				bson.E{
					Key:   "from",
					Value: opt.From,
				},
				bson.E{
					Key:   "startWith",
					Value: opt.StartWith,
				},
				bson.E{
					Key:   "connectFromField",
					Value: opt.ConnectFromField,
				},
				bson.E{
					Key:   "connectToField",
					Value: opt.ConnectToField,
				},
				bson.E{
					Key:   "depthField",
					Value: opt.DepthField,
				},
				bson.E{
					Key:   "as",
					Value: opt.As,
				},
			},
		},
	}
}

func MongoDateToString(field, format string) bson.D {
	return bson.D{
		{
			Key: "$dateToString",
			Value: bson.D{
				{
					Key:   "format",
					Value: format,
				},
				{
					Key:   "date",
					Value: field,
				},
			},
		},
	}
}
