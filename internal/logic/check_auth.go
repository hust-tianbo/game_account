package logic

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/hust-tianbo/game_account/client"

	"github.com/hust-tianbo/game_account/internal/model"
	"github.com/hust-tianbo/go_lib/log"
	"github.com/jinzhu/gorm"
)

const (
	RetSuccess       = 0
	RetNotValidCode  = -10000 // 微信code校验失败
	RetInternalError = -10001 // 内部异常
)

var TokenValidDuration = 5 * time.Minute // token有效时长5分钟

type CheckAuthReq struct {
	PersonID      string `json:"personid"`
	Code          string `json:"code"`           // 平台返回的code码
	InternalToken string `json:"internal_token"` // 如果已经有内部票据，则携带
}

type CheckAuthRsp struct {
	Ret           int    `json:"ret"`            // 错误码
	Msg           string `json:"msg"`            // 错误信息
	InternalToken string `json:"internal_token"` // 内部票据
	PersonID      string `json:"personid"`       // 内部id
}

// 内部票据
func IsInternalTokenValid(req *CheckAuthReq) bool {
	// 没有内部账号，则直接刷新票据
	if req.InternalToken == "" || req.PersonID == "" {
		log.Debugf("[IsInternalTokenValid]internal account empty:%+v", req)
		return false
	}

	var accountEle model.AccountToken
	dbRes := db.Table(model.AccountTokenTable).Where(
		&model.AccountToken{InternalToken: req.InternalToken, PersonID: req.PersonID}).First(&accountEle)

	// 如果没有查到则需要更新票据
	if dbRes.Error != nil {
		log.Errorf("[IsInternalTokenValid]query account failed:%+v,%+v", req, dbRes.Error)
		return false
	}

	// 如果token生成时间不在有效期内，则认为token过期，需要重新刷新票据
	if accountEle.InternalTokenTime.Add(TokenValidDuration).Before(time.Now()) {
		log.Debugf("[IsInternalTokenValid]internal token invalid:%+v,%+v", req, accountEle.InternalTokenTime)
		return false
	}
	return true
}

// 生成用户id
func GenePersonid() string {
	time := time.Now()
	randInt := rand.Intn(1000)
	return fmt.Sprintf("%+v%+03d", time.Unix(), randInt)
}

// 生成内部票据
func GeneInternalToken() string {
	result := make([]byte, 6)
	rand.Read(result)
	return hex.EncodeToString(result)
}

// 客户端是否携带有效的code
func IsCodeValid(req *CheckAuthReq) (string, string, string, bool) {
	if req.Code == "" {
		log.Errorf("[IsCodeValid]code is invalid:%+v", req)
		return "", "", "", false
	}

	session, sessionErr := client.CodeToSession(req.Code)
	if sessionErr != nil {
		log.Errorf("[IsCodeValid]CodeToSession failed:%+v,%+v", req, sessionErr)
		return "", "", "", false
	}

	log.Debugf("[IsCodeValid]code is valid:%+v,%+v", req, session)
	return session.Openid, session.SessionKey, session.Unionid, true
}

func CheckAuth(req CheckAuthReq) CheckAuthRsp {
	// 校验内部登录态是否正常，如果在有效期内，则直接返回
	if IsInternalTokenValid(&req) {
		return CheckAuthRsp{Ret: RetSuccess, InternalToken: req.InternalToken, PersonID: req.PersonID}
	}
	// 根据code换取票据，如果没有code，则提示错误
	openid, sessionKey, unionID, isValid := IsCodeValid(&req)
	if !isValid {
		return CheckAuthRsp{
			Ret: RetNotValidCode,
		}
	}

	// 生成新的内部票据
	internalToken := GeneInternalToken()
	nowTime := time.Now()
	var personId string

	var ele model.AccountToken
	var dbRes *gorm.DB
	dbRes = db.Table(model.AccountTokenTable).Where(&model.AccountToken{Openid: openid, SessionKey: sessionKey}).First(&ele)
	if dbRes.Error != nil && !dbRes.RecordNotFound() {
		log.Errorf("[CheckAuth]read table failed:%+v,%+v", req, dbRes.Error)
		return CheckAuthRsp{
			Ret: RetInternalError,
		}
	}
	if dbRes.RecordNotFound() {
		// 检查openid,sessionkey是否已经存在，不存在则新建用户
		personId = GenePersonid()
		dbRes = db.Table(model.AccountTokenTable).Create(&model.AccountToken{
			Openid:            openid,
			SessionKey:        sessionKey,
			InternalToken:     internalToken,
			InternalTokenTime: nowTime,
			CTime:             nowTime,
			MTime:             nowTime,
			PersonID:          personId,
			Unionid:           unionID,
		})
	} else {
		personId = ele.PersonID
		// 如果已经存在用户，则直接更新内部票据
		dbRes = db.Table(model.AccountTokenTable).
			Where(&model.AccountToken{Openid: openid, SessionKey: sessionKey}).Update(map[string]interface{}{
			"internal_token": internalToken, "internal_token_time": nowTime, "m_time": nowTime})
	}

	if dbRes.Error != nil || dbRes.RowsAffected != 1 {
		log.Errorf("[CheckAuth]update table failed:%+v,%+v", req, dbRes.Error)
		return CheckAuthRsp{
			Ret: RetInternalError,
		}
	}

	log.Debugf("[CheckAuth]update table success:%+v", req)
	// 生成内部登录态并返回
	return CheckAuthRsp{
		Ret:           RetSuccess,
		InternalToken: internalToken,
		PersonID:      personId,
	}
}
