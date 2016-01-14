package com.xinou.volunteer.controller;

import com.jfinal.core.Controller;
import com.xinou.volunteer.dao.Activities;
import com.xinou.volunteer.utils.IdUtil;

import java.util.List;
import java.util.Map;


/**
 *
 * 活动操作
 * Created by shizhida on 16/1/13.
 */
public class ActivitiesController extends Controller {


    public void launchActivities (){
        Activities activities = new Activities();
        activities.put("id", IdUtil.getTimeMillisSequence());
        activities.put("title",getPara("title"));
        activities.put("require",getPara("require"));
        activities.put("time",getPara("time"));
        activities.put("size",getPara("size"));
        activities.put("organizer",getPara("organizer"));
        activities.put("context",getPara("context"));
        activities.put("num",0);
        activities.put("activities_time",getPara("activities_time"));
        activities.put("type", getPara("type"));
        boolean b = activities.save();
        if (!b){
            render("/resources/launch_error.html");
        }else {
            render("/resources/launch_ok.html");
        }
    }
    public void launchjsp() {
        render("/resources/activities_launch.html");
    }


    public void activitiesList() {
        render("/resources/activities_type.html");
    }

    public void findAll() {
        List<Activities> list = Activities.dao.find("select * from activities order by id desc");
        renderJson(list);
    }
    public void findActivitiesById() {
        renderJson(Activities.dao.findById(getPara("id")));
    }

    public void findActivitiesByType() {
        List<Activities> list = Activities.dao.find("select * from activities where type=? order by id desc",getPara("type"));
        renderJson(list);
    }
    public void findAllList() {
        render("/resources/activities.html");
    }



}
