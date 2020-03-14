package gowe

import (
	"context"
	"errors"
)

//模板消息
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html

//WxMTemplateMsgSend 发送模板消息
func WxMTemplateMsgSend(ctx context.Context, params map[string]interface{}) (*APIResult, error) {
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/message/template/send?access_token=" + wxMpConfig.getAccessToken(ctx)
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgSetIndustry 设置所属行业
func WxMpTemplateMsgSetIndustry(ctx context.Context, industryId1 string, industryId2 string) (*APIResult, error) {
	if len(industryId1) < 1 || len(industryId2) < 1 {
		return nil, errors.New("industry_id1或者industry_id2不能为空")
	}
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/template/api_set_industry?access_token=" + wxMpConfig.getAccessToken(ctx)

	params := make(map[string]interface{})
	params["industry_id1"] = industryId1
	params["industry_id2"] = industryId2
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgGetIndustry 获取帐号设置的行业信息
func WxMpTemplateMsgGetIndustry(ctx context.Context) (*APIResult, error) {
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/template/get_industry?access_token=" + wxMpConfig.getAccessToken(ctx)
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgGetTemplateId 获得模板ID  templateIdShort 模板库中模板的编号，有“TM**”和“OPENTMTM**”等形式
//从行业模板库选择模板到帐号后台,获得模板ID的过程可在微信公众平台后台完成.为方便第三方开发者,提供通过接口调用的方式来获取模板ID
func WxMpTemplateMsgGetTemplateId(ctx context.Context, templateIdShort string) (*APIResult, error) {
	if len(templateIdShort) < 1 {
		return nil, errors.New("templateIdShort不能为空")
	}
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/template/api_set_industry?access_token=" + wxMpConfig.getAccessToken(ctx)

	params := make(map[string]interface{})
	params["template_id_short"] = templateIdShort
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgGetAllTemplate 获取模板列表
func WxMpTemplateMsgGetAllTemplate(ctx context.Context) (*APIResult, error) {
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/template/get_all_private_template?access_token=" + wxMpConfig.getAccessToken(ctx)
	resultMap, errMap := httpGetResultMap(apiurl)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}

//WxMpTemplateMsgDelPrivateTemplate  删除模板
func WxMpTemplateMsgDelPrivateTemplate(ctx context.Context, templateId string) (*APIResult, error) {
	if len(templateId) < 1 {
		return nil, errors.New("templateId不能为空")
	}
	wxMpConfig, errWxMpConfig := getWxMpConfig(ctx)
	if errWxMpConfig != nil {
		return nil, errWxMpConfig
	}
	apiurl := WxmpApiUrl + "/cgi-bin/template/del_private_template?access_token=" + wxMpConfig.getAccessToken(ctx)

	params := make(map[string]interface{})
	params["template_id"] = templateId
	resultMap, errMap := httpPostResultMap(apiurl, params)
	if errMap != nil {
		return nil, errMap
	}
	apiResult := newAPIResult(resultMap)
	return &apiResult, nil
}
