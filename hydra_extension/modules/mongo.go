package modules

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"reflect"
	"time"
	"unicode/utf8"
	"yym/hydra_extension/hydra"
)

var (
	mongoRegistry *bsoncodec.Registry
)

func init() {
	hydra.RegisterType("MongoCollection", &MongoCollection{})
	hydra.RegisterType("MongoClient", NewMongoClient)
	builder := bson.NewRegistryBuilder()
	tBytes := reflect.TypeOf([]byte(nil))
	builder.RegisterTypeMapEntry(bsontype.Binary, tBytes)
	builder.RegisterDecoder(tBytes, bsoncodec.ValueDecoderFunc(func(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
		if !val.CanSet() || val.Type() != tBytes {
			return bsoncodec.ValueDecoderError{Name: "BinaryDecodeValue", Types: []reflect.Type{tBytes}, Received: val}
		}

		if vr.Type() != bsontype.Binary {
			return fmt.Errorf("cannot decode %v into a Binary", vr.Type())
		}

		data, _, err := vr.ReadBinary()
		if err != nil {
			return err
		}

		val.SetBytes(data)
		return nil
	}))

	builder.RegisterEncoder(tBytes, bsoncodec.ValueEncoderFunc(func(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
		if !val.IsValid() || val.Type() != tBytes {
			return bsoncodec.ValueEncoderError{Name: "ByteSliceEncodeValue", Types: []reflect.Type{tBytes}, Received: val}
		}
		if val.IsNil() {
			return vw.WriteNull()
		}
		data := val.Interface().([]byte)
		if utf8.Valid(data) {
			return vw.WriteString(string(data))
		} else {
			return vw.WriteBinary(data)
		}

	}))
	mongoRegistry = builder.Build()

}

type MongoClient struct {
	client  *mongo.Client
	timeOut time.Duration
}

type MongoCollection struct {
	col     *mongo.Collection
	timeOut time.Duration
}

func NewMongoClient(uri string) (*MongoClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri), &options.ClientOptions{
		Registry: mongoRegistry,
	})
	if err != nil {
		return nil, err
	}
	return &MongoClient{client: client, timeOut: time.Second * 10}, nil
}

func (self *MongoClient) defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), self.timeOut)
	return ctx
}

func (self *MongoClient) SetTimeOut(timeout hydra.Duration) {
	self.timeOut = timeout.Value()
}

func (self *MongoClient) Connect() error {
	return self.client.Connect(self.defaultContext())
}

func (self *MongoClient) Collection(db, col string) *MongoCollection {
	return &MongoCollection{
		col:     self.client.Database(db).Collection(col),
		timeOut: self.timeOut,
	}
}

func (self *MongoClient) Ping() error {
	return self.client.Ping(self.defaultContext(), readpref.Primary())
}

func (self *MongoClient) ListDatabase() ([]string, error) {
	return self.client.ListDatabaseNames(self.defaultContext(), map[string]interface{}{})
}

func (self *MongoCollection) SetTimeout(timeout hydra.Duration) {
	self.timeOut = timeout.Value()
}

func (self *MongoCollection) defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), self.timeOut)
	return ctx
}

func (self *MongoCollection) FindOne(filter interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	ret := bson.M{}
	if filter == nil {
		filter = bson.M{}
	}
	err := self.col.FindOne(self.defaultContext(), filter, opts...).Decode(&ret)
	if err == nil {
		return ret, nil
	} else {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
}

func (self *MongoCollection) Find(filter interface{}, opts ...*options.FindOptions) (interface{}, error) {
	ret := []bson.M{}
	if filter == nil {
		filter = bson.M{}
	}

	cur, err := self.col.Find(self.defaultContext(), filter, opts...)
	if err == nil {
		err = cur.All(nil, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		if err == mongo.ErrNoDocuments {
			return ret, nil
		}
		return nil, err
	}
}

func (self *MongoCollection) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) (interface{}, error) {
	if filter == nil {
		filter = bson.M{}
	}
	ret := bson.M{}
	err := self.col.FindOneAndDelete(self.defaultContext(), filter, opts...).Decode(&ret)
	if err == nil {
		return ret, nil
	} else {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
}

func (self *MongoCollection) FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (interface{}, error) {
	if filter == nil {
		filter = bson.M{}
	}

	ret := bson.M{}
	err := self.col.FindOneAndUpdate(self.defaultContext(), filter, update, opts...).Decode(&ret)
	if err == nil {
		return ret, nil
	} else {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
}

func (self *MongoCollection) FindOneAndReplace(filter interface{}, update interface{}, opts ...*options.FindOneAndReplaceOptions) (interface{}, error) {
	if filter == nil {
		filter = bson.M{}
	}
	ret := bson.M{}
	err := self.col.FindOneAndReplace(self.defaultContext(), filter, update, opts...).Decode(&ret)
	if err == nil {
		return ret, nil
	} else {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
}

func (self *MongoCollection) InsertOne(doc map[string]interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	var err error
	if id, ok := doc["_id"]; ok {
		if bytes.Equal(id.([]byte), []byte("sequence()")) {
			doc["_id"], err = self.SequenceId()
			if err != nil {
				return nil, err
			}
		}
	}
	result, err := self.col.InsertOne(self.defaultContext(), doc, opts...)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (self *MongoCollection) Insert(docs []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error) {
	var err error
	for _, value := range docs {
		doc := value.(map[string]interface{})
		if id, ok := doc["_id"]; ok {
			if bytes.Equal(id.([]byte), []byte("sequence()")) {
				doc["_id"], err = self.SequenceId()
				if err != nil {
					return nil, err
				}
			}
		}
	}
	result, err := self.col.InsertMany(self.defaultContext(), docs, opts...)
	if err != nil {
		return nil, err
	}
	ret := make([]interface{}, len(result.InsertedIDs))
	for key, value := range result.InsertedIDs {
		ret[key] = value
	}
	return ret, nil
}

func (self *MongoCollection) SetIndexes(docs []mongo.IndexModel) error {
	_, err := self.col.Indexes().DropAll(self.defaultContext())
	if de, ok := err.(driver.Error); ok {
		if de.Code != 26{
			return err
		}
	}
	_, err = self.col.Indexes().CreateMany(self.defaultContext(), docs)
	if err != nil {
		return err
	}
	return nil
}

func (self *MongoCollection) GetIndexes() ([]bson.M, error) {
	indexes, err := self.col.Indexes().List(self.defaultContext())
	if err != nil {
		return nil, err
	}
	ret := []bson.M{}
	err = indexes.All(self.defaultContext(), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *MongoCollection) DeleteOne(filter interface{}) (int64, error) {
	result, err := self.col.DeleteOne(self.defaultContext(), filter)
	var count int64 = 0
	if result != nil {
		count = result.DeletedCount
	}
	return count, err
}

func (self *MongoCollection) Delete(filter interface{}) (int64, error) {
	result, err := self.col.DeleteMany(self.defaultContext(), filter)
	var count int64 = 0
	if result != nil {
		count = result.DeletedCount
	}
	return count, err
}

func (self *MongoCollection) Count(filter interface{}) (int64, error) {
	count, err := self.col.CountDocuments(self.defaultContext(), filter)
	return count, err
}

func (self *MongoCollection) SequenceId() (int64, error) {
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)
	opt.SetUpsert(true)
	result := self.col.Database().Collection("sequence").FindOneAndUpdate(self.defaultContext(), bson.M{"_id": self.col.Name()}, bson.M{
		"$set": bson.M{"_id": self.col.Name()},
		"$inc": bson.M{"seq": 1},
	}, opt)
	if result.Err() != nil {
		return 0, result.Err()
	}
	type Seq struct {
		Seq int64 `bson:"seq"`
	}
	seq := &Seq{}
	err := result.Decode(&seq)
	if err != nil {
		return 0, err
	}

	return seq.Seq, nil
}

func (self *MongoCollection) Aggregate(doc interface{}, opts ...*options.AggregateOptions) (interface{}, error) {
	ret := []bson.M{}
	cursor, err := self.col.Aggregate(self.defaultContext(), doc, opts...)
	if err == nil {
		err = cursor.All(nil, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		if err == mongo.ErrNoDocuments {
			return ret, nil
		}
		return nil, err
	}
}

func (self *MongoCollection) Drop() error {
	return self.col.Drop(self.defaultContext())
}

func (self *MongoCollection) Update(filter interface{}, doc interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return self.col.UpdateMany(self.defaultContext(), filter, doc, opts...)
}

func (self *MongoCollection) UpdateOne(filter interface{}, doc interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return self.col.UpdateOne(self.defaultContext(), filter, doc, opts...)
}
