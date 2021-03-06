// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: log_count_index.proto

package models

import (
	context "context"
	fmt "fmt"
	
	_ "github.com/infobloxopen/protoc-gen-gorm/options"
	math "math"

	gorm2 "github.com/infobloxopen/atlas-app-toolkit/gorm"
	errors1 "github.com/infobloxopen/protoc-gen-gorm/errors"
	gorm1 "github.com/jinzhu/gorm"
	field_mask1 "google.golang.org/genproto/protobuf/field_mask"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = fmt.Errorf
var _ = math.Inf

type LogCountIndexORM struct {
	LogIndex        uint64 `gorm:"primary_key"`
	TransactionHash string `gorm:"primary_key"`
}

// TableName overrides the default tablename generated by GORM
func (LogCountIndexORM) TableName() string {
	return "log_count_indices"
}

// ToORM runs the BeforeToORM hook if present, converts the fields of this
// object to ORM format, runs the AfterToORM hook, then returns the ORM object
func (m *LogCountIndex) ToORM(ctx context.Context) (LogCountIndexORM, error) {
	to := LogCountIndexORM{}
	var err error
	if prehook, ok := interface{}(m).(LogCountIndexWithBeforeToORM); ok {
		if err = prehook.BeforeToORM(ctx, &to); err != nil {
			return to, err
		}
	}
	to.TransactionHash = m.TransactionHash
	to.LogIndex = m.LogIndex
	if posthook, ok := interface{}(m).(LogCountIndexWithAfterToORM); ok {
		err = posthook.AfterToORM(ctx, &to)
	}
	return to, err
}

// ToPB runs the BeforeToPB hook if present, converts the fields of this
// object to PB format, runs the AfterToPB hook, then returns the PB object
func (m *LogCountIndexORM) ToPB(ctx context.Context) (LogCountIndex, error) {
	to := LogCountIndex{}
	var err error
	if prehook, ok := interface{}(m).(LogCountIndexWithBeforeToPB); ok {
		if err = prehook.BeforeToPB(ctx, &to); err != nil {
			return to, err
		}
	}
	to.TransactionHash = m.TransactionHash
	to.LogIndex = m.LogIndex
	if posthook, ok := interface{}(m).(LogCountIndexWithAfterToPB); ok {
		err = posthook.AfterToPB(ctx, &to)
	}
	return to, err
}

// The following are interfaces you can implement for special behavior during ORM/PB conversions
// of type LogCountIndex the arg will be the target, the caller the one being converted from

// LogCountIndexBeforeToORM called before default ToORM code
type LogCountIndexWithBeforeToORM interface {
	BeforeToORM(context.Context, *LogCountIndexORM) error
}

// LogCountIndexAfterToORM called after default ToORM code
type LogCountIndexWithAfterToORM interface {
	AfterToORM(context.Context, *LogCountIndexORM) error
}

// LogCountIndexBeforeToPB called before default ToPB code
type LogCountIndexWithBeforeToPB interface {
	BeforeToPB(context.Context, *LogCountIndex) error
}

// LogCountIndexAfterToPB called after default ToPB code
type LogCountIndexWithAfterToPB interface {
	AfterToPB(context.Context, *LogCountIndex) error
}

// DefaultCreateLogCountIndex executes a basic gorm create call
func DefaultCreateLogCountIndex(ctx context.Context, in *LogCountIndex, db *gorm1.DB) (*LogCountIndex, error) {
	if in == nil {
		return nil, errors1.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(LogCountIndexORMWithBeforeCreate_); ok {
		if db, err = hook.BeforeCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	if err = db.Create(&ormObj).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(LogCountIndexORMWithAfterCreate_); ok {
		if err = hook.AfterCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormObj.ToPB(ctx)
	return &pbResponse, err
}

type LogCountIndexORMWithBeforeCreate_ interface {
	BeforeCreate_(context.Context, *gorm1.DB) (*gorm1.DB, error)
}
type LogCountIndexORMWithAfterCreate_ interface {
	AfterCreate_(context.Context, *gorm1.DB) error
}

// DefaultApplyFieldMaskLogCountIndex patches an pbObject with patcher according to a field mask.
func DefaultApplyFieldMaskLogCountIndex(ctx context.Context, patchee *LogCountIndex, patcher *LogCountIndex, updateMask *field_mask1.FieldMask, prefix string, db *gorm1.DB) (*LogCountIndex, error) {
	if patcher == nil {
		return nil, nil
	} else if patchee == nil {
		return nil, errors1.NilArgumentError
	}
	var err error
	for _, f := range updateMask.Paths {
		if f == prefix+"TransactionHash" {
			patchee.TransactionHash = patcher.TransactionHash
			continue
		}
		if f == prefix+"LogIndex" {
			patchee.LogIndex = patcher.LogIndex
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	return patchee, nil
}

// DefaultListLogCountIndex executes a gorm list call
func DefaultListLogCountIndex(ctx context.Context, db *gorm1.DB) ([]*LogCountIndex, error) {
	in := LogCountIndex{}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(LogCountIndexORMWithBeforeListApplyQuery); ok {
		if db, err = hook.BeforeListApplyQuery(ctx, db); err != nil {
			return nil, err
		}
	}
	db, err = gorm2.ApplyCollectionOperators(ctx, db, &LogCountIndexORM{}, &LogCountIndex{}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(LogCountIndexORMWithBeforeListFind); ok {
		if db, err = hook.BeforeListFind(ctx, db); err != nil {
			return nil, err
		}
	}
	db = db.Where(&ormObj)
	db = db.Order("transaction_hash")
	ormResponse := []LogCountIndexORM{}
	if err := db.Find(&ormResponse).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(LogCountIndexORMWithAfterListFind); ok {
		if err = hook.AfterListFind(ctx, db, &ormResponse); err != nil {
			return nil, err
		}
	}
	pbResponse := []*LogCountIndex{}
	for _, responseEntry := range ormResponse {
		temp, err := responseEntry.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		pbResponse = append(pbResponse, &temp)
	}
	return pbResponse, nil
}

type LogCountIndexORMWithBeforeListApplyQuery interface {
	BeforeListApplyQuery(context.Context, *gorm1.DB) (*gorm1.DB, error)
}
type LogCountIndexORMWithBeforeListFind interface {
	BeforeListFind(context.Context, *gorm1.DB) (*gorm1.DB, error)
}
type LogCountIndexORMWithAfterListFind interface {
	AfterListFind(context.Context, *gorm1.DB, *[]LogCountIndexORM) error
}
