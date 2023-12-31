package api

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryadavDeqode/dq-vault/config"
	"github.com/ryadavDeqode/dq-vault/test/unit_test/mocks"
)

func TestPathSaveUserConsent(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mocks.NewMockStorage(ctrl)

	tErr := "test error"
	s.EXPECT().Get(context.Background(), config.StorageBasePath+"test").Return(&logical.StorageEntry{}, errors.New(tErr))
	s.EXPECT().List(context.Background(), config.StorageBasePath).Return([]string{"test"}, nil).AnyTimes()
	s.EXPECT().Put(context.Background(), gomock.Any()).Return(nil).AnyTimes()
	b := backend{}
	req := logical.Request{}

	req.Storage = s

	mpGet := MPatchGet("test")

	res, err := b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err.Error() != tErr {
		t.Error("expected test error, received - ", res, err)
	}

	s.EXPECT().Get(context.Background(), config.StorageBasePath+"test").Return(&logical.StorageEntry{}, nil)
	mpdj := MPatchDecodeJSON(nil)

	res, err = b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err != nil {
		t.Error(" error wasn't expected, received - ", err)
	} else {
		if res.Data["status"].(bool) {
			t.Error(" unexpected value of status,expected false, received - ", res)
		}
	}

	mpGet.Unpatch()
	mpGet = MPatchGet("MNEMONICS")
	mpjwt := MPatchVerifyJWTSignature(true, "")
	mpnc := MPatchNewClient(nil)
	mpgetps := MPatchGetPubSub("test", nil)

	s.EXPECT().Get(context.Background(), config.StorageBasePath+"MNEMONICS").Return(&logical.StorageEntry{}, nil)
	res, err = b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err != nil {
		t.Error(" error wasn't expected, received - ", err)
	} else {
		if !res.Data["status"].(bool) {
			t.Error(" unexpected value of status,expected true, received - ", res)
		}
	}

	mpStorageEntryJson := MPatchEntryJSON(errors.New(tErr))

	s.EXPECT().Get(context.Background(), config.StorageBasePath+"MNEMONICS").Return(&logical.StorageEntry{}, nil)
	res, err = b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err == nil {
		t.Error("expected error, received ", res, err)
	}

	fmt.Println("unPached ")

	// mpjwt.Unpatch()
	// mpnc.Unpatch()
	// mpgetps.Unpatch()

	// mpjwt = MPatchVerifyJWTSignature(true, "")
	// mpnc = MPatchNewClient(nil)
	// mpgetps = MPatchGetPubSub("test", nil)

	mpGet.Unpatch()
	mpStorageEntryJson.Unpatch()
	mpGet = MPatchGet("PRIVATE_KEY")

	s.EXPECT().Get(context.Background(), config.StorageBasePath+"PRIVATE_KEY").Return(&logical.StorageEntry{}, nil)
	res, err = b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err != nil {
		t.Error(" error wasn't expected, received - ", err)
	} else {
		if !res.Data["status"].(bool) {
			t.Error(" unexpected value of status,expected true, received - ", res)
		}
	}

	mpGet.Unpatch()
	mpGet = MPatchGet("test")

	s.EXPECT().Get(context.Background(), config.StorageBasePath+"test").Return(&logical.StorageEntry{}, nil).AnyTimes()

	res, err = b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err != nil {
		t.Error(" error wasn't expected, received - ", err)
	} else {
		if res.Data["status"].(bool) {
			t.Error(" unexpected value of status,expected false, received - ", res)
		}
	}

	mpdj.Unpatch()

	mpdj = MPatchDecodeJSON(errors.New("tErr"))
	res, err = b.pathSaveUserConsent(context.Background(), &req, &framework.FieldData{})

	if err == nil {
		t.Error("expected error, received ", res, err)
	}

	mpStorageEntryJson.Unpatch()
	mpGet.Unpatch()
	mpjwt.Unpatch()
	mpnc.Unpatch()
	mpgetps.Unpatch()
	mpdj.Unpatch()
}
