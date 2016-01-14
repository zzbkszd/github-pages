package com.xinou.volunteer.controller;

import com.jfinal.core.Controller;
import com.jfinal.kit.PropKit;
import com.jfinal.weixin.sdk.api.*;
import com.jfinal.weixin.sdk.utils.JsonUtils;
import com.xinou.volunteer.utils.HttpsUtil;

import java.io.*;
import java.util.HashMap;
import java.util.Map;

/**
 * 测试及管理操作使用
 * Created by shizhida on 16/1/13.
 */
public class AdminController extends Controller {

    /**
     * 必须先设置ApiConfig才能调用微信API
     * @return
     */
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
     * 查看所有资源
     */
    public void getAllResource(){
        ApiConfigKit.setThreadLocalApiConfig(getApiConfig());

        ApiResult result = MediaApi.batchGetMaterial(MediaApi.MediaType.NEWS,0,20);
        InputStream in = MediaApi.getMaterial("AYk0hq2czp9yVI0q8SQsRMPBa6QlgnPcnLvCbIRF0A0");
        StringBuffer buffer = new StringBuffer();
        InputStreamReader inputStreamReader = null;
        try {
            inputStreamReader = new InputStreamReader(in, "utf-8");
            BufferedReader bufferedReader = new BufferedReader(inputStreamReader);

            String str = null;
            while ((str = bufferedReader.readLine()) != null) {
                buffer.append(str);
            }
            renderText(buffer.toString());
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        }

    }
}
