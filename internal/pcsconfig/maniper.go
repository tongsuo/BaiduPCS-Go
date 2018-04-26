package pcsconfig

const (
	opDelete = "delete"
	opSwitch = "switch"
	opGet    = "get"
)

func (c *PCSConfig) manipUser(op string, baiduBase *BaiduBase) (*Baidu, error) {
	// empty baiduBase
	if baiduBase == nil || (baiduBase.UID == 0 && baiduBase.Name == "") {
		return nil, ErrBaiduUserNotFound
	}
	if len(c.baiduUserList) == 0 {
		return nil, ErrNoSuchBaiduUser
	}

	for k, user := range c.baiduUserList {
		if user == nil {
			continue
		}

		switch {
		case baiduBase.UID != 0 && baiduBase.Name != "":
			if user.UID == baiduBase.UID && user.Name == baiduBase.Name {
				goto handle
			}
			continue
		case baiduBase.UID == 0 && baiduBase.Name != "":
			if user.Name == baiduBase.Name {
				goto handle
			}
			continue
		case baiduBase.UID != 0 && baiduBase.Name == "":
			if user.UID == baiduBase.UID {
				goto handle
			}
			continue
		default:
			continue
		}
		// unreachable zone

	handle:
		switch op {
		case opSwitch:
			c.setupNewUser(user)
		case opDelete:
			c.baiduUserList = append(c.baiduUserList[:k], c.baiduUserList[k+1:]...)

			// 修改 正在使用的 百度帐号
			// 如果要删除的帐号为当前登录的帐号, 则设置当前登录帐号为列表中第一个帐号
			if c.baiduActiveUID == user.UID {
				if len(c.baiduUserList) != 0 {
					c.setupNewUser(c.baiduUserList[0])
				} else {
					c.baiduActiveUID = 0
				}
			}
		case opGet:
			// do nothing
		default:
			// do nothing
		}
		return user, nil
	}

	return nil, ErrBaiduUserNotFound
}

//setupNewUser 从已有用户中, 设置新的当前登录用户
func (c *PCSConfig) setupNewUser(user *Baidu) {
	if user == nil {
		return
	}
	c.baiduActiveUID = user.UID
	c.activeUser = user
	c.pcs = user.BaiduPCS()
}

// SwitchUser 切换用户, 返回切换成功的用户
func (c *PCSConfig) SwitchUser(baiduBase *BaiduBase) (*Baidu, error) {
	return c.manipUser(opSwitch, baiduBase)
}

// DeleteUser 删除用户, 返回删除成功的用户
func (c *PCSConfig) DeleteUser(baiduBase *BaiduBase) (*Baidu, error) {
	return c.manipUser(opDelete, baiduBase)
}

// GetBaiduUser 获取百度用户信息
func (c *PCSConfig) GetBaiduUser(baidubase *BaiduBase) (*Baidu, error) {
	return c.manipUser(opGet, baidubase)
}

// CheckBaiduUserExist 检查百度用户是否存在于已登录列表
func (c *PCSConfig) CheckBaiduUserExist(baidubase *BaiduBase) bool {
	_, err := c.manipUser("", baidubase)
	return err == nil
}

// SetupUserByBDUSS 设置百度 bduss, ptoken, stoken 并保存
func (c *PCSConfig) SetupUserByBDUSS(bduss, ptoken, stoken string) (baidu *Baidu, err error) {
	b, err := NewUserInfoByBDUSS(bduss)
	if err != nil {
		return nil, err
	}

	c.DeleteUser(&BaiduBase{
		UID: b.UID,
	}) // 删除旧的信息

	b.PTOKEN = ptoken
	b.STOKEN = stoken

	c.baiduUserList = append(c.baiduUserList, b)

	// 自动切换用户
	c.setupNewUser(b)
	return b, nil
}

// SetAppID 设置app_id
func (c *PCSConfig) SetAppID(appID int) {
	c.appID = appID
}

// SetCacheSize 设置cache_size, 下载缓存
func (c *PCSConfig) SetCacheSize(cacheSize int) {
	if c.cacheSize > 1024 {
		c.cacheSize = cacheSize
	}
}

// SetMaxParallel 设置max_parallel, 下载最大并发量
func (c *PCSConfig) SetMaxParallel(maxParallel int) {
	if c.maxParallel > 0 {
		c.maxParallel = maxParallel
	}
}

// SetUserAgent 设置User-Agent
func (c *PCSConfig) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// SetSaveDir 设置下载保存路径
func (c *PCSConfig) SetSaveDir(saveDir string) {
	c.saveDir = saveDir
}

// SetEnableHTTPS 设置是否启用https
func (c *PCSConfig) SetEnableHTTPS(https bool) {
	c.enableHTTPS = https
}
