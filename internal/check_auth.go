package internal

const (
	RetSuccess      = 0
	RetNotValidCode = -10000
)

type CheckAuthRes struct {
	Ret           int    // 错误码
	Msg           string // 错误信息
	InternalToken string // 内部票据
}

// 内部票据
func IsInternalTokenValid() bool {
	return true
}

// 客户端是否携带有效的code
func IsCodeValid() bool {
	return true
}

func GetInternalToken() string {
	return "abcdefg"
}

func CheckAuth() CheckAuthRes {
	// 校验内部登录态是否正常，如果在有效期内，则直接返回
	if IsInternalTokenValid() {
		return CheckAuthRes{Ret: RetSuccess, InternalToken: "abcdefg"}
	}
	// 根据code换取票据，如果没有code，则提示错误
	if !IsCodeValid() {
		return CheckAuthRes{
			Ret: RetNotValidCode,
		}
	}

	// 生成内部登录态并返回
	return CheckAuthRes{
		Ret:           RetSuccess,
		InternalToken: GetInternalToken(),
	}
}
