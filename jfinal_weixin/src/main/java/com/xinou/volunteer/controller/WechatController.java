package com.xinou.volunteer.controller;

import com.jfinal.kit.PropKit;
import com.jfinal.weixin.sdk.api.AccessTokenApi;
import com.jfinal.weixin.sdk.api.ApiConfig;
import com.jfinal.weixin.sdk.api.MediaFile;
import com.jfinal.weixin.sdk.jfinal.MsgControllerAdapter;
import com.jfinal.weixin.sdk.msg.in.InTextMsg;
import com.jfinal.weixin.sdk.msg.in.event.InFollowEvent;
import com.jfinal.weixin.sdk.msg.in.event.InMenuEvent;
import com.jfinal.weixin.sdk.msg.out.News;
import com.jfinal.weixin.sdk.msg.out.OutCustomMsg;
import com.jfinal.weixin.sdk.msg.out.OutNewsMsg;
import com.jfinal.weixin.sdk.msg.out.OutTextMsg;
import com.jfinal.weixin.sdk.utils.HttpUtils;
import com.jfinal.weixin.sdk.utils.IOUtils;
import com.jfinal.weixin.sdk.utils.JsonUtils;
import com.xinou.volunteer.manager.SourceManager;
import com.xinou.volunteer.utils.HttpsUtil;

import java.io.InputStream;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 微信主要的API，可以重写更多的方法以应对不同的事件
 * Created by shizhida on 16/1/13.
 */
public class WechatController extends MsgControllerAdapter {
    /**
     * 加入关注的事件
     * @param inFollowEvent
     */
    @Override
    protected void processInFollowEvent(InFollowEvent inFollowEvent) {
        String e = inFollowEvent.getEvent();
        if(InFollowEvent.EVENT_INFOLLOW_SUBSCRIBE.equals(e)){
            render(new OutTextMsg().setContent("欢迎加入,感谢您对亦志愿服务的支持与关注！"));
        }
    }

    @Override
    public ApiConfig getApiConfig() {
        ApiConfig ac = new ApiConfig();

        // 配置微信 API 相关常量
        ac.setToken(PropKit.get("token"));
        ac.setAppId(PropKit.get("appId"));
        ac.setAppSecret(PropKit.get("appSecret"));

        /**
         *  是否对消息进行加密，对应于微信平台的消息加解密方式：
         *  1：true进行加密且必须配置 encodingAesKey
         *  2：false采用明文模式，同时也支持混合模式
         */
        ac.setEncryptMessage(PropKit.getBoolean("encryptMessage", false));
        ac.setEncodingAesKey(PropKit.get("encodingAesKey", "setting it in config file"));
        return ac;
    }

    /**
     * 接收文本消息的事件
     * @param inTextMsg
     */
    @Override
    protected void processInTextMsg(InTextMsg inTextMsg) {
        OutTextMsg outTextMsg = new OutTextMsg(inTextMsg);
        outTextMsg.setContent("接收到了文字消息");
        render(outTextMsg);
    }

    /**
     * 获取素材
     * @param type
     * @param offset
     * @param count
     * @return
     */
    public Map<String,Object> getResource(String type,int offset,int count){

        String at = AccessTokenApi.getAccessTokenStr();//PropKit.get("access_token");
        String url = PropKit.get("source_url");
        url = url.replace("ACCESS_TOKEN", at);
        if (null != at) {
            // 调用接口创建菜单

            Map<String, Object> params = new HashMap<>();
            params.put("type", type);
            params.put("offset", offset);
            params.put("count", count);

            String jsonParam = JsonUtils.toJson(params);

//            InputStream in = HttpUtils.download(url, jsonParam);
            Map<String,Object> jsonObject = HttpsUtil.httpRequest(url,"POST",jsonParam);
            System.out.print(jsonObject);
            return jsonObject;
        }
        return null;

    }

    /**
     * 自定义菜单的点击事件
     * @param inMenuEvent
     */
    @Override
    protected void processInMenuEvent(InMenuEvent inMenuEvent) {

        String eventKey = inMenuEvent.getEventKey();
        List<News> news;
        switch (eventKey){
            case "13":
                news = SourceManager.getNews("主题宣传");
                render(new OutNewsMsg(inMenuEvent).addNews(news));
                break;
            case "14":
                news = SourceManager.getNews("志愿者风采");
                render(new OutNewsMsg(inMenuEvent).addNews(news));
                break;
            case "21":
                news = SourceManager.getNews("课程纵览");
                render(new OutNewsMsg(inMenuEvent).addNews(news));
                break;
            case "22":
                news = SourceManager.getNews("师资团队");
                render(new OutNewsMsg(inMenuEvent).addNews(news));
                break;
            case "23":
                news = SourceManager.getNews("品牌宣传");
                render(new OutNewsMsg(inMenuEvent).addNews(news));
                break;
            case "35":
                render(new OutTextMsg(inMenuEvent).setContent("联系电话：010-67883216\n联系QQ：35380767"));
                break;
            default:
                render(new OutTextMsg(inMenuEvent).setContent(eventKey+" 被点击！"));
                break;
        }

    }


}
