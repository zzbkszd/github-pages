package com.xinou.volunteer.manager;

import com.jfinal.weixin.sdk.api.ApiResult;
import com.jfinal.weixin.sdk.api.MediaApi;
import com.jfinal.weixin.sdk.msg.out.News;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 *
 * Created by shizhida on 16/1/13.
 */
public class SourceManager {

    /**
     * 根据Author获取新闻素材，暂时没有添加图片url
     * @param author
     * @return
     */
    public static List<News> getNews(String author){

        List<News> news_list = new ArrayList<>();
        ApiResult result = MediaApi.batchGetMaterial(MediaApi.MediaType.NEWS, 0, 20);
        List<Map<String,Object>> items = result.getList("item");
        for (Map<String, Object> item : items) {
            Map<String,Object> contents = (Map<String, Object>) item.get("content");
            List<Map> news_items = (List<Map>) contents.get("news_item");
            for (Map news_item : news_items) {
                if(author.equals(news_item.get("author"))){
                    news_list.add(new News(
                            (String)news_item.get("title"),
                            "",//简介
                            (String)news_item.get("thumb_url"),//图片url
                            (String)news_item.get("url")));
                }
            }
        }
        return news_list;
    }
}
